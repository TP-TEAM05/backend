package database

import (
	"recofiit/models"

	"github.com/getsentry/sentry-go"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func GetModels() []interface{} {
	return []interface{}{
		&models.Car{},
		&models.Controller{},
		&models.Firmware{},
		&models.Measurement{},
		&models.Sensor{},
		&models.Session{},
		&models.CarController{},
		&models.CarSession{},
		&models.CarSessionController{},
		&models.ControllerInstance{},
	}
}

func AutoMigrate() {
	db, db_err := GetDB().DB()

	if db_err != nil {
		panic(db_err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	m, m_err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)

	if m_err != nil {
		panic(m_err)
	}

	m.Up()
}
