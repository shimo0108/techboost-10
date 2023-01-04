package api

import (
	csvMock "hello-world/pkg/csv/mock"
	"hello-world/pkg/hospital/entity"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type mocks struct {
	sqlMock sqlmock.Sqlmock
	csv     *csvMock.MockCSVClient
}

func NewDbMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	mockDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	return mockDB, mock, err
}

func newMocks(ctrl *gomock.Controller, sqlMock sqlmock.Sqlmock) *mocks {
	return &mocks{
		sqlMock: sqlMock,
		csv:     csvMock.NewMockCSVClient(ctrl),
	}
}

func Test_hospitalService_ListHospitalsByMunicipality(t *testing.T) {
	name := "テスト"
	phoneNumber := "0000000000"
	address := "テスト区テスト町"
	postalCode := "000-0000"
	municipality := "渋谷区"
	timestamp := time.Now()
	db, sqlMock, err := NewDbMock()
	require.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(*testing.T, *mocks)
		searchMap map[string]string
		expect    entity.Hospitals
		wantErr   bool
	}{
		{
			name:    "failed to search map not found",
			setup:   func(t *testing.T, mocks *mocks) {},
			expect:  nil,
			wantErr: true,
		},
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks) {
				rows := sqlmock.NewRows([]string{"name", "phone_number", "address", "postal_code", "municipality", "created_at", "updated_at"}).
					AddRow("テスト", "0000000000", "テスト区テスト町", "000-0000", "渋谷区", timestamp, timestamp)
				mocks.sqlMock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "hospitals" WHERE municipality = $1`)).
					WithArgs("渋谷区").
					WillReturnRows(rows)
			},
			searchMap: map[string]string{
				"municipality": "渋谷区",
			},
			expect: entity.Hospitals{
				{
					Name:         name,
					PhoneNumber:  phoneNumber,
					Address:      address,
					PostalCode:   postalCode,
					Municipality: municipality,
					CreatedAt:    timestamp,
					UpdatedAt:    timestamp,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mocks := newMocks(ctrl, sqlMock)
			params := &Params{
				DB: db,
			}
			tt.setup(t, mocks)
			service := NewHospitalService(params)
			actual, err := service.ListHospitalsByMunicipality(tt.searchMap)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func Test_hospitalService_UpsertHospitalInformation(t *testing.T) {
	name := "テスト"
	phoneNumber := "0000000000"
	address := "テスト区テスト町"
	postalCode := "000-0000"
	municipality := "渋谷区"
	timestamp := time.Now()
	db, sqlMock, err := NewDbMock()
	require.NoError(t, err)
	tests := []struct {
		name    string
		setup   func(*testing.T, *mocks)
		want    entity.Hospitals
		expect  entity.Hospitals
		wantErr bool
	}{
		{
			name: "failed to csv error",
			setup: func(t *testing.T, mocks *mocks) {
				mocks.csv.EXPECT().ListHospitalsInformation().Return(nil, assert.AnError)
			},
			expect:  nil,
			wantErr: true,
		},
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks) {
				mocks.csv.EXPECT().ListHospitalsInformation().Return(entity.Hospitals{
					{
						Name:         "テスト",
						PhoneNumber:  "00000000000",
						Address:      "テスト区テスト町",
						PostalCode:   "000-0000",
						Municipality: "渋谷区",
						CreatedAt:    timestamp,
						UpdatedAt:    timestamp,
					},
				}, nil)
				mocks.sqlMock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "hospitals" ("name","phone_number","address","postal_code","municipality","created_at","updated_at") 
						VALUES ($1,$2,$3,$4,$5,$6,$7') ON CONFLICT DO NOTHING`)).
					WithArgs(name, phoneNumber, address, postalCode, municipality, timestamp, timestamp).
					WillReturnRows(
						sqlmock.NewRows([]string{"name", "phone_number", "address", "postal_code", "municipality", "created_at", "updated_at"}).AddRow(name, phoneNumber, address, postalCode, municipality, timestamp, timestamp))
			},
			expect: entity.Hospitals{
				{
					Name:         "テスト",
					PhoneNumber:  "00000000000",
					Address:      "テスト区テスト町",
					PostalCode:   "000-0000",
					Municipality: "渋谷区",
					CreatedAt:    timestamp,
					UpdatedAt:    timestamp,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mocks := newMocks(ctrl, sqlMock)
			params := &Params{
				DB:  db,
				CSV: mocks.csv,
			}
			tt.setup(t, mocks)
			service := NewHospitalService(params)
			actual, err := service.UpsertHospitalInformation()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
