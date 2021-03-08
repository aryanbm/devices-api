package main

import (
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
// In this case we define a mock struct to be used in our unit tests of getDevice()
type stubDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

// Mocking DynamoDB GetItem function to use in our tests
func (m *stubDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	// Create a mock input with "id1" value as Item id key
	mockInput := (&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("id1"),
			},
		},
	})

	// Returns a mock device data if input matches mock input
	if reflect.DeepEqual(input, mockInput) {
		mockDevice := Device{
			Id:          "id1",
			DeviceModel: "/devicemodels/id1",
			Name:        "Sensor",
			Note:        "Testing a sensor",
			Serial:      "A020000102",
		}
		// Marshaling Device to a dynamodbattribute Map
		mockOutputItem, _ := dynamodbattribute.MarshalMap(mockDevice)

		return &dynamodb.GetItemOutput{Item: mockOutputItem}, nil
	}
	// Return an empty GetItemOutput if the device Id is not "id1"
	return &dynamodb.GetItemOutput{}, nil
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
			// Test that the handler responds the device data
			// When Dynamodb finds the device ID
			Request:    events.APIGatewayProxyRequest{PathParameters: map[string]string{"id": "id1"}},
			Expect:     "{\"id\":\"/devices/id1\",\"deviceModel\":\"/devicemodels/id1\",\"name\":\"Sensor\",\"note\":\"Testing a sensor\",\"serial\":\"A020000102\"}",
			StatusCode: 200,
			Err:        nil,
		},
		{
			// Test that the handler responds with the device data
			// When Dynamodb cannot find the device ID
			Request:    events.APIGatewayProxyRequest{PathParameters: map[string]string{"id": "INVALID_ID"}},
			Expect:     "{\"message\":\"Device with id (INVALID_ID) not founded\",\"statusCode\":404}",
			StatusCode: 404,
			Err:        nil,
		},
		{
			// Test that the handler responds Bad Request 400
			// when no id is provided in the request body
			Request:    events.APIGatewayProxyRequest{},
			Expect:     "{\"message\":\"Device identifiers cannot be empty\",\"statusCode\":400}",
			StatusCode: 400,
			Err:        nil,
		},
	}

	for _, test := range tests {
		response, err := Handler(test.Request)
		assert.IsType(t, test.Err, err)
		assert.Equal(t, test.Expect, response.Body)
		assert.Equal(t, test.StatusCode, response.StatusCode)
	}
}
