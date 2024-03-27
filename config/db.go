package config

import (
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/exception"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewPostgres(dbSource string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{})
	exception.FatalIfNeeded(err, "cannot connect to db")

	return db
}
