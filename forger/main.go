package main

import (
	"forger/gita"
	"forger/imager"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if strings.Contains(request.Path, "/gita") {
		return gita.GitaHandler(request), nil
	}

	if strings.Contains(request.Path, "/imager") {
		return imager.ImagerHandler(request), nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "No Path Found",
		StatusCode: http.StatusInternalServerError,
	}, nil

}

func main() {

	// svc := dynamodb.New(db.DB())

	// migrations.CreateUserActivityTable(svc)

	lambda.Start(handler)

}
