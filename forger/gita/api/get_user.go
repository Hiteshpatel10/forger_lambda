package api

import (
	"forger/gita/models"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetUser(request events.APIGatewayProxyRequest, svc *dynamodb.DynamoDB) events.APIGatewayProxyResponse {

	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error scanning DynamoDB table: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "un authorised user")
	}
	query := &dynamodb.QueryInput{
		TableName:              aws.String("User"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(email),
			},
		},
	}

	// Perform the query operation
	result, err := svc.Query(query)
	if err != nil {
		log.Printf("Error scanning DynamoDB table: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to fetch chapter: DynamoDB scan error")
	}

	if len(result.Items) == 0 {
		log.Printf("Chapter with number %s not found", "")
		return responseBuilder(0, nil, "Not Found", "Chapter not found")
	}

	var user models.User
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		log.Printf("Failed to unmarshal item: %v", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to parse chapter data: Unmarshal error")
	}

	return responseBuilder(1, user, "Success", "")
}
