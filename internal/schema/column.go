package schema

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type Column struct {
	Name         string
	Type         ColumnType
	IsPrimaryKey bool
}

func NewColumn(name string, typ ColumnType, isPK bool) (*Column, error) {
	colInfo := &Column{
		Name:         name,
		Type:         typ,
		IsPrimaryKey: isPK,
	}
	return colInfo, colInfo.Validate()
}

func (c *Column) Validate() error {
	if c.Name == "" {
		return errors.New("column name can't be empty")
	}

	if len(c.Name) >= math.MaxUint16 {
		return errors.New("column name length exceeds max size")
	}

	if !c.Type.IsValid() {
		return errors.New("invalid type specified for column")
	}

	return nil
}

func (c *Column) Serialize() ([]byte, error) {
	buf := make([]byte, 0)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(c.Name)))
	buf = append(buf, []byte(c.Name)...)
	buf = append(buf, byte(c.Type))
	if c.IsPrimaryKey {
		buf = append(buf, 1)
	} else {
		buf = append(buf, 0)
	}
	return buf, nil
}

func (c *Column) Deserialize(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	nCnt := int(binary.BigEndian.Uint16(data))

	if len(data) < 2+nCnt+2 {
		return 0, io.ErrUnexpectedEOF
	}
	c.Name = string(data[2 : 2+nCnt])
	c.Type = ColumnType(data[2+nCnt])
	c.IsPrimaryKey = data[2+nCnt+1] == 1
	return 2 + nCnt + 2, nil
}
