package util

import (
	"gobit-demo/internal/pagination"

	"gorm.io/gorm"
)

type GormUtil struct{}

func (s *GormUtil) Pagination(p *pagination.Pagination) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize)
	}
}

func (s *GormUtil) Exists(g *gorm.DB) (bool, error) {
	var result bool
	if err := g.Select("true").Scan(&result).Error; err != nil {
		return false, err
	}
	return result, nil
}
