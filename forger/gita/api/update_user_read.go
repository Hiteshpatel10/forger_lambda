package api

import (
	"fmt"
	"forger/gita/models"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func UpdateUserRead(request events.APIGatewayProxyRequest, svc *dynamodb.DynamoDB) events.APIGatewayProxyResponse {

	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error extracting email: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	updateRead, err := decodeAndUnmarshal[models.UpdateRead](request)
	if err != nil {
		log.Printf("Error decoding and unmarshalling request body: %s", err)
		return responseBuilder(0, nil, "Bad Request", "Failed to parse request body")
	}

	// Fetch current user data from DynamoDB
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("User"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		log.Printf("Error fetching user data: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to fetch user data")
	}

	if result.Item == nil {
		return responseBuilder(0, nil, "Not Found", "User not found")
	}

	// Unmarshal the user's data into a User struct
	var user models.User
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		log.Printf("Error unmarshalling user data: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to unmarshal user data")
	}

	user.UpdatedAt = time.Now().Format(time.RFC3339)
	for i := range user.Reads {
		if user.Reads[i].Chapter == updateRead.ChapterNo {
			// Check if the verse number is already present
			verseExists := false
			for _, verse := range user.Reads[i].Verses {
				if verse == updateRead.VerseNo {
					verseExists = true
					break
				}
			}

			// Append the verse number if it does not exist
			if !verseExists {
				user.Reads[i].Verses = append(user.Reads[i].Verses, updateRead.VerseNo)
				sort.Ints(user.Reads[i].Verses)
				user.LastRead = fmt.Sprintf("BG%d.%d", updateRead.ChapterNo, updateRead.VerseNo)
				user.Reads[i].Progress = len(user.Reads[i].Verses) * 100 / models.GitaChapters[i]
			}
			break
		}
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Printf("Error marshalling updated user data: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to marshal updated user data")
	}

	// Update the 'Item' map with the email attribute
	av["email"] = &dynamodb.AttributeValue{
		S: aws.String(user.Email),
	}

	// Create PutItem input to update the DynamoDB record
	input := &dynamodb.PutItemInput{
		TableName: aws.String("User"),
		Item:      av,
	}

	// Put item into DynamoDB to update the user's record
	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Error updating user data: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to update user data")
	}

	// Return successful response
	return responseBuilder(1, user, "Success", "User reads updated successfully")
}
