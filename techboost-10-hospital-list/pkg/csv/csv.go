package csv

import (
	"hello-world/pkg/hospital/entity"
	"net/http"

	"github.com/gocarina/gocsv"
)

type CSVClient interface {
	ListHospitalsInformation() (entity.Hospitals, error)
}

type csvClient struct{}

func NewCSVClient() CSVClient {
	return &csvClient{}
}

func (c *csvClient) ListHospitalsInformation() (entity.Hospitals, error) {
	const csvURL = "https://www.opendata.metro.tokyo.lg.jp/fukushihoken/130001_shinryoukensa.csv"
	resp, err := http.Get(csvURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var hospitals entity.Hospitals
	if err := gocsv.Unmarshal(resp.Body, &hospitals); err != nil {
		return nil, err
	}
	return hospitals, nil
}
