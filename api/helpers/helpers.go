package helpers

import (
	"errors"
	"log"

	"cekkustomer.com/pkg/aws"
)

var ErrPermission = errors.New("permission denied")

func CheckAccountAdmin(username *string) error {
	resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*username)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	targetValue := "admin"
	found := false
	for _, value := range resp {
		if value == targetValue {
			found = true
			break
		}
	}

	if !found {
		return ErrPermission
	} else {
		log.Printf("Found %s\n", targetValue)
	}

	return nil
}
