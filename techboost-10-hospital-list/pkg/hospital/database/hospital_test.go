package database

import (
	"hello-world/pkg/hospital/entity"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDbMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	mockDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	return mockDB, mock, err
}

func TestUpsertHospitals(t *testing.T) {
	name := "テスト"
	phoneNumber := "0000000000"
	address := "テスト区テスト町"
	postalCode := "000-0000"
	municipality := "渋谷区"
	timestamp := time.Now()
	db, mock, err := NewDbMock()
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(*testing.T, sqlmock.Sqlmock)
		hospitals entity.Hospitals
	}{
		{
			name: "success",
			setup: func(t *testing.T, dbmocks sqlmock.Sqlmock) {
				dbmocks.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "hospitals" ("name","phone_number","address","postal_code","municipality","created_at","updated_at") 
						VALUES ($1,$2,$3,$4,$5,$6,$7') ON CONFLICT DO NOTHING`)).
					WithArgs(name, phoneNumber, address, postalCode, municipality, timestamp, timestamp).
					WillReturnRows(
						sqlmock.NewRows([]string{"name", "phone_number", "address", "postal_code", "municipality", "created_at", "updated_at"}).AddRow(name, phoneNumber, address, postalCode, municipality, timestamp, timestamp))
			},
			hospitals: entity.Hospitals{
				{
					Name:         name,
					PhoneNumber:  phoneNumber,
					Address:      address,
					PostalCode:   postalCode,
					Municipality: municipality,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, mock)
			UpsertHospitals(db, tt.hospitals)
		})
	}
}

func TestMultiGetHospitalsByMunicipality(t *testing.T) {
	db, mock, err := NewDbMock()
	require.NoError(t, err)
	now := time.Now()
	tests := []struct {
		name         string
		setup        func(*testing.T, sqlmock.Sqlmock)
		expect       entity.Hospitals
		municipality string
	}{
		{
			name: "success",
			setup: func(t *testing.T, dbmocks sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "phone_number", "address", "postal_code", "municipality", "created_at", "updated_at"}).
					AddRow("テスト", "00000000000", "テスト区テスト町", "000-0000", "渋谷区", now, now)
				dbmocks.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "hospitals" WHERE municipality = $1`)).
					WithArgs("渋谷区").
					WillReturnRows(rows)
			},
			expect: entity.Hospitals{
				{
					Name:         "テスト",
					PhoneNumber:  "00000000000",
					Address:      "テスト区テスト町",
					PostalCode:   "000-0000",
					Municipality: "渋谷区",
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			municipality: "渋谷区",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, mock)
			actual := MultiGetHospitalsByMunicipality(db, tt.municipality)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
