package btree

import (
	"reflect"
	"testing"

	"github.com/harishtpj/indiansql/internal/page"
	"github.com/harishtpj/indiansql/internal/pager"
)

func newTestTree(t *testing.T) *BTree {
	t.Helper()

	p, err := pager.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		p.Close()
	})

	pg, err := p.GetPage(0)
	if err != nil {
		t.Fatal(err)
	}

	page.InitPage(pg.Data, page.PageTypeLeaf)

	return NewTree(p, 0)
}

func populateTree(t *testing.T, tree *BTree, keys ...uint64) {
	t.Helper()

	for _, k := range keys {
		if err := tree.Insert(k, []byte{byte(k)}); err != nil {
			t.Fatalf("insert %d: %v", k, err)
		}
	}
}

func collectKeys(tree *BTree) ([]uint64, error) {
	cur, err := NewCursor(tree)
	if err != nil {
		return nil, err
	}

	keys := make([]uint64, 0)

	for !cur.IsFinished() {
		k, err := cur.Key()
		if err != nil {
			return nil, err
		}

		keys = append(keys, k)
		cur.Next()
	}

	return keys, nil
}

func requireKeys(t *testing.T, tree *BTree, want []uint64) {
	t.Helper()

	got, err := collectKeys(tree)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("keys = %v, want %v", got, want)
	}
}

func TestCursorDeleteFirst(t *testing.T) {
	tree := newTestTree(t)
	populateTree(t, tree, 1, 2, 3, 4, 5)

	cur, err := NewCursor(tree)
	if err != nil {
		t.Fatal(err)
	}

	if err := cur.Delete(); err != nil {
		t.Fatal(err)
	}

	k, err := cur.Key()
	if err != nil {
		t.Fatal(err)
	}

	if k != 2 {
		t.Fatalf("cursor at %d, want 2", k)
	}

	requireKeys(t, tree, []uint64{2, 3, 4, 5})
}

func TestCursorDeleteMiddle(t *testing.T) {
	tree := newTestTree(t)
	populateTree(t, tree, 1, 2, 3, 4, 5)

	cur, err := NewCursor(tree)
	if err != nil {
		t.Fatal(err)
	}

	cur.Next() // 2
	cur.Next() // 3

	if err := cur.Delete(); err != nil {
		t.Fatal(err)
	}

	k, err := cur.Key()
	if err != nil {
		t.Fatal(err)
	}

	if k != 4 {
		t.Fatalf("cursor at %d, want 4", k)
	}

	requireKeys(t, tree, []uint64{1, 2, 4, 5})
}

func TestCursorDeleteLast(t *testing.T) {
	tree := newTestTree(t)
	populateTree(t, tree, 1, 2, 3)

	cur, err := NewCursor(tree)
	if err != nil {
		t.Fatal(err)
	}

	cur.Next()
	cur.Next() // 3

	if err := cur.Delete(); err != nil {
		t.Fatal(err)
	}

	if !cur.IsFinished() {
		t.Fatal("cursor should be finished")
	}

	requireKeys(t, tree, []uint64{1, 2})
}

func TestCursorDeleteAll(t *testing.T) {
	tree := newTestTree(t)
	populateTree(t, tree, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9)

	cur, err := NewCursor(tree)
	if err != nil {
		t.Fatal(err)
	}

	for !cur.IsFinished() {
		if err := cur.Delete(); err != nil {
			t.Fatal(err)
		}
	}

	requireKeys(t, tree, []uint64{})

	cnt, err := cur.Count()
	if err != nil {
		t.Fatal(err)
	}

	if cnt != 0 {
		t.Fatalf("count = %d, want 0", cnt)
	}
}

func TestCursorDeleteCount(t *testing.T) {
	tree := newTestTree(t)
	populateTree(t, tree, 0, 1, 2, 3, 4)

	cur, err := NewCursor(tree)
	if err != nil {
		t.Fatal(err)
	}

	if err := cur.Delete(); err != nil {
		t.Fatal(err)
	}

	if err := cur.Delete(); err != nil {
		t.Fatal(err)
	}

	requireKeys(t, tree, []uint64{2, 3, 4})

	cnt, err := cur.Count()
	if err != nil {
		t.Fatal(err)
	}

	if cnt != 3 {
		t.Fatalf("count = %d, want 3", cnt)
	}
}

func TestBTreeDeleteMissingKey(t *testing.T) {
	tree := newTestTree(t)
	populateTree(t, tree, 1)

	if err := tree.Delete(2); err == nil {
		t.Fatal("expected error deleting missing key")
	}

	requireKeys(t, tree, []uint64{1})
}
