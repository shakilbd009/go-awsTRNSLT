package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/translate"
)

func main() {
	lambda.Start(translt)
}

func translt(reQ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := external.LoadDefaultAWSConfig(external.WithDefaultRegion("us-east-2"))
	if err != nil {
		log.Fatalln(err)
	}
	svc := translate.New(cfg)
	query := reQ.QueryStringParameters
	if v, ok := query["text"]; ok {
		req := svc.TranslateTextRequest(&translate.TranslateTextInput{
			SourceLanguageCode: aws.String("auto"),
			TargetLanguageCode: aws.String("en"),
			Text:               aws.String(v),
		})
		resp, err := req.Send(context.Background())
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: resp.String(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       err.Error(),
	}, nil
}
