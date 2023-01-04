package main

import (
	"encoding/json"
	"hello-world/pkg/csv"
	"hello-world/pkg/hospital/api"
	"hello-world/pkg/hospital/database"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := database.NewDB()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	csvClient := csv.NewCSVClient()
	params := &api.Params{
		DB:  db,
		CSV: csvClient,
	}
	service := api.NewHospitalService(params)
	hospitals, err := service.UpsertHospitalInformation()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	msg, err := json.Marshal(hospitals)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       string(msg),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
