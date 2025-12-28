package system

import (
	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterSystemCollections 注册 /api/system/collections 相关接口
func RegisterSystemCollections(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/system/collections", func(e *core.RequestEvent) error {
			allCollections, err := app.FindAllCollections()

			if err != nil {
				return err
			}

			return tools.JSONSuccess(e, allCollections)
		})

		se.Router.GET("/api/system/collection/{collectionName}", func(e *core.RequestEvent) error {
			collection, err := app.FindCollectionByNameOrId(e.Request.PathValue("collectionName"))

			if err != nil {
				return err
			}

			return tools.JSONSuccess(e, collection)
		})

		return se.Next()
	})
}
