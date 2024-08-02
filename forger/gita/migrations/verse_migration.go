package migrations

import (
	"encoding/json"
	"fmt"
	"forger/gita/models"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func parseAndTransformSlokaData(data []byte) (*models.Verse, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	sloka := &models.Verse{
		ID:              raw["_id"].(string),
		Chapter:         int(raw["chapter"].(float64)),
		Verse:           int(raw["verse"].(float64)),
		Slok:            raw["slok"].(string),
		Transliteration: raw["transliteration"].(string),
	}

	commentFields := []string{"tej", "siva", "purohit", "chinmay", "san", "adi", "gambir", "madhav", "anand", "rams", "raman", "abhinav", "sankar", "jaya", "vallabh", "ms", "srid", "dhan", "venkat", "puru", "neel"}

	for _, field := range commentFields {
		if raw[field] != nil {
			commentData := raw[field].(map[string]interface{})
			comment := models.VerseComment{
				Author: commentData["author"].(string),
			}
			languages := []models.VerseLanguage{}
			for _, lang := range []string{"ht", "et", "ec", "hc", "sc"} {
				if text, exists := commentData[lang]; exists {
					languages = append(languages, models.VerseLanguage{
						Language: lang,
						Text:     text.(string),
					})
				}
			}
			comment.Languages = languages
			sloka.Comments = append(sloka.Comments, comment)
		}
	}

	return sloka, nil
}

func createVerseTable(svc *dynamodb.DynamoDB, tableName string) error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		return err
	}

	// Wait until the table is created
	err = svc.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return err
	}

	return nil
}

func ProcessSloks(svc *dynamodb.DynamoDB) {
	tableName := "Verses"
	fmt.Println("Starting to process sloks...")

	exists, err := checkTableExists(svc, tableName)
	if err != nil {
		log.Fatalf("Failed to check if table exists: %s", err)
	}

	if !exists {
		err = createVerseTable(svc, tableName)
		if err != nil {
			log.Fatalf("Failed to create table: %s", err)
		}
	}

	path := "/Users/user/Documents/hp/aws-local/dataset/slok"

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// Read the JSON file
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Parse and transform the JSON data
			sloka, err := parseAndTransformSlokaData(data)
			if err != nil {
				return err
			}

			av, err := dynamodbattribute.MarshalMap(sloka)
			if err != nil {
				log.Fatalf("Got error marshalling new sloka item: %s", err)
			}

			av["ID"] = &dynamodb.AttributeValue{
				S: aws.String(sloka.ID),
			}
			input := &dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item:      av,
			}

			// Insert the item into DynamoDB
			_, err = svc.PutItem(input)
			if err != nil {
				log.Fatalf("Error inserting item into DynamoDB: %s", err)
			}

			fmt.Printf("Processed JSON file: %s\n", path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", path, err)
		return
	}
	fmt.Println("Finished processing sloks.")
}
