package usecase

import (
	"errors"

	"github.com/tbtec/vdlg-login/internal/dto"

	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	//"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
)

// var (
// 	ErrorInvalidCredentials = xerrors.NewBusinessError("TLL-LOGIN-001", "Invalid Credentials")
// )

type UscLogin struct {
}

func NewUseCaseLogin() *UscLogin {
	return &UscLogin{}
}

var (
	userPoolID   = "us-east-1_tRk65cj3j"        // User Pool ID
	clientID     = "15v4gl95tci1ajqj7sro466im3" // App Client ID
	region       = "us-east-1"
	clientSecret = "154iful2mel5oljgmshpsh1aeoihr1jubiikvjqo8hlstp2n80sb"
)

func generateSecretHash(clientID, clientSecret, username string) string {
	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func authenticateUser(username, password, secretHash string, sess *session.Session) (*cognitoidentityprovider.AdminInitiateAuthOutput, error) {

	svc := cognitoidentityprovider.New(sess)

	//slog.Info("Iniciando autenticação no Cognito", "username", username, "userPoolID", userPoolID, "clientID", clientID)
	authParams := map[string]*string{
		"USERNAME":    aws.String(username),
		"PASSWORD":    aws.String(password),
		"SECRET_HASH": aws.String(secretHash),
	}
	//slog.Info("Parâmetros de autenticação", "authParams", authParams)

	authInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:       aws.String("ADMIN_NO_SRP_AUTH"),
		UserPoolId:     aws.String(userPoolID),
		ClientId:       aws.String(clientID),
		AuthParameters: authParams,
	}
	//slog.Info("Chamando AdminInitiateAuth", "authInput", authInput)
	authOutput, err := svc.AdminInitiateAuth(authInput)
	if err != nil {
		//slog.Error("Erro ao chamar AdminInitiateAuth", "error", err)
		log.Fatalf("Erro na autenticação: %s\n", err)
		return authOutput, errors.New("invalid credentials")
	}
	//slog.Info("AdminInitiateAuth chamado com sucesso", "authOutput", authOutput)
	return authOutput, nil
}

func validateChallenge(username, password, secretHash string, sess *session.Session, authOutput *cognitoidentityprovider.AdminInitiateAuthOutput) (dto.Login, error) {

	//slog.Info("Validando desafio de autenticação", "challengeName", *authOutput.ChallengeName)
	if *authOutput.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		cognitoClient := cognitoidentityprovider.New(sess)
		//slog.Info("Desafio NEW_PASSWORD_REQUIRED encontrado, respondendo ao desafio")
		input := &cognitoidentityprovider.RespondToAuthChallengeInput{
			Session:       aws.String(*authOutput.Session),
			ClientId:      aws.String(clientID),
			ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
			ChallengeResponses: map[string]*string{
				"NEW_PASSWORD": aws.String(password),
				"USERNAME":     aws.String(username),
				"SECRET_HASH":  aws.String(secretHash),
			},
		}
		//slog.Info("Preparando resposta ao desafio", "input", input)

		_, err := cognitoClient.RespondToAuthChallenge(input)
		if err != nil {
			//slog.Error("Erro ao responder ao desafio NEW_PASSWORD_REQUIRED", "error", err)
			return dto.Login{}, err
		}
		//slog.Info("Resposta ao desafio NEW_PASSWORD_REQUIRED enviada com sucesso")
		return dto.Login{
			AccessToken: "Primeiro login efetuado",
		}, nil
	}
	//slog.Info("Nenhum desafio NEW_PASSWORD_REQUIRED encontrado, retornando token de acesso")
	return dto.Login{
		AccessToken: "Challenge não encontrado",
	}, nil

}

func (u *UscLogin) Login(loginRequest dto.LoginRequest) (dto.Login, error) {

	username := loginRequest.DocumentNumber
	password := loginRequest.Password

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatal("Erro ao criar sessão AWS:", err)
	}
	//slog.Info("Sessão AWS criada com sucesso", "region", region)
	secretHash := generateSecretHash(clientID, clientSecret, username)
	//slog.Info("Secret hash", "secretHash", secretHash)
	authOutput, err := authenticateUser(username, password, secretHash, sess)
	//slog.Info("Resultado da autenticação", "authOutput", authOutput)
	if err != nil {
		//slog.Error("Erro na autenticação", "error", err)
		log.Fatalf("Erro na autenticação: %s\n", err)
		return dto.Login{}, err
	}

	if authOutput.Session != nil {
		//slog.Info("Sessão encontrada na resposta", "session", *authOutput.Session)
		return validateChallenge(username, password, secretHash, sess, authOutput)
	} else {
		slog.Info("Nenhuma sessão encontrada na resposta")
	}
	//slog.Info("Retornando token de acesso", "accessToken", *authOutput.AuthenticationResult.AccessToken)
	return dto.Login{
		AccessToken: *authOutput.AuthenticationResult.AccessToken,
	}, nil

}
