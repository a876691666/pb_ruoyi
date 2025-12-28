package tools

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// AfterCreateTempMap 用于在创建阶段暂存关联ID的内存映射
// key: __temp_key, value: 需要在创建成功后落库的ID列表
type AfterCreateTempMap = map[string][]string

// ProcessAfterCreateTempIds 通用的“创建成功后处理暂存关联”逻辑
//   - e: PocketBase 的 RecordEvent
//   - tempMap: 暂存映射，通常是包级 map[string][]string
//   - collectionName: 目标关联表集合名称，例如 "role_menu"、"user_post"、"user_role"
//   - set: 对新建记录进行字段赋值的回调，例如设置外键和关联ID
//     示例：
//     ProcessAfterCreateTempIds(e, tempRoleMenus, "role_menu", func(nr, parent *core.Record, id string){
//     nr.Set("role_id", parent.Id)
//     nr.Set("menu_id", id)
//     })
func ProcessAfterCreateTempIds(
	e *core.RecordEvent,
	tempMap AfterCreateTempMap,
	collectionName string,
	set func(newRec *core.Record, parent *core.Record, id string),
) {
	if e == nil || e.Record == nil || e.App == nil {
		return
	}

	tempKey := e.Record.GetString("__temp_key")
	if tempKey == "" {
		return
	}

	ids, ok := tempMap[tempKey]
	if !ok || len(ids) == 0 {
		return
	}

	coll, err := e.App.FindCollectionByNameOrId(collectionName)
	if err != nil {
		return
	}

	for _, id := range ids {
		newRecord := core.NewRecord(coll)
		set(newRecord, e.Record, id)
		_ = e.App.Save(newRecord)
	}

	delete(tempMap, tempKey)
}

// EnsureTempKeyForRequest 确保在创建阶段存在 __temp_key，并返回该值
func EnsureTempKeyForRequest(e *core.RecordRequestEvent) string {
	tempKey := e.Record.GetString("__temp_key")
	if tempKey == "" {
		tempKey = core.GenerateDefaultRandomId()
		e.Record.Set("__temp_key", tempKey)
	}
	return tempKey
}

// CacheIdsForCreate 在创建阶段（记录尚未持久化，无ID）缓存需要在创建成功后写入的IDs
// 仅当 headerName 对应请求头为 "true" 时生效
func CacheIdsForCreate(e *core.RecordRequestEvent, headerName string, tempMap AfterCreateTempMap, ids []string) {
	if e == nil || e.Request == nil || e.Record == nil || len(ids) == 0 {
		return
	}
	if e.Request.Header.Get(headerName) != "true" {
		return
	}
	if e.Record.Id != "" { // 非创建阶段
		return
	}
	tempKey := EnsureTempKeyForRequest(e)
	tempMap[tempKey] = ids
}

// ReplaceJoinTableForUpdate 在更新阶段替换关联表（先删后插）
// 仅当 headerName 对应请求头为 "true" 且记录已持久化（有ID）时生效
// filterExp 示例："role_id={:role_id}" 或 "user.id={:user_id}"
// params 需包含用于 filter 的 recordID 值
// set 用于设置新记录的关联字段
func ReplaceJoinTableForUpdate(
	e *core.RecordRequestEvent,
	headerName string,
	joinCollection string,
	filterExp string,
	params dbx.Params,
	ids []string,
	set func(newRec *core.Record, parentId string, id string),
) {
	if e == nil || e.Request == nil || e.Record == nil || e.App == nil {
		return
	}
	if e.Request.Header.Get(headerName) != "true" {
		return
	}
	if e.Record.Id == "" { // 创建阶段不做更新逻辑
		return
	}

	// 删除旧关联
	records, _ := e.App.FindRecordsByFilter(joinCollection, filterExp, "", 999, 0, params)
	for _, r := range records {
		_ = e.App.Delete(r)
	}

	coll, err := e.App.FindCollectionByNameOrId(joinCollection)
	if err != nil {
		return
	}
	for _, id := range ids {
		newRec := core.NewRecord(coll)
		set(newRec, e.Record.Id, id)
		_ = e.App.Save(newRec)
	}
}
