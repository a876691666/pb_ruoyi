package menu

import "sort"

func uniqueMenus(in []Menu) []Menu {
	seen := make(map[string]struct{}, len(in))
	out := make([]Menu, 0, len(in))
	for _, m := range in {
		if _, ok := seen[m.ID]; ok {
			continue
		}
		seen[m.ID] = struct{}{}
		out = append(out, m)
	}
	return out
}

func buildMenuTree(list []Menu) []*MenuNode {
	byMenuID := make(map[string]*MenuNode, len(list))
	roots := make([]*MenuNode, 0)

	for i := range list {
		m := list[i]
		byMenuID[m.ID] = &MenuNode{Menu: m}
	}

	for _, n := range byMenuID {
		if n.ParentID == "0" {
			roots = append(roots, n)
			continue
		}
		if p, ok := byMenuID[n.ParentID]; ok {
			p.Children = append(p.Children, n)
		} else {
			roots = append(roots, n)
		}
	}

	var sortRec func(nodes []*MenuNode)
	sortRec = func(nodes []*MenuNode) {
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
