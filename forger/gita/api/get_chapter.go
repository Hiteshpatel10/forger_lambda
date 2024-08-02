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

func GetChapter(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	svc := dynamodb.New(db.DB())

	// Extract the chapter number from the request path
	parts := strings.Split(request.Path, "/")
	chapterNumber := parts[len(parts)-1]

	// Set up scan input parameters
	input := &dynamodb.ScanInput{
		TableName:        aws.String("ChaptersTable"),
		FilterExpression: aws.String("chapter_number = :chapter_number"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":chapter_number": {
				N: aws.String(chapterNumber),
			},
		},
	}

	// Perform the scan operation
	result, err := svc.Scan(input)
	if err != nil {
		log.Printf("Error scanning DynamoDB table: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to fetch chapter: DynamoDB scan error")
	}

	if len(result.Items) == 0 {
		log.Printf("Chapter with number %s not found", chapterNumber)
		return responseBuilder(0, nil, "Not Found", "Chapter not found")
	}

	// Unmarshal the result into a Chapter struct
	var chapter models.Chapter
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &chapter)
	if err != nil {
		log.Printf("Failed to unmarshal item: %v", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to parse chapter data: Unmarshal error")
	}

	return responseBuilder(1, chapter, "Success", "")
}
