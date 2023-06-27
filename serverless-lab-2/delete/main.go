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

type APIResponse struct {
	ResponseId      string `json:"responseId"`
	ResponseTime    string `json:"responseTime"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type Request struct {
	RequestID   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Data        struct {
		Username string `json:"username"`
	} `json:"data"`
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

	var req Request

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

	// Kiểm tra sự tồn tại duy nhất của username
	exists, err := checkUsernameExists(conn, req.Data.Username)
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
	if !exists {
		response := APIResponse{
			ResponseId:      req.RequestID,
			ResponseTime:    time.Now().String(),
			ResponseCode:    "ERROR",
			ResponseMessage: "Username does not exist",
		}
		responseJSON, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJSON),
		}, nil
	}

	// Xóa người dùng
	err = deleteUser(conn, req.Data.Username)
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

	response := APIResponse{
		ResponseId:      req.RequestID,
		ResponseTime:    time.Now().String(),
		ResponseCode:    "SUCCESS",
		ResponseMessage: "User deleted successfully",
	}
	responseJSON, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
	}, nil
}

func checkUsernameExists(conn *pgx.Conn, username string) (bool, error) {
	var count int
	err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func deleteUser(conn *pgx.Conn, username string) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM users WHERE username = $1", username)
	return err
}

func main() {
	lambda.Start(HandleRequest)
}
