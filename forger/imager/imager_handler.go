package imager

import (
	"forger/imager/api"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func ImagerHandler(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	if request.Path == "/imager/image-generator" {
		return api.BuildIcon(request)

	}

	return events.APIGatewayProxyResponse{
		Body:       "No Gita Path Found",
		StatusCode: http.StatusInternalServerError,
	}
}
