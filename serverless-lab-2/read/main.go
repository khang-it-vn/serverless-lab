package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v4"
)

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type APIRequest struct {
	RequestID   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Data        struct {
		Username string `json:"username"`
	} `json:"data"`
}

type APIResponse struct {
	ResponseId      string `json:"responseId"`
	ResponseTime    string `json:"responseTime"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Data            User   `json:"data"`
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dbConfig := DBConfig{
		Host:     "postgres.cjkfitk009d7.ap-southeast-1.rds.amazonaws.com",
		Port:     "5432",
		Username: "postgres",
		Password: "Xinchao123",
		DBName:   "postgres",
	}

	conn, err := pgx.Connect(context.Background(), "postgres://"+dbConfig.Username+":"+dbConfig.Password+"@"+dbConfig.Host+":"+dbConfig.Port+"/"+dbConfig.DBName)
	if err != nil {
		log.Fatal(err)
	}

	var req APIRequest

	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		response := APIResponse{
			ResponseId:      req.RequestID,
			ResponseTime:    time.Now().String(),
			ResponseCode:    "ERROR",
			ResponseMessage: "Invalid request body",
		}
		responseJSON, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJSON),
		}, nil
	}

	// Lấy thông tin user theo username
	user, err := getUserByUsername(conn, req.Data.Username)
	if err != nil {
		response := APIResponse{
			ResponseId:      req.RequestID,
			ResponseTime:    time.Now().String(),
			ResponseCode:    "ERROR",
			ResponseMessage: "Internal server error",
		}
		responseJSON, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string(responseJSON),
		}, nil
	}
	if user == nil {
		response := APIResponse{
			ResponseId:      req.RequestID,
			ResponseTime:    time.Now().String(),
			ResponseCode:    "ERROR",
			ResponseMessage: "User not found",
		}
		responseJSON, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJSON),
		}, nil
	}

	response := APIResponse{
		ResponseId:      req.RequestID,
		ResponseTime:    time.Now().String(),
		ResponseCode:    "SUCCESS",
		ResponseMessage: "User detail retrieved successfully",
		Data: User{
			Username: user.Username,
			Name:     user.Name,
			Phone:    user.Phone,
		},
	}
	responseJSON, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
	}, nil
}

func getUserByUsername(conn *pgx.Conn, username string) (*User, error) {
	var user User
	err := conn.QueryRow(context.Background(), "SELECT username, name, phone FROM users WHERE username = $1", username).Scan(&user.Username, &user.Name, &user.Phone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func main() {
	lambda.Start(HandleRequest)
}
