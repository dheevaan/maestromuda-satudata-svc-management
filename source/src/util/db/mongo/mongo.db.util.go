package db

import (
	"data-management/src/model"
	"strings"
)

func GetSortValue(requestPagination model.Request_Pagination) (order int) {
	order = -1
	if strings.ToLower(requestPagination.Order) == "asc" {
		order = 1
	}
	return
}

func GetSkipAndLimit(requestPagination model.Request_Pagination) (skip int64, limit int64) {
	skip = requestPagination.Page
	if skip > 0 {
		skip--
	}
	skip *= requestPagination.Size

	limit = requestPagination.Size
	if limit == 0 {
		limit = 10
	}
	return
}
