package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type Request struct {
	RequestId   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Data        struct {
		PlaintText string `json:"plaintText"`
		SecretKey  string `json:"secretKey"`
	} `json:"data"`
}

type Response struct {
	ResponseId   string `json:"responseId"`
	ResponseTime string `json:"responseTime"`
	Data         struct {
		Signature string `json:"signature"`
	} `json:"data"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req Request
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Invalid request body", StatusCode: 400}, nil
	}

	// Generate UUID for responseId
	responseId := uuid.New().String()

	// Generate time in RFC3339 format for responseTime
	responseTime := time.Now().UTC().Format(time.RFC3339)

	// Generate signature using HMAC-SHA256
	h := hmac.New(sha256.New, []byte(req.Data.SecretKey))
	h.Write([]byte(req.Data.PlaintText))
	signature := hex.EncodeToString(h.Sum(nil))

	// Create response object
	res := Response{
		ResponseId:   responseId,
		ResponseTime: responseTime,
		Data: struct {
			Signature string `json:"signature"`
		}{
			Signature: signature,
		},
	}

	// Convert response object to JSON
	resJson, err := json.Marshal(res)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Error creating response", StatusCode: 500}, nil
	}

	// Return response
	return events.APIGatewayProxyResponse{Body: string(resJson), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
