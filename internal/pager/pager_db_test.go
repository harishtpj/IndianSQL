package pager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/harishtpj/indiansql/internal/consts"
	"github.com/harishtpj/indiansql/internal/page"
)

func TestDatabaseHeaderPersistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	p, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}

	pg, err := p.GetPage(0)
	if err != nil {
		t.Fatal(err)
	}

	err = page.InitDBHeader(pg.Data, consts.PageSize)
	if err != nil {
		t.Fatal(err)
	}

	pg.Dirty = true

	if err := p.Close(); err != nil {
		t.Fatal(err)
	}

	p2, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}

	pg2, err := p2.GetPage(0)
	if err != nil {
		t.Fatal(err)
	}

	h, err := page.DecodeDBHeader(pg2.Data)
	if err != nil {
		t.Fatal(err)
	}

	if h.PageSize != consts.PageSize {
		t.Fatalf("header not persisted")
	}

	p2.Close()

	_, err = os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
}
