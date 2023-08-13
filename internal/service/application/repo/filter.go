package repo

import (
	"github.com/nash567/GoSentinel/internal/service/application/model"
	dbHelper "github.com/nash567/GoSentinel/pkg/db/helper"
)

const (
	applicationTable string = "applications"
)

func buildQuery(filter *model.Filter) (string, []interface{}) {
	sqlQuery := getApplication
	if filter == nil {
		return sqlQuery, nil
	}
	query, params := buildFilter(filter)

	if query != "" && filter.Email != nil {
		sqlQuery += " WHERE " + query
	}

	return sqlQuery, params
}
func buildFilter(filter *model.Filter) (string, []interface{}) {
	if filter == nil {
		return "", nil
	}

	f := &dbHelper.Filters{}

	if filter.Email != nil {
		f.AppendInFilter(applicationTable, "email", toInterfaceArr(filter.Email)...)
	}

	if filter.ID != nil {
		f.AppendInFilter(applicationTable, "id", toInterfaceArr(filter.ID)...)
	}
	return f.Query(dbHelper.LogicalOperatorAnd, false), f.Params()
}
func toInterfaceArr[T int | string](v []T) []interface{} {
	out := make([]interface{}, 0, len(v))
	for _, i := range v {
		out = append(out, i)
	}
	return out
}
