package menu

// Menu 模型映射（尽量贴合 SQL 定义）。
type Menu struct {
	ID         string `db:"id" json:"id"`
	MenuName   string `db:"menu_name" json:"menu_name"`
	ParentID   string `db:"parent_id" json:"parent_id"`
	OrderNum   int64  `db:"order_num" json:"order_num"`
	Path       string `db:"path" json:"path"`
	Component  string `db:"component" json:"component"`
	QueryParam string `db:"query_param" json:"query_param"`
	IsFrame    string `db:"is_frame" json:"is_frame"`
	IsCache    string `db:"is_cache" json:"is_cache"`
	MenuType   string `db:"menu_type" json:"menu_type"`
	Visible    string `db:"visible" json:"visible"`
	Status     string `db:"status" json:"status"`
	Perms      string `db:"perms" json:"perms"`
	Icon       string `db:"icon" json:"icon"`
}

// MenuNode 菜单树节点结构体
type MenuNode struct {
	Menu
	Children []*MenuNode `json:"children,omitempty"`
}

// RouteMeta 前端路由元信息结构体
type RouteMeta struct {
	Title   string  `json:"title"`
	Icon    string  `json:"icon,omitempty"`
	NoCache bool    `json:"noCache"`
	Link    *string `json:"link"`
}

// Router 前端路由结构体
type Router struct {
	AlwaysShow bool      `json:"alwaysShow"`
	Children   []Router  `json:"children,omitempty"`
	Component  string    `json:"component"`
	Hidden     bool      `json:"hidden"`
	Meta       RouteMeta `json:"meta"`
	ActiveMenu *string   `json:"activeMenu"`
	Icon       string    `json:"icon"`
	Link       *string   `json:"link"`
	NoCache    bool      `json:"noCache"`
	Title      string    `json:"title"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Redirect   string    `json:"redirect,omitempty"`
}
