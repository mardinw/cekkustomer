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

func (c *AwsCognito) ResendConfirmationCode(email string) (string, error) {
	input := &cognito.ResendConfirmationCodeInput{
		ClientId: aws.String(c.cognitoClientId),
		Username: aws.String(email),
	}

	secretHash := computeSecretHash(c.cognitoClientSecret, email, c.cognitoClientId)
	input.SecretHash = aws.String(secretHash)

	code, err := c.cognitoClient.ResendConfirmationCode(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return *code.CodeDeliveryDetails.Destination, nil
}

func (c *AwsCognito) GetUsername(token string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: &token,
	}
	result, err := c.cognitoClient.GetUser(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *AwsCognito) UpdateUserAttributes(email, attribute, value string) error {
	inputAttributes := []types.AttributeType{
		{
			Name:  aws.String(attribute),
			Value: aws.String(value),
		},
	}

	input := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId:     aws.String(c.cognitoPoolId),
		Username:       aws.String(email),
		UserAttributes: inputAttributes,
	}

	_, err := c.cognitoClient.AdminUpdateUserAttributes(context.TODO(), input)

	return err
}

func (c *AwsCognito) CheckUserAttributes(email string) ([]types.AttributeType, error) {

	input := &cognito.AdminGetUserInput{
		Username:   aws.String(email),
		UserPoolId: &c.cognitoPoolId,
	}

	output, err := c.cognitoClient.AdminGetUser(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return output.UserAttributes, nil
}

func (c *AwsCognito) ForgotPassword(email string) (string, error) {
	input := &cognito.ForgotPasswordInput{
		ClientId: aws.String(c.cognitoClientId),
		Username: aws.String(email),
	}

	secretHash := computeSecretHash(c.cognitoClientSecret, email, c.cognitoClientId)
	input.SecretHash = aws.String(secretHash)

	result, err := c.cognitoClient.ForgotPassword(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return *result.CodeDeliveryDetails.Destination, nil
}

func (c *AwsCognito) ResetPassword(email, password, code string) (*cognito.ConfirmForgotPasswordOutput, error) {
	input := &cognito.ConfirmForgotPasswordInput{
		ClientId:         aws.String(c.cognitoClientId),
		Username:         aws.String(email),
		Password:         aws.String(password),
		ConfirmationCode: aws.String(code),
	}

	secretHash := computeSecretHash(c.cognitoClientSecret, email, c.cognitoClientId)
	input.SecretHash = aws.String(secretHash)

	result, err := c.cognitoClient.ConfirmForgotPassword(context.TODO(), input)
	if err != nil {
		if strings.Contains(err.Error(), "CodeMismatchException") {
			err = errors.New("kode verifikasi tidak cocok, silahkan request kode konfirmasi kembali")
		}

		if strings.Contains(err.Error(), "LimitExceededException") {
			err = errors.New("maksium pengulangan password hanya sampai 3x, silahkan coba lagi 1 jam berikutnya")
		}

		if strings.Contains(err.Error(), "ExpiredCodeException") {
			err = errors.New("kode verifikasi gagal, silahkan request lupa password kembali")
		}

		if strings.Contains(err.Error(), "InvalidParameterException") {
			err = errors.New("maaf, parameter tidak bisa dipakai, silahkan ubah parameternya")
		}
		return nil, err
	}

	return result, err
}

func (c *AwsCognito) AddUserToGroup(userName, groupName string) error {
	input := &cognito.AdminAddUserToGroupInput{
		UserPoolId: aws.String(c.cognitoPoolId),
		Username:   aws.String(userName),
		GroupName:  aws.String(groupName),
	}

	_, err := c.cognitoClient.AdminAddUserToGroup(context.TODO(), input)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			err = errors.New("group tidak ada")
		}

		return err
	}

	return nil
}

func (c *AwsCognito) CheckUserInGroup(userName string) ([]string, error) {
	input := &cognito.AdminListGroupsForUserInput{
		Username:   aws.String(userName),
		UserPoolId: aws.String(c.cognitoPoolId),
	}

	resp, err := c.cognitoClient.AdminListGroupsForUser(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	groupNames := make([]string, len(resp.Groups))
	for i, group := range resp.Groups {
		groupNames[i] = *group.GroupName
	}

	return groupNames, nil
}

func (c *AwsCognito) ListGroup() ([]string, error) {
	input := &cognito.ListGroupsInput{
		UserPoolId: aws.String(c.cognitoPoolId),
	}

	resp, err := c.cognitoClient.ListGroups(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	groupNames := make([]string, len(resp.Groups))
	for i, group := range resp.Groups {
		groupNames[i] = *group.GroupName
	}
	return groupNames, nil
}
