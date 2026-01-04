package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"pocketbase-ruoyi/api/monitor"
	"pocketbase-ruoyi/api/system/menu"
	"pocketbase-ruoyi/tools"
	"strings"
	"time"

	// 使用简单的行解析以避免引入额外依赖

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

var collWhitelist, permWhitelist = loadRBACWhitelist()

// RegisterRBAC registers a router middleware that computes a permission
// identifier for collection record routes and attaches any custom checks.
// It mirrors the logic previously in main.go.
func RegisterRBAC(app *pocketbase.PocketBase) {

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.BindFunc(func(e *core.RequestEvent) error {
			// 如果是 GET 请求 且 路径不以 /api 开头，直接跳过鉴权和日志，返回下一处理器
			if e.Request != nil && e.Request.Method == "GET" {
				p := e.Request.URL.Path
				if !strings.HasPrefix(p, "/api") {
					return e.Next()
				}
			}

			if e.Auth != nil && e.Auth.IsSuperuser() {
				return e.Next()
			}
			// 统一执行 e.Next()，并在末尾统一记录日志和返回
			var blockErr error

			path := e.Request.URL.Path
			trimmed := strings.TrimSuffix(path, "/")
			parts := strings.Split(trimmed, "/")

			var collectionName, action, perm string

			// 根据7种情况组合权限标识：list/add/query/edit/remove/export/import
			// 支持的路由形态：
			// - GET    /api/collections/{collection}/records           -> list (列表)
			// - POST   /api/collections/{collection}/records           -> add  (新增)
			// - GET    /api/collections/{collection}/records/{id}      -> query(详情)
			// - PATCH  /api/collections/{collection}/records/{id}      -> edit (修改)
			// - DELETE /api/collections/{collection}/records/{id}      -> remove(删除)
			// - GET    /api/collections/{collection}/export            -> export(导出)
			// - POST   /api/collections/{collection}/import            -> import(导入)

			// 识别集合相关路由：records、export、import
			if len(parts) >= 5 && parts[1] == "api" && parts[2] == "collections" {
				collectionName = parts[3]
				method := e.Request.Method
				// /api/collections/{collection}/... (第五段为资源类型或记录标识)
				if len(parts) == 5 { // 无记录ID的情况
					res := parts[4]
					switch {
					case res == "records" && method == "GET":
						action = "query"
					case res == "records" && method == "POST":
						action = "add"
					case res == "export" && method == "GET":
						action = "export"
					case res == "import" && method == "POST":
						action = "import"
					}
				} else if len(parts) == 6 { // 记录ID相关操作 /records/{id}
					if parts[4] == "records" {
						switch method {
						case "GET":
							action = "query"
						case "PATCH":
							action = "edit"
						case "DELETE":
							action = "remove"
						}
					}
				}
			}

			// 解析集合真实名称；若不存在集合则不做 RBAC 限制
			if collectionName != "" {
				if coll, err := e.App.FindCachedCollectionByNameOrId(collectionName); err == nil && coll != nil {
					collectionName = coll.Name
				} else {
					collectionName = ""
				}
			}

			if collectionName != "" && action != "" {
				perm = collectionName + ":" + action
			}

			// 非超级管理员/应用管理员才进行权限判断
			if !(IsSuperuser(e) || IsAdminByApp(e)) {
				// 白名单集合直接放行（不设置 blockErr）
				if collectionName != "" {
					if _, ok := collWhitelist[collectionName]; !ok {
						// 校验用户权限
						userID := ""
						if e.Auth != nil {
							userID = e.Auth.Id
						}
						// 使用公共函数执行基于 userID 的权限检测（collectionName 的提取仍在调用处完成）
						if err := EnsureUserHasPermission(e, userID, perm); err != nil {
							blockErr = err
						}
					}
				}
			}

			operParam := extractOperParam(e, 2048)
			// 统一调用下一处理器并记录结果
			start := time.Now()
			var nextErr error
			if blockErr == nil {
				nextErr = e.Next()
			} else {
				nextErr = blockErr
			}

			if perm == "" && path != "" {
				perm = path
			}
			// 异步/尽力而为地写入操作日志（失败不影响主流程）
			_ = monitor.RecordOperLog(e, monitor.OperLogInput{
				Title:         buildOperTitle(collectionName, action, path),
				BusinessType:  mapBusinessType(action),
				OperatorType:  "1", // 后台用户
				Status:        mapStatus(nextErr),
				Method:        perm,
				RequestMethod: e.Request.Method,
				OperParam:     operParam,
				ErrorMsg:      errorMsg(nextErr, 800),
				CostTime:      time.Since(start).Milliseconds(),
			})

			return nextErr
		})
		return se.Next()
	})
}

func mapBusinessType(action string) string {
	// 业务类型编码：0=其它 1=新增 2=修改 3=删除 4=查列表 5=查详情 6=接口调用 10=导出 11=导入
	switch strings.ToLower(action) {
	case "add":
		return "1"
	case "edit":
		return "2"
	case "remove":
		return "3"
	case "list":
		return "4"
	case "query":
		return "5"
	case "import":
		return "11"
	case "export":
		return "10"
	default:
		return "6"
	}
}

func mapStatus(err error) string {
	if err == nil {
		return "0" // 成功
	}
	return "1" // 失败
}

func buildOperTitle(collectionName, action, path string) string {
	if collectionName != "" && action != "" {
		return collectionName + " 表进行 " + actionName(action) + " 操作"
	}
	if path != "" {
		return "访问接口 " + path
	}
	return "访问"
}

