package btree

import (
	"testing"

	"github.com/harishtpj/indiansql/internal/page"
)

func TestBTreeFindLeafSingleNode(t *testing.T) {
	p := newTestPager(t)

	n := createLeafNode(t, p, 0)

	err := n.page.InsertCell(&page.Cell{
		Key:   10,
		Value: []byte("x"),
	})
	if err != nil {
		t.Fatal(err)
	}

	tree := NewTree(p, 0)

	leaf, pos, found, err := tree.FindLeaf(10)
	if err != nil {
		t.Fatal(err)
	}

	if !leaf.IsLeaf() {
		t.Fatal("expected leaf")
	}

	if !found || pos != 0 {
		t.Fatalf("wrong search result")
	}
}

func TestBTreeTraversal(t *testing.T) {
	p := newTestPager(t)

	left := createLeafNode(t, p, 1)
	right := createLeafNode(t, p, 2)
	root := createInternalNode(t, p, 0)

	h := root.page.Header()
	h.RightChild = 2
	root.page.SetHeader(h)

	ic := &page.InternalCell{
		Key:   50,
		Value: 1, // left child page
	}

	if err := root.page.InsertCell(ic); err != nil {
		t.Fatal(err)
	}

	tree := NewTree(p, 0)

	leaf, _, _, err := tree.FindLeaf(25)
	if err != nil {
		t.Fatal(err)
	}

	if leaf.pageID != left.pageID {
		t.Fatal("should go to left child")
	}

	leaf, _, _, err = tree.FindLeaf(75)
	if err != nil {
		t.Fatal(err)
	}

	if leaf.pageID != right.pageID {
		t.Fatal("should go to right child")
	}
}
