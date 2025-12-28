package tenant

import (
	"math/rand"
	"pocketbase-ruoyi/tools"
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterTenant 注册
func RegisterTenant(app *pocketbase.PocketBase) {
	app.OnRecordCreateRequest().BindFunc(autoTenantID)

	// 当创建 tenant 时，补全 tenant_id，并在创建成功后完成角色/部门/用户/字典/配置的初始化
	app.OnRecordCreateRequest("tenant").BindFunc(ensureTenantID)
	app.OnRecordAfterCreateSuccess("tenant").BindFunc(afterTenantCreated)

	// 当更新租户套餐时，联动更新绑定该套餐的租户管理员角色的菜单（role_menu）
	app.OnRecordAfterUpdateSuccess("tenant_package").BindFunc(afterTenantPackageUpdated)

	// 当删除 tenant 时，做联动清理
	// 1) 执行期钩子：可阻止删除默认租户，并预清理强依赖的关联表
	app.OnRecordDeleteExecute("tenant").BindFunc(beforeTenantDelete)
	// 2) 删除成功后：清理该租户下的角色/部门/用户/字典/配置等数据
	app.OnRecordAfterDeleteSuccess("tenant").BindFunc(afterTenantDeleted)
}

func autoTenantID(e *core.RecordRequestEvent) error {
	collectionNameOrID := e.Request.PathValue("collection")
	collection, _ := e.App.FindCollectionByNameOrId(collectionNameOrID)
	collection.Fields.FieldNames()

	isTenantField := false
	for _, field := range collection.Fields {
		if field.GetName() == "tenant_id" {
			isTenantField = true
			break
		}
	}

	if !isTenantField {
		return e.Next()
	}

	tenantID := e.Auth.Get("tenant_id")

	if tenantID != nil {
		e.Record.Set("tenant_id", tenantID)
	}

	return e.Next()
}

type tenantCreateRequest struct {
	AccountCount    int    `json:"account_count"`
	CompanyName     string `json:"company_name"`
	ContactPhone    string `json:"contact_phone"`
	ContactUserName string `json:"contact_user_name"`
	ExpireTime      string `json:"expire_time"`
	PackageID       string `json:"package_id"`
	Password        string `json:"password"`
	Username        string `json:"username"`
}

// 默认租户编号（用于克隆字典和配置）
const defaultTenantID = "000000"

// ensureTenantID 在创建租户时确保 id 存在且唯一
func ensureTenantID(e *core.RecordRequestEvent) error {
	if e.Record == nil {
		return e.Next()
	}
	// 若已有 id 则跳过
	if tid := strings.TrimSpace(e.Record.GetString("id")); tid != "" {
		return e.Next()
	}
	req, _, _ := tools.ParseBody[tenantCreateRequest](e.RequestEvent.Request)
	// 检测user_name重名
	record, _ := e.App.FindFirstRecordByFilter("users", "user_name={:uname}", dbx.Params{"uname": req.Username})
	if record != nil {
		return apis.NewBadRequestError("租户管理员用户名已存在，请更换后重试", nil)
	}

	// 读取现有所有 id
	type row struct {
		TenantID string `db:"id"`
	}
	var rows []row
	_ = e.App.DB().Select("id").From("tenant").All(&rows)
	exists := map[string]struct{}{}
	for _, r := range rows {
		exists[r.TenantID] = struct{}{}
	}

	id := generateTenantID(exists)

	// 生成不重复的 6 位数字租户号（首位不为 0）
	e.Record.Set("id", id)

	e.App.Store().Set("tenant_id_generated_for_"+id, req)

	return e.Next()
}

// afterTenantCreated 在租户创建成功后，初始化角色、部门、管理员用户、字典与配置
func afterTenantCreated(e *core.RecordEvent) error {
	if e.Record == nil {
		return e.Next()
	}

	tenantID := e.Record.GetString("id")
	if tenantID == "" {
		return e.Next()
	}

	// 1) 创建租户管理员角色（基于套餐的菜单）
	roleID := ""
	if rid, err := createTenantAdminRole(e, tenantID); err == nil {
		roleID = rid
	}

	// 2) 创建部门：公司名作为部门名称，父级为 0
	deptID := ""
	if did, err := createRootDeptForTenant(e, tenantID, guessCompanyName(e)); err == nil {
		deptID = did
	}

	// 3) 角色与部门关联
	if roleID != "" && deptID != "" {
		_ = createJoinRecord(e, "role_dept", func(nr *core.Record) {
			nr.Set("role", roleID)
			nr.Set("dept", deptID)
		})
	}

	// 4) 创建管理员用户，并设为部门负责人（若存在 leader 字段）
	userID := ""
	if uid, err := createAdminUserForTenant(e, tenantID, deptID); err == nil {
		userID = uid
		// 部门负责人
		if deptID != "" && hasField(e.App, "dept", "leader") {
			if deptRec, err := e.App.FindRecordById("dept", deptID); err == nil && deptRec != nil {
				deptRec.Set("leader", userID)
				_ = e.App.Save(deptRec)
			}
		}
	}

	// 5) 用户-角色关联
	if userID != "" && roleID != "" {
		_ = createJoinRecord(e, "user_role", func(nr *core.Record) {
			nr.Set("user", userID)
			nr.Set("role", roleID)
		})
	}

	// 6) 克隆默认租户字典与配置
	_ = cloneByTenant(e, "dict_type", defaultTenantID, tenantID)
	_ = cloneByTenant(e, "dict_data", defaultTenantID, tenantID)
	_ = cloneByTenant(e, "config", defaultTenantID, tenantID)

	return e.Next()
}
func getCacheReq(e *core.RecordEvent, tenantID string) (tenantCreateRequest, error) {
	record := e.App.Store().Get("tenant_id_generated_for_" + tenantID)
	if record == nil {
		return tenantCreateRequest{}, apis.NewBadRequestError("未获取到缓存记录", nil)
	}

	switch m := record.(type) {
	case tenantCreateRequest:
		return m, nil
	case *tenantCreateRequest:
		return *m, nil
	default:
		return tenantCreateRequest{}, apis.NewBadRequestError("缓存记录类型错误", nil)
	}
}

// createTenantAdminRole 基于套餐创建租户管理员角色与角色菜单
func createTenantAdminRole(e *core.RecordEvent, tenantID string) (string, error) {
	record, err := getCacheReq(e, tenantID)
	if err != nil {
		return "", err
	}
	// 读取套餐ID：优先 package_id，其次 package
	pkgID := strings.TrimSpace(record.PackageID)
	if pkgID == "" {
		// 无套餐信息时，仍然创建一个空角色
		return createRoleWithMenus(e, tenantID, []string{})
	}

	pkg, err := e.App.FindRecordById("tenant_package", pkgID)
	if err != nil || pkg == nil {
		return createRoleWithMenus(e, tenantID, []string{})
	}

	menuIds := pkg.GetStringSlice("menu_ids")

	return createRoleWithMenus(e, tenantID, menuIds)
}

func createRoleWithMenus(e *core.RecordEvent, tenantID string, menuIds []string) (string, error) {
	coll, err := e.App.FindCollectionByNameOrId("role")
	if err != nil {
		return "", err
	}
	nr := core.NewRecord(coll)
	nr.Set("id", core.GenerateDefaultRandomId())
	nr.Set("tenant_id", tenantID)
	nr.Set("role_name", "租户管理员")
	nr.Set("role_key", "admin")
	nr.Set("role_sort", 1)
	nr.Set("status", "0") // 正常

	if err := e.App.Save(nr); err != nil {
		return "", err
	}
	roleID := nr.Id

	// 创建角色菜单关联
	for _, mid := range menuIds {
		if mid == "" {
			continue
		}
		_ = createJoinRecord(e, "role_menu", func(jr *core.Record) {
			jr.Set("role", roleID)
			jr.Set("menu", mid)
		})
	}
	return roleID, nil
}

// createRootDeptForTenant 创建根部门（父级为 0）
func createRootDeptForTenant(e *core.RecordEvent, tenantID, deptName string) (string, error) {
	coll, err := e.App.FindCollectionByNameOrId("dept")
	if err != nil {
		return "", err
	}
	nr := core.NewRecord(coll)
	nr.Set("id", core.GenerateDefaultRandomId())
	nr.Set("tenant_id", tenantID)
	nr.Set("dept_name", deptName)
	nr.Set("parent_id", "0")
	if err := e.App.Save(nr); err != nil {
		return "", err
	}
	// ancestors 将由 dept 钩子自动纠正
	return nr.Id, nil
}

// createAdminUserForTenant 创建系统用户（管理员）
func createAdminUserForTenant(e *core.RecordEvent, tenantID, deptID string) (string, error) {
	coll, err := e.App.FindCollectionByNameOrId("users")
	if err != nil {
		return "", err
	}
	record, err := getCacheReq(e, tenantID)
	if err != nil {
		return "", err
	}
	nr := core.NewRecord(coll)

	username := record.Username
	nickname := record.ContactUserName
	rawPwd := record.Password
	phonenumber := record.ContactPhone

	nr.Set("tenant_id", tenantID)
	nr.Set("user_name", username)
	nr.Set("user_type", "admin")
	nr.Set("nick_name", nickname)
	nr.Set("phonenumber", phonenumber)
	nr.Set("dept_id", deptID)
	nr.Set("status", "0") // 正常
	nr.SetPassword(rawPwd)

	if err := e.App.Save(nr); err != nil {
		return "", err
	}
	return nr.Id, nil
}

// cloneByTenant 从 sourceTenantID 克隆到 targetTenantID（忽略 id/created/updated）
func cloneByTenant(e *core.RecordEvent, collName, sourceTenantID, targetTenantID string) error {
	// 读取源记录
	records, err := e.App.FindRecordsByFilter(collName, "tenant_id={:tid}", "", 10000, 0, dbx.Params{"tid": sourceTenantID})
	if err != nil || len(records) == 0 {
		return nil
	}
	coll, err := e.App.FindCollectionByNameOrId(collName)
	if err != nil {
		return nil
	}
	// 需要复制的字段列表
	fields := coll.Fields.FieldNames()
	// 排除字段
	exclude := map[string]struct{}{"id": {}, "created": {}, "updated": {}}

	for _, src := range records {
		nr := core.NewRecord(coll)
		for _, f := range fields {
			if _, skip := exclude[f]; skip {
				continue
			}
			if f == "tenant_id" {
				nr.Set("tenant_id", targetTenantID)
				continue
			}
			// 复制原值
			nr.Set(f, src.Get(f))
		}
		_ = e.App.Save(nr)
	}
	return nil
}

// createJoinRecord 工具：创建关联表记录
func createJoinRecord(e *core.RecordEvent, collName string, set func(nr *core.Record)) error {
	coll, err := e.App.FindCollectionByNameOrId(collName)
	if err != nil {
		return err
	}
	nr := core.NewRecord(coll)
	set(nr)
	return e.App.Save(nr)
}

// generateTenantID 生成不重复的 6 位数字租户编号
func generateTenantID(exists map[string]struct{}) string {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		v := 100000 + rand.Intn(900000) // 100000-999999
		s := strconv.Itoa(v)
		if _, ok := exists[s]; !ok {
			return s
		}
	}
	// 兜底：时间戳后 6 位
	s := strconv.Itoa(int(time.Now().UnixNano()%900000 + 100000))
	return s
}

