package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var dbSess *session.Session

func init() {
	var err error
	dbSess, err = session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-south-1"),
		},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
	}
}

func DB() *session.Session {
	return dbSess
}

func PrintDBSession() {
	fmt.Println(dbSess)
}
