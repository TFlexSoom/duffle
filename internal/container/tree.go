package container

type Tree[V any] interface {
	Length() uint
	GetValue() V
	SetValue(V) Tree[V]
	IsLeaf() bool
	AddChild(V) Tree[V]
	GetChild(int) Tree[V]
	GetChildren() []Tree[V]
	GetChildrenData() []V
	AllData() []V
}

func AddChildren[V any](self Tree[V], other Tree[V]) Tree[V] {
	subTree := self.AddChild(other.GetValue())
	for _, child := range other.GetChildren() {
		AddChildren[V](subTree, child)
	}

	return self
}

func TreeCollector[V any](
	lt Tree[V],
	appender func(Tree[V]),
	nextItem func(int) Tree[V],
	isLeftRightAdd bool,
) []V {
	appender(lt)
	result := make([]V, 0, 1024)

	for count := 1; count > 0; count-- {
		curNode := nextItem(count)

		result = append(result, curNode.GetValue())
		if curNode.IsLeaf() {
			continue
		}

		if isLeftRightAdd {
			length := len(curNode.GetChildren())
			for i := 0; i < length; i++ {
				appender(curNode.GetChild(i))
				count++
			}
		} else {
			for i := len(curNode.GetChildren()) - 1; i >= 0; i-- {
				appender(curNode.GetChild(i))
				count++
			}
		}
	}

	return result
}

func LeftDepthFirst[V any](lt Tree[V]) []V {
	stack := make([]Tree[V], 0, 1024)
	return TreeCollector[V](
		lt,
		func(node Tree[V]) {
			stack = append(stack, node)
		},
		func(count int) Tree[V] {
			result := stack[count-1]
			stack = stack[0 : count-1]
			return result
		},
		false,
	)
}

func RightDepthFirst[V any](lt Tree[V]) []V {
	stack := make([]Tree[V], 0, 1024)
	return TreeCollector[V](
		lt,
		func(node Tree[V]) {
			stack = append(stack, node)
		},
		func(count int) Tree[V] {
			result := stack[count-1]
			stack = stack[0 : count-1]
			return result
		},
		true,
	)
}

func LeftBreadthFirst[V any](lt Tree[V]) []V {
	queue := make([]Tree[V], 0, 1024)
	return TreeCollector[V](
		lt,
		func(node Tree[V]) {
			queue = append(queue, node)
		},
		func(count int) Tree[V] {
			result := queue[0]
			queue = queue[1:count]
			return result
		},
		true,
	)
}

func RightBreadthFirst[V any](lt Tree[V]) []V {
	queue := make([]Tree[V], 0, 1024)
	return TreeCollector[V](
		lt,
		func(node Tree[V]) {
			queue = append(queue, node)
		},
		func(count int) Tree[V] {
			result := queue[0]
			queue = queue[1:count]
			return result
		},
		false,
	)
}
