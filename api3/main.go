package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type Request struct {
	NeedEncode string `json:"needEncode"`
	NeedDecode string `json:"needDecode"`
}

type Response struct {
	ResponseId   string    `json:"responseId"`
	ResponseTime time.Time `json:"responseTime"`
	Data         Data      `json:"data"`
}

type Data struct {
	OutEncode string `json:"outEncode"`
	OutDecode string `json:"outDecode"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req Request
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Invalid request body", StatusCode: 400}, nil
	}

	// Encode
	encoded := base64.StdEncoding.EncodeToString([]byte(req.NeedEncode))

	// Decode
	decoded, err := base64.StdEncoding.DecodeString(req.NeedDecode)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Invalid base64 string", StatusCode: 400}, nil
	}

	res := Response{
		ResponseId:   uuid.New().String(),
		ResponseTime: time.Now().UTC(),
		Data: Data{
			OutEncode: encoded,
			OutDecode: string(decoded),
		},
	}

	resBody, err := json.Marshal(res)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Error marshalling response", StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(resBody), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
