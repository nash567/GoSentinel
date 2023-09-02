package helper

import (
	"fmt"
	"strconv"
	"strings"
)

// sql operators
type Operator string

const (
	OperatorIn  Operator = " in "
	OperatorOr  Operator = " or "
	OperatorAnd Operator = " and "
)

type LogicalOperator int32

const (
	LogicalOperatorAnd LogicalOperator = iota
	LogicalOperatorOr
)

//go:generate enumer -type=LogicalOperator -text -json -transform=snake -output=enum_logical_operator_gen.go

func (o Operator) ToString() string {
	return string(o)
}

type Filters struct {
	clause          []string      // clause eg ->> columnName in(?,?,?,>)
	params          []interface{} // params will be the values of ?
	queryParamCount int           // count of next query param
}

func (f *Filters) Query(log LogicalOperator, appendWhere bool) string {
	q := ""
	if len(f.clause) != 0 {
		operator := OperatorAnd
		if log == LogicalOperatorOr {
			operator = OperatorOr
		}

		if appendWhere {
			q += " WHERE "
		}
		q = fmt.Sprintf("%s(%s)", q, strings.Join(f.clause, operator.ToString()))
	}

	return q
}

func (f *Filters) Params() []interface{} {
	return f.params
}

func (f *Filters) buildMultiParamCondition(column string, operator Operator, params ...interface{}) {
	if len(params) == 0 {
		return
	}
	q := ""
	for i := 0; i < len(params); i++ {
		if f.queryParamCount == 0 {
			f.queryParamCount = 1

		}
		q += fmt.Sprintf(",$%s", strconv.Itoa(f.queryParamCount))

		f.queryParamCount++
	}
	if q[0] == ',' {
		f.clause = append(f.clause, column+operator.ToString()+"("+q[1:]+")")
	} else {
		f.clause = append(f.clause, column+operator.ToString()+"("+q+")")
	}
	f.params = append(f.params, params...)
}

func (f *Filters) AppendInFilter(table, col string, params ...interface{}) {
	if len(params) > 0 {
		f.buildMultiParamCondition(getColumnName(table, col), OperatorIn, params...)
	}
}

func getColumnName(table, column string) string {
	if strings.TrimSpace(table) != "" {
		return fmt.Sprintf("%s.%s", table, column)
	}
	return column
}
