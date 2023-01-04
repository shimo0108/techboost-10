package main

import (
	"encoding/json"
	"fmt"
	"hello-world/pkg/hospital/api"
	"hello-world/pkg/hospital/database"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := database.NewDB()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "db error",
		}, fmt.Errorf("failed to db error: %w", err)
	}
	params := &api.Params{
		DB: db,
	}

	service := api.NewHospitalService(params)
	hospitals, err := service.ListHospitalsByMunicipality(request.PathParameters)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "hospitals error",
		}, fmt.Errorf("failed to list hospitals by municipality: %w", err)
	}
	if len(hospitals) == 0 {
		return events.APIGatewayProxyResponse{
			Body: "hospitals is not found",
		}, nil
	}
	json, err := json.Marshal(hospitals)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "json error",
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       string(json),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
