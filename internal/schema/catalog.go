package schema

import (
	"encoding/binary"
	"errors"
	"io"
	"maps"
	"slices"
)

type Catalog struct {
	tables         map[string]*Table
	nextRootPageID uint32
}

func NewCatalog() *Catalog {
	return &Catalog{
		tables:         make(map[string]*Table),
		nextRootPageID: 1,
	}
}

func (c *Catalog) CreateTable(table *Table) error {
	if err := table.Validate(); err != nil {
		return err
	}

	for tName := range c.tables {
		if tName == table.Name {
			return errors.New("table already exists: " + tName)
		}
	}

	table.RootPageID = c.GetNextRootPageID()
	c.nextRootPageID++
	c.tables[table.Name] = table
	return nil
}

func (c *Catalog) GetTable(name string) *Table {
	return c.tables[name]
}

func (c *Catalog) TableExists(name string) bool {
	_, ok := c.tables[name]
	return ok
}

func (c *Catalog) ListTables() []*Table {
	return slices.Collect(maps.Values(c.tables))
}

func (c *Catalog) TableCount() int {
	return len(c.tables)
}

func (c *Catalog) GetNextRootPageID() uint32 {
	return c.nextRootPageID
}

func (c *Catalog) Serialize() ([]byte, error) {
	buf := make([]byte, 0)
	buf = binary.BigEndian.AppendUint32(buf, uint32(c.TableCount()))
	for _, tab := range c.tables {
		tabBuf, err := tab.Serialize()
		if err != nil {
			return nil, err
		}
		buf = append(buf, tabBuf...)
	}
	return buf, nil
}

func (c *Catalog) Deserialize(data []byte) error {
	if len(data) < 4 {
		return io.ErrUnexpectedEOF
	}
	nTabs := binary.BigEndian.Uint32(data)
	c.tables = make(map[string]*Table)
	c.nextRootPageID = 0
	offset := 4
	for range nTabs {
		tab := &Table{}
		o, err := tab.Deserialize(data[offset:])
		if err != nil {
			return err
		}
		c.tables[tab.Name] = tab
		c.nextRootPageID = max(c.nextRootPageID, tab.RootPageID)
		offset += o
	}
	c.nextRootPageID++
	return nil
}

func (c *Catalog) Clear() {
	c.tables = make(map[string]*Table)
	c.nextRootPageID = 1
}
