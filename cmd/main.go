package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tbtec/vdlg-login/internal/controller"
	"github.com/tbtec/vdlg-login/internal/dto"
)

func routerReq(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slog.Info("received a request", "path", req.Path, "method", req.HTTPMethod)

	loginController := controller.NewLoginController()

	if req.Path == "/login" && req.HTTPMethod == "POST" {
		return loginController.Login(req)
	}

	errorMessage := dto.ErrorMessage{Error: dto.Error{Description: "Method not allowed"}}
	body, err := json.Marshal(errorMessage)
	if err != nil {
		slog.Error("error while trying to marshal the response", "error", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 405,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	log := slog.New(handler)

	slog.SetDefault(log)

	lambda.Start(routerReq)
}
