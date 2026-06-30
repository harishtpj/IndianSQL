package btree

import (
	"fmt"

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

func (t *BTree) Delete(key uint64) error {
	leaf, pos, found, err := t.FindLeaf(key)
	if err != nil {
		return err
	}

	if !found {
		return apperrors.ErrKeyNotFound
	}

	return leaf.DeleteLeafCell(pos)
}

func (t *BTree) splitRootLeaf(root *Node) error {
	oldRoot := t.root
	root, err := LoadNode(t.pager, oldRoot)
	if err != nil {
		return fmt.Errorf("could not split root: %w", err)
	}

	rightPageId := t.pager.PageCount()
	split, err := root.SplitLeaf(rightPageId)
	if err != nil {
		return fmt.Errorf("could not split root: %w", err)
	}

	newRootPageId := t.pager.PageCount()
	newRtPg, err := t.pager.GetPage(newRootPageId)
	if err != nil {
		return fmt.Errorf("could not split root: %w", err)
	}

	page.InitPage(newRtPg.Data, page.PageTypeInternal)
	newRoot, err := LoadNode(t.pager, newRootPageId)
	if err != nil {
		return fmt.Errorf("could not split root: %w", err)
	}
	if err := newRoot.InternalInsert(split.Separator, oldRoot, split.RightPage); err != nil {
		return fmt.Errorf("could not split root: %w", err)
	}

	t.root = newRootPageId
	return nil
}
