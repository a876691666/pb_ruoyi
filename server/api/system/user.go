package system

import (
	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// PasswordUpdatePayload 请求体结构
type PasswordUpdatePayload struct {
	OldPassword string `json:"oldPassword" form:"oldPassword"`
	NewPassword string `json:"newPassword" form:"newPassword"`
}

// PasswordResetPayload 重置密码请求体结构
type PasswordResetPayload struct {
	ID       string `json:"id" form:"id"`
	Password string `json:"password" form:"password"`
}

// SimpleUser 用户最小信息集
type SimpleUser struct {
	ID          string `db:"id" json:"id"`
	UserName    string `db:"user_name" json:"user_name"`
	NickName    string `db:"nick_name" json:"nick_name"`
	Email       string `db:"email" json:"email"`
	Phonenumber string `db:"phonenumber" json:"phonenumber"`
	DeptID      int64  `db:"dept_id" json:"dept_id"`
	Status      string `db:"status" json:"status"`
	DelFlag     string `db:"del_flag" json:"del_flag"`
	Avatar      string `db:"avatar" json:"avatar"`
}
type syncUserReq struct {
	RoleIds []string `json:"role_ids" form:"role_ids"`
	PostIds []string `json:"post_ids" form:"post_ids"`
}

// 暂存创建用户阶段的角色和岗位，待创建成功后再落库（参考 role.go 的实现）
var tempUserRoles = map[string][]string{}
var tempUserPosts = map[string][]string{}

func syncUser(e *core.RecordRequestEvent) error {
	payload := &syncUserReq{}
	e.BindBody(payload)

	if e.Request.Header.Get("X-Post") == "true" {
		tools.CacheIdsForCreate(e, "X-Post", tempUserPosts, payload.PostIds)
	}
	tools.ReplaceJoinTableForUpdate(
		e,
		"X-Post",
		"user_post",
		"user={:user}",
		dbx.Params{"user": e.Record.Id},
		payload.PostIds,
		func(nr *core.Record, userId string, postID string) {
			nr.Set("user", userId)
			nr.Set("post", postID)
		},
	)

	if e.Request.Header.Get("X-Role") == "true" {
		tools.CacheIdsForCreate(e, "X-Role", tempUserRoles, payload.RoleIds)
	}
	tools.ReplaceJoinTableForUpdate(
		e,
		"X-Role",
		"user_role",
		"user={:user}",
		dbx.Params{"user": e.Record.Id},
		payload.RoleIds,
		func(nr *core.Record, userId string, roleID string) {
			nr.Set("user", userId)
			nr.Set("role", roleID)
		},
	)

	if e.Request.Header.Get("X-Dept") == "true" {
		record, _ := e.App.FindRecordById("dept", e.Record.GetString("dept_id"))

		if record != nil {
			e.Record.Set("dept_name", record.GetString("dept_name"))
		}
	}

	e.Record.SetEmailVisibility(true)

	return e.Next()
}

// 创建成功后处理暂存的角色与岗位关联
func syncUserAfter(e *core.RecordEvent) error {
	// 处理岗位关联 user_post
	tools.ProcessAfterCreateTempIds(e, tempUserPosts, "user_post", func(nr, parent *core.Record, postID string) {
		nr.Set("user", parent.Id)
		nr.Set("post", postID)
	})

	// 处理角色关联 user_role
	tools.ProcessAfterCreateTempIds(e, tempUserRoles, "user_role", func(nr, parent *core.Record, roleID string) {
		nr.Set("user", parent.Id)
		nr.Set("role", roleID)
	})

	return e.Next()
}

// RegisterSystemUserProfile 注册用户相关接口（例如修改密码）
func RegisterSystemUserProfile(app *pocketbase.PocketBase) {

	app.OnRecordCreateRequest("users").BindFunc(syncUser)
	app.OnRecordUpdateRequest("users").BindFunc(syncUser)
	app.OnRecordAfterCreateSuccess("users").BindFunc(syncUserAfter)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 按部门查询用户列表
		se.Router.GET("/api/system/user/list/dept/{dept_id}", func(e *core.RequestEvent) error {
			ri, err := e.RequestInfo()
			if err != nil {
				return e.BadRequestError("请求信息错误", err)
			}
			if ri.Auth == nil {
				return e.UnauthorizedError("未登录或无权限", nil)
			}

			deptID := e.Request.PathValue("dept_id")
			if deptID == "" {
				return e.BadRequestError("缺少部门ID", nil)
			}

			var users []SimpleUser
			q := e.App.DB().
				Select("id", "user_name", "nick_name", "email", "phonenumber", "dept_id", "status", "del_flag", "avatar").
				From("users").
				Where(dbx.HashExp{"id": deptID}).
				AndWhere(dbx.HashExp{"del_flag": "0"}).
				OrderBy("created ASC")

			if err := q.All(&users); err != nil {
				return e.InternalServerError("查询用户失败", err)
			}

			return tools.JSONSuccess(e, users)
		})

		se.Router.PUT("/api/system/user/resetPwd", func(e *core.RequestEvent) error {
			ri, err := e.RequestInfo()
			if err != nil {
				return e.BadRequestError("请求信息错误", err)
			}
			if ri.Auth == nil {
				return e.UnauthorizedError("未登录或无权限", nil)
			}
			var payload PasswordResetPayload

			if err := e.BindBody(&payload); err != nil {
				return e.BadRequestError("无效的请求体", err)
			}

			record, err := e.App.FindRecordById(ri.Auth.Collection(), payload.ID)
			if err != nil {
				return e.InternalServerError("查找用户失败", err)
			}

			record.SetPassword(payload.Password)

			if err := e.App.Save(record); err != nil {
				return e.InternalServerError("保存新密码失败", err)
			}

			return tools.JSONSuccess(e, nil)
		})

		se.Router.PUT("/api/system/user/profile/updatePwd", func(e *core.RequestEvent) error {
			ri, err := e.RequestInfo()
			if err != nil {
				return e.BadRequestError("请求信息错误", err)
			}
			if ri.Auth == nil {
				return e.UnauthorizedError("未登录或无权限", nil)
			}

			var payload PasswordUpdatePayload
			if err := e.BindBody(&payload); err != nil {
				return e.BadRequestError("无效的请求体", err)
			}

			if payload.OldPassword == "" || payload.NewPassword == "" {
				return e.BadRequestError("旧密码和新密码为必填字段", nil)
			}

			record, err := e.App.FindRecordById(ri.Auth.Collection(), ri.Auth.Id)
			if err != nil {
				return e.InternalServerError("查找用户失败", err)
			}

			if !record.ValidatePassword(payload.OldPassword) {
				return e.BadRequestError("旧密码不正确", nil)
			}

			record.SetPassword(payload.NewPassword)

			if err := e.App.Save(record); err != nil {
				return e.InternalServerError("保存新密码失败", err)
			}

			return tools.JSONSuccess(e, nil)
		})

		return se.Next()
	})
}
