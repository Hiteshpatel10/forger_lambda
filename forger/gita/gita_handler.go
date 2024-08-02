package gita

import (
	"forger/db"
	"forger/gita/api"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GitaHandler(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	svc := dynamodb.New(db.DB())

	if strings.Contains(request.Path, "/gita/createUser") {
		return api.CreateUser(request, svc)
	}

	if strings.Contains(request.Path, "/gita/user") {
		return api.GetUser(request, svc)
	}

	if strings.Contains(request.Path, "/gita/updateRead") {
		return api.UpdateUserRead(request, svc)
	}

	if strings.Contains(request.Path, "/gita/chapter") {
		return api.GetChapter(request)
	}

	if strings.Contains(request.Path, "/gita/verse") {
		return api.GetVerse(request)
	}

	return events.APIGatewayProxyResponse{
		Body:       "No Gita Path Found",
		StatusCode: http.StatusInternalServerError,
	}
}
