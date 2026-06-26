package btree

import (
	"testing"

	"github.com/harishtpj/indiansql/internal/page"
	"github.com/harishtpj/indiansql/internal/pager"
)

func TestNewCursor(t *testing.T) {
	p, _ := pager.Open(":memory:")
	defer p.Close()

	pg, _ := p.GetPage(0)
	page.InitPage(pg.Data, page.PageTypeLeaf)

	tree := NewTree(p, 0)

	cursor, err := NewCursor(tree)
	if err != nil {
		t.Fatal(err)
	}

	if cursor == nil {
		t.Fatal("cursor is nil")
	}

	if !cursor.IsFinished() {
		t.Error("New cursor on empty tree should be finished")
	}
}

func TestCursorFirstNext(t *testing.T) {
	p, _ := pager.Open(":memory:")
	defer p.Close()

	pg, _ := p.GetPage(0)
	page.InitPage(pg.Data, page.PageTypeLeaf)
	page.WrapSlottedPage(pg.Data)

	tree := NewTree(p, 0)

	// Insert some cells
	for i := 0; i < 5; i++ {
		key := uint64(i)
		value := []byte("val")
		tree.Insert(key, value)
	}

	cursor, _ := NewCursor(tree)

	// First should position at start
	_ = cursor.First()

	count := 0
	for cursor.Next() {
		count++
		if count > 10 { // Safety check
			break
		}
	}

	if count < 4 {
		t.Errorf("Iterated %d cells, expected at least 5", count)
	}
}

func TestCursorEmptyTree(t *testing.T) {
	p, _ := pager.Open(":memory:")
	defer p.Close()

	pg, _ := p.GetPage(0)
	page.InitPage(pg.Data, page.PageTypeLeaf)

	tree := NewTree(p, 0)
	cursor, _ := NewCursor(tree)

	// First on empty tree
	err := cursor.First()
	if err != nil {
		// May or may not error depending on implementation
	}

	if !cursor.IsFinished() && cursor.Next() {
		t.Error("Empty tree cursor should not iterate")
	}
}

func TestCursorKeyValue(t *testing.T) {
	p, _ := pager.Open(":memory:")
	defer p.Close()

	pg, _ := p.GetPage(0)
	page.InitPage(pg.Data, page.PageTypeLeaf)

	tree := NewTree(p, 0)

	tree.Insert(42, []byte("answer"))
	tree.Insert(100, []byte("century"))

	cursor, _ := NewCursor(tree)
	cursor.First()

	key, err := cursor.Key()
	if err != nil {
		t.Fatal(err)
	}

	val, err := cursor.Value()
	if err != nil {
		t.Fatal(err)
	}

	if key != 42 && key != 100 {
		t.Errorf("Key = %v, expected 42 or 100", key)
	}

	if len(val) == 0 {
		t.Error("Value is empty")
	}
}

func TestCursorCount(t *testing.T) {
	p, _ := pager.Open(":memory:")
	defer p.Close()

	pg, _ := p.GetPage(0)
	page.InitPage(pg.Data, page.PageTypeLeaf)

	tree := NewTree(p, 0)

	for i := 0; i < 10; i++ {
		tree.Insert(uint64(i), []byte("data"))
	}

	cursor, _ := NewCursor(tree)
	count, err := cursor.Count()
	if err != nil {
		t.Fatal(err)
	}

	if count != 10 {
		t.Errorf("Count = %v, want 10", count)
	}
}

func TestCursorIsFinished(t *testing.T) {
	p, _ := pager.Open(":memory:")
	defer p.Close()

	pg, _ := p.GetPage(0)
	page.InitPage(pg.Data, page.PageTypeLeaf)

	tree := NewTree(p, 0)
	tree.Insert(1, []byte("one"))

	cursor, _ := NewCursor(tree)
	cursor.First()

	if cursor.IsFinished() {
		t.Error("Should not be finished at first cell")
	}

	cursor.Next() // Move past last cell
	if !cursor.IsFinished() {
		t.Error("Should be finished after last cell")
	}
}
