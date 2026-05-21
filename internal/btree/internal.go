package btree

import (
	"fmt"

	"github.com/harishtpj/indiansql/internal/apperrors"
	"github.com/harishtpj/indiansql/internal/page"
)

func (t *BTree) FindLeaf(key uint64) (*Node, int, bool, error) {
	childId := t.root

	for {
		node, err := LoadNode(t.pager, childId)
		if err != nil {
			return nil, 0, false, err
		}

		pos, found, err := node.FindPosition(key)
		if err != nil {
			return nil, 0, false, err
		}

		if node.IsLeaf() {
			return node, pos, found, nil
		}

		childId, err = node.ChildAt(pos)
		if err != nil {
			return nil, 0, false, err
		}
	}
}

func (n *Node) InternalInsert(key uint64, leftChild, rightChild uint32) error {
	if !n.IsInternal() {
		return apperrors.ErrInvalidPageType
	}

	if n.CellCount() != 0 {
		panic("internal insert: cannot insert into already used node")
	}

	cell := &page.InternalCell{
		Key:   key,
		Child: leftChild,
	}
	if err := n.page.AppendCell(cell); err != nil {
		return fmt.Errorf("internal insert: %w", err)
	}
	pgHdr := n.Header()
	pgHdr.RightChild = rightChild
	n.SetHeader(pgHdr)
	return nil
}

func (n *Node) InternalCellAt(idx uint16) (*page.InternalCell, error) {
	if !n.IsInternal() {
		return nil, apperrors.ErrInvalidPageType
	}

	return page.DecodeInternalCell(n.page.GetCellRaw(idx))
}
