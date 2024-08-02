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

func transformMeaningAndSummary(original map[string]string) []models.ChapterText {
	var transformed []models.ChapterText
	for lang, text := range original {
		transformed = append(transformed, models.ChapterText{
			Language: lang,
			Text:     text,
		})
	}
	return transformed
}

func createChapterTable(svc *dynamodb.DynamoDB, tableName string) error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ChapterNumber"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ChapterNumber"),
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

	err = svc.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return err
	}

	return nil
}

func ProcessChapters(svc *dynamodb.DynamoDB) {
	tableName := "ChaptersTable"
	fmt.Print("stating")
	exists, err := checkTableExists(svc, tableName)
	if err != nil {
		log.Fatalf("Failed to check if table exists: %s", err)
	}

	if !exists {
		err = createChapterTable(svc, tableName)
		if err != nil {
			log.Fatalf("Failed to create table: %s", err)
		}
	}

	var path = "/Users/user/Documents/hp/aws-local/dataset/chapter"

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

			// Temporary struct to unmarshal the original format
			var temp struct {
				ChapterNumber   int               `json:"chapter_number"`
				VersesCount     int               `json:"verses_count"`
				Name            string            `json:"name"`
				Translation     string            `json:"translation"`
				Transliteration string            `json:"transliteration"`
				Meaning         map[string]string `json:"meaning"`
				Summary         map[string]string `json:"summary"`
			}

			err = json.Unmarshal(data, &temp)
			if err != nil {
				return err
			}

			// Transform meaning and summary to the new format
			chapter := models.Chapter{
				ChapterNumber:   temp.ChapterNumber,
				VersesCount:     temp.VersesCount,
				Name:            temp.Name,
				Translation:     temp.Translation,
				Transliteration: temp.Transliteration,
				Meaning:         transformMeaningAndSummary(temp.Meaning),
				Summary:         transformMeaningAndSummary(temp.Summary),
			}

			// Ensure ChapterNumber is set correctly
			if chapter.ChapterNumber == 0 {
				return fmt.Errorf("ChapterNumber is required and cannot be zero")
			}

			av, err := dynamodbattribute.MarshalMap(chapter)
			if err != nil {
				log.Fatalf("Got error marshalling new chapter item: %s", err)
			}

			av["ChapterNumber"] = &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%d", chapter.ChapterNumber)),
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

			fmt.Print("inserted successfully", chapter.ChapterNumber)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", path, err)
		return
	}
}
