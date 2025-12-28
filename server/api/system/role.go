package system

import (
	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// Role 角色结构体
type Role struct {
	TenantID          string `db:"tenant_id" json:"tenant_id"`
	ID                string `db:"id" json:"id"`
	RoleKey           string `db:"role_key" json:"role_key"`
	RoleSort          int    `db:"role_sort" json:"role_sort"`
	DataScope         string `db:"data_scope" json:"data_scope"`
	MenuCheckStrictly bool   `db:"menu_check_strictly" json:"menu_check_strictly"`
	DeptCheckStrictly bool   `db:"dept_check_strictly" json:"dept_check_strictly"`
	Status            string `db:"status" json:"status"`     // select values: "0","1"
	DelFlag           string `db:"del_flag" json:"del_flag"` // select values: "0","1"
	CreateDept        string `db:"create_dept" json:"create_dept"`
	CreateBy          string `db:"create_by" json:"create_by"`
	CreateTime        string `db:"create_time" json:"create_time"` // RFC3339 timestamp string
	UpdateBy          string `db:"update_by" json:"update_by"`
	UpdateTime        string `db:"update_time" json:"update_time"` // RFC3339 timestamp string
	Remark            string `db:"remark" json:"remark"`
}

// GetAllRolesKeyByUser 获取用户的所有角色标识
func GetAllRolesKeyByUser(e *core.RequestEvent, userID string) []string {
	rows := GetAllRolesByUser(e, userID)

	var roleKeys []string
	for _, row := range rows {
		roleKeys = append(roleKeys, row.RoleKey)
	}
	return roleKeys
}

// GetAllRolesByUser 获取用户的所有角色标识
func GetAllRolesByUser(e *core.RequestEvent, userID string) []Role {
	rows := []Role{}

	err := e.App.DB().Select("role.*").From("role").
		InnerJoin("user_role as ur", dbx.NewExp("ur.role = role.id")).
		Where(dbx.HashExp{"ur.user": userID}).
		OrderBy("role.data_scope ASC").
		All(&rows)

	if err != nil {
		return []Role{}
	}

	return rows
}

// RegisterSystemRole 注册 /api/system/role 相关接口
func RegisterSystemRole(app *pocketbase.PocketBase) {
	app.OnRecordCreateRequest("role").BindFunc(syncRole)
	app.OnRecordAfterCreateSuccess("role").BindFunc(syncRoleAfter)

	app.OnRecordUpdateRequest("role").BindFunc(syncRole)
}

type syncRoleReq struct {
	MenuIds []string `json:"menu_ids" form:"menu_ids"`
	DeptIds []string `json:"dept_ids" form:"dept_ids"`
}

var tempRoleMenus = map[string][]string{}

func syncRole(e *core.RecordRequestEvent) error {
	payload := &syncRoleReq{}
	e.BindBody(payload)

	if e.Request.Header.Get("X-Menu") == "true" {
		tools.CacheIdsForCreate(e, "X-Menu", tempRoleMenus, payload.MenuIds)
	}
	tools.ReplaceJoinTableForUpdate(
		e,
		"X-Menu",
		"role_menu",
		"role={:role}",
		dbx.Params{"role": e.Record.Id},
		payload.MenuIds,
		func(nr *core.Record, roleId string, menuID string) {
			nr.Set("role", roleId)
			nr.Set("menu", menuID)
		},
	)

	if e.Record.Id == "" {
		e.Record.Set("data_scope", "1")
	}

	tools.ReplaceJoinTableForUpdate(
		e,
		"X-Dept",
		"role_dept",
		"role={:role}",
		dbx.Params{"role": e.Record.Id},
		payload.DeptIds,
		func(nr *core.Record, roleId string, deptID string) {
			nr.Set("role", roleId)
			nr.Set("dept", deptID)
		},
	)

	return e.Next()
}

func syncRoleAfter(e *core.RecordEvent) error {
	tools.ProcessAfterCreateTempIds(e, tempRoleMenus, "role_menu", func(nr, parent *core.Record, menuID string) {
		nr.Set("role", parent.Id)
		nr.Set("menu", menuID)
	})

	return e.Next()
}
