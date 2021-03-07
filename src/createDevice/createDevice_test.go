package main

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

// Package dynamodbiface provides an interface to enable
// mocking the Amazon DynamoDB service client for testing our code.
// In this case we define a mock struct to be used in our unit tests of createDevice()
type stubDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

// Mocking DynamoDB PutItem function to use in our tests
// We defined a state that we can mock an error from DynamoDB PutItem
func (m *stubDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	// mock a device with invalid inputs
	device := Device{
		Id:          "INVALID_VALUE",
		DeviceModel: "INVALID_VALUE",
		Name:        "INVALID_VALUE",
		Note:        "INVALID_VALUE",
		Serial:      "INVALID_VALUE",
	}

	mockDevice, _ := dynamodbattribute.MarshalMap(device)

	mockInput := &dynamodb.PutItemInput{
		Item:      mockDevice,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	// Returns an error if input parameters are invalid
	if reflect.DeepEqual(input, mockInput) {
		return &dynamodb.PutItemOutput{}, errors.New("ValidationException: One or more parameter values were invalid:xxxxxx")
	}

	return &dynamodb.PutItemOutput{}, nil
}

func TestHandler(t *testing.T) {

	// Put the mocked DynamoDB API in place of the main DynamoDB API.
	mDynamodb := &stubDynamoDBClient{}
	svc = mDynamodb

	// Test Cases
	tests := []struct {
		Request    events.APIGatewayProxyRequest
		Expect     string
		StatusCode int
		Err        error
	}{

		{
			// Test that the handler gets correct inputs
			// and it will be added to the DynamoDB successfully.
			Request:    events.APIGatewayProxyRequest{Body: "{\"id\":\"/devices/id1\",\"deviceModel\":\"/devicemodels/id1\",\"name\":\"Sensor\",\"note\":\"Testing a sensor\",\"serial\":\"A020000102\"}"},
			Expect:     "{\"message\":\"device added to DynamoDB\",\"statusCode\":201}",
			StatusCode: 201,
			Err:        nil,
		},
		{
			// Test that the handler responds Bad Request 400
			// when some parameters are missing
			Request:    events.APIGatewayProxyRequest{Body: "{\"name\":\"Sensor\",\"note\":\"Testing a sensor\",\"serial\":\"A020000102\"}"},
			Expect:     "{\"message\":\"{Id, DeviceModel} are missing\",\"statusCode\":400}",
			StatusCode: 400,
			Err:        nil,
		},
		{
			// Test that the handler responds internal error for invalid input
			// for example, if the hash key exceeds its maximum size limit of 2048 bytes.
			Request:    events.APIGatewayProxyRequest{Body: "{\"id\":\"INVALID_VALUE\",\"deviceModel\":\"INVALID_VALUE\",\"name\":\"INVALID_VALUE\",\"note\":\"INVALID_VALUE\",\"serial\":\"INVALID_VALUE\"}"},
			Expect:     "{\"message\":\"DynamoDB PutItem method returned an error: ValidationException: One or more parameter values were invalid:xxxxxx\",\"statusCode\":500}",
			StatusCode: 500,
			Err:        nil,
		},
	}

	for _, test := range tests {
		response, err := Handler(test.Request)
		assert.IsType(t, test.Err, err)
		assert.Equal(t, test.StatusCode, response.StatusCode)
		assert.Equal(t, test.Expect, response.Body)
	}
}
