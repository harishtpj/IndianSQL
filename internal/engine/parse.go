package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/harishtpj/indiansql/internal/row"
	"github.com/harishtpj/indiansql/internal/schema"
)

func parseCreateColumns(cols []string) ([]*schema.Column, error) {
	columns := make([]*schema.Column, 0, len(cols))

	for _, c := range cols {
		parts := strings.Split(c, ":")

		if len(parts) < 2 || len(parts) > 3 {
			return nil, fmt.Errorf("invalid column definition: %q", c)
		}

		col := &schema.Column{
			Name: parts[0],
		}

		switch strings.ToLower(parts[1]) {
		case "numeric", "num":
			col.Type = schema.ColumnTypeNumeric

		case "varchar", "text", "string":
			col.Type = schema.ColumnTypeVarchar

		default:
			return nil, fmt.Errorf("unknown type %q", parts[1])
		}

		if len(parts) == 3 {
			if strings.ToLower(parts[2]) != "pk" {
				return nil, fmt.Errorf("unknown modifier %q", parts[2])
			}
			col.IsPrimaryKey = true
		}

		columns = append(columns, col)
	}

	return columns, nil
}

func parseInsertValues(table *schema.Table, vals []string) ([]row.Value, error) {
	if len(vals) != table.ColumnCount() {
		return nil, fmt.Errorf(
			"expected %d values, got %d",
			table.ColumnCount(),
			len(vals),
		)
	}

	values := make([]row.Value, 0, len(vals))

	for i, str := range vals {
		switch table.Columns[i].Type {

		case schema.ColumnTypeNumeric:
			n, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"column %q expects NUMERIC",
					table.Columns[i].Name,
				)
			}

			values = append(values, row.NewNumericValue(n))

		case schema.ColumnTypeVarchar:
			values = append(values, row.NewVarcharValue(str))

		default:
			return nil, fmt.Errorf(
				"unsupported datatype for column %q",
				table.Columns[i].Name,
			)
		}
	}

	return values, nil
}

func PrintRows(table *schema.Table, rows []*row.Row) {
	for _, col := range table.Columns {
		fmt.Printf("%-20s", col.Name)
	}
	fmt.Println()

	for _, col := range table.Columns {
		fmt.Printf("%-20s", strings.Repeat("-", len(col.Name)))
	}
	fmt.Println()

	for _, r := range rows {
		for _, v := range r.Values {
			fmt.Printf("%-20s", (*v).String())
		}
		fmt.Println()
	}

	fmt.Printf("\n%d row(s)\n", len(rows))
}
