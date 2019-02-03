package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Response events.APIGatewayProxyResponse

type Message struct {
	Message string `json:"message"`
}

func Handler(req events.APIGatewayProxyRequest) (Response, error) {
	name, ok := req.QueryStringParameters["name"]
	if !ok {
		name = "John Doe"
	}

	msg := &Message{Message: "Hello, " + name}

	log.Infoln(msg)

	m, err := jsoniter.MarshalToString(msg)
	if err != nil {
		log.Debugf("Failed to marshal %s", msg)
		return Response{
			StatusCode: http.StatusBadRequest,
		}, err
	}

	resp := Response{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: m,
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
