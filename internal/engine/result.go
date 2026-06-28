package engine

import (
	"github.com/harishtpj/indiansql/internal/row"
	"github.com/harishtpj/indiansql/internal/schema"
)

type Result interface {
	isResult()
}

type EmptyResult struct{}

type MessageResult struct {
	Message string
}

type TableResult struct {
	Table *schema.Table
	Rows  []*row.Row
}

type SchemaResult struct {
	Table *schema.Table
}

type TablesResult struct {
	Tables []*schema.Table
}

type ExitResult struct{}

func (*MessageResult) isResult() {}
func (*TableResult) isResult()   {}
func (*SchemaResult) isResult()  {}
func (*TablesResult) isResult()  {}
func (*EmptyResult) isResult()   {}
func (*ExitResult) isResult()    {}
