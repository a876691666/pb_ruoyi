package menu

import (
	"fmt"
	"strings"
)

func buildRouters(nodes []*MenuNode) []Router {
	out := make([]Router, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, nodeToRoute(n))
	}
	return out
}

func nodeToRoute(n *MenuNode) Router {
	hidden := n.Visible == "1"
	noCache := n.IsCache == "1"

	var linkPtr *string
	if n.IsFrame == "1" && strings.TrimSpace(n.Path) != "" {
		l := n.Path
		linkPtr = &l
	}

	var activeMenu *string = nil

	r := Router{
		AlwaysShow: len(n.Children) > 0,
		Component:  deriveComponent(n),
		Hidden:     hidden,
		Meta: RouteMeta{
			Title:   n.MenuName,
			Icon:    n.Icon,
			NoCache: noCache,
			Link:    linkPtr,
		},
		ActiveMenu: activeMenu,
		Icon:       n.Icon,
		Link:       linkPtr,
		NoCache:    noCache,
		Title:      n.MenuName,
		Name:       deriveName(n.MenuName, n.ID),
		Path:       derivePath(n),
	}

	if len(n.Children) > 0 {
		r.Redirect = "noRedirect"
		r.Children = buildRouters(n.Children)
	}

	return r
}

func deriveComponent(n *MenuNode) string {
	if n.MenuType == "M" {
		if n.ParentID == "0" {
			return "Layout"
		}
		return "ParentView"
	}
	if strings.TrimSpace(n.Component) != "" {
		return n.Component
	}
	return "ParentView"
}

func derivePath(n *MenuNode) string {
	p := strings.TrimSpace(n.Path)
	if n.ParentID == "0" {
		if p == "" || p[0] != '/' {
			return "/" + strings.TrimPrefix(p, "/")
		}
		return p
	}
	return strings.TrimPrefix(p, "/")
}

func deriveName(menuName string, id string) string {
	base := strings.TrimSpace(menuName)
	if base == "" {
		base = "Menu"
	}
	return fmt.Sprintf("%s%s", base, id)
}
