package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v4"
	"github.com/google/uuid"
	"fmt"
	"strings"
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

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
	Data        User   `json:"data"`
	Signature 	string `json:"signature"`
}

const SUCCESS  = "SUCCESS"
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req Request

	err := json.Unmarshal([]byte(request.Body), &req)

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

	secretKey := "%secretKey%"
	plainText := req.RequestID + req.Data.Phone + req.Data.Username + secretKey
	signature := generateHMAC(secretKey, plainText )

	resultCompareSignature := strings.Compare(signature, req.Signature)

	if resultCompareSignature != 0{
		response := APIResponse{
			ResponseId:      req.RequestID,
			ResponseTime:    time.Now().String(),
			ResponseCode:    "ERROR",
			ResponseMessage: "Invalid signature",
		}
		responseJSON, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJSON),
		}, nil
	}


	phone, err := strconv.Atoi(req.Data.Phone)
	
	if err != nil {
		response := APIResponse{
			ResponseId:      req.RequestID,
			ResponseTime:    time.Now().String(),
			ResponseCode:    "ERROR",
			ResponseMessage: "WRONG PHONE NUMBER",
		}
		responseJSON, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJSON),
		}, nil
	}
	stateString := checkPhone(phone)

	result := strings.Compare(stateString, SUCCESS)
	if result != 0{
		response := APIResponse{
			ResponseId:      req.RequestID,
			ResponseTime:    time.Now().String(),
			ResponseCode:    "ERROR",
			ResponseMessage: "WRONG PHONE NUMBER",
		}
		responseJSON, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJSON),
		}, nil
	}
	
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
	if exists {
		response := APIResponse{
			ResponseId:      req.RequestID,
			ResponseTime:    time.Now().String(),
			ResponseCode:    "ERROR",
			ResponseMessage: "Username already exists",
		}
		responseJSON, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseJSON),
		}, nil
	}

	// Tạo người dùng mới
	err = createUser(conn, req.Data)
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
		ResponseMessage: "User created successfully",
	}
	responseJSON, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
	}, nil
}
func generateHMAC(key, message string) string {
	// Chuyển đổi khóa (key) và thông điệp (message) thành byte slices
	keyBytes := []byte(key)
	messageBytes := []byte(message)

	// Tạo một đối tượng HMAC-SHA256 với khóa
	h := hmac.New(sha256.New, keyBytes)

	// Thêm thông điệp vào HMAC
	h.Write(messageBytes)

	// Tính toán mã hash HMAC-SHA256
	hash := h.Sum(nil)

	// Chuyển đổi mã hash thành chuỗi hex
	hashString := hex.EncodeToString(hash)

	return hashString
}

func checkPhone(phone int) string{
	url := "https://1g1zcrwqhj.execute-api.ap-southeast-1.amazonaws.com/dev/testapi"
	requestID := uuid.New().String()
	requestBody := []byte(fmt.Sprintf(`{
		"requestId": "%s",
		"data": {
			"value": %d
		}
	}`,requestID, phone))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("INTERNAL SERVER ERROR:", err)
		return "INTERNAL SERVER ERROR";
	}

	req.Header.Set("x-api-key", "B5d4JtTU8u1ggV8gp7OF88gcCGxZls6T3f5PYZSa")
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("CALL API ERROR:", err)
		return "CALL API ERROR"
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HANDLE DATA ERROR:", err)
		return "HANDLE DATA ERROR"
	}

	fmt.Println("CALL API SUCCESS:", string(responseBody))
	return SUCCESS
}

func checkUsernameExists(conn *pgx.Conn, username string) (bool, error) {
	var count int
	err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createUser(conn *pgx.Conn, user User) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO users (username, name, phone) VALUES ($1, $2, $3)", user.Username, user.Name, user.Phone)
	return err
}

func main() {
	lambda.Start(HandleRequest)
}

