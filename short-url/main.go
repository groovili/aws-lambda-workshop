package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

type Response events.APIGatewayProxyResponse

var db *dynamodb.DynamoDB

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	var resp Response
	var err error

	switch req.HTTPMethod {
	case http.MethodPost:
		resp, err = PostHandler(ctx, req)
	case http.MethodGet:
		resp, err = GetHandler(ctx, req)
	default:
		resp = Response{
			StatusCode: http.StatusForbidden,
		}
		err = nil
	}

	return resp, err
}

func PostHandler(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	URL := req.Body
	short, ok := req.PathParameters["short_url"]
	if !ok {
		return Response{
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	_, err := db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DB_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"alias": {S: aws.String(short)},
			"url":   {S: aws.String(URL)},
		}})
	if err != nil {
		log.Debug(err)

		return Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return Response{
		StatusCode: http.StatusOK,
		Body:       "Short URL - " + short + " added",
	}, nil
}

func GetHandler(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	short, ok := req.PathParameters["short_url"]
	if !ok {
		return Response{
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	out, err := db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DB_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"alias": {S: aws.String(short)},
		},
	})
	if err != nil {
		log.Debug(err)

		return Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	originalURL, ok := out.Item["url"]
	if err != nil {
		log.Debug(err)

		return Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return Response{
		StatusCode: http.StatusOK,
		Body:       "Original URL - " + aws.StringValue(originalURL.S),
	}, nil
}

func main() {
	// Create a new AWS session and fail immediately on error
	sess := session.Must(session.NewSession())
	// Create the DynamoDB client
	db = dynamodb.New(sess)

	lambda.Start(Handler)
}
