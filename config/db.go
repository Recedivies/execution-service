package config

import (
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/exception"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewPostgres(dbSource string, env string) *gorm.DB {
	var gormLogger logger.Interface = logger.Default.LogMode(logger.Info)
	if env != "development" {
		// For other environments (e.g., production), use a silent logger
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{
		Logger: gormLogger,
	})
	exception.FatalIfNeeded(err, "cannot connect to db")

	return db
}
