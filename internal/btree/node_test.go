package btree

import (
	"testing"

	"github.com/harishtpj/indiansql/internal/page"
)

func TestNodeLoadAndType(t *testing.T) {
	p := newTestPager(t)

	n := createLeafNode(t, p, 0)

	if !n.IsLeaf() {
		t.Fatal("expected leaf node")
	}

	if n.IsInternal() {
		t.Fatal("should not be internal")
	}
}

func TestNodeKeyAtLeaf(t *testing.T) {
	p := newTestPager(t)

	n := createLeafNode(t, p, 0)

	cell := &page.Cell{
		Key:   10,
		Value: []byte("a"),
	}

	if err := n.page.AppendCell(cell); err != nil {
		t.Fatal(err)
	}

	key, err := n.KeyAt(0)
	if err != nil {
		t.Fatal(err)
	}

	if key != 10 {
		t.Fatalf("expected key 10 got %d", key)
	}
}

func TestNodeFindPosition(t *testing.T) {
	p := newTestPager(t)

	n := createLeafNode(t, p, 0)

	keys := []uint64{10, 20, 30}

	for _, k := range keys {
		cell := &page.Cell{
			Key:   k,
			Value: []byte("x"),
		}

		if err := n.page.AppendCell(cell); err != nil {
			t.Fatal(err)
		}
	}

	pos, found, err := n.FindPosition(20)
	if err != nil {
		t.Fatal(err)
	}

	if !found || pos != 1 {
		t.Fatalf("expected pos=1 found=true")
	}

	pos, found, err = n.FindPosition(25)
	if err != nil {
		t.Fatal(err)
	}

	if found || pos != 2 {
		t.Fatalf("expected insert position 2")
	}
}

func TestNodeChildAt(t *testing.T) {
	p := newTestPager(t)

	n := createInternalNode(t, p, 0)

	hdr := n.page.Header()
	hdr.RightChild = 99
	n.page.SetHeader(hdr)

	ic := &page.InternalCell{
		Key:   50,
		Value: 123,
	}

	if err := n.page.AppendCell(ic); err != nil {
		t.Fatal(err)
	}

	child, err := n.ChildAt(0)
	if err != nil {
		t.Fatal(err)
	}

	if child != 123 {
		t.Fatalf("expected child 123 got %d", child)
	}

	right, err := n.ChildAt(1)
	if err != nil {
		t.Fatal(err)
	}

	if right != 99 {
		t.Fatalf("expected right child 99")
	}
}

func TestNodeInsertLeafCellSorted(t *testing.T) {
	p := newTestPager(t)

	n := createLeafNode(t, p, 0)

	keys := []uint64{30, 10, 20}

	for _, k := range keys {
		pos, found, err := n.FindPosition(k)
		if err != nil {
			t.Fatal(err)
		}

		if found {
			t.Fatal("duplicate key")
		}

		err = n.InsertLeafCell(pos, &page.Cell{
			Key:   k,
			Value: []byte("x"),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	expected := []uint64{10, 20, 30}

	for i, want := range expected {
		got, err := n.KeyAt(i)
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
