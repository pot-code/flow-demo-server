package util

import (
	"gobit-demo/internal/pagination"

	"gorm.io/gorm"
)

type gormWrap struct {
	*gorm.DB
}

func NewGormWrap(db *gorm.DB) *gormWrap {
	return &gormWrap{db}
}

func (g *gormWrap) Paginate(p *pagination.Pagination) *gorm.DB {
	return g.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize)
}

func (g *gormWrap) Exists() (bool, error) {
	var result bool
	if err := g.Select("true").Scan(&result).Error; err != nil {
		return false, err
	}
	return result, nil
}
