package auth

import (
	// Use simple line parsing to avoid introducing extra dependencies
	"fmt"
	"os"
	"path/filepath"
	"pocketbase-ruoyi/api/system"
	"pocketbase-ruoyi/tools"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterDataScope registers a router middleware that computes a permission
// identifier for collection record routes and attaches any custom checks.
// It mirrors the logic previously in main.go.
func RegisterDataScope(app *pocketbase.PocketBase) {

	// Load the data_scope whitelist (return an empty whitelist if the file doesn't exist)
	collWhitelist := loadDataScopeWhitelist()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.BindFunc(func(e *core.RequestEvent) error {

			if e.Auth == nil || e.Auth.IsSuperuser() {
				return e.Next()
			}

			collectionName := e.Request.PathValue("collection")
			method := e.Request.Method

			if !(method == "GET" && collectionName != "") {
				return e.Next()
			}

			collection, err := e.App.FindCachedCollectionByNameOrId(collectionName)
			if err != nil || collection == nil {
				return e.Next()
			}

			collectionName = collection.Name

			// If the collection is in the data_scope whitelist, skip appending the data_scope filter
			if collectionName != "" {
				if _, ok := collWhitelist[collectionName]; ok {
					return e.Next()
				}
			}

			tenantIDField := collection.Fields.GetByName("tenant_id")
			createDeptField := collection.Fields.GetByName("create_dept")
			createByField := collection.Fields.GetByName("create_by")

			userTenantID := tools.GetUserTenant(e)
			userDeptID := e.Auth.GetString("dept_id")

			query := e.Request.URL.Query()

			if tenantIDField != nil {
				if userTenantID == "" {
					return fmt.Errorf("User tenant information is missing; cannot access tenant data")
				}

				oldFilter := query.Get("filter")
				tenantFilter := fmt.Sprintf("%s=\"%s\"", tenantIDField.GetName(), userTenantID)
				// Preserve existing query params; only update the filter value.
				if oldFilter == "" {
					query.Set("filter", tenantFilter)
				} else {
					query.Set("filter", fmt.Sprintf("(%s) && %s", oldFilter, tenantFilter))
				}
				e.Request.URL.RawQuery = query.Encode()
			}

			if !IsSuperuserByApp(e) {
				return e.Next()
			}

			roles := system.GetAllRolesByUser(e, e.Auth.Id)

			filters := []string{}

			// Split each role's data permission logic into a separate function for reuse (e.g., dataScope == "6")
			for _, role := range roles {
				clause, stop := filterClauseForRole(e, e.App, role.ID, role.DataScope, getFieldName(createDeptField), getFieldName(createByField), userDeptID)
				if stop {
					// All data access: clear and stop
					filters = []string{}
					break
				}
				if clause != "" {
					filters = append(filters, clause)
				}
			}

			if len(filters) > 0 {
				if userDeptID == "" {
					return fmt.Errorf("User department information is missing; cannot access department data")
				}

				oldFilter := query.Get("filter")
				dsFilter := strings.Join(func() []string { return wrapEach(filters) }(), " || ")
				createDeptFilter := fmt.Sprintf("(%s)", dsFilter)
				// Merge with existing filter while preserving other query parameters.
				if oldFilter == "" {
					query.Set("filter", createDeptFilter)
				} else {
					query.Set("filter", fmt.Sprintf("(%s) && %s", oldFilter, createDeptFilter))
				}
				e.Request.URL.RawQuery = query.Encode()
			}

			return e.Next()
		})

		return se.Next()
	})
}

// loadDataScopeWhitelist loads the list of collections that should skip data_scope from config/data_scope_whitelist.yml
// Returns a map for quick lookup. If the file is missing or parsing fails, returns an empty map (doesn't affect default behavior).
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

// --- helpers for data scope filtering ---

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

// filterClauseForRole returns a filter clause based on the role's dataScope (with or without parentheses), and whether processing should stop ("1" means full access)
func filterClauseForRole(e *core.RequestEvent, app core.App, roleID, dataScope, createDeptFieldName, createByFieldName, userDeptID string) (string, bool) {
	switch dataScope {
	case "1":
		// All data access
		return "", true
	case "2":
		// Custom data access
		if createDeptFieldName != "" {
			parts := partsFromRoleDept(app, createDeptFieldName, roleID)
			if len(parts) > 0 {
				return "(" + strings.Join(parts, " || ") + ")", false
			}
		}
	case "3":
		// Department-only data access
		if createDeptFieldName != "" && userDeptID != "" {
			return fmt.Sprintf("%s=\"%s\"", createDeptFieldName, userDeptID), false
		}
	case "4":
		// Department and sub-departments data access
		if createDeptFieldName != "" {
			ids := getDeptIDsIncludingChildren(e, userDeptID)
			parts := partsFromIDs(createDeptFieldName, ids)
			if len(parts) > 0 {
				return "(" + strings.Join(parts, " || ") + ")", false
			}
		}
	case "5":
		// Owner-only data access
		if createByFieldName != "" {
			return fmt.Sprintf("%s=\"%s\"", createByFieldName, e.Auth.Id), false
		}
	case "6":
		// Combined: prefer merging custom access with department and sub-departments access
		parts := []string{}
		if createDeptFieldName != "" {
			// Custom
			parts = append(parts, partsFromRoleDept(app, createDeptFieldName, roleID)...)
			// Department and sub-departments
			ids := getDeptIDsIncludingChildren(e, userDeptID)
			parts = append(parts, partsFromIDs(createDeptFieldName, ids)...)
		}
		if len(parts) > 0 {
			return "(" + strings.Join(parts, " || ") + ")", false
		}
	}
	return "", false
}

// partsFromRoleDept reads the dept list associated with a role from the role_dept table and constructs parts in the form field="id"
func partsFromRoleDept(app core.App, fieldName, roleID string) []string {
	if fieldName == "" || roleID == "" {
		return nil
	}
	depts, _ := app.FindRecordsByFilter("role_dept", "role={:roleID}", "", 999, 0, dbx.Params{"roleID": roleID})
	parts := make([]string, 0, len(depts))
	for _, d := range depts {
		deptID := d.GetString("dept")
		parts = append(parts, fmt.Sprintf("%s=\"%s\"", fieldName, deptID))
	}
	return parts
}

// partsFromIDs constructs parts in the form field="id" based on a list of ids
func partsFromIDs(fieldName string, ids []string) []string {
	if fieldName == "" || len(ids) == 0 {
		return nil
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%s=\"%s\"", fieldName, id))
	}
	return parts
}

// getDeptIDsIncludingChildren returns the list of ids including the user's department and its child departments
func getDeptIDsIncludingChildren(e *core.RequestEvent, userDeptID string) []string {
	ids := []string{}
	if userDeptID == "" {
		return ids
	}
	ids = append(ids, userDeptID)
	depts := system.GetChildrensByDeptID(e, userDeptID)
	for _, d := range depts {
		ids = append(ids, d.ID)
	}
	return ids
}

// wrapEach adds a pair of parentheses around each filter element to facilitate joining
func wrapEach(filters []string) []string {
	out := make([]string, len(filters))
	for i, f := range filters {
		out[i] = "(" + f + ")"
	}
	return out
}
