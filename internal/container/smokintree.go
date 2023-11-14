package container

// Capped at depth 64
// Could be more with a Binary Tree Bookmarks Field
type SmokinTree64[V any] struct {
	Bookmarks map[int64]int
	Continues []int
	Pages     []V
	Capacity  int
}

func NewSmokinTree[V any]() Tree[V] {

}

func NewSmokinTreeCap[V any]() Tree[V] {

}

func (st *SmokinTree64) Length() uint           {}
func (st *SmokinTree64) GetValue() V            {}
func (st *SmokinTree64) SetValue(V) Tree[V]     {}
func (st *SmokinTree64) IsLeaf() bool           {}
func (st *SmokinTree64) AddChild(V) Tree[V]     {}
func (st *SmokinTree64) GetChild(int) Tree[V]   {}
func (st *SmokinTree64) GetChildren() []V       {}
func (st *SmokinTree64) LeftDepthFirst() []V    {}
func (st *SmokinTree64) RightDepthFirst() []V   {}
func (st *SmokinTree64) LeftBreadthFirst() []V  {}
func (st *SmokinTree64) RightBreadthFirst() []V {}
