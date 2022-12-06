# Receipt Processor - My Solution

## Approach

- Create a `Receipt` struct to hold the data from the inbound JSON
  - Bind the JSON to the `Receipt` struct
  - Perform validation on the inbound JSON
- Create a `Receipts` struct to hold a map of `ReceiptPoints` structs
  - Create a `ReceiptPoints` struct to hold the points for a receipt
  - Efficient lookup of points for a receipt by ID
  - Calculate the points for a receipt during processing
    - Makes GET requests to `/receipts/{id}/points` more efficient
    - Slightly less performant POST requests to `/receipts/process`, but better to do processing alongside validation


## Assumptions

- No persistence layer is required (for this exercise)
- Negative prices, totals, and points are not inbound/outbound from the API
  - Have error handling to cover these cases

## Dependencies

- Go version 1.18.3 darwin/amd64
- Gin Web Framework
- Other dependencies are listed in `go.mod`

## Execution

I've opted to use Docker to run the application. This allows for a consistent environment across all platforms.

### Build

- Run `docker build -t receipt-processor-service .` at the root.

### Run

- Run `docker run -dp 8080:8080 --name receipt-rest-server receipt-processor-service` to start the service

## Test Environment

- Make sure you're on the most up to date build: `docker build -t receipt-processor-service:latest .`
- Run the command `go test -v ./main` at the root to run the tests.

## Discussion

how to handle duplicate receipts? - initially thought this should not be allowed.

- Idea: generate unique ID based on complete receipt body, and check if ID already exists in map
- Reconsideration: duplicate receipts should be allowed, since identical transactions can occur
  - For example, two customers can purchase at different registers but still have the same time, store, and items
  - additionally, having an unbounded id length based on the whole receipt body could be problematic
- Solution: generate unique ID based on date, time, retailer, and customerID, and check if ID already exists in map
  - We would need to know the CustomerID to determine if a receipt is a duplicate

## Conclusion
