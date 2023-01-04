package api

import (
	"errors"
	"hello-world/pkg/csv"
	"hello-world/pkg/hospital/database"
	"hello-world/pkg/hospital/entity"
	"net/url"

	"gorm.io/gorm"
)

type Params struct {
	DB  *gorm.DB
	CSV csv.CSVClient
}

type hospitalService struct {
	db  *gorm.DB
	csv csv.CSVClient
}

type hospitalServiceServer interface {
	ListHospitalsByMunicipality(searchMap map[string]string) (entity.Hospitals, error)
	UpsertHospitalInformation() (entity.Hospitals, error)
}

func NewHospitalService(p *Params) hospitalServiceServer {
	return &hospitalService{
		db:  p.DB,
		csv: p.CSV,
	}
}

func (h *hospitalService) ListHospitalsByMunicipality(searchMap map[string]string) (entity.Hospitals, error) {
	const key = "municipality"
	municipality, ok := searchMap[key]
	if !ok {
		return nil, errors.New("search params is not found")
	}
	municipality, err := url.QueryUnescape(municipality)
	if err != nil {
		return nil, err
	}
	hospitals := database.MultiGetHospitalsByMunicipality(h.db, municipality)
	return hospitals, nil
}

func (h *hospitalService) UpsertHospitalInformation() (entity.Hospitals, error) {
	hospitals, err := h.csv.ListHospitalsInformation()
	if err != nil {
		return nil, err
	}
	database.UpsertHospitals(h.db, hospitals)
	return hospitals, nil
}
