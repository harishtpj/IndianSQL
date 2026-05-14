package btree

import (
	"testing"

	"github.com/harishtpj/indiansql/internal/page"
	"github.com/harishtpj/indiansql/internal/pager"
)

func newTestPager(t *testing.T) *pager.Pager {
	t.Helper()

	p, err := pager.Open(":memory:")
	if err != nil {
		t.Fatalf("open pager: %v", err)
	}

	return p
}

func createLeafNode(t *testing.T, p *pager.Pager, id uint32) *Node {
	t.Helper()

	pg, err := p.GetPage(id)
	if err != nil {
		t.Fatal(err)
	}

	page.NewSlottedPage(pg.Data, page.PageTypeLeaf)

	n, err := LoadNode(p, id)
	if err != nil {
		t.Fatal(err)
	}

	return n
}

func createInternalNode(t *testing.T, p *pager.Pager, id uint32) *Node {
	t.Helper()

	pg, err := p.GetPage(id)
	if err != nil {
		t.Fatal(err)
	}

	sp := page.NewSlottedPage(pg.Data, page.PageTypeInternal)

	h := sp.Header()
	h.RightChild = 0
	sp.SetHeader(h)

	n, err := LoadNode(p, id)
	if err != nil {
		t.Fatal(err)
	}

	return n
}
