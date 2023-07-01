package main

import (
	"context"
	"encoding/json"

	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

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

type Request struct {
	RequestID   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Data        User   `json:"data"`
	Signature 	string `json:"signature"`
}
type RequestBody struct {
	RequestID    string `json:"requestId"`
	RequestTime  string `json:"requestTime"`
	Data         User   `json:"data"`
}

const SUCCESS  = "SUCCESS"
const SERVER_ERROR = "SERVER_ERROR"

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// format request to json
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

	// verify signature
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


	// verify phone number
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

	// create user
	stateCreateUser := strings.Compare(createUserApiLab2(req), SUCCESS)
	
	response := APIResponse{
		ResponseId:      req.RequestID,
		ResponseTime:    time.Now().String(),
		ResponseCode:    "SUCCESS",
		ResponseMessage: "USER CREATED SUCCESSFULLY",
	}
	if stateCreateUser != 0{
		response := APIResponse{
			ResponseId:      req.RequestID,
			ResponseTime:    time.Now().String(),
			ResponseCode:    "ERROR",
			ResponseMessage: "USER CREATION FAILED",
		}
		responseJSON, _ := json.Marshal(response)

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(responseJSON),
		}, nil
	}
	responseJSON, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
	}, nil
}

func createUserApiLab2(data Request ) string{
	// Tạo request body
	requestBody := RequestBody{
		RequestID:   data.RequestID,
		RequestTime: data.RequestTime,
		Data: data.Data,
	}

	// Chuyển đổi request body thành JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("FORMAT ERROR:", err)
		return SERVER_ERROR
	}

	// Tạo HTTP request
	url := "https://2xqobkgcwa.execute-api.us-east-1.amazonaws.com/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("INTERNAL SERVER ERROR:", err)
		return SERVER_ERROR 
	}

	// Set header và content type
	req.Header.Set("x-api-key", "B5d4JtTU8u1ggV8gp7OF88gcCGxZls6T3f5PYZSa")
	req.Header.Set("Content-Type", "application/json")

	// Gửi HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("REQUEST ERROR:", err)
		return SERVER_ERROR
	}
	defer resp.Body.Close()

	// Đọc và in ra kết quả trả về
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("FORMAT ERROR:", err)
		return SERVER_ERROR
	}

	fmt.Println("SUCCESS:", response)

	return SUCCESS
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

// checkPhone function to verify phone number
func checkPhone(phone int) string{
	url := "https://1g1zcrwqhj.execute-api.ap-southeast-1.amazonaws.com/dev/testapi"
	requestID := uuid.New().String()
	requestBody := []byte(fmt.Sprintf(`{
		"requestId": "%s",
		"data": {
			"value": %d
		}
	}`,requestID, phone))

	// create http request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("INTERNAL SERVER ERROR:", err)
		return "INTERNAL SERVER ERROR";
	}

	// set headers
	req.Header.Set("x-api-key", "B5d4JtTU8u1ggV8gp7OF88gcCGxZls6T3f5PYZSa")
	req.Header.Set("Content-Type", "text/plain")

	// send http request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("CALL API ERROR:", err)
		return "CALL API ERROR"
	}
	defer resp.Body.Close()

	// read response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HANDLE DATA ERROR:", err)
		return "HANDLE DATA ERROR"
	}

	fmt.Println("CALL API SUCCESS:", string(responseBody))
	return SUCCESS
}

func main() {
	lambda.Start(HandleRequest)
}

