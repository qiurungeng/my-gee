package gee

import "strings"

//路由前缀树节点
//      GET           POST   ...
//       ↓              ↓
//      nil            nil
//     /   \          /
//  /hello /docs    /login
//           \
//            /*
type node struct {
	// 根节点到本节点的路由
	pattern string
	// 本节点的路由
	part string
	// 本节点的子节点
	children []*node
	// 是否模糊匹配，part 含有 : 或 * 时为true
	isWild bool
}


// 返回第一个匹配成功的节点
// 用于插入
func (n node) matchChild(part string) *node {
	for _, child := range n.children{
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}


// 返回子节点中所有匹配成功的节点
// 用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}


// 插入功能
// 递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个
func (n *node) insert(pattern string, parts []string, height int)  {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}


// 查询功能
// 同样也是递归查询每一层的节点
// 退出规则是，匹配到了*，匹配失败，或者匹配到了第len(parts)层节点。
func (n *node) search(parts []string, height int) *node {
	// 已匹配到最后一层节点 或 当前节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	// 子节点中所有匹配的节点
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}


