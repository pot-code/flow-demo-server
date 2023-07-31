package db

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormClient(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
}
