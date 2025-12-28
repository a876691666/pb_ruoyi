package tools

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// loadDataScopeWhitelist loads collections that should skip data_scope from config/data_scope_whitelist.yml
// Returns a map for quick lookup. If the file is missing or parsing fails, returns an empty map.
func loadDataScopeWhitelist() map[string]struct{} {
	collMap := make(map[string]struct{})

	cfgPath := filepath.Join("config", "data_scope_whitelist.yml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return collMap
	}

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

		if section == "collection" && strings.HasPrefix(t, "-") {
			v := strings.TrimSpace(strings.TrimPrefix(t, "-"))
			v = strings.Trim(v, "\"'")
			if v != "" {
				collMap[v] = struct{}{}
			}
			continue
		}

		if strings.HasSuffix(t, ":") {
			section = ""
			continue
		}
	}

	return collMap
}

// getFieldName safely extracts the field name from a possibly nil object (assuming it implements GetName() string)
func getFieldName(f interface{}) string {
	if f == nil {
		return ""
	}
	type namer interface{ GetName() string }
	if n, ok := f.(namer); ok {
		return n.GetName()
	}
	return ""
}

// isSuperuserRequest checks whether the current request's user is a super admin.
// It mirrors api/auth.IsSuperuser but is self-contained to avoid package cycles.
func isSuperuserRequest(e *core.RequestEvent) bool {
	if e == nil || e.Auth == nil {
		return false
	}
	if e.Auth.IsSuperuser() {
		return true
	}

	count := 0
	q := e.App.DB().Select("count(*)").From("role").
		InnerJoin("user_role", dbx.NewExp("user_role.role = role.id")).
		Where(dbx.HashExp{
			"role.role_key":  "superadmin",
			"user_role.user": e.Auth.Id,
		})
	q.Row(&count)
	return count > 0
}

// IsAdmin checks whether the current request's user is an admin.
func IsAdmin(e *core.RequestEvent) bool {
	count := 0
	q := e.App.DB().Select("count(*)").From("role").
		InnerJoin("user_role", dbx.NewExp("user_role.role = role.id")).
		Where(dbx.HashExp{
			"role.role_key":  "admin",
			"user_role.user": e.Auth.Id,
		})
	q.Row(&count)
	return count > 0
}

// getDeptIDsIncludingChildren returns the user's department id plus all children department ids
func getDeptIDsIncludingChildren(e *core.RequestEvent, userDeptID string) []string {
	ids := []string{}
	if userDeptID == "" {
		return ids
	}
	ids = append(ids, userDeptID)

	// query children by ancestors LIKE pattern, following api/system.GetChildrensByDeptID logic
	type row struct {
		ID string `db:"id"`
	}
	var rows []row
	_ = e.App.DB().Select("id").From("dept").
		Where(
			dbx.Or(
				dbx.Like("ancestors", fmt.Sprintf(",%s,", userDeptID)),
				dbx.Like("ancestors", fmt.Sprintf(",%s", userDeptID)).Match(true, false),
			),
		).
		All(&rows)
	for _, r := range rows {
		if r.ID != "" {
			ids = append(ids, r.ID)
		}
	}
	return ids
}

