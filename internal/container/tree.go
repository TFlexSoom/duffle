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
