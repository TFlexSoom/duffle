package container

type GraphTree[V any] struct {
	CurrentIndex  uint64
	Relationships *map[uint64][]uint64
	Data          *[]V
	BreadthAlloc  uint
}

func NewGraphTree[V any]() Tree[V] {
	return NewGraphTreeCap[V](8, 16)
}

func NewGraphTreeCap[V any](depth uint, breadth uint) Tree[V] {
	relations := make(map[uint64][]uint64, depth*breadth)
	data := make([]V, 1, breadth*depth)
	return &GraphTree[V]{
		CurrentIndex:  0,
		Relationships: &relations,
		Data:          &data,
		BreadthAlloc:  breadth,
	}
}

func (st *GraphTree[V]) Length() uint {
	return uint(len(*st.Data))
}

func (st *GraphTree[V]) GetValue() V {
	return (*st.Data)[st.CurrentIndex]
}

func (st *GraphTree[V]) SetValue(value V) Tree[V] {
	(*st.Data)[st.CurrentIndex] = value
	return st
}

func (st *GraphTree[V]) IsLeaf() bool {
	return (*st.Relationships)[st.CurrentIndex] == nil ||
		len((*st.Relationships)[st.CurrentIndex]) == 0
}

func (st *GraphTree[V]) getOrAddRelationship(parent uint64, child uint64) {
	relation := (*st.Relationships)[parent]
	if relation == nil {
		relation = make([]uint64, 0, st.BreadthAlloc)
	}

	(*st.Relationships)[parent] = append(relation, child)
}

func (st *GraphTree[V]) AddChild(value V) Tree[V] {
	indexOfChild := st.Length()
	*st.Data = append((*st.Data), value)
	st.getOrAddRelationship(st.CurrentIndex, uint64(indexOfChild))

	return st
}

func (st *GraphTree[V]) GetChild(index int) Tree[V] {
	relation := (*st.Relationships)[st.CurrentIndex]
	if relation == nil {
		return nil
	} else if len(relation) <= index {
		return nil
	}

	// Shallow Copy
	return &GraphTree[V]{
		CurrentIndex:  (*st.Relationships)[st.CurrentIndex][index],
		Relationships: st.Relationships,
		Data:          st.Data,
		BreadthAlloc:  st.BreadthAlloc,
	}
}

func (st *GraphTree[V]) GetChildren() []V {
	relation := (*st.Relationships)[st.CurrentIndex]
	if relation == nil {
		return nil
	}

	result := make([]V, 0, st.BreadthAlloc)

	for _, childIndex := range relation {
		result = append(result, (*st.Data)[childIndex])
	}

	return result
}

func (st *GraphTree[V]) AllData() []V {
	return *st.Data
}
