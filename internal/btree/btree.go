package btree

import (
	"errors"
	"fmt"

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
		fmt.Printf("Checking key: %d from root: %d\n", key, childId)
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
		return errors.New("duplicate key")
	}

	cell := &page.Cell{
		Key:   key,
		Value: value,
	}

	if !leaf.page.HasSpace(cell.Size()) {
		return errors.New("insufficient space! leaf split required")
	}
	return leaf.InsertLeafCell(pos, cell)
}
