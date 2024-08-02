package migrations

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func checkTableExists(svc *dynamodb.DynamoDB, tableName string) (bool, error) {
	_, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		if isResourceNotFoundException(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isResourceNotFoundException(err error) bool {
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
		return true
	}
	return false
}
