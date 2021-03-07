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

	// Get the device {id} parameter
	deviceId := request.PathParameters["id"]

	// If the get parameter is empty
	if deviceId == "" {
		return Response{Body: (&ResponseMessage{"Device identifiers cannot be empty", 400}).json(), StatusCode: 400}, nil
	}

	// Add "/devices/" to the beginning of id parameter
	if !strings.HasPrefix(deviceId, "/devices/") {
		deviceId = "/devices/" + deviceId
	}

	// Get device from DynamoDB
	device, err := getDevice(deviceId)

	// If getDevice() returns an error
	if err != nil {
		return Response{Body: err.json(), StatusCode: err.StatusCode}, nil
	}

	// Converting Device type to JSON
	deviceString, _ := json.Marshal(device)

	return Response{Body: string(deviceString), StatusCode: 200}, nil
}

func getDevice(deviceId string) (_device *Device, _error *ResponseMessage) {
	// search for deviceId in DynammoDB Table
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(deviceId),
			},
		},
	})
	// If DynamoDB GetItem returns an error
	if err != nil {
		return nil, &ResponseMessage{err.Error(), 500}
	}

	// Device {id} is not founded
	if result.Item == nil {
		return nil, &ResponseMessage{fmt.Sprintf("Device with id (%v) not founded", deviceId), 404}
	}

	device := Device{}
	// Unmarshalling the query result into a Device type
	err = dynamodbattribute.UnmarshalMap(result.Item, &device)
	if err != nil {
		return nil, &ResponseMessage{fmt.Sprintf("Failed to unmarshal Record, %v", err), 500}
	}

	// Everything is okay, so we can return the device.
	return &device, nil
}

// Generate a JSON string from ResponseMessage
func (e *ResponseMessage) json() string {
	json, _ := json.Marshal(e)
	return string(json)
}

func main() {
	lambda.Start(Handler)
}
