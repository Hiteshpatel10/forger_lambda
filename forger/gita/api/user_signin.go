package api

import (
	"forger/gita/models"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CreateUserBody struct {
	DisplayName *string `json:"display_name,omitempty"`
	ProfileURL  *string `json:"sam ,omitempty"`
	FCMToken    *string `json:"fcm_token,omitempty"`
}

func CreateUser(request events.APIGatewayProxyRequest, svc *dynamodb.DynamoDB) events.APIGatewayProxyResponse {
	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error extracting email: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	createUserBody, err := decodeAndUnmarshal[CreateUserBody](request)
	if err != nil {
		log.Printf("Error decoding and unmarshalling request body: %s", err)
		return responseBuilder(0, nil, "Bad Request", "Failed to parse request body")
	}

	// Set up query input parameters
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
		log.Printf("Error querying DynamoDB table: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to query DynamoDB table")
	}

	if len(result.Items) > 0 {
		var user models.User
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
		if err != nil {
			log.Printf("Failed to unmarshal item: %v", err)
			return responseBuilder(0, nil, "Internal Server Error", "Failed to parse user data: Unmarshal error")
		}
		log.Printf("User with email %s already exists", email)
		return responseBuilder(1, user, "User Already Exists", "User with this email already exists")
	}

	var reads []models.Read
	for i := 1; i <= 18; i++ {
		read := models.Read{
			Chapter:  i,
			Verses:   []int{},
			Progress: 0,
		}
		reads = append(reads, read)
	}

	user := models.User{
		Email:       email,
		FCMToken:    createUserBody.FCMToken,
		DisplayName: createUserBody.DisplayName,
		ProfileURL:  createUserBody.ProfileURL,
		Reads:       reads,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Fatalf("Got error marshalling new user item: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to marshal new user data")
	}

	// Create PutItem input
	input := &dynamodb.PutItemInput{
		TableName: aws.String("User"),
		Item:      av,
	}

	// Put item into DynamoDB
	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to put item into DynamoDB")
	}

	// Return successful response
	return responseBuilder(1, user, "Success", "User created successfully")
}
