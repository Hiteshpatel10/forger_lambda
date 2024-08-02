package api

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// responseBuilder builds a response for the API Gateway
func responseBuilder(status int, result interface{}, message, errorMessage string, header ...map[string]string) events.APIGatewayProxyResponse {
	response := map[string]interface{}{
		"status":  status,
		"result":  result,
		"message": message,
	}

	if status == 0 {
		response["error_message"] = errorMessage
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(responseBody),
		}
	}

	var headers map[string]string
	if len(header) > 0 && header[0] != nil {
		headers = header[0]
	} else {
		headers = map[string]string{
			"Content-Type": "application/json",
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(responseBody),
	}
}
