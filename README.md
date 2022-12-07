# Receipt Processor - Joel Yoshiya Foster

This is a webservice that processes receipts and returns the number of points earned. The service is built with Golang and the Gin Web Framework, and containerized with Docker. Tests are written with the Go testing package, as well as Testify and Dockertest. The service conforms to the API described in `api.yml`.


## Approach

- Create a `Receipt` struct to hold the data from the inbound JSON
  - Bind the JSON to the `Receipt` struct
  - Perform validation on the inbound JSON
  - Process points post-validation
- Create a `Receipts` struct to hold a map of `ReceiptPoints` structs
  - Create a `ReceiptPoints` struct to hold the points for a receipt
    - keeps a receipt and its points tightly coupled in the same location
  - A unique ID generated per valid receipt, with O(1) lookup

Why process points during receipt processing?

- Makes GET requests to `/receipts/{id}/points` more performant (assuming GET requests are more frequent than POST requests)
- Slightly less performant POST requests to `/receipts/process`, but safer to process points immediately after validation.

## Assumptions

- No persistence layer is required (for this exercise)
- Negative prices, totals, and points are not inbound/outbound from the API
  - Have error handling to cover these cases
- Assuming points tied to a receipt are immutable.

## Dependencies

- Go version: go1.18.3 darwin/amd64
- Gin Web Framework
- Other dependencies are listed in `go.mod`

## API

### Endpoint: Process Receipts

- Path: `/receipts/process`
- Method: `POST`
- Payload: Receipt JSON
- Response: JSON containing an id for the receipt.

Description:

Takes in a JSON receipt (see example in the example directory) and returns a JSON object with an ID specifying the id.

The ID returned is the ID that should be passed into `/receipts/{id}/points` to get the number of points the receipt
was awarded.

How many points should be earned are defined by the rules below.

Example Response:

```json
{"id": "7fb1377b-b223-49d9-a31a-5a02701dd310"}
```

### Endpoint: Get Points

- Path: `/receipts/{id}/points`
- Method: `GET`
- Response: A JSON object containing the number of points awarded.

A simple Getter endpoint that looks up the receipt by the ID and returns an object specifying the points awarded.

Example Response:

```json
{ "points": 32 }
```

## Execution

I've opted to use Docker to run the application. This allows for a consistent environment across all platforms.

### Build

A Dockerfile is included in the root of the project. To build the image, run the following command:

- Run `docker build -t receipt-processor-service .` at the root of the project.

### Run The Service

- Run `docker run -dp 8080:8080 --name receipt-rest-server receipt-processor-service` to start the service

## Test Environment

- Tests are found in the main directory
- Note: Testing will require Go (version specified in Dockerfile recommended) and Docker to be installed on your machine.
  - Testing is primarily for my personal use, assuming engineers will be using their own means to test the service within a Docker container.
- Make sure you're on the most up to date build: `docker build -t receipt-processor-service:latest .`
- Run the command `go test -v ./main` at the root to run the tests.
- Tests are written with the Go testing package, as well as Testify and Dockertest.

### Test Classes

- `service_test.go` - tests the service layer
- `service_docker_test.go` - tests the service layer within a docker container

### Test Cases

- `TestProcessReceipts` - tests the `/receipts/process` endpoint
  - In the docker test, this also tests the GetPoints endpoint, since test data is isolated (cannot split into two tests)
- `TestGetPoints` - tests the `/receipts/{id}/points` endpoint
- `TestProcessReceipts_Bad_*` - tests the process endpoint with bad data, including empty and invalid receipt data
- `TestGetPoints_Bad_*` - tests the get points endpoint with bad data, including empty and invalid ids

## Discussion

how to handle duplicate receipts? - initially thought this should not be allowed.

- Idea: generate unique ID based on complete receipt body, and check if ID already exists in map
- Reconsideration: duplicate receipts should be allowed, since identical transactions can occur
  - For example, two customers can purchase at different registers but still have the same time, store, and items
  - additionally, having an unbounded id length based on the whole receipt body could be problematic
- Solution: generate unique ID based on date, time, retailer, and customerID, and check if ID already exists in map
  - We would need to know the CustomerID to determine if a receipt is a duplicate

## Conclusion

- Despite not being a seasoned Go developer, I've been learning Go recently and am enjoying it! This exercise showed me how elegant Go can be when implementing a REST API.
- I've enjoyed working on this exercise, and I look forward to hearing your feedback!
