package container

type LinkedTreeNode[V any] struct {
	Value    V
	Children []LinkedTreeNode[V]
}

func NewLinkedTree[V any]() Tree[V] {
	return NewLinkedTreeCap[V](0, 0)
}

func NewLinkedTreeCap[V any](branches uint, depth uint) Tree[V] {
	result := LinkedTreeNode[V]{}

	if branches == 0 || depth == 0 {
		return &result
	}

	newLinkedTreeCapRecursive[V](&result, branches, depth)

	return &result
}

func newLinkedTreeCapRecursive[V any](node *LinkedTreeNode[V], branches uint, depth uint) {
	var depthIter uint = 0
	var branchIter uint
	for ; depthIter < depth; depthIter++ {
		node.Children = make([]LinkedTreeNode[V], branches)
		for branchIter = 0; branchIter < branches; branchIter++ {
			newLinkedTreeCapRecursive[V](&node.Children[branchIter], branches, depth-1)
		}
	}
}

func (lt *LinkedTreeNode[V]) Length() uint {
	return uint(len(LeftBreadthFirst(lt)))
}

func (lt *LinkedTreeNode[V]) GetValue() V {
	return lt.Value
}

func (lt *LinkedTreeNode[V]) SetValue(value V) Tree[V] {
	lt.Value = value
	return lt
}

func (lt *LinkedTreeNode[V]) IsLeaf() bool {
	return len(lt.Children) == 0
}

func (lt *LinkedTreeNode[V]) AddChild(childVal V) Tree[V] {
	if lt.Children == nil {
		lt.Children = []LinkedTreeNode[V]{
			{
				Value:    childVal,
				Children: nil,
			},
		}
	} else {
		lt.Children = append(lt.Children, LinkedTreeNode[V]{
			Value:    childVal,
			Children: nil,
		})
	}

	return lt
}

func (lt *LinkedTreeNode[V]) GetChild(index int) Tree[V] {
	if lt.Children == nil {
		return nil
	}

	return &lt.Children[index]
}

func (lt *LinkedTreeNode[V]) GetChildren() []Tree[V] {
	if lt.Children == nil {
		return []Tree[V]{}
	}

	result := make([]Tree[V], 0, len(lt.Children))

	for _, child := range lt.Children {
		result = append(result, (&child))
	}

	return result
}

func (lt *LinkedTreeNode[V]) GetChildrenData() []V {
	if lt.Children == nil {
		return []V{}
	}

	result := make([]V, 0, len(lt.Children))
	for _, child := range lt.Children {
		result = append(result, child.Value)
	}

	return result
}

func (lt *LinkedTreeNode[V]) AllData() []V {
	return LeftBreadthFirst(lt)
}
