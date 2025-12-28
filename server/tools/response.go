package tools

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

// JSONSuccess 返回一个标准化的成功响应
// 使用示例：
//
//	return tools.JSONSuccess(e, data)
func JSONSuccess(e *core.RequestEvent, data any) error {
	return e.JSON(http.StatusOK, map[string]any{
		"code": http.StatusOK,
		"msg":  "操作成功",
		"data": data,
	})
}