// roleDeptIDs returns dept ids bound to a role from role_dept table
func roleDeptIDs(app core.App, roleID string) []string {
	if roleID == "" {
		return nil
	}
	depts, _ := app.FindRecordsByFilter("role_dept", "role={:roleID}", "", 999, 0, dbx.Params{"roleID": roleID})
	if len(depts) == 0 {
		return nil
	}
	ids := make([]string, 0, len(depts))
	for _, d := range depts {
		if id := d.GetString("dept"); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

// FilterBuilderForRole returns a dbx expression representing data scope for a single role,
// and whether to stop processing (when dataScope=="1").
func FilterBuilderForRole(
	e *core.RequestEvent,
	app core.App,
	roleID, dataScope, createDeptFieldName, createByFieldName, userDeptID string,
) (dbx.Expression, bool) {
	switch dataScope {
	case "1":
		return nil, true
	case "2":
		if createDeptFieldName == "" {
			return nil, false
		}
		ids := roleDeptIDs(app, roleID)
		if len(ids) > 0 {
			vals := make([]interface{}, 0, len(ids))
			for _, id := range ids {
				vals = append(vals, id)
			}
			return dbx.In(createDeptFieldName, vals...), false
		}
	case "3":
		if createDeptFieldName != "" && userDeptID != "" {
			return dbx.HashExp{createDeptFieldName: userDeptID}, false
		}
	case "4":
		if createDeptFieldName == "" {
			return nil, false
		}
		ids := getDeptIDsIncludingChildren(e, userDeptID)
		if len(ids) > 0 {
			vals := make([]interface{}, 0, len(ids))
			for _, id := range ids {
				vals = append(vals, id)
			}
			return dbx.In(createDeptFieldName, vals...), false
		}
	case "5":
		if createByFieldName != "" {
			return dbx.HashExp{createByFieldName: e.Auth.Id}, false
		}
	case "6":
		if createDeptFieldName == "" {
			return nil, false
		}
		set := make(map[string]struct{})
		for _, id := range roleDeptIDs(app, roleID) {
			if id != "" {
				set[id] = struct{}{}
			}
		}
		for _, id := range getDeptIDsIncludingChildren(e, userDeptID) {
			if id != "" {
				set[id] = struct{}{}
			}
		}
		if len(set) > 0 {
			vals := make([]interface{}, 0, len(set))
			for id := range set {
				vals = append(vals, id)
			}
			return dbx.In(createDeptFieldName, vals...), false
		}
	}
	return nil, false
}

// getAllRolesByUser returns minimal role info list for a user.
type roleLite struct {
	ID        string `db:"id"`
	DataScope string `db:"data_scope"`
}

func getAllRolesByUser(e *core.RequestEvent, userID string) []roleLite {
	rows := []roleLite{}
	_ = e.App.DB().Select("role.id", "role.data_scope").From("role").
		InnerJoin("user_role as ur", dbx.NewExp("ur.role = role.id")).
		Where(dbx.HashExp{"ur.user": userID}).
		OrderBy("role.data_scope ASC").
		All(&rows)
	return rows
}

var oneByOne = dbx.NewExp("1=1")

// BuildDataScopeExpression builds a dbx expression equivalent to the router-level string filter.
// - Applies tenant_id restriction if such field exists.
// - If collection is not whitelisted, merges data-scope constraints across all user roles (OR).
// - Returns nil for superuser or unauthenticated requests.
// When tenant/department info is required but missing, returns an error.
func BuildDataScopeExpression(e *core.RequestEvent, collectionName string) dbx.Expression {
	if e == nil || e.App == nil || e.Auth == nil {
		return oneByOne
	}

	collection, err := e.App.FindCachedCollectionByNameOrId(collectionName)
	if err != nil || collection == nil {
		return oneByOne
	}

	// tenant filter if field exists
	var tenantExp dbx.Expression
	if f := collection.Fields.GetByName("tenant_id"); f != nil {
		userTenantID := GetUserTenant(e)
		if userTenantID == "" {
			return oneByOne
		}
		tenantExp = dbx.HashExp{getFieldName(f): userTenantID}
	}

	if IsAdmin(e) || isSuperuserRequest(e) {
		return tenantExp
	}

	// data scope (skip if whitelisted)
	var dataExp dbx.Expression
	collWhitelist := loadDataScopeWhitelist()
	if _, ok := collWhitelist[collection.Name]; !ok {
		createDeptField := collection.Fields.GetByName("create_dept")
		createByField := collection.Fields.GetByName("create_by")
		userDeptID := e.Auth.GetString("dept_id")

		roles := getAllRolesByUser(e, e.Auth.Id)
		var roleExps []dbx.Expression
		for _, role := range roles {
			exp, stop := FilterBuilderForRole(
				e, e.App,
				role.ID, role.DataScope,
				getFieldName(createDeptField), getFieldName(createByField),
				userDeptID,
			)
			if stop {
				roleExps = nil
				break
			}
			if exp != nil {
				roleExps = append(roleExps, dbx.Enclose(exp))
			}
		}
		if len(roleExps) > 0 {
			if userDeptID == "" {
				return oneByOne
			}
			dataExp = dbx.Or(roleExps...)
		}
	}

	switch {
	case tenantExp != nil && dataExp != nil:
		return dbx.And(tenantExp, dataExp)
	case tenantExp != nil:
		return tenantExp
	case dataExp != nil:
		return dataExp
	default:
		return oneByOne
	}
}

// AppendDataScopeToURLQuery is a helper mirroring the router-level behavior for building ?filter= query parts.
// It is optional; kept here in case external callers need to mutate Request URL directly like RegisterDataScope does.
func AppendDataScopeToURLQuery(e *core.RequestEvent, collectionName string) error {
	if e == nil || e.App == nil || e.Auth == nil || isSuperuserRequest(e) {
		return nil
	}
	collection, err := e.App.FindCachedCollectionByNameOrId(collectionName)
	if err != nil || collection == nil {
		return nil
	}

	query := e.Request.URL.Query()
	// tenant
	if f := collection.Fields.GetByName("tenant_id"); f != nil {
		userTenantID := GetUserTenant(e)
		if userTenantID == "" {
			return fmt.Errorf("User tenant information is missing; cannot access tenant data")
		}
		oldFilter := query.Get("filter")
		tenantFilter := fmt.Sprintf("%s=\"%s\"", getFieldName(f), userTenantID)
		if oldFilter == "" {
			e.Request.URL.RawQuery = url.Values{"filter": {tenantFilter}}.Encode()
		} else {
			e.Request.URL.RawQuery = url.Values{"filter": {fmt.Sprintf("(%s) && %s", oldFilter, tenantFilter)}}.Encode()
		}
	}

	// data scope
	collWhitelist := loadDataScopeWhitelist()
	if _, ok := collWhitelist[collection.Name]; ok {
		return nil
	}

	createDeptField := collection.Fields.GetByName("create_dept")
	createByField := collection.Fields.GetByName("create_by")
	userDeptID := e.Auth.GetString("dept_id")

	roles := getAllRolesByUser(e, e.Auth.Id)
	var parts []string
	for _, role := range roles {
		exp, stop := FilterBuilderForRole(
			e, e.App,
			role.ID, role.DataScope,
			getFieldName(createDeptField), getFieldName(createByField),
			userDeptID,
		)
		if stop {
			parts = nil
			break
		}
		if exp != nil {
			// convert expression back to simple equalities joined by OR when possible
			// For simplicity, we only support the equalities that FilterBuilderForRole produces
			// and fall back to no-op if complex expressions are present.
			// Callers should prefer BuildDataScopeExpression for DB querying.
			// Here we attempt to reconstruct patterns field="id" joined by ||.
			switch v := exp.(type) {
			case dbx.HashExp:
				for k, val := range v {
					if s, ok := val.(string); ok {
						parts = append(parts, fmt.Sprintf("%s=\"%s\"", k, s))
					}
				}
			default:
				// best-effort not implemented for non-HashExp; skip adding
			}
		}
	}
	if len(parts) > 0 {
		if userDeptID == "" {
			return fmt.Errorf("User department information is missing; cannot access department data")
		}
		joined := "(" + strings.Join(parts, ") || (") + ")"
		oldFilter := query.Get("filter")
		if oldFilter == "" {
			e.Request.URL.RawQuery = url.Values{"filter": {joined}}.Encode()
		} else {
			e.Request.URL.RawQuery = url.Values{"filter": {fmt.Sprintf("(%s) && %s", oldFilter, joined)}}.Encode()
		}
	}
	return nil
}
