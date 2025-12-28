package system

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/xuri/excelize/v2"
)

// RegisterSystemImportExport 注册集合级别的 Excel 导入/导出接口：
// GET  /api/collections/{collection}/export -> 导出
// POST /api/collections/{collection}/import -> 导入
// 注意：RBAC 已在 auth/rbac.go 中处理权限标识 system:{collection}:export / import
func RegisterSystemImportExport(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// 导出 Excel
		se.Router.GET("/api/collections/{collection}/export", func(e *core.RequestEvent) error {
			collName := e.Request.PathValue("collection")
			coll, err := e.App.FindCachedCollectionByNameOrId(collName)
			if err != nil || coll == nil {
				return e.BadRequestError("集合不存在", err)
			}

			// 读取查询参数 filter 和 sort，用于导出时的筛选与排序
			// - filter: PocketBase 过滤表达式，如 name~"abc" && status=true
			// - sort:   排序字段，支持前缀 '-' 表示倒序，如 -created
			q := e.Request.URL.Query()
			filter := q.Get("filter")
			if filter == "" { // 兼容可能的拼写错误 fiter
				filter = q.Get("fiter")
			}
			sort := q.Get("sort")

			// limit & page 简单分页控制，避免一次性拉取过多数据
			limit := parsePositiveInt(q.Get("limit"), 1000)
			page := parsePositiveInt(q.Get("page"), 1)
			if limit > 5000 { // 强制上限
				limit = 5000
			}
			offset := (page - 1) * limit

			records, err := e.App.FindRecordsByFilter(coll.Name, filter, sort, limit, offset)
			if err != nil {
				return e.InternalServerError("查询记录失败", err)
			}

			f := excelize.NewFile()
			sheet := f.GetSheetName(0)

			// 收集字段名（使用集合 schema 中的字段 + id + 创建/更新时间）
			fieldNames := collectExportFields(coll)
			for i, name := range fieldNames {
				cell, _ := excelize.CoordinatesToCellName(i+1, 1)
				_ = f.SetCellValue(sheet, cell, name)
			}

			for rIdx, rec := range records {
				for cIdx, name := range fieldNames {
					cell, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+2)
					_ = f.SetCellValue(sheet, cell, rec.Get(name))
				}
			}

			buf, err := f.WriteToBuffer()
			if err != nil {
				return e.InternalServerError("生成 Excel 失败", err)
			}

			filename := fmt.Sprintf("%s_export.xlsx", coll.Name)
			e.Response.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			e.Response.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
			_, _ = e.Response.Write(buf.Bytes())
			return nil
		})

		// 导入 Excel
		se.Router.POST("/api/collections/{collection}/import", func(e *core.RequestEvent) error {
			collName := e.Request.PathValue("collection")
			coll, err := e.App.FindCachedCollectionByNameOrId(collName)
			if err != nil || coll == nil {
				return e.BadRequestError("集合不存在", err)
			}

			file, header, err := e.Request.FormFile("file")
			if err != nil {
				return e.BadRequestError("缺少上传文件字段 file", err)
			}
			defer file.Close()

			if header != nil && header.Size > 20*1024*1024 { // 20MB 限制
				return e.BadRequestError("文件过大，限制 20MB", nil)
			}

			excel, err := excelize.OpenReader(file)
			if err != nil {
				return e.BadRequestError("解析 Excel 失败", err)
			}
			defer excel.Close()

			sheet := excel.GetSheetName(0)
			rows, err := excel.GetRows(sheet)
			if err != nil || len(rows) == 0 {
				return e.BadRequestError("Excel 无可用数据", err)
			}

			headers := normalizeHeaders(rows[0])
			if len(headers) == 0 {
				return e.BadRequestError("首行表头为空", nil)
			}

			validFields := buildFieldSet(coll)
			imported := 0
			failed := 0
			errs := []string{}

			for r := 1; r < len(rows); r++ {
				row := rows[r]
				rec := core.NewRecord(coll)
				for cIdx, h := range headers {
					if h == "" || !validFields[h] { // 跳过无效字段
						continue
					}
					val := ""
					if cIdx < len(row) {
						val = row[cIdx]
					}
					rec.Set(h, val)
				}
				if err := e.App.Save(rec); err != nil {
					failed++
					if len(errs) < 10 { // 只保留前 10 条错误信息
						errs = append(errs, fmt.Sprintf("第 %d 行: %v", r+1, err))
					}
					continue
				}
				imported++
			}

			return tools.JSONSuccess(e, map[string]any{
				"imported": imported,
				"failed":   failed,
				"errors":   errs,
			})
		})

		return se.Next()
	})
}

func parsePositiveInt(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil || i <= 0 {
		return def
	}
	return i
}

// collectExportFields 组装导出字段
func collectExportFields(coll *core.Collection) []string {
	seen := map[string]struct{}{}
	fields := []string{}
	for _, f := range coll.Fields {
		name := f.GetName()
		if name == "" {
			continue
		}
		if _, ok := seen[name]; !ok {
			fields = append(fields, name)
			seen[name] = struct{}{}
		}
	}
	return fields
}

func normalizeHeaders(hdrs []string) []string {
	out := make([]string, len(hdrs))
	for i, h := range hdrs {
		out[i] = strings.TrimSpace(h)
	}
	return out
}

func buildFieldSet(coll *core.Collection) map[string]bool {
	m := map[string]bool{}
	for _, f := range coll.Fields {
		name := f.GetName()
		if name != "" {
			m[name] = true
		}
	}
	return m
}

// helper: 读取 multipart 文件（当前使用 FormFile 简化，此方法留作后续扩展）
func readUploadedFile(fh *multipart.FileHeader) ([]byte, error) {
	if fh == nil {
		return nil, errors.New("file header nil")
	}
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
