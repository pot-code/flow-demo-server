package util

import (
	"fmt"
	"gobit-demo/internal/pagination"

	"gorm.io/gorm"
)

type existenceResult struct {
	Exist bool
}

func GormCheckExistence(db *gorm.DB, queryFn func(tx *gorm.DB) *gorm.DB) (bool, error) {
	var result existenceResult
	sql := db.ToSQL(queryFn)
	if err := db.Raw(fmt.Sprintf("SELECT EXISTS(%s) as exist", sql)).Scan(&result).Error; err != nil {
		return false, fmt.Errorf("select exists: %w", err)
	}
	return result.Exist, nil
}

type gormQuery[T any] interface {
	Limit(int) T
	Offset(int) T
}

func GormPaginator[Q gormQuery[Q]](
	query Q,
	pagination *pagination.Pagination,
) Q {
	return query.Limit(pagination.PageSize).Offset((pagination.Page - 1) * pagination.PageSize)
}
