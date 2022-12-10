# Discussion of Exercise

## Overview

Overally I enjoyed doing this exercise! I'm getting more comfortable with REST apis, so I decided to challenge myself and write the program in Go. Learning testing with a Go package was also a new experience for me. I'm looking forward to hearing your feedback! I've included a few discussion points below on some of the decisions I made and some of the things I would have liked to do differently.

## Discussion Points

### Gin

I decided to use the Gin web framework upon hearing it works well to make performant and succintly written Go web applications. I found it easy to use and reduced boilerplate relating to routing and error handling. I also found it easy to use the `shouldBindJSON` function to bind incoming JSON to a struct, and the `c.JSON` function to return JSON responses.

### Data Storage

#### Hashmap

I chose to store data (including receipt data and points for a receipt) in structs in memory, accessed via a Hashmap with the key being the unique ID generated upon processing of a receipt. I chose this approach since I anticipated that the `Get Points` endpoint would see equal to or more traffic than the `Process Receipt` endpoint, so making retrieval of points constant time via a Hashmap would make the request performant. I also chose to create a struct to hold a `<Receipt, Points>` pair in order to form a tight coupling between a specific receipt and the points processed for that receipt. Although the receipt data was not accessed via another route, I anticipated that in a production environment, there may be routes added that access or mutate receipt data, so I left the receipt data in the struct.

#### Concerns

My concern with this approach is that a unique ID is generated for each valid inbound receipt, which means that identical receipt data could receive different IDs, resulting in many duplicates in the hashmap. I've discussed this in more depth on the discussion section of the main README.

### Binding incoming JSON and Validation

#### Binding

I used Gin's binding functionality to de-serialize incoming JSON data. I found this approach favorable since it leverages Go's unmarshalling feature, inferring from the struct field's JSON tag what JSON data to bind to. I discovered that if a JSON value is not appropriate for a given type, the Bind function will throw an error, making type validation succinct.

#### Validation

For application-specific validation, I called a validator function on the Receipt struct formed, checking for valid date, time, and other field formats, including checking the length of items as well as negative totals and prices. This is because all the receipt fields are used to process points, so I decided to have a no tolerance policy during validation, since any flexibility would lead to unpredictable point outputs. Therefore the internal func for processing points assumes a completely valid Receipt struct.

#### Validation - Concerns

Although I believe `validateReceipt` does a satisfactory job of validation coverage, using Gin's validation features could've saved time and possible performed a more effective validation. With a simple `binding:"required"` tag, Gin would've checked for empty fields and thrown a `400: bad request` error. A formatting standard could also be included in the tag, making validation of dates and times relatively easy.

I decided to use my own validation function since I wanted to be able to check for the length of the items array, as well as negative totals and prices as part of the validation process. I also used the  `shouldBindJSON` variant of the `BindJSON` function, such that I can return a custom error message that conforms to the api standard.

### Testing

#### Unit Testing

I used a unit testing approach to test my program. Tests included those that check for appropriate responses to valid and invalid inputs to the `/receipts/process` route, as well as those to the `/receipts/{id}/points` route respectively.

#### Unit Testing - Concerns

Since tests to the `/receipts/{id}/points` route rely on already created resources server side, I stored the id of the receipt created in the test for the `/receipts/process` route in a global variable, and used it in the tests for the `/receipts/{id}/points` route. I'm not sure if there was a more elegant way to do this, but I found this approach to work for my purposes.

I also believe it would have been more appropriate to use [table-driven tests](https://pkg.go.dev/testing#hdr-Subtests_and_Sub_benchmarks) in order to reduce code duplication, as I was writing a unique JSON body for each test. However, I found it difficult to imagine how to include various types of JSON inputs without code duplication, so I decided to leave it as is.

#### Testing Within A Docker Container

I also chose to run the same tests within a Docker container specified built from the included Dockerfile, but ran into the issue of each Unit Test having its own memory/stack space, and therefore not being able to access global vars reliably. Therefore, to validate the Get Points route, I would run both routes within a single unit test to ensure shared stated. I'm not sure if it was necessary to run identical tests within a Docker container, nor if this adjustment to the tests can even qualify as a unit test, but I wanted to try it out, and wanted to be able to run unit tests within a container.

#### Testing - Coverage

My code coverage with my current code tests is as shown:
`ok      github.com/joelyoshiya/Fetch_Rewards_Backend_Coding_Challenge/main      8.058s  coverage: 91.2% of statements`

After viewing the uncovered code, I discovered that it was mainly boolean return statements (which are safe and predictable), as well as the `main` function, which I did not test, since it is mainly library functions and is indirectly covered by every other test.

## Conclusion

As far as the scalability of this project, I would expand the program to use a database, as well as have a more flexible validator that can understand and validate a wider range of JSON inputs, as well as have a more modular points calculator that can be changed on the fly. I would also reformat the testing library to include table-driven tests, and to include benchmarking.
