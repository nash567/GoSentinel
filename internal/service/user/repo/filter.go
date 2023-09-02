package repo

import (
	"github.com/nash567/GoSentinel/internal/service/user/model"
	dbHelper "github.com/nash567/GoSentinel/pkg/db/helper"
)

const (
	usersTable = "users"
)

func buildFilter(filter *model.Filter) (string, []interface{}) {
	if filter == nil {
		return "", nil
	}

	f := &dbHelper.Filters{}

	if len(filter.Email) > 0 {
		f.AppendInFilter(usersTable, "email", toInterfaceArr((filter.Email))...)
	}
	if len(filter.ID) > 0 {
		f.AppendInFilter(usersTable, "id", toInterfaceArr(filter.ID)...)
	}
	return f.Query(dbHelper.LogicalOperatorAnd, true), f.Params()
}
func toInterfaceArr[T int | string](v []T) []interface{} {
	out := make([]interface{}, 0, len(v))
	for _, i := range v {
		out = append(out, i)
	}
	return out
}
