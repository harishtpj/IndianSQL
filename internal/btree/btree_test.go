package btree

import (
	"errors"
	"testing"

	"github.com/harishtpj/indiansql/internal/apperrors"
	"github.com/harishtpj/indiansql/internal/page"
)

func TestBTreeFindLeafSingleNode(t *testing.T) {
	p := newTestPager(t)

	n := createLeafNode(t, p, 0)

	err := n.page.AppendCell(&page.Cell{
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
		Child: 1, // left child page
	}

	if err := root.page.AppendCell(ic); err != nil {
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

func TestBTreeInsertSorted(t *testing.T) {
	p := newTestPager(t)

	createLeafNode(t, p, 0)

	tree := NewTree(p, 0)

	keys := []uint64{30, 10, 20}

	for _, k := range keys {
		err := tree.Insert(k, []byte("x"))
		if err != nil {
			t.Fatal(err)
		}
	}

	leaf, _, _, err := tree.FindLeaf(20)
	if err != nil {
		t.Fatal(err)
	}

	expected := []uint64{10, 20, 30}

	for i, want := range expected {
		got, err := leaf.KeyAt(i)
		if err != nil {
			t.Fatal(err)
		}

		if got != want {
			t.Fatalf(
				"wrong key at index %d: got=%d want=%d",
				i,
				got,
				want,
			)
		}
	}
}

func TestBTreeInsertDuplicate(t *testing.T) {
	p := newTestPager(t)

	createLeafNode(t, p, 0)

	tree := NewTree(p, 0)

	err := tree.Insert(10, []byte("a"))
	if err != nil {
		t.Fatal(err)
	}

	err = tree.Insert(10, []byte("b"))
	if err == nil {
		t.Fatal("expected duplicate key error")
	}
	if !errors.Is(err, apperrors.ErrDuplicateKey) {
		t.Fatalf("expected ErrDuplicateKey, got %v", err)
	}
}
