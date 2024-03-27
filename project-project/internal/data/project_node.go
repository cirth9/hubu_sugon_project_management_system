package data

import "strings"

type ProjectNode struct {
	Id       int64
	Node     string
	Title    string
	IsMenu   int
	IsLogin  int
	IsAuth   int
	CreateAt int64
}

func (*ProjectNode) TableName() string {
	return "ms_project_node"
}

type ProjectNodeTree struct {
	Id       int64
	Node     string
	Title    string
	IsMenu   int
	IsLogin  int
	IsAuth   int
	Pnode    string
	Children []*ProjectNodeTree
}

//node节点  project  project/account

func ToNodeTreeList(list []*ProjectNode) []*ProjectNodeTree {
	var roots []*ProjectNodeTree
	for _, v := range list {
		paths := strings.Split(v.Node, "/")
		if len(paths) == 1 {
			//根节点
			root := &ProjectNodeTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeTree{},
			}
			roots = append(roots, root)
		}
	}
	for _, v := range roots {
		addChild(list, v, 2)
	}
	return roots
}

func addChild(list []*ProjectNode, root *ProjectNodeTree, level int) {
	//root project  project/account
	for _, v := range list {
		paths := strings.Split(v.Node, "/")
		if len(paths) == level && strings.HasPrefix(v.Node, root.Node+"/") {
			child := &ProjectNodeTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeTree{},
			}
			root.Children = append(root.Children, child)
		}
	}
	for _, v := range root.Children {
		addChild(list, v, level+1)
	}
}

type ProjectNodeAuthTree struct {
	Id       int64
	Node     string
	Title    string
	IsMenu   int
	IsLogin  int
	IsAuth   int
	Pnode    string
	Key      string
	Checked  bool
	Children []*ProjectNodeAuthTree
}

func ToAuthNodeTreeList(list []*ProjectNode, checkedList []string) []*ProjectNodeAuthTree {
	checkedMap := make(map[string]struct{})
	for _, v := range checkedList {
		checkedMap[v] = struct{}{}
	}
	var roots []*ProjectNodeAuthTree
	for _, v := range list {
		paths := strings.Split(v.Node, "/")
		if len(paths) == 1 {
			checked := false
			if _, ok := checkedMap[v.Node]; ok {
				checked = true
			}
			//根节点
			root := &ProjectNodeAuthTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeAuthTree{},
				Checked:  checked,
				Key:      v.Node,
			}
			roots = append(roots, root)
		}
	}
	for _, v := range roots {
		addAuthNodeChild(list, v, 2, checkedMap)
	}
	return roots
}

func addAuthNodeChild(list []*ProjectNode, root *ProjectNodeAuthTree, level int, checkedMap map[string]struct{}) {
	for _, v := range list {
		if strings.HasPrefix(v.Node, root.Node+"/") && len(strings.Split(v.Node, "/")) == level {
			//此根节点子节点
			checked := false
			if _, ok := checkedMap[v.Node]; ok {
				checked = true
			}

			child := &ProjectNodeAuthTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeAuthTree{},
				Checked:  checked,
				Key:      v.Node,
			}
			root.Children = append(root.Children, child)
		}
	}
	for _, v := range root.Children {
		addAuthNodeChild(list, v, level+1, checkedMap)
	}
}
