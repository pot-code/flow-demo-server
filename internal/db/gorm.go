package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormClient(db *sql.DB, logger zerolog.Logger) *gorm.DB {
	gd, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: NewGormLogger(logger),
	})
	if err != nil {
		panic(fmt.Errorf("error connecting mysql database: %w", err))
	}
	return gd
}

type gormLogger struct {
	l zerolog.Logger
}

func NewGormLogger(l zerolog.Logger) logger.Interface {
	return &gormLogger{l}
}

func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &gormLogger{
		l: g.l.Level(parseGormLoggerLevel(level)),
	}
}

func (g *gormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	g.l.Info().Msgf(msg, data...)
}

func (g *gormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	g.l.Warn().Msgf(msg, data...)
}

func (g *gormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	// g.l.Error().Msgf(msg, data...)
}

func (g *gormLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	g.l.Trace().
		Dur("elapsed", elapsed).
		Str("sql", sql).
		Int64("rows", rows).
		Msg("")
}

func parseGormLoggerLevel(l logger.LogLevel) zerolog.Level {
	switch l {
	case logger.Info:
		return zerolog.InfoLevel
	case logger.Warn:
		return zerolog.WarnLevel
	case logger.Error:
		return zerolog.ErrorLevel
	case logger.Silent:
		return zerolog.NoLevel
	}
	panic("unknown log level")
}
