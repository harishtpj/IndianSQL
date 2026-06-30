package btree

import (
	"errors"
	"fmt"

	"github.com/harishtpj/indiansql/internal/apperrors"
	"github.com/harishtpj/indiansql/internal/page"
)

func (n *Node) InsertLeafCell(pos int, cell *page.Cell) error {
	if !n.IsLeaf() {
		return errors.New("cannot insert leaf cell into internal node")
	}

	return n.page.InsertCell(uint16(pos), cell)
}

func (n *Node) DeleteLeafCell(pos int) error {
	if !n.IsLeaf() {
		return errors.New("cannot delete leaf cell from internal node")
	}

	cells, err := n.LeafCells()
	if err != nil {
		return err
	}

	if pos < 0 || pos >= len(cells) {
		return errors.New("invalid cell position")
	}

	copy(cells[pos:], cells[pos+1:])
	cells = cells[:len(cells)-1]
	return n.WriteLeafCells(cells)
}

func (n *Node) LeafCells() ([]*page.Cell, error) {
	if !n.IsLeaf() {
		return nil, apperrors.ErrInvalidPageType
	}

	pgHdr := n.Header()
	arr := make([]*page.Cell, pgHdr.CellCount)
	for i := range arr {
		rawCell := n.page.GetCellRaw(uint16(i))
		cell, err := page.DecodeCell(rawCell)
		if err != nil {
			return nil, fmt.Errorf("decode leaf cell %d: %w", i, err)
		}
		arr[i] = cell
	}
	return arr, nil
}

func (n *Node) ResetLeaf() {
	if !n.IsLeaf() {
		panic("ResetLeaf called on non-leaf node")
	}

	hdr := n.Header()
	parent := hdr.Parent()
	next := hdr.GetNextLeaf()
	prev := hdr.GetPrevLeaf()

	clear(n.page.Data())
	page.InitPage(n.page.Data(), page.PageTypeLeaf)

	hdr = n.Header()
	hdr.SetParent(parent)
	hdr.SetNextLeaf(next)
	hdr.SetPrevLeaf(prev)
	n.SetHeader(hdr)
}

func (n *Node) WriteLeafCells(cells []*page.Cell) error {
	if !n.IsLeaf() {
		return apperrors.ErrInvalidPageType
	}
	n.ResetLeaf()
	for i, cell := range cells {
		if err := n.page.AppendCell(cell); err != nil {
			return fmt.Errorf("rewrite leaf cell %d: %w", i, err)
		}
	}
	return nil
}

func (n *Node) SplitLeaf(newPageId uint32) (*SplitResult, error) {
	if !n.IsLeaf() {
		return nil, apperrors.ErrInvalidPageType
	}

	cells, err := n.LeafCells()
	if err != nil {
		return nil, err
	}

	if len(cells) < 2 {
		return nil, apperrors.ErrCantSplit
	}

	mid := len(cells) / 2
	left := cells[:mid]
	right := cells[mid:]

	oldHdr := n.Header()
	oldNext := oldHdr.GetNextLeaf()
	oldPrev := oldHdr.GetPrevLeaf()
	parent := oldHdr.Parent()

	if err := n.WriteLeafCells(left); err != nil {
		return nil, err
	}

	leftHdr := n.Header()
	leftHdr.SetParent(parent)
	leftHdr.SetNextLeaf(newPageId)
	leftHdr.SetPrevLeaf(oldPrev)
	n.SetHeader(leftHdr)

	rtNode, err := LoadNode(n.pager, newPageId)
	if err != nil {
		return nil, err
	}
	page.InitPage(rtNode.page.Data(), page.PageTypeLeaf)

	rtHdr := rtNode.Header()
	rtHdr.SetParent(parent)
	rtHdr.SetNextLeaf(oldNext)
	rtHdr.SetPrevLeaf(n.pageID)
	rtNode.SetHeader(rtHdr)

	if err := rtNode.WriteLeafCells(right); err != nil {
		return nil, err
	}

	if oldNext != 0 {
		nxtNode, err := LoadNode(n.pager, oldNext)
		if err != nil {
			return nil, err
		}

		hdr := nxtNode.Header()
		hdr.SetPrevLeaf(newPageId)
		nxtNode.SetHeader(hdr)
	}

	return &SplitResult{
		Separator: right[0].Key,
		RightPage: newPageId,
	}, nil
}
