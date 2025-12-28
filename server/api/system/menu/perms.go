package menu

import (
	"sort"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// GetAllPermissionsByUser 获取指定用户的所有权限标识
func GetAllPermissionsByUser(e *core.RequestEvent, userID string) []string {
	var rows []Menu

	q := e.App.DB().Select("m.perms").From("menu as m").
		InnerJoin("role_menu as rm", dbx.NewExp("rm.menu = m.id")).
		InnerJoin("user_role as ur", dbx.NewExp("ur.role = rm.role")).
		Where(dbx.HashExp{"ur.user": userID})

	if err := q.All(&rows); err != nil {
		// 如果查询失败，返回空切片
		return []string{}
	}

	return extractPerms(rows)
}

// GetAllPermissionsByRole 根据角色对应的菜单列表获取所有权限标识
func GetAllPermissionsByRole(e *core.RequestEvent, roleID string) []string {
	var rows []Menu

	q := e.App.DB().Select("m.perms").From("menu as m").
		InnerJoin("role_menu as rm", dbx.NewExp("rm.menu = m.id")).
		Where(dbx.HashExp{"rm.role": roleID})

	if err := q.All(&rows); err != nil {
		return []string{}
	}

	return extractPerms(rows)
}

// extractPerms 将查询到的 Menu 列表中的 perms 字段拆分、去重并返回排序后的权限切片
func extractPerms(rows []Menu) []string {
	permSet := make(map[string]struct{})
	for _, r := range rows {
		p := strings.TrimSpace(r.Perms)
		if p == "" {
			continue
		}
		parts := strings.Split(p, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			permSet[part] = struct{}{}
		}
	}

	out := make([]string, 0, len(permSet))
	for k := range permSet {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
