package aws

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type AwsCognito struct {
	cognitoClient       *cognito.Client
	cognitoClientId     string
	cognitoClientSecret string
	cognitoPoolId       string
}

func NewCognitoClient(config *aws.Config, clientId, clientSecret, poolId string) *AwsCognito {
	client := cognito.NewFromConfig(*config)

	return &AwsCognito{
		cognitoClient:       client,
		cognitoClientId:     clientId,
		cognitoClientSecret: clientSecret,
		cognitoPoolId:       poolId,
	}
}

func computeSecretHash(clientSecret, userName, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(userName + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *AwsCognito) SignUp(email, password string) (string, error) {
	input := &cognito.SignUpInput{
		ClientId: aws.String(c.cognitoClientId),
		Username: aws.String(email),
		Password: aws.String(password),
	}

	secretHash := computeSecretHash(c.cognitoClientSecret, email, c.cognitoClientId)

	input.SecretHash = aws.String(secretHash)

	user, err := c.cognitoClient.SignUp(context.TODO(), input)
	if err != nil {
		if strings.Contains(err.Error(), "UsernameExistsException") {
			err = errors.New("akun dengan email yang akan didaftarkan telah ada, jika lupa password silahkan klik lupa password")
		}

		if strings.Contains(err.Error(), "InvalidParameterException") {
			err = errors.New("username harus berupa email")
		}

		return "", err
	}

	return *user.CodeDeliveryDetails.Destination, nil
}

func (c *AwsCognito) ConfirmSignUp(email, code string) (*cognito.ConfirmSignUpOutput, error) {
	input := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String(c.cognitoClientId),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
	}

	secretHash := computeSecretHash(c.cognitoClientSecret, email, c.cognitoClientId)
	input.SecretHash = aws.String(secretHash)

	result, err := c.cognitoClient.ConfirmSignUp(context.TODO(), input)
	if err != nil {
		if strings.Contains(err.Error(), "ExpiredCodeException") {
			err = errors.New("kode telah expire, silahkan request kode konfirmasi kembali")
		}

		if strings.Contains(err.Error(), "CodeMismatchException") {
			err = errors.New("kode verifikasi gagal, silahkan cek kembali pesan di email anda")
		}
		return nil, err
	}

	return result, nil
}

func (c *AwsCognito) SignIn(email, password string) (*types.AuthenticationResultType, error) {
	secretHash := computeSecretHash(c.cognitoClientSecret, email, c.cognitoClientId)

	input := &cognito.InitiateAuthInput{
		ClientId: aws.String(c.cognitoClientId),
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		},
	}

	result, err := c.cognitoClient.InitiateAuth(context.TODO(), input)
	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			err = errors.New("nama pengguna atau password kurang tepat")
		}

		if strings.Contains(err.Error(), "UserNotConfirmedException") {
			err = errors.New("akun belum terkonfirmasi, silahkan cek kode konfirmasi didalm email")
		}

		if strings.Contains(err.Error(), "InvalidParameterException") {
			err = errors.New("parameter tidak valid")
		}
		return nil, err
	}

	return result.AuthenticationResult, nil
}

func (c *AwsCognito) SignOut(userName string) error {
	input := &cognito.AdminUserGlobalSignOutInput{
		UserPoolId: aws.String(c.cognitoPoolId),
		Username:   aws.String(userName),
	}

	_, err := c.cognitoClient.AdminUserGlobalSignOut(context.TODO(), input)
	if err != nil {
		log.Println("failed to perform global sign out:", err.Error())
		return err
	}

	return nil
}
