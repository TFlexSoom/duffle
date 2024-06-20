package dfa

// import (
// 	"github.com/tflexsoom/duffle/internal/container"
// )

// type DfaGraph[V any] struct {
// 	CurrentIndex  uint64
// 	Relationships *map[uint64](map[rune]uint64)
// 	Data          *[]container.Optional[V]
// 	BreadthAlloc  uint
// }

// func NewDfaGraph[V any]() *DfaGraph[V] {
// 	return NewDfaGraphCap[V](8, 16)
// }

// func NewDfaGraphCap[V any](depth uint, breadth uint) *DfaGraph[V] {
// 	relations := make(map[uint64](map[rune]uint64), depth*breadth)
// 	data := make([]container.Optional[V], 1, breadth*depth)
// 	return &DfaGraph[V]{
// 		CurrentIndex:  0,
// 		Relationships: &relations,
// 		Data:          &data,
// 		BreadthAlloc:  breadth,
// 	}
// }

// func (st *DfaGraph[V]) Length() uint {
// 	return uint(len(*st.Data))
// }

// func (st *DfaGraph[V]) GetValue() V {
// 	return (*st.Data)[st.CurrentIndex].MustGet()
// }

// func (st *DfaGraph[V]) SetValue(value V) *DfaGraph[V] {
// 	(*st.Data)[st.CurrentIndex] = container.NewPointerOf[V](value)
// 	return st
// }

// func (st *DfaGraph[V]) IsLeaf() bool {
// 	return (*st.Relationships)[st.CurrentIndex] == nil ||
// 		len((*st.Relationships)[st.CurrentIndex]) == 0
// }

// func (st *DfaGraph[V]) getOrAddRelationship(parent uint64, edge rune, child uint64) {
// 	relation := (*st.Relationships)[parent]
// 	if relation == nil {
// 		relation = make(map[rune]uint64, st.BreadthAlloc)
// 	}

// 	relation[edge] = child
// 	(*st.Relationships)[parent] = relation
// }

// func (st *DfaGraph[V]) AddChild(edge rune, value V) *DfaGraph[V] {
// 	indexOfChild := st.Length()
// 	*st.Data = append((*st.Data), container.NewPointerOf[V](value))
// 	st.getOrAddRelationship(st.CurrentIndex, edge, uint64(indexOfChild))

// 	return st
// }

// func (st *DfaGraph[V]) GetChild(index int) *DfaGraph[V] {
// 	relation := (*st.Relationships)[st.CurrentIndex]
// 	if relation == nil {
// 		return nil
// 	} else if len(relation) <= index {
// 		return nil
// 	}

// 	// Shallow Copy
// 	return &DfaGraph[V]{
// 		CurrentIndex:  (*st.Relationships)[st.CurrentIndex][0],
// 		Relationships: st.Relationships,
// 		Data:          st.Data,
// 		BreadthAlloc:  st.BreadthAlloc,
// 	}
// }

// func (st *DfaGraph[V]) GetChildren() []*DfaGraph[V] {
// 	relation := (*st.Relationships)[st.CurrentIndex]
// 	if relation == nil {
// 		return nil
// 	}

// 	result := make([]*DfaGraph[V], 0, st.BreadthAlloc)

// 	for _, childIndex := range relation {
// 		result = append(result, st.GetChild(int(childIndex)))
// 	}

// 	return result
// }
