package btree

import (
	"errors"

	"github.com/harishtpj/indiansql/internal/page"
)

type Cursor struct {
	tree     *BTree
	pageID   uint32
	cellIdx  int
	finished bool
}

func NewCursor(tree *BTree) (*Cursor, error) {
	if tree == nil {
		return nil, errors.New("tree is nil")
	}

	c := &Cursor{tree: tree}
	err := c.First()
	return c, err
}

func (c *Cursor) First() error {
	c.pageID = c.tree.root

	for {
		node, err := LoadNode(c.tree.pager, c.pageID)
		if err != nil {
			return err
		}

		if node.IsLeaf() {
			c.finished = node.CellCount() == 0
			break
		}

		c.pageID, err = node.ChildAt(0)
		if err != nil {
			return err
		}
	}
	c.cellIdx = 0
	return nil
}

func (c *Cursor) Next() bool {
	if c.finished {
		return false
	}

	node, err := LoadNode(c.tree.pager, c.pageID)
	if err != nil {
		c.finished = true
		return false
	}

	if c.cellIdx+1 < node.CellCount() {
		c.cellIdx++
		return true
	}

	nxtPageId := node.Header().NextLeaf
	if nxtPageId == 0 {
		c.finished = true
		return false
	}

	nxtNode, err := LoadNode(c.tree.pager, nxtPageId)
	if err != nil {
		c.finished = true
		return false
	}

	c.pageID = nxtPageId
	c.cellIdx = 0

	if nxtNode.CellCount() == 0 {
		c.finished = true
		return false
	}

	return true
}

func (c *Cursor) Key() (uint64, error) {
	curCell, err := c.currentCell()
	if err != nil {
		return 0, err
	}
	return curCell.Key, nil
}

func (c *Cursor) Value() ([]byte, error) {
	curCell, err := c.currentCell()
	if err != nil {
		return nil, err
	}
	return curCell.Value, nil
}

func (c *Cursor) IsFinished() bool {
	return c.finished
}

func (c *Cursor) Count() (int, error) {
	cur, err := NewCursor(c.tree)
	if err != nil {
		return 0, err
	}

	cnt := 0
	for !cur.finished {
		cnt++
		cur.Next()
	}
	return cnt, nil
}

func (c *Cursor) Delete() error {
	key, err := c.Key()
	if err != nil {
		return err
	}

	if err := c.tree.Delete(key); err != nil {
		return err
	}

	node, err := LoadNode(c.tree.pager, c.pageID)
	if err != nil {
		return err
	}

	switch {
	case node.CellCount() == 0:
		if node.Header().NextLeaf == 0 {
			c.finished = true
			return nil
		}

		c.pageID = node.Header().NextLeaf
		c.cellIdx = 0

		next, err := LoadNode(c.tree.pager, c.pageID)
		if err != nil {
			return err
		}

		if next.CellCount() == 0 {
			c.finished = true
		}

	case c.cellIdx >= node.CellCount():
		if node.Header().NextLeaf == 0 {
			c.finished = true
			return nil
		}

		c.pageID = node.Header().NextLeaf
		c.cellIdx = 0
	}

	return nil
}

func (c *Cursor) currentCell() (*page.Cell, error) {
	if c.finished {
		return nil, errors.New("cursor is iterated")
	}

	node, err := LoadNode(c.tree.pager, c.pageID)
	if err != nil {
		return nil, err
	}

	cells, err := node.LeafCells()
	if err != nil {
		return nil, err
	}

	return cells[c.cellIdx], nil
}
