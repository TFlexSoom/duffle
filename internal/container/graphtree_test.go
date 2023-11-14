package container

import (
	"testing"
)

func TestOneGraphTree(t *testing.T) {
	tree := NewGraphTreeCap[int](1, 1)

	tree.SetValue(1)

	if tree.GetValue() != 1 {
		t.Errorf("unexpected value: %v", tree.GetValue())
	}
}

func TestGraphTreeBreadth(t *testing.T) {
	tree := NewGraphTree[int]()

	tree.SetValue(1)
	tree.AddChild(3).AddChild(2)
	tree.GetChild(1).AddChild(5).AddChild(4)
	tree.GetChild(0).AddChild(6)

	breadthResult := RightBreadthFirst(tree)
	if !equals(breadthResult, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("unexpected right breadth first got : %v", breadthResult)
	}

	tree = NewGraphTree[int]()
	tree.SetValue(1)
	tree.AddChild(2).AddChild(3).AddChild(4)
	tree.GetChild(0).AddChild(5)
	tree.GetChild(1).AddChild(6)

	breadthResult = LeftBreadthFirst(tree)
	if !equals(breadthResult, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("unexpected left breadth first got : %v", breadthResult)
	}
}

func TestGraphTreeDepth(t *testing.T) {
	tree := NewGraphTree[int]()

	tree.SetValue(1)
	tree.AddChild(2).AddChild(4)
	tree.GetChild(0).AddChild(3)
	tree.GetChild(1).AddChild(5).AddChild(6)

	depthResult := LeftDepthFirst(tree)
	if !equals(depthResult, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("unexpected left depth first got : %v", depthResult)
	}

	tree = NewGraphTree[int]()
	tree.SetValue(1)
	tree.AddChild(6).AddChild(3).AddChild(2)
	tree.GetChild(1).AddChild(5).AddChild(4)

	depthResult = RightDepthFirst(tree)
	if !equals(depthResult, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("unexpected right depth first got : %v", depthResult)
	}
}

func TestGraphTreeMaxDepthSpeed(t *testing.T) {
	tree := NewGraphTree[int]()

	expected := make([]int, 0, 100)

	treeIter := tree
	for i := 0; i < 1_000_000; i++ {
		expected = append(expected, i)
		treeIter.SetValue(i)
		treeIter.AddChild(0)
		treeIter = treeIter.GetChild(0)
	}

	// Trailing 0 okay
	expected = append(expected, 0)

	treeResult := tree.AllData()
	if !equals(treeResult, expected) {
		t.Errorf("unexpected left depth first got : %v", treeResult)
	}
}