// 小工具函数
func hasField(app core.App, collName, fieldName string) bool {
	coll, err := app.FindCollectionByNameOrId(collName)
	if err != nil {
		return false
	}
	for _, f := range coll.Fields {
		if f.GetName() == fieldName {
			return true
		}
	}
	return false
}

func splitIDs(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// guessCompanyName 猜测公司名字段
func guessCompanyName(e *core.RecordEvent) string {
	n := firstNonEmpty(
		e.Record.GetString("company_name"),
		e.Record.GetString("tenant_name"),
		e.Record.GetString("name"),
	)
	if n == "" {
		return "默认部门"
	}
	return n
}

// ------------------------ 删除租户：联动处理 ------------------------

// beforeTenantDelete 在执行删除 tenant 前的保护与预清理
func beforeTenantDelete(e *core.RecordEvent) error {
	if e == nil || e.App == nil || e.Record == nil {
		return e.Next()
	}

	// 保护默认租户，阻止删除
	tid := strings.TrimSpace(e.Record.GetString("id"))
	if tid == defaultTenantID {
		return apis.NewBadRequestError("默认租户不允许删除", nil)
	}

	// 预清理：尽量删除依赖于租户下角色/部门/用户的关联表，避免外键/业务约束
	// 这里仅删除纯关联表记录，实体表在 afterTenantDeleted 中逐个删除以触发对应钩子
	// 根据 id 先查出相关角色、部门、用户ID，再批量清理关联表

	// 角色ID
	roleIDs := mustIDsByTenant(e, "role", tid)
	// 部门ID
	deptIDs := mustIDsByTenant(e, "dept", tid)
	// 用户ID
	userIDs := mustIDsByTenant(e, "users", tid)

	// 清理关联表：user_role, role_menu, role_dept
	if len(userIDs) > 0 {
		_, _ = e.App.DB().Delete("user_role", dbx.In("user", userIDs...)).Execute()
	}
	if len(roleIDs) > 0 {
		_, _ = e.App.DB().Delete("user_role", dbx.In("role", roleIDs...)).Execute()
		_, _ = e.App.DB().Delete("role_menu", dbx.In("role", roleIDs...)).Execute()
		_, _ = e.App.DB().Delete("role_dept", dbx.In("role", roleIDs...)).Execute()
	}
	if len(deptIDs) > 0 {
		_, _ = e.App.DB().Delete("role_dept", dbx.In("dept", deptIDs...)).Execute()
	}

	return e.Next()
}

// afterTenantDeleted 在删除 tenant 成功后，清理该租户下的实体数据
func afterTenantDeleted(e *core.RecordEvent) error {
	if e == nil || e.App == nil || e.Record == nil {
		return e.Next()
	}

	tid := strings.TrimSpace(e.Record.GetString("id"))
	if tid == "" {
		return e.Next()
	}

	// 删除顺序：实体表优先，关联表在 beforeTenantDelete 已预清理（这里仍然兜底）

	// 1) 删除用户
	_ = deleteRecordsByTenant(e, "users", tid)
	// 2) 删除角色
	_ = deleteRecordsByTenant(e, "role", tid)
	// 3) 删除部门（包含层级，逐条删以触发其自身钩子）
	_ = deleteRecordsByTenant(e, "dept", tid)
	// 4) 删除字典与配置
	_ = deleteRecordsByTenant(e, "dict_data", tid)
	_ = deleteRecordsByTenant(e, "dict_type", tid)
	_ = deleteRecordsByTenant(e, "config", tid)

	return e.Next()
}

// deleteRecordsByTenant 查找并删除指定 collection 下 id=tid 的所有记录
func deleteRecordsByTenant(e *core.RecordEvent, collName, tid string) error {
	_, err := e.App.DB().Delete(collName, dbx.HashExp{"tenant_id": tid}).Execute()
	if err != nil {
		return err
	}
	return nil
}

// mustIDsByTenant 返回 collection 下 id=tid 的记录ID列表（忽略错误）
func mustIDsByTenant(e *core.RecordEvent, collName, tid string) []any {
	recs, err := e.App.FindRecordsByFilter(collName, "tenant_id={:tid}", "", 10000, 0, dbx.Params{"tid": tid})
	if err != nil || len(recs) == 0 {
		return nil
	}
	ids := make([]any, 0, len(recs))
	for _, r := range recs {
		ids = append(ids, r.Id)
	}
	return ids
}

// ------------------------ 套餐更新：联动菜单同步 ------------------------

// afterTenantPackageUpdated 在更新 tenant_package 成功后：
// 同步所有关联该套餐的租户的管理员角色菜单（role_menu）为套餐的最新 menu_ids
func afterTenantPackageUpdated(e *core.RecordEvent) error {
	if e == nil || e.App == nil || e.Record == nil {
		return e.Next()
	}

	pkgID := strings.TrimSpace(e.Record.Id)
	if pkgID == "" {
		return e.Next()
	}

	// 读取套餐中的最新菜单ID集合
	menuIDs := e.Record.GetStringSlice("menu_ids")

	// 查找绑定此套餐的所有租户
	tenants, _ := e.App.FindRecordsByFilter("tenant", "package_id={:pid}", "", 10000, 0, dbx.Params{"pid": pkgID})

	if len(tenants) == 0 {
		return e.Next()
	}

	// 遍历租户：定位其管理员角色（role_key=admin），重建其 role_menu 关联
	for _, t := range tenants {
		tid := t.GetString("id")
		if strings.TrimSpace(tid) == "" {
			continue
		}

		// 查找该租户的管理员角色
		role, err := e.App.FindFirstRecordByFilter("role", "tenant_id={:tid} && role_key='admin'", dbx.Params{"tid": tid})
		if err != nil {
			continue
		}
		roleID := role.Id

		// 清空原有 role_menu 关联
		_, _ = e.App.DB().Delete("role_menu", dbx.HashExp{"role": roleID}).Execute()

		// 重新写入当前套餐的菜单集合
		for _, mid := range menuIDs {
			m := strings.TrimSpace(mid)
			if m == "" {
				continue
			}
			_ = createJoinRecord(e, "role_menu", func(jr *core.Record) {
				jr.Set("role", roleID)
				jr.Set("menu", m)
			})
		}
	}

	return e.Next()
}
