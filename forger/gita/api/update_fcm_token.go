package api

import (
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func UpdateFCMToken(request events.APIGatewayProxyRequest, svc *dynamodb.DynamoDB) events.APIGatewayProxyResponse {
	type RequestBody struct {
		FCMToken string `json:"fcm_token"`
	}

	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error extracting email: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	body, err := decodeAndUnmarshal[RequestBody](request)
	if err != nil {
		log.Printf("Error decoding and unmarshalling request body: %s", err)
		return responseBuilder(0, nil, "Bad Request", "Failed to parse request body")
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fcm_token": {
				S: aws.String(body.FCMToken),
			},
			":updated_at": {
				S: aws.String(time.Now().Format(time.RFC3339)), // or any other desired format
			},
		},
		TableName:        aws.String("User"),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set fcm_token = :fcm_token, updated_at = :updated_at"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Printf("Error updating FCM token: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to update FCM token")
	}

	return responseBuilder(1, nil, "Success", "Token Updated Successfully")
}
