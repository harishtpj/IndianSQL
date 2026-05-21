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
	clear(n.page.Data())
	page.InitPage(n.page.Data(), page.PageTypeLeaf)
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
	if err := n.WriteLeafCells(left); err != nil {
		return nil, err
	}

	rtNode, err := LoadNode(n.pager, newPageId)
	if err != nil {
		return nil, err
	}
	page.InitPage(rtNode.page.Data(), page.PageTypeLeaf)
	if err := rtNode.WriteLeafCells(right); err != nil {
		return nil, err
	}

	return &SplitResult{
		Separator: right[0].Key,
		RightPage: newPageId,
	}, nil
}
