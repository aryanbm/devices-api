# Simple Restful API om AWS with Go Language

An example of how to implement a Restful API on AWS by utilizing the following tech stack:
* [Serverless Framework](https://serverless.com)
* [Go Language](https://golang.org)
* [AWS API Gateway](https://aws.amazon.com/api-gateway/)
* [AWS Lambda](https://aws.amazon.com/lambda/)
* [AWS DynamoDB](https://aws.amazon.com/dynamodb/)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Prerequisites

##### 1. Go Language
First you need to install Go language on your system,
On this website you can find installation instructions based on your system OS:
https://golang.org/doc/install

##### 2. Go Dep
Go dep helps us manage our dependencies,
you can install go-dep on your system with the following command

```
sudo apt-get install go-dep
```
##### 3. Serverless framework
Install via npm:
```
npm install -g serverless
```

### Installing and Deployment

1. Clone the repo
```
git clone https://github.com/aryanbm/devices-api
```
2. Install dependencies
```
dep ensure -v
```
3. Build go files
```
env GOOS=linux go build -ldflags="-s -w" -o bin/getDeviceBin src/getDevice/getDevice.go
env GOOS=linux go build -ldflags="-s -w" -o bin/createDeviceBin src/createDevice/createDevice.go
```
or you can just use `make build`

4. Deploy your code to AWS account
```
serverless deploy --verbose
```
or
`make deploy`

## Running the tests
Test covarages are available in cover.html files

You can test `getDevice` and `createDevice` unit tests with the following command
```
go test ./src/getDevice/
go test ./src/createDevice/
```

## Testing API

#### GetDevice:

```
curl -X GET https://<api-gateway-url>/api/devices/id1
```



#### CreateDevice:
```
curl -X POST https://<api-gateway-url>/api/devices/ -H "Content-Type: application/json" --data-binary @- <<DATA
{
    "id": "/devices/id1",
    "deviceModel": "/devicesmodels/id1",
    "name": "Sensor",
    "note": "Testing a sensor.",
    "serial": "A020000102"
}
DATA
```