func actionName(action string) string {
	switch action {
	case "list":
		return "列表查询"
	case "add":
		return "新增"
	case "edit":
		return "修改"
	case "remove":
		return "删除"
	case "query":
		return "详情查询"
	default:
		return action
	}
}

func extractOperParam(e *core.RequestEvent, max int) string {
	if e == nil || e.Request == nil {
		return ""
	}

	// 解析 query 参数为字典（仅取首个值），而不是原始 a=1&b=2 字符串
	queryVals := e.Request.URL.Query()
	queryMap := make(map[string]string, len(queryVals))

	// 读取 Body 原文
	bodyStr := ""
	var bodyJSON any // 若可解析为 JSON，则保存结构化对象
	_, raw, err := tools.ParseBody[any](e.Request)
	if err == nil && len(raw) > 0 {
		bodyStr = string(raw)
		// 优先尝试解析为 JSON；成功则使用解析后的对象，不再做字符串截断
		var parsed any
		if jsonErr := json.Unmarshal(raw, &parsed); jsonErr == nil {
			bodyJSON = parsed
		}
	}

	// 计算每字段可用最大长度（简单平分给 query 中的 value 与整体 body）
	var per int
	if max > 0 {
		per = max / 2
		if per < 1 {
			per = max
		}
	}

	for k, vs := range queryVals {
		if len(vs) == 0 {
			continue
		}
		v := vs[0]
		if per > 0 && len(v) > per { // 对单个值截断，保持键完整
			v = v[:per]
		}
		queryMap[k] = v
	}

	// 仅在未成功解析 JSON 时对原始字符串截断
	if bodyJSON == nil && per > 0 && len(bodyStr) > per {
		bodyStr = bodyStr[:per]
	}

	// 若解析成功，则 body 字段为结构化 JSON；否则为可能被截断的字符串
	var bodyField any
	if bodyJSON != nil {
		bodyField = bodyJSON
	} else {
		bodyField = bodyStr
	}

	obj := map[string]any{
		"query": queryMap,
		"body":  bodyField,
	}
	b, err := json.Marshal(obj)
	if err != nil {
		// 回退：使用原始原 query 字符串 + body 拼接
		fallbackQuery := e.Request.URL.RawQuery
		combined := fallbackQuery
		if bodyStr != "" {
			if combined != "" {
				combined += "&"
			}
			combined += bodyStr
		}
		if max > 0 && len(combined) > max {
			return combined[:max]
		}
		return combined
	}

	// 不对整体 JSON 再截断，保持合法性；调用方可按需限制长度
	return string(b)
}

func errorMsg(err error, max int) string {
	if err == nil {
		return ""
	}
	s := err.Error()
	if len(s) > max {
		return s[:max]
	}
	return s
}

// loadRBACWhitelist 从 config/rbac_whitelist.yml 加载白名单的 collection 列表和单一权限白名单。
// 返回两个 map 便于快速查验。文件不存在或解析错误时返回空 map（即不影响默认行为）。
func loadRBACWhitelist() (map[string]struct{}, map[string]struct{}) {
	collMap := make(map[string]struct{})
	permMap := make(map[string]struct{})

	// 默认相对工作目录的 config 文件
	cfgPath := filepath.Join("config", "rbac_whitelist.yml")

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		// 文件不存在或无法读取 -> 返回空白名单
		return collMap, permMap
	}

	// 简单行解析 yaml 格式：
	// collectionList:
	//   - collection_a
	// permissionList:
	//   - system:tenant:list
	lines := strings.Split(string(data), "\n")
	section := ""
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}

		if strings.HasPrefix(t, "collectionList:") {
			section = "collection"
			continue
		}
		if strings.HasPrefix(t, "permissionList:") {
			section = "permission"
			continue
		}

		// 已进入某个列表部分，查找以 '-' 开头的项
		if section == "collection" && strings.HasPrefix(t, "-") {
			v := strings.TrimSpace(strings.TrimPrefix(t, "-"))
			v = strings.Trim(v, "\"'")
			if v != "" {
				collMap[v] = struct{}{}
			}
			continue
		}
		if section == "permission" && strings.HasPrefix(t, "-") {
			v := strings.TrimSpace(strings.TrimPrefix(t, "-"))
			v = strings.Trim(v, "\"'")
			if v != "" {
				permMap[v] = struct{}{}
			}
			continue
		}

		// 如果遇到新的键（以冒号结束），退出当前列表解析
		if strings.HasSuffix(t, ":") {
			section = ""
			continue
		}
	}

	return collMap, permMap
}

// EnsureUserHasPermission 封装用户权限检测逻辑：
// - 若用户拥有通配权限 "*:*:*" 则放行
// - 若目标 perm 为空则视为允许
// - 若用户包含 perm 则放行
// - 若 perm 在白名单中也放行
// 否则返回一个权限不足的错误
func EnsureUserHasPermission(e *core.RequestEvent, userID string, perm string) error {
	permissions := menu.GetAllPermissionsByUser(e, userID)

	for _, p := range permissions {
		if p == "*:*:*" {
			return nil
		}
	}

	if perm == "" {
		return nil
	}

	for _, p := range permissions {
		if p == perm {
			return nil
		}
	}

	if _, ok := permWhitelist[perm]; ok {
		return nil
	}

	return apis.NewBadRequestError("权限不足", nil)
}

// RBAC is a middleware that checks if the authenticated user has the specified permission.
// It relies on EnsureUserHasPermission for the actual check so other packages can reuse it.
func RBAC(permission string) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		// 获取 userID（尽量一致的提取逻辑）
		userID := ""
		if e.Auth != nil {
			userID = e.Auth.Id
		}

		if err := EnsureUserHasPermission(e, userID, permission); err != nil {
			return err
		}
		return e.Next()
	}
}
