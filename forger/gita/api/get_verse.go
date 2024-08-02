package api

import (
	"forger/db"
	"forger/gita/models"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetVerse(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	svc := dynamodb.New(db.DB())

	// Extract the verse ID from the request path
	parts := strings.Split(request.Path, "/")
	verseID := parts[len(parts)-1]

	// Set up query input parameters
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Verses"),
		KeyConditionExpression: aws.String("ID = :verse_id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":verse_id": {
				S: aws.String(verseID),
			},
		},
	}

	// Perform the query operation
	result, err := svc.Query(input)
	if err != nil {
		log.Printf("Error querying DynamoDB: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to fetch verse: DynamoDB query error")
	}

	if len(result.Items) == 0 {
		log.Printf("Verse with ID %s not found", verseID)
		return responseBuilder(0, nil, "Not Found", "Verse not found")
	}

	// Unmarshal the result into a Verse struct
	var verse models.Verse
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &verse)
	if err != nil {
		log.Printf("Failed to unmarshal item: %v", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to fetch verse: Unmarshal error")
	}

	return responseBuilder(1, verse, "Success", "")
}
