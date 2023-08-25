package orm

import (
	"gobit-demo/internal/pagination"

	"gorm.io/gorm"
)

func Pagination(p *pagination.Pagination) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize)
	}
}

func Exists(g *gorm.DB) (bool, error) {
	var result bool
	if err := g.Select("true").Scan(&result).Error; err != nil {
		return false, err
	}
	return result, nil
}
