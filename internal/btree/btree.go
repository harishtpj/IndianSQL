package btree

import (
	"github.com/harishtpj/indiansql/internal/apperrors"
	"github.com/harishtpj/indiansql/internal/page"
	"github.com/harishtpj/indiansql/internal/pager"
)

type BTree struct {
	pager *pager.Pager
	root  uint32
}

func NewTree(p *pager.Pager, rt uint32) *BTree {
	return &BTree{
		pager: p,
		root:  rt,
	}
}

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

func (t *BTree) Insert(key uint64, value []byte) error {
	leaf, pos, found, err := t.FindLeaf(key)
	if err != nil {
		return err
	}

	if found {
		return apperrors.ErrDuplicateKey
	}

	cell := &page.Cell{
		Key:   key,
		Value: value,
	}

	if !leaf.HasSpaceFor(cell.Size()) {
		return apperrors.ErrNodeFull
	}
	return leaf.InsertLeafCell(pos, cell)
}
