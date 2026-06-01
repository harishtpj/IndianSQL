package schema

type ColumnType uint8

const (
	ColumnTypeNumeric ColumnType = iota
	ColumnTypeVarchar
	MaxColumnType
)

func (ct ColumnType) String() string {
	switch ct {
	case ColumnTypeNumeric:
		return "NUMERIC"
	case ColumnTypeVarchar:
		return "VARCHAR"
	default:
		return "UNKNOWN"
	}
}

func (ct ColumnType) IsValid() bool {
	return ct < MaxColumnType
}
