package btree

import (
	"errors"

	"github.com/harishtpj/indiansql/internal/page"
	"github.com/harishtpj/indiansql/internal/pager"
)

type Node struct {
	pager  *pager.Pager
	pageID uint32
	page   *page.SlottedPage
}

func LoadNode(p *pager.Pager, id uint32) (*Node, error) {
	pg, err := p.GetPage(id)
	if err != nil {
		return nil, err
	}

	return &Node{
		pager:  p,
		pageID: id,
		page:   page.WrapSlottedPage(pg.Data),
	}, nil
}

func (n *Node) IsLeaf() bool {
	return n.page.Header().Type == page.PageTypeLeaf
}

func (n *Node) IsInternal() bool {
	return n.page.Header().Type == page.PageTypeInternal
}

func (n *Node) CellCount() int {
	return int(n.page.Header().CellCount)
}

func (n *Node) KeyAt(idx int) (uint64, error) {
	raw := n.page.GetCellRaw(uint16(idx))

	switch {
	case n.IsLeaf():
		cell, err := page.DecodeCell(raw)
		if err != nil {
			return 0, err
		}
		return cell.Key, nil
	case n.IsInternal():
		cell, err := page.DecodeInternalCell(raw)
		if err != nil {
			return 0, err
		}
		return cell.Key, nil
	}

	return 0, errors.New("invalid cell type")
}

func (n *Node) FindPosition(key uint64) (int, bool, error) {
	l, h := 0, n.CellCount()

	for l < h {
		mid := (l + h) / 2
		iKey, err := n.KeyAt(mid)
		if err != nil {
			return 0, false, err
		} else if iKey == key {
			return mid, true, nil
		} else if key < iKey {
			h = mid
		} else {
			l = mid + 1
		}
	}

	return l, false, nil
}

func (n *Node) ChildAt(idx int) (uint32, error) {
	if !n.IsInternal() {
		return 0, errors.New("leaf node has no children")
	}

	hdr := n.page.Header()
	cnt := int(hdr.CellCount)

	if idx < 0 || idx > cnt {
		return 0, errors.New("child index is out of bounds")
	}

	if idx == cnt {
		return hdr.RightChild, nil
	}

	ic, err := page.DecodeInternalCell(n.page.GetCellRaw(uint16(idx)))
	if err != nil {
		return 0, err
	}
	return ic.Value, nil
}

func (n *Node) InsertLeafCell(pos int, cell *page.Cell) error {
	if !n.IsLeaf() {
		return errors.New("cannot insert leaf cell into internal node")
	}

	return n.page.InsertCell(uint16(pos), cell)
}
