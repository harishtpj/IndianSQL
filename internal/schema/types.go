package schema

type ColumnType uint8

const (
	ColumnTypeInteger ColumnType = iota
	ColumnTypeVarchar
	ColumnTypeBoolean
	MaxColumnType
)

func (ct ColumnType) String() string {
	switch ct {
	case ColumnTypeInteger:
		return "INTEGER"
	case ColumnTypeVarchar:
		return "VARCHAR"
	case ColumnTypeBoolean:
		return "BOOLEAN"
	default:
		return "UNKNOWN"
	}
}

func (ct ColumnType) IsValid() bool {
	return ct < MaxColumnType
}
