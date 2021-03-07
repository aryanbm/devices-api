package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Device struct {
	Id          string `json:"id"`
	DeviceModel string `json:"deviceModel"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Serial      string `json:"serial"`
}

type ResponseMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type Response events.APIGatewayProxyResponse

// Define global DynamoDBAPI varieble, so we can change it in the unit test
var svc dynamodbiface.DynamoDBAPI

func init() {
	//Set up a session to be used by the SDK to load
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)
}
func Handler(request events.APIGatewayProxyRequest) (Response, error) {

	// Device is used to construct the request from the client
	device := Device{
		Id:          "",
		DeviceModel: "",
		Name:        "",
		Note:        "",
		Serial:      "",
	}

	// Unmarshal the Request body
	// return (Bad Request 400) if error
	err := json.Unmarshal([]byte(request.Body), &device)
	if err != nil {
		return Response{Body: (&ResponseMessage{err.Error(), 400}).json(), StatusCode: 400}, nil
	}

	// Trim "/devices/" from id parameter
	device.Id = strings.TrimPrefix(device.Id, "/devices/")

	// Check if payloads are valid
	// return missing payloads if error
	missingPayloads, isValid := checkPayloads(&device)
	if !isValid {
		return Response{Body: (&ResponseMessage{missingPayloads, 400}).json(), StatusCode: 400}, nil
	}

	/*
		TODO:
		Dynammodb will update other parameters if the device id already exists in the table.
		but we can check if this device id already exists or not (with getDevice function)
	*/

	// Add device to DynamoDB Table
	// return createDevice() function message if error
	_err := createDevice(&device)
	if _err != nil {
		return Response{Body: _err.json(), StatusCode: _err.StatusCode}, nil
	}

	// The device was submitted successfully
	return Response{Body: (&ResponseMessage{"device added to DynamoDB", 201}).json(), StatusCode: 201}, nil
}

func createDevice(device *Device) *ResponseMessage {

	// Marshal Device to DynamoDB item
	dbItem, err := dynamodbattribute.MarshalMap(device)
	if err != nil {
		return &ResponseMessage{fmt.Sprintf("Error marshalling item: %v", err.Error()), 500}
	}

	// Build DynamoDB PutItem input
	input := &dynamodb.PutItemInput{
		Item:      dbItem,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	// Asking DynamoDB client to put the item
	_, err = svc.PutItem(input)

	// Checking if PutItem has error
	if err != nil {
		return &ResponseMessage{fmt.Sprintf("DynamoDB PutItem method returned an error: %v", err.Error()), 500}
	}

	// Everything is okay, so we can return function without error
	return nil
}

// Check payloads to make sure they are not missing
func checkPayloads(device *Device) (message string, isValid bool) {
	var mp []string // Missing payloads

	if device.Id == "" {
		mp = append(mp, "Id") // Id is missing
	}
	if device.DeviceModel == "" {
		mp = append(mp, "DeviceModel") // DeviceModel is missing
	}
	if device.Name == "" {
		mp = append(mp, "Name") // Name is missing
	}
	if device.Note == "" {
		mp = append(mp, "Note") // Note is missing
	}
	if device.Serial == "" {
		mp = append(mp, "Serial") // Serial is missing
	}

	// If any of the payloads are missing
	if len(mp) > 0 {
		if len(mp) == 1 {
			return fmt.Sprintf("{%v} is missing", mp[0]), false
		}

		return fmt.Sprintf("{%v} are missing", strings.Join(mp, ", ")), false
	}

	return "the inputs are valid", true
}

// Generate a JSON string from ResponseMessage
func (e *ResponseMessage) json() string {
	json, _ := json.Marshal(e)
	return string(json)
}

func main() {
	lambda.Start(Handler)
}
