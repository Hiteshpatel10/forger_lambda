package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func responseBuilder(status int, result interface{}, message, errorMessage string) events.APIGatewayProxyResponse {
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
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(responseBody),
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}
}

func headerHandler(headers map[string]string) (string, error) {

	authHeader := headers["Authorization"]
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	email := tokenParts[1]

	// Validate email format
	if !isValidEmail(email) {
		return "", errors.New("invalid email format")
	}

	return email, nil
}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func decodeAndUnmarshal[T any](request events.APIGatewayProxyRequest) (T, error) {
	var targetStruct T

	// Decode the base64-encoded body
	decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return targetStruct, fmt.Errorf("error decoding base64 body: %w", err)
	}

	// Unmarshal the decoded body into the target struct
	err = json.Unmarshal(decodedBody, &targetStruct)
	if err != nil {
		return targetStruct, fmt.Errorf("error unmarshalling JSON body to struct: %w", err)
	}

	return targetStruct, nil
}
