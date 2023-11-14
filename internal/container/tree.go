package container

type Tree[V any] interface {
	Length() uint
	GetValue() V
	SetValue(V) Tree[V]
	IsLeaf() bool
	AddChild(V) Tree[V]
	GetChild(int) Tree[V]
	GetChildren() []V
	LeftDepthFirst() []V
	RightDepthFirst() []V
	LeftBreadthFirst() []V
	RightBreadthFirst() []V
}

type LinkedTreeNode[V any] struct {
	Value    V
	Children []LinkedTreeNode[V]
}

func NewLinkedTree[V any]() LinkedTreeNode[V] {
	return NewLinkedTreeCap[V](0, 0)
}

func NewLinkedTreeCap[V any](branches uint, depth uint) LinkedTreeNode[V] {
	result := LinkedTreeNode[V]{}

	if branches == 0 || depth == 0 {
		return result
	}

	newLinkedTreeCapRecursive[V](&result, branches, depth)

	return result
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
	return uint(len(lt.LeftBreadthFirst()))
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

func (lt *LinkedTreeNode[V]) GetChildren() []V {
	if lt.Children == nil {
		return []V{}
	}

	result := make([]V, 0, len(lt.Children))
	for _, child := range lt.Children {
		result = append(result, child.Value)
	}

	return result
}

func (lt *LinkedTreeNode[V]) collector(
	appender func(*[]LinkedTreeNode[V], *LinkedTreeNode[V]),
	nextItem func(*[]LinkedTreeNode[V], int) *LinkedTreeNode[V],
	isLeftRightAdd bool,
) []V {
	data := make([]LinkedTreeNode[V], 0, 1024)
	appender(&data, lt)
	result := make([]V, 0, 1024)

	for count := 1; count > 0; count-- {
		curNode := nextItem(&data, count)

		result = append(result, curNode.Value)
		if curNode.IsLeaf() {
			continue
		}

		if isLeftRightAdd {
			length := len(curNode.Children)
			for i := 0; i < length; i++ {
				appender(&data, &curNode.Children[i])
				count++
			}
		} else {
			for i := len(curNode.Children) - 1; i >= 0; i-- {
				appender(&data, &curNode.Children[i])
				count++
			}
		}
	}

	return result
}

func (lt LinkedTreeNode[V]) LeftDepthFirst() []V {
	return lt.collector(
		func(stack *[]LinkedTreeNode[V], node *LinkedTreeNode[V]) {
			*stack = append(*stack, *node)
		},
		func(stack *[]LinkedTreeNode[V], count int) *LinkedTreeNode[V] {
			result := (*stack)[count-1]
			(*stack) = (*stack)[0 : count-1]
			return &result
		},
		false,
	)
}

func (lt LinkedTreeNode[V]) RightDepthFirst() []V {
	return lt.collector(
		func(stack *[]LinkedTreeNode[V], node *LinkedTreeNode[V]) {
			*stack = append(*stack, *node)
		},
		func(stack *[]LinkedTreeNode[V], count int) *LinkedTreeNode[V] {
			result := (*stack)[count-1]
			(*stack) = (*stack)[0 : count-1]
			return &result
		},
		true,
	)
}

func (lt LinkedTreeNode[V]) LeftBreadthFirst() []V {
	return lt.collector(
		func(queue *[]LinkedTreeNode[V], node *LinkedTreeNode[V]) {
			*queue = append(*queue, *node)
		},
		func(queue *[]LinkedTreeNode[V], count int) *LinkedTreeNode[V] {
			result := (*queue)[0]
			(*queue) = (*queue)[1:count]
			return &result
		},
		true,
	)
}

func (lt LinkedTreeNode[V]) RightBreadthFirst() []V {
	return lt.collector(
		func(queue *[]LinkedTreeNode[V], node *LinkedTreeNode[V]) {
			*queue = append(*queue, *node)
		},
		func(queue *[]LinkedTreeNode[V], count int) *LinkedTreeNode[V] {
			result := (*queue)[0]
			(*queue) = (*queue)[1:count]
			return &result
		},
		false,
	)
}
