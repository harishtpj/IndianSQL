package btree

import (
	"errors"
	"testing"

	"github.com/harishtpj/indiansql/internal/apperrors"
	"github.com/harishtpj/indiansql/internal/page"
)

func fillLeafToCapacity(
	t *testing.T,
	n *Node,
	start uint64,
) uint64 {
	t.Helper()

	key := start

	for {
		err := n.page.AppendCell(&page.Cell{
			Key:   key,
			Value: []byte("x"),
		})

		if err != nil {
			if errors.Is(err, apperrors.ErrNodeFull) {
				return key
			}

			t.Fatalf("append key %d: %v", key, err)
		}

		key++
	}
}

func TestSplitLeaf(t *testing.T) {
	p := newTestPager(t)

	rootPageID := uint32(0)

	rootPg, err := p.GetPage(rootPageID)
	if err != nil {
		t.Fatal(err)
	}

	page.InitPage(rootPg.Data, page.PageTypeLeaf)

	root, err := LoadNode(p, rootPageID)
	if err != nil {
		t.Fatal(err)
	}

	nextKey := fillLeafToCapacity(t, root, 1)

	if nextKey <= 1 {
		t.Fatal("expected multiple inserted keys")
	}

	rightPageID := p.PageCount()

	_, err = p.GetPage(rightPageID)
	if err != nil {
		t.Fatal(err)
	}

	split, err := root.SplitLeaf(rightPageID)
	if err != nil {
		t.Fatal(err)
	}

	right, err := LoadNode(p, rightPageID)
	if err != nil {
		t.Fatal(err)
	}

	leftCells, err := root.LeafCells()
	if err != nil {
		t.Fatal(err)
	}

	rightCells, err := right.LeafCells()
	if err != nil {
		t.Fatal(err)
	}

	if len(leftCells) == 0 {
		t.Fatal("left page empty after split")
	}

	if len(rightCells) == 0 {
		t.Fatal("right page empty after split")
	}

	// separator == first key in right child
	if split.Separator != rightCells[0].Key {
		t.Fatalf(
			"separator mismatch: got=%d want=%d",
			split.Separator,
			rightCells[0].Key,
		)
	}

	// left max < right min
	leftMax := leftCells[len(leftCells)-1].Key
	rightMin := rightCells[0].Key

	if leftMax >= rightMin {
		t.Fatalf(
			"invalid split ordering: leftMax=%d rightMin=%d",
			leftMax,
			rightMin,
		)
	}
}

func TestSplitLeafPreservesAllKeys(t *testing.T) {
	p := newTestPager(t)

	rootPageID := uint32(0)

	rootPg, err := p.GetPage(rootPageID)
	if err != nil {
		t.Fatal(err)
	}

	page.InitPage(rootPg.Data, page.PageTypeLeaf)

	root, err := LoadNode(p, rootPageID)
	if err != nil {
		t.Fatal(err)
	}

	inserted := make(map[uint64]bool)

	var key uint64 = 1

	for {
		err := root.page.AppendCell(&page.Cell{
			Key:   key,
			Value: []byte("x"),
		})

		if err != nil {
			if errors.Is(err, apperrors.ErrNodeFull) {
				break
			}

			t.Fatal(err)
		}

		inserted[key] = true
		key++
	}

	rightPageID := p.PageCount()

	_, err = p.GetPage(rightPageID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = root.SplitLeaf(rightPageID)
	if err != nil {
		t.Fatal(err)
	}

	rightNode, err := LoadNode(p, rightPageID)
	if err != nil {
		t.Fatal(err)
	}

	leftCells, err := root.LeafCells()
	if err != nil {
		t.Fatal(err)
	}

	rightCells, err := rightNode.LeafCells()
	if err != nil {
		t.Fatal(err)
	}

	found := make(map[uint64]bool)

	for _, c := range leftCells {
		found[c.Key] = true
	}

	for _, c := range rightCells {
		found[c.Key] = true
	}

	if len(found) != len(inserted) {
		t.Fatalf(
			"key count mismatch after split: got=%d want=%d",
			len(found),
			len(inserted),
		)
	}

	for k := range inserted {
		if !found[k] {
			t.Fatalf("missing key after split: %d", k)
		}
	}
}

func TestSplitRootLeaf(t *testing.T) {
	p := newTestPager(t)

	rootPageID := uint32(0)

	rootPg, err := p.GetPage(rootPageID)
	if err != nil {
		t.Fatal(err)
	}

	page.InitPage(rootPg.Data, page.PageTypeLeaf)

	tree := &BTree{
		pager: p,
		root:  rootPageID,
	}

	root, err := LoadNode(p, rootPageID)
	if err != nil {
		t.Fatal(err)
	}

	fillLeafToCapacity(t, root, 1)

	err = tree.splitRootLeaf(root)
	if err != nil {
		t.Fatal(err)
	}

	if tree.root == rootPageID {
		t.Fatal("root page did not change after split")
	}

	newRoot, err := LoadNode(p, tree.root)
	if err != nil {
		t.Fatal(err)
	}

	if !newRoot.IsInternal() {
		t.Fatal("new root is not internal")
	}

	if newRoot.CellCount() != 1 {
		t.Fatalf(
			"unexpected root cell count: got=%d want=1",
			newRoot.CellCount(),
		)
	}
}

func TestRootSplitPreservesOrdering(t *testing.T) {
	p := newTestPager(t)

	rootPageID := uint32(0)

	rootPg, err := p.GetPage(rootPageID)
	if err != nil {
		t.Fatal(err)
	}

	page.InitPage(rootPg.Data, page.PageTypeLeaf)

	tree := &BTree{
		pager: p,
		root:  rootPageID,
	}

	root, err := LoadNode(p, rootPageID)
	if err != nil {
		t.Fatal(err)
	}

	fillLeafToCapacity(t, root, 1)

	err = tree.splitRootLeaf(root)
	if err != nil {
		t.Fatal(err)
	}

	internalRoot, err := LoadNode(p, tree.root)
	if err != nil {
		t.Fatal(err)
	}

	internalCell, err := internalRoot.InternalCellAt(0)
	if err != nil {
		t.Fatal(err)
	}

	leftNode, err := LoadNode(p, internalCell.Child)
	if err != nil {
		t.Fatal(err)
	}

	rightNode, err := LoadNode(
		p,
		internalRoot.Header().RightChild,
	)
	if err != nil {
		t.Fatal(err)
	}

	leftCells, err := leftNode.LeafCells()
	if err != nil {
		t.Fatal(err)
	}

	rightCells, err := rightNode.LeafCells()
	if err != nil {
		t.Fatal(err)
	}

	if len(leftCells) == 0 {
		t.Fatal("left child empty")
	}

	if len(rightCells) == 0 {
		t.Fatal("right child empty")
	}

	leftMax := leftCells[len(leftCells)-1].Key
	rightMin := rightCells[0].Key

	if leftMax >= rightMin {
		t.Fatalf(
			"cross-page ordering broken: leftMax=%d rightMin=%d",
			leftMax,
			rightMin,
		)
	}
}
