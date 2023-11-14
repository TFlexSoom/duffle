package container

import (
	"testing"
)

func equals(a []int, b []int) bool {
	aLen := len(a)
	bLen := len(b)
	if aLen != bLen {
		return false
	}

	for i := 0; i < aLen; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestOneTree(t *testing.T) {
	tree := NewLinkedTreeCap[int](1, 1)

	tree.SetValue(1)

	if tree.GetValue() != 1 {
		t.Errorf("unexpected value: %v", tree.GetValue())
	}
}

func TestBreadth(t *testing.T) {
	tree := NewLinkedTree[int]()

	tree.SetValue(1)
	tree.AddChild(3).AddChild(2)
	tree.GetChild(1).AddChild(5).AddChild(4)
	tree.GetChild(0).AddChild(6)

	breadthResult := tree.RightBreadthFirst()
	if !equals(breadthResult, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("unexpected right breadth first got : %v", breadthResult)
	}

	tree = NewLinkedTree[int]()
	tree.SetValue(1)
	tree.AddChild(2).AddChild(3).AddChild(4)
	tree.GetChild(0).AddChild(5)
	tree.GetChild(1).AddChild(6)

	breadthResult = tree.LeftBreadthFirst()
	if !equals(breadthResult, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("unexpected left breadth first got : %v", breadthResult)
	}
}

func TestDepth(t *testing.T) {
	tree := NewLinkedTree[int]()

	tree.SetValue(1)
	tree.AddChild(2).AddChild(4)
	tree.GetChild(0).AddChild(3)
	tree.GetChild(1).AddChild(5).AddChild(6)

	depthResult := tree.LeftDepthFirst()
	if !equals(depthResult, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("unexpected left depth first got : %v", depthResult)
	}

	tree = NewLinkedTree[int]()
	tree.SetValue(1)
	tree.AddChild(6).AddChild(3).AddChild(2)
	tree.GetChild(1).AddChild(5).AddChild(4)

	depthResult = tree.RightDepthFirst()
	if !equals(depthResult, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("unexpected right depth first got : %v", depthResult)
	}
}
