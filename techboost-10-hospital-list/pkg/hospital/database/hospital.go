package database

import (
	"fmt"
	"hello-world/pkg/hospital/entity"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// postgres接続情報
	dbFormat = "host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable"
)

func NewDB() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf(dbFormat, host, user, password, name)
	return gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
}

func MultiGetHospitalsByMunicipality(db *gorm.DB, municipality string) entity.Hospitals {
	var hospitals entity.Hospitals
	db.Where("municipality = ?", municipality).Find(&hospitals)
	return hospitals
}

func UpsertHospitals(db *gorm.DB, hospitals entity.Hospitals) {
	db.Debug().Clauses(clause.OnConflict{DoNothing: true}).Create(&hospitals)
}
