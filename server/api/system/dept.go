package system

import (
	"fmt"
	"sort"
	"strings"

	"pocketbase-ruoyi/tools"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// Dept 部门模型（贴合表结构的最小字段集）
type Dept struct {
	ID        string `db:"id" json:"id"`
	DeptName  string `db:"dept_name" json:"dept_name"`
	ParentID  string `db:"parent_id" json:"parent_id"`
	Ancestors string `db:"ancestors" json:"ancestors"`
	OrderNum  int64  `db:"order_num" json:"order_num"`
	Status    string `db:"status" json:"status"`
	DelFlag   string `db:"del_flag" json:"del_flag"`
}

// DeptNode 部门树节点
type DeptNode struct {
	Dept
	Children []*DeptNode `json:"children,omitempty"`
}

// GetChildrensByDeptID 获取部门的所有子部门
func GetChildrensByDeptID(e *core.RequestEvent, deptID string) []Dept {
	rows := []Dept{}

	err := e.App.DB().Select("*").From("dept").
		Where(
			dbx.Or(
				dbx.Like("ancestors", fmt.Sprintf(",%s,", deptID)),
				dbx.Like("ancestors", fmt.Sprintf(",%s", deptID)).Match(true, false),
			),
		).
		All(&rows)

	if err != nil {
		return []Dept{}
	}

	return rows
}

// RegisterSystemDept 注册 /api/system/user/deptTree 接口
func RegisterSystemDept(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/system/user/deptTree", func(e *core.RequestEvent) error {
			ri, err := e.RequestInfo()
			if err != nil {
				return e.BadRequestError("请求信息错误", err)
			}
			if ri.Auth == nil {
				return e.UnauthorizedError("未登录或无权限", nil)
			}

			// 查询有效部门（status=0 且 del_flag=0）
			var depts []Dept
			q := e.App.DB().
				Select("id", "dept_name", "parent_id", "order_num", "status", "del_flag").
				From("dept").
				Where(dbx.HashExp{"status": "0"}).
				AndWhere(tools.BuildDataScopeExpression(e, "dept")).
				OrderBy("parent_id ASC", "order_num ASC", "id ASC")

			if err := q.All(&depts); err != nil {
				return e.InternalServerError("查询部门失败", err)
			}

			// 去重（按 PocketBase 系统 id）
			depts = uniqueDepts(depts)

			// 构建树并返回
			tree := buildDeptTree(depts)
			return tools.JSONSuccess(e, tree)
		})
		return se.Next()
	})

	// 新增：创建成功后，基于 parent_id 设置当前记录 ancestors，并更新直接子节点的 ancestors
	app.OnRecordAfterCreateSuccess("dept").BindFunc(func(e *core.RecordEvent) error {
		if err := setAncestorsAndChildren(e.App, e.Record); err != nil {
			return err
		}
		return e.Next()
	})

	// 新增：更新成功后，同步 ancestors 与直接子节点 ancestors
	app.OnRecordAfterUpdateSuccess("dept").BindFunc(func(e *core.RecordEvent) error {
		if err := setAncestorsAndChildren(e.App, e.Record); err != nil {
			return err
		}
		return e.Next()
	})
}

// setAncestorsAndChildren 根据记录的 parent_id 计算并设置自身 ancestors，同时更新所有直接子部门的 ancestors
func setAncestorsAndChildren(app core.App, rec *core.Record) error {
	// 1) 计算自身 ancestors
	parentID := strings.TrimSpace(rec.GetString("parent_id"))
	var newAncestors string

	if parentID == "" || parentID == "0" {
		newAncestors = "0"
	} else {
		parent, err := app.FindRecordById("dept", parentID)
		if err != nil {
			// 若父级不存在，降级为根（0）
			newAncestors = "0"
		} else {
			pa := strings.Trim(parent.GetString("ancestors"), ",")
			if pa == "" {
				pa = "0"
			}
			newAncestors = pa + "," + parent.Id
		}
	}

	if rec.GetString("ancestors") != newAncestors {
		rec.Set("ancestors", newAncestors)
		if err := app.Save(rec); err != nil {
			return err
		}
	}

	// 2) 更新直接子部门的 ancestors = 本记录 ancestors + "," + 本记录 id
	children := []*core.Record{}
	if err := app.RecordQuery("dept").Where(dbx.HashExp{"parent_id": rec.Id}).All(&children); err != nil {
		return err
	}
	for _, ch := range children {
		ch.Set("ancestors", newAncestors+","+rec.Id)
		if err := app.Save(ch); err != nil {
			return err
		}
	}

	return nil
}

func uniqueDepts(in []Dept) []Dept {
	seen := make(map[string]struct{}, len(in))
	out := make([]Dept, 0, len(in))
	for _, d := range in {
		if _, ok := seen[d.ID]; ok {
			continue
		}
		seen[d.ID] = struct{}{}
		out = append(out, d)
	}
	return out
}

func buildDeptTree(list []Dept) []*DeptNode {
	byID := make(map[string]*DeptNode, len(list))
	roots := make([]*DeptNode, 0)

	// 索引
	for i := range list {
		d := list[i]
		byID[d.ID] = &DeptNode{Dept: d}
	}

	// 组装父子
	for _, n := range byID {
		if n.ParentID == "0" {
			roots = append(roots, n)
			continue
		}
		if p, ok := byID[n.ParentID]; ok {
			p.Children = append(p.Children, n)
		} else {
			// 无父节点时作为根节点
			roots = append(roots, n)
		}
	}

	// 递归排序：order_num 升序，同序按 dept_id 升序
	var sortRec func(nodes []*DeptNode)
	sortRec = func(nodes []*DeptNode) {
		sort.SliceStable(nodes, func(i, j int) bool {
			if nodes[i].OrderNum == nodes[j].OrderNum {
				return nodes[i].ID < nodes[j].ID
			}
			return nodes[i].OrderNum < nodes[j].OrderNum
		})
		for _, ch := range nodes {
			if len(ch.Children) > 0 {
				sortRec(ch.Children)
			}
		}
	}
	sortRec(roots)

	return roots
}
