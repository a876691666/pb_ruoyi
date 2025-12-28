package auth

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// IsSuperuser 检查当前请求的用户是否为超级管理员
func IsSuperuser(e *core.RequestEvent) bool {

	if e.Auth == nil {
		return false
	}

	if e.Auth.IsSuperuser() {
		return true
	}

	return IsSuperuserByApp(e)
}

// IsSuperuserByApp 检查指定应用中的用户是否为超级管理员
func IsSuperuserByApp(e *core.RequestEvent) bool {

	count := 0
	query := e.App.DB().Select("count(*)").From("role").
		InnerJoin("user_role", dbx.NewExp("user_role.role = role.id")).
		Where(dbx.HashExp{
			"role.role_key":  "superadmin",
			"user_role.user": e.Auth.Id,
		})

	query.Row(&count)

	return count > 0
}

// IsAdminByApp 检查指定应用中的用户是否为管理员
func IsAdminByApp(e *core.RequestEvent) bool {

	if e.Auth == nil {
		return false
	}

	count := 0
	query := e.App.DB().Select("count(*)").From("role").
		InnerJoin("user_role", dbx.NewExp("user_role.role = role.id")).
		Where(dbx.HashExp{
			"role.role_key":  "admin",
			"user_role.user": e.Auth.Id,
		})

	query.Row(&count)

	return count > 0
}
