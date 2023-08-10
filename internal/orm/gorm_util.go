package orm

import (
	"gobit-demo/internal/pagination"

	"gorm.io/gorm"
)

type gormWrapper struct {
	*gorm.DB
}

func NewGormWrapper(db *gorm.DB) *gormWrapper {
	return &gormWrapper{db}
}

func (g *gormWrapper) Paginate(p *pagination.Pagination) *gorm.DB {
	return g.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize)
}

func (g *gormWrapper) Exists() (bool, error) {
	var result bool
	if err := g.Select("true").Scan(&result).Error; err != nil {
		return false, err
	}
	return result, nil
}
