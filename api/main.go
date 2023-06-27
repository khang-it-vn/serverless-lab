package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	RequestID   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Data        struct {
		Value1 int `json:"value1"`
		Value2 int `json:"value2"`
	} `json:"data"`
}

type Response struct {
	ResponseID   string `json:"responseId"`
	ResponseTime string `json:"responseTime"`
	Data         struct {
		Sum int `json:"sum"`
	} `json:"data"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req Request
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshalling request: %v", err)
		return events.APIGatewayProxyResponse{}, err
	}

	sum := req.Data.Value1 + req.Data.Value2

	resp := Response{
		ResponseID:   req.RequestID,
		ResponseTime: time.Now().Format(time.RFC3339),
		Data: struct {
			Sum int `json:"sum"`
		}{
			Sum: sum,
		},
	}

	responseBody, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	lambda.Start(Handler)
}
