package pager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/harishtpj/indiansql/internal/consts"
)

func TestOpenMemoryPager(t *testing.T) {
	p, err := Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory pager: %v", err)
	}

	if p.nPages != 0 {
		t.Fatalf("expected 0 pages, got %d", p.nPages)
	}

	if p.file == nil {
		t.Fatal("expected storage to be initialized")
	}

	err = p.Close()
	if err != nil {
		t.Fatalf("close failed: %v", err)
	}
}

func TestGetPageCreatesNewPage(t *testing.T) {
	p, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	page, err := p.GetPage(0)
	if err != nil {
		t.Fatalf("GetPage failed: %v", err)
	}

	if page == nil {
		t.Fatal("expected page, got nil")
	}

	if len(page.Data) != consts.PageSize {
		t.Fatalf("expected page size %d, got %d", consts.PageSize, len(page.Data))
	}

	if p.nPages != 1 {
		t.Fatalf("expected nPages=1, got %d", p.nPages)
	}
}

func TestGetPageCache(t *testing.T) {
	p, _ := Open(":memory:")
	defer p.Close()

	p1, _ := p.GetPage(0)
	p2, _ := p.GetPage(0)

	if p1 != p2 {
		t.Fatal("pager did not cache page")
	}
}

func TestPagerFlushPage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	p, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}

	page, err := p.GetPage(0)
	if err != nil {
		t.Fatal(err)
	}

	copy(page.Data[:4], []byte("TEST"))
	page.Dirty = true

	err = p.FlushPage(0)
	if err != nil {
		t.Fatalf("FlushPage failed: %v", err)
	}

	if page.Dirty {
		t.Fatal("page should not be dirty after flush")
	}

	p.Close()

	p2, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer p2.Close()

	page2, err := p2.GetPage(0)
	if err != nil {
		t.Fatal(err)
	}

	if string(page2.Data[:4]) != "TEST" {
		t.Fatalf("expected persisted data, got %q", page2.Data[:4])
	}
}

func TestFlushAll(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "flushall.db")

	p, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}

	for i := uint32(0); i < 3; i++ {
		pg, _ := p.GetPage(i)
		pg.Data[0] = byte(i)
		pg.Dirty = true
	}

	err = p.FlushAll()
	if err != nil {
		t.Fatalf("FlushAll failed: %v", err)
	}

	for _, pg := range p.pages {
		if pg.Dirty {
			t.Fatal("page still dirty after FlushAll")
		}
	}

	p.Close()

	// reopen and verify data
	p2, _ := Open(path)
	defer p2.Close()

	for i := uint32(0); i < 3; i++ {
		pg, _ := p2.GetPage(i)
		if pg.Data[0] != byte(i) {
			t.Fatalf("data mismatch on page %d", i)
		}
	}
}

func TestCloseFlushesDirtyPages(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "closeflush.db")

	p, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}

	pg, _ := p.GetPage(0)
	copy(pg.Data[:5], []byte("HELLO"))
	pg.Dirty = true

	err = p.Close()
	if err != nil {
		t.Fatalf("close failed: %v", err)
	}

	// reopen to confirm flush happened
	p2, _ := Open(path)
	defer p2.Close()

	pg2, _ := p2.GetPage(0)

	if string(pg2.Data[:5]) != "HELLO" {
		t.Fatal("Close() did not flush dirty pages")
	}
}

func TestPagerFileGrowth(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "growth.db")

	p, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}

	pg, err := p.GetPage(5)
	if err != nil {
		t.Fatal(err)
	}

	if p.nPages != 6 {
		t.Fatalf("expected nPages=6 got %d", p.nPages)
	}
	pg.Dirty = true

	if err := p.Close(); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	expected := int64(6 * consts.PageSize)
	if info.Size() != expected {
		t.Fatalf("expected file size %d got %d", expected, info.Size())
	}
}
