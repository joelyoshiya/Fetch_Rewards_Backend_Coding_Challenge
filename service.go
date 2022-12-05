package main

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// define structs
type Item struct {
	ShortDescription string  `json:"shortDescription"` // from client
	Price            float64 `json:"price"`            // from client
}

type Items struct {
	Items []Item `json:"items"` // from client
}

// struct representing inbound receipt
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []struct {
		ShortDescription string `json:"shortDescription"`
		Price            string `json:"price"`
	} `json:"items"`
	Total string `json:"total"`
}

// struct representing Receipt Points pair
type ReceiptPoints struct {
	Receipt Receipt `json:"receipt"`
	Points  int     `json:"points"`
}

// struct representing Receipts
type Receipts struct {
	// store a map of receipts, points pairs accessed via ID
	ReceiptsMap map[string]ReceiptPoints `json:"receipts"`
}

// constructor for receipts
func NewReceipts() *Receipts {
	var rs Receipts
	rs.ReceiptsMap = make(map[string]ReceiptPoints)
	return &rs
}

// global receipts object
var rs = NewReceipts() // pointer to receipt object, in place of persisting data struct

// setup router
func setupRouter() *gin.Engine {
	r := gin.Default()
	// define routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// process receipt
	r.POST("/receipts/process", processReceipt)
	// get points
	r.GET("/receipts/:id/points", getPoints)
	return r
}

// main function - start server
func main() {
	r := setupRouter()
	// run server
	r.Run(":8080")
}

// Internal functions

// Calculate points for receipt - based on ruleset given
func processPoints(r Receipt) int {
	// process points

	// referring to: https://gosamples.dev/remove-non-alphanumeric/
	// 1 point for every alphanumeric character in the retailer name.
	// clean retailer name for non alphanumeric characters
	// define regex for alphanumeric characters
	nonAlphaNumericRegex := regexp.MustCompile("[^a-zA-Z0-9]+")
	// replace non alphanumeric characters with empty string
	retailer := nonAlphaNumericRegex.ReplaceAllString(r.Retailer, "")
	// count alphanumeric characters
	retailerPoints := len(retailer)

	// 50 points if the total is a round dollar amount with no cents.
	// 25 points if the total is a multiple of 0.25
	// parse total to float
	total, _ := strconv.ParseFloat(r.Total, 64)
	// check if total is a round dollar amount
	totalPoints := 0
	if total == float64(int(total)) { // check if total is a round dollar amount
		totalPoints += 50
	}
	if total == float64(int(total*4))/4 { // check if total is a multiple of 0.25
		totalPoints += 25
	}

	// 5 points for every two items on the receipt.
	// count items
	itemCount := len(r.Items)
	// calculate points
	itemCountPoints := (itemCount / 2) * 5

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	itemPoints := 0
	for _, item := range r.Items {
		// trim the item description
		itemDesc := strings.Trim(item.ShortDescription, " ")
		// check if trimmmed length of item description is a multiple of 3
		if len(itemDesc)%3 == 0 {
			// parse price to float
			price, _ := strconv.ParseFloat(item.Price, 64)
			// calculate points
			itemPoints += int(math.Ceil(price * 0.2)) // add to item points for each item
		}
	}
	// 6 points if the day in the purchase date is odd.
	// parse purchase date to int
	datePoints := 0
	// clean purchase date for whitespace and dashes
	purchaseDate := strings.ReplaceAll(r.PurchaseDate, " ", "")
	purchaseDate = strings.ReplaceAll(purchaseDate, "-", "")
	// parse purchase date to int
	purchaseDateInt, _ := strconv.Atoi(purchaseDate)
	// check if day is odd
	if purchaseDateInt%2 != 0 {
		datePoints += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm
	timePoints := 0
	// clean purchase time for whitespace and colons
	purchaseTime := strings.ReplaceAll(r.PurchaseTime, " ", "")
	purchaseTime = strings.ReplaceAll(purchaseTime, ":", "")
	// parse purchase time to int
	purchaseTimeInt, _ := strconv.Atoi(purchaseTime)
	// check if time is between 2:00pm and 4:00pm
	if purchaseTimeInt > 1400 && purchaseTimeInt < 1600 {
		timePoints += 10
	}

	// print all points to see if values are correct
	// fmt.Println("retailerPoints: ", retailerPoints)
	// fmt.Println("totalPoints: ", totalPoints)
	// fmt.Println("itemCountPoints: ", itemCountPoints)
	// fmt.Println("itemPoints: ", itemPoints)
	// fmt.Println("datePoints: ", datePoints)
	// fmt.Println("timePoints: ", timePoints)

	return retailerPoints + totalPoints + itemCountPoints + itemPoints + datePoints + timePoints
}

// Route Functions

// Path: /receipts/process
// Method: POST
// Payload: Receipt JSON
// Response: JSON containing an id for the receipt.
// Description:
// Takes in a JSON receipt (see example in the example directory) and returns a JSON object with an ID generated by your code.
// The ID returned is the ID that should be passed into /receipts/{id}/points to get the number of points the receipt was awarded.
// How many points should be earned are defined by the rules in the README.
// Reminder: Data does not need to survive an application restart. This is to allow you to use in-memory solutions to track any data generated by this endpoint.
func processReceipt(c *gin.Context) {
	// New items, receipt objects
	// var i Items
	var r Receipt
	// var s ServerReceipt

	// generate ID
	id := uuid.New().String()

	// bind JSON to receipt object
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// process points
	points := processPoints(r)

	// create a ReceiptPoints object and add to receipts map
	(*rs).ReceiptsMap[id] = ReceiptPoints{r, points}

	// TODO write logic to determine bad receipt, return error

	// return status created and receipt ID
	c.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}

// Path: /receipts/{id}/points
// Method: GET
// Response: A JSON object containing the number of points awarded.
// A simple Getter endpoint that looks up the receipt by the ID and returns an object specifying the points awarded.
func getPoints(c *gin.Context) {
	// get ID
	id := c.Param("id")
	// get receipt object with ID from receipts
	rp, present := (*rs).ReceiptsMap[id]
	if !present {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "receipt not found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"points": rp.Points})
	}

}
