package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/harishtpj/indiansql/internal/manager"
	"github.com/harishtpj/indiansql/internal/row"
	"github.com/harishtpj/indiansql/internal/schema"
	"github.com/xwb1989/sqlparser"
)

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

func CompilePredicate(expr sqlparser.Expr) (manager.Predicate, error) {
	switch e := expr.(type) {

	case *sqlparser.AndExpr:
		left, err := CompilePredicate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := CompilePredicate(e.Right)
		if err != nil {
			return nil, err
		}

		return func(r *row.Row) (bool, error) {
			ok, err := left(r)
			if err != nil || !ok {
				return ok, err
			}

			return right(r)
		}, nil

	case *sqlparser.OrExpr:
		left, err := CompilePredicate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := CompilePredicate(e.Right)
		if err != nil {
			return nil, err
		}

		return func(r *row.Row) (bool, error) {
			ok, err := left(r)
			if err != nil {
				return false, err
			}

			if ok {
				return true, nil
			}

			return right(r)
		}, nil

	case *sqlparser.ParenExpr:
		return CompilePredicate(e.Expr)

	case *sqlparser.ComparisonExpr:
		return compileComparison(e)

	default:
		return nil, fmt.Errorf("unsupported WHERE expression %T", e)
	}
}

func compileComparison(expr *sqlparser.ComparisonExpr) (manager.Predicate, error) {
	leftCol, ok := expr.Left.(*sqlparser.ColName)
	if !ok {
		return nil, fmt.Errorf("left operand must be a column")
	}

	column := leftCol.Name.String()
	op := expr.Operator

	return func(r *row.Row) (bool, error) {
		value, err := r.GetValueByName(column)
		if err != nil {
			return false, err
		}

		cmp, err := compareValues(value, expr.Right)
		if err != nil {
			return false, err
		}

		switch op {

		case "=":
			return cmp == 0, nil

		case "!=", "<>":
			return cmp != 0, nil

		case "<":
			return cmp < 0, nil

		case "<=":
			return cmp <= 0, nil

		case ">":
			return cmp > 0, nil

		case ">=":
			return cmp >= 0, nil

		default:
			return false, fmt.Errorf("unsupported operator %v", op)
		}
	}, nil
}

func compareValues(val row.Value, expr sqlparser.Expr) (int, error) {
	switch v := val.(type) {

	case *row.IntegerValue:
		switch e := expr.(type) {

		case *sqlparser.SQLVal:
			switch e.Type {

			case sqlparser.IntVal:
				rhs, err := strconv.ParseInt(string(e.Val), 10, 64)
				if err != nil {
					return 0, err
				}

				switch {
				case v.GetInt64() < rhs:
					return -1, nil
				case v.GetInt64() > rhs:
					return 1, nil
				default:
					return 0, nil
				}

			case sqlparser.FloatVal:
				rhs, err := strconv.ParseFloat(string(e.Val), 64)
				if err != nil {
					return 0, err
				}

				lhs := float64(v.GetInt64())

				switch {
				case lhs < rhs:
					return -1, nil
				case lhs > rhs:
					return 1, nil
				default:
					return 0, nil
				}
			}
		}

	case *row.NumericValue:
		switch e := expr.(type) {

		case *sqlparser.SQLVal:
			switch e.Type {

			case sqlparser.IntVal:
				rhs, err := strconv.ParseInt(string(e.Val), 10, 64)
				if err != nil {
					return 0, err
				}

				lhs := v.GetFloat()

				switch {
				case lhs < float64(rhs):
					return -1, nil
				case lhs > float64(rhs):
					return 1, nil
				default:
					return 0, nil
				}

			case sqlparser.FloatVal:
				rhs, err := strconv.ParseFloat(string(e.Val), 64)
				if err != nil {
					return 0, err
				}

				switch {
				case v.GetFloat() < rhs:
					return -1, nil
				case v.GetFloat() > rhs:
					return 1, nil
				default:
					return 0, nil
				}
			}
		}

	case *row.VarcharValue:
		sqlVal, ok := expr.(*sqlparser.SQLVal)
		if !ok || sqlVal.Type != sqlparser.StrVal {
			return 0, fmt.Errorf("cannot compare VARCHAR with %T", expr)
		}

		return strings.Compare(v.GetString(), string(sqlVal.Val)), nil

	case *row.BoolValue:
		boolVal, ok := expr.(sqlparser.BoolVal)
		if !ok {
			return 0, fmt.Errorf("cannot compare BOOLEAN with %T", expr)
		}

		lhs := v.GetBool()
		rhs := bool(boolVal)

		switch {
		case lhs == rhs:
			return 0, nil
		case !lhs && rhs:
			return -1, nil
		default:
			return 1, nil
		}
	}

	return 0, fmt.Errorf("cannot compare %T with %T", val, expr)
}
