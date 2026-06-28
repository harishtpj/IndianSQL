package schema

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
)

type Table struct {
	Name            string
	Columns         []*Column
	PrimaryKeyIndex int
	RootPageID      uint32
}

func NewTable(name string, columns []*Column) (*Table, error) {
	tab := &Table{
		Name:    name,
		Columns: columns,
	}

	for idx, val := range columns {
		if val.IsPrimaryKey {
			tab.PrimaryKeyIndex = idx
		}
	}

	return tab, tab.Validate()
}

func (t *Table) Validate() error {
	if t.Name == "" {
		return errors.New("table name can't be empty")
	}

	if len(t.Name) >= math.MaxUint16 {
		return errors.New("table name length exceeds max size")
	}

	if len(t.Columns) == 0 {
		return errors.New("table with zero column(s) is not allowed")
	}

	if len(t.Columns) >= math.MaxUint16 {
		return errors.New("too many columns in table (exceeded max)")
	}

	hasPK := false
	for i, col := range t.Columns {
		if col.IsPrimaryKey {
			if hasPK {
				return errors.New("table already has a Primary Key")
			} else if i != t.PrimaryKeyIndex {
				return errors.New("table's primary key index mismatch")
			} else {
				hasPK = true
			}
		}

		if err := col.Validate(); err != nil {
			return err
		}
	}

	if !hasPK {
		return errors.New("table must have one Primary Key")
	}

	return nil
}

func (t *Table) GetColumn(name string) *Column {
	for _, col := range t.Columns {
		if col.Name == name {
			return col
		}
	}
	return nil
}

func (t *Table) GetColumnIndex(name string) int {
	for i, col := range t.Columns {
		if col.Name == name {
			return i
		}
	}
	return -1
}

func (t *Table) GetPrimaryKeyColumn() *Column {
	return t.Columns[t.PrimaryKeyIndex]
}

func (t *Table) ColumnCount() int {
	return len(t.Columns)
}

func (t *Table) String() string {
	var repr strings.Builder
	fmt.Fprintf(&repr, "Table: %s(", t.Name)
	for i, col := range t.Columns {
		fmt.Fprint(&repr, col.Name, " ", col.Type.String())
		if i == t.PrimaryKeyIndex {
			repr.WriteString(" PRIMARY KEY")
		}

		if i+1 != t.ColumnCount() {
			repr.WriteString(", ")
		}
	}
	return repr.String() + ")"
}

func (t *Table) Serialize() ([]byte, error) {
	buf := make([]byte, 0)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(t.Name)))
	buf = append(buf, []byte(t.Name)...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(t.ColumnCount()))
	for _, col := range t.Columns {
		colBuf, err := col.Serialize()
		if err != nil {
			return nil, err
		}
		buf = append(buf, colBuf...)
	}
	buf = binary.BigEndian.AppendUint32(buf, uint32(t.RootPageID))
	return buf, nil
}

func (t *Table) Deserialize(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	nCnt := int(binary.BigEndian.Uint16(data))
	offset := 2 + nCnt + 2
	if len(data) < offset {
		return 0, io.ErrUnexpectedEOF
	}
	t.Name = string(data[2 : nCnt+2])
	nCols := int(binary.BigEndian.Uint16(data[2+nCnt:]))
	t.Columns = make([]*Column, nCols)
	for i := range nCols {
		t.Columns[i] = &Column{}
		o, err := t.Columns[i].Deserialize(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += o
	}

	if len(data)-offset < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	t.RootPageID = binary.BigEndian.Uint32(data[offset:])
	return offset + 4, nil
}
