package migrations

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateUserActivityTable(svc *dynamodb.DynamoDB) error {
	tableName := "UserActivity"
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("email"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("date"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("email"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("date"),
				KeyType:       aws.String("RANGE"),
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
