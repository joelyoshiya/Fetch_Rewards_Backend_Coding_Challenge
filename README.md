# Receipt Processor - My Solution

## Approach

## Assumptions

- No persistence layer is required (for this exercise)
- Negative prices, totals, and points are not inbound/outbound from the API
  - Have error handling to cover these cases

## Dependencies

- Go version 1.18.3 darwin/amd64
- Gin Web Framework
- Other dependencies are listed in `go.mod`

## Test Environment

Run tests via `go test` in the `/main` directory

## Discussion

how to handle duplicate receipts? - initially thought this should not be allowed.

- Idea: generate unique ID based on complete receipt body, and check if ID already exists in map
- Reconsideration: duplicate receipts should be allowed, since identical transactions can occur
  - For example, two customers can purchase at different registers but still have the same time, store, and items
  - additionally, having an unbounded id length based on the whole receipt body could be problematic
- Solution: generate unique ID based on date, time, retailer, and customerID, and check if ID already exists in map
  - We would need to know the CustomerID to determine if a receipt is a duplicate

## Conclusion
