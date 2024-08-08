package api

import (
	"fmt"
	"forger/db"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const tableName = "UserActivity"

// RequestBody structure for incoming request

// UpdateUserActivity updates the user's activity for a specific date
func UpdateUserActivity(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	type RequestBody struct {
		ChapterNo int `json:"chapter_no"`
		VerseNo   int `json:"verse_no"`
	}
	svc := dynamodb.New(db.DB())

	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error extracting email: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	body, err := decodeAndUnmarshal[RequestBody](request)
	if err != nil {
		log.Printf("Error decoding and unmarshalling request body: %s", err)
		return responseBuilder(0, request.Body, "Bad Request", err.Error())
	}

	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")

	updateExpression := "SET #activity = list_append(if_not_exists(#activity, :empty_list), :new_activity)"
	activity := map[string]*dynamodb.AttributeValue{
		"chapter_no": {N: aws.String(fmt.Sprintf("%d", body.ChapterNo))},
		"verse_no":   {N: aws.String(fmt.Sprintf("%d", body.VerseNo))},
	}
	conditionExpression := "not contains(#activity, :existing_activity)"
	existingActivity := map[string]*dynamodb.AttributeValue{
		"chapter_no": {N: aws.String(fmt.Sprintf("%d", body.ChapterNo))},
		"verse_no":   {N: aws.String(fmt.Sprintf("%d", body.VerseNo))},
	}

	updateParams := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(email)},
			"date":  {S: aws.String(formattedDate)},
		},
		UpdateExpression:    aws.String(updateExpression),
		ConditionExpression: aws.String(conditionExpression),
		ExpressionAttributeNames: map[string]*string{
			"#activity": aws.String("activity"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":new_activity":      {L: []*dynamodb.AttributeValue{{M: activity}}},
			":empty_list":        {L: []*dynamodb.AttributeValue{}},
			":existing_activity": {M: existingActivity},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err = svc.UpdateItem(updateParams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			log.Printf("Activity already exists: %s", err)
			return responseBuilder(1, nil, "No Operation", "Activity already exists")
		}
		log.Printf("Error updating item in DynamoDB: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", err.Error())
	}

	return responseBuilder(1, nil, "Success", "User activity updated successfully")
}
