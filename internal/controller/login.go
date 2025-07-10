package controller

import (
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tbtec/vdlg-login/internal/dto"
	"github.com/tbtec/vdlg-login/internal/usecase"
)

type LoginController struct {
	UscLogin *usecase.UscLogin
}

func NewLoginController() *LoginController {
	return &LoginController{
		UscLogin: usecase.NewUseCaseLogin(),
	}
}

func (ctl *LoginController) Login(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginRequest dto.LoginRequest

	if err := json.Unmarshal([]byte(req.Body), &loginRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string("Invalid Body"),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	login, err := ctl.UscLogin.Login(loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       string("Invalid credentials"),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	body, err := json.Marshal(login)
	if err != nil {
		slog.Error("error while trying to marshal the response", "error", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}
