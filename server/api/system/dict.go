package system

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// DictData 字典数据结构体
type DictData struct {
	TenantID string `db:"tenant_id" json:"tenant_id"`
	ID       string `db:"id" json:"id"`
}

// RegisterSystemDict 注册 /api/system/role 相关接口
func RegisterSystemDict(app *pocketbase.PocketBase) {
	// app.OnRecordCreateRequest("dict_data").BindFunc(syncDictDataCreate)
}

type syncDictDataReq struct {
}

func syncDictDataCreate(e *core.RecordRequestEvent) error {
	payload := &syncDictDataReq{}
	e.BindBody(payload)

	tenantID := e.Auth.Get("tenant_id")
	e.Record.Set("tenant_id", tenantID)

	return e.Next()
}
