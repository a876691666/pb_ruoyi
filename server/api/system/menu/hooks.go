package menu

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func syncDeleteRoleMenu(e *core.RecordEvent) error {
	if e == nil || e.App == nil || e.Record == nil {
		return nil
	}

	e.App.DB().Delete("role_menu", dbx.HashExp{"menu": e.Record.Id}).Execute()

	return e.Next()
}

func syncMenuDeleteAfter(e *core.RecordEvent) error {
	// 安全检查
	if e == nil || e.App == nil || e.Record == nil {
		return nil
	}

	children, _ := e.App.FindRecordsByFilter("menu", "parent_id={:parent_id}", "", 999, 0, dbx.Params{"parent_id": e.Record.Id})

	for _, child := range children {
		e.App.Delete(child)
	}

	return e.Next()
}
