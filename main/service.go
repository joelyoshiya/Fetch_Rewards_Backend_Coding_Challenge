package main

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Struct definitions & constructors

// Struct representing inbound receipt
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

// Struct representing Receipt Points pair - used for storing receipts/points pairs
type ReceiptPoints struct {
	Receipt Receipt `json:"receipt"`
	Points  int     `json:"points"`
}

// Struct representing Receipts - internal storage of receipts/points
type Receipts struct {
	// store a map of receipts, points pairs accessed via ID
	ReceiptsMap map[string]ReceiptPoints `json:"receipts"`
}

// Constructor for Receipts
func NewReceipts() *Receipts {
	var rs Receipts
	rs.ReceiptsMap = make(map[string]ReceiptPoints)
	return &rs
}

// Internal data

// Global receipts object - in place of persisting data
var rs = NewReceipts() // pointer to Receipts object

// Internal functions - not exported

// Setup router
func setupRouter() *gin.Engine {
	r := gin.Default()
	// define routes
	r.POST("/receipts/process", processReceipt)
	r.GET("/receipts/:id/points", getPoints)
	return r
}

// Validate receipt - make sure all fields are populated and valid
// Any invalid fields will result in an invalid receipt
func validateReceipt(r Receipt) bool {
	// check if all fields populated
	if r.Retailer == "" || r.Total == "" || r.PurchaseDate == "" || r.PurchaseTime == "" || r.Items == nil {
		return false
	}
	// check if purchase date is valid
	_, err := time.Parse("2006-01-02", r.PurchaseDate)
	if err != nil {
		return false
	}
	// check if purchase time is valid
	_, err = time.Parse("15:04", r.PurchaseTime)
	if err != nil {
		return false
	}
	// check if total is valid
	_, err = strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return false
	}
	// check if total is negative
	total, _ := strconv.ParseFloat(r.Total, 64)
	if total < 0 {
		return false
	}
	// check if r.Items meets minimum length requirement of 1
	if r.Items != nil && len(r.Items) < 1 {
		return false
	}
	// check if bad data in r.Items
	for _, item := range r.Items {
		// check for empty vals
		if item.ShortDescription == "" || item.Price == "" {
			return false
		}
		// check for invalid price
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return false
		}
		// check that price is above 0
		if price < 0 {
			return false
		}

	}
	return true
}

// Calculate points for receipt - based on ruleset given
// Assumes a valid receipt is passed in
func processPoints(r Receipt) int {
	// 1 point for every alphanumeric character in the retailer name.
	// define regex for alphanumeric characters - referring to: https://gosamples.dev/remove-non-alphanumeric/
	nonAlphaNumericRegex := regexp.MustCompile("[^a-zA-Z0-9]+")
	// replace non alphanumeric characters with empty string
	retailer := nonAlphaNumericRegex.ReplaceAllString(r.Retailer, "")
	// count alphanumeric characters
	retailerPoints := len(retailer)

	// 50 points if the total is a round dollar amount with no cents. 25 points if the total is a multiple of 0.25
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
	itemCount := len(r.Items)
	itemCountPoints := (itemCount / 2) * 5

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	itemPoints := 0
	for _, item := range r.Items {
		itemDesc := strings.Trim(item.ShortDescription, " ")
		// check if trimmmed length of item description is a multiple of 3
		if len(itemDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			// calculate points
			itemPoints += int(math.Ceil(price * 0.2)) // add to item points for each item
		}
	}

	// 6 points if the day in the purchase date is odd.
	// parse purchase date to int
	datePoints := 0
	// define a non numeric regex
	nonNumericRegex := regexp.MustCompile("[^0-9]+")
	// replace non numeric characters with empty string
	purchaseDate := nonNumericRegex.ReplaceAllString(r.PurchaseDate, "")
	purchaseDateInt, _ := strconv.Atoi(purchaseDate)
	// check if day is odd
	if purchaseDateInt%2 != 0 {
		datePoints += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm
	timePoints := 0
	purchaseTime := nonNumericRegex.ReplaceAllString(r.PurchaseTime, "")
	purchaseTimeInt, _ := strconv.Atoi(purchaseTime)
	// check if time is between 2:00pm and 4:00pm (after and before 1400 and 1600)
	if purchaseTimeInt > 1400 && purchaseTimeInt < 1600 {
		timePoints += 10
	}
	return retailerPoints + totalPoints + itemCountPoints + itemPoints + datePoints + timePoints
}

// Internal Route Functions

// Path: /receipts/process
// Method: POST
// Payload: Receipt JSON
// Response: JSON containing an id for the receipt.
// Description: Takes in a JSON receipt (see example in the example directory) and returns a JSON object with an ID generated by your code.
func processReceipt(c *gin.Context) {
	var r Receipt // to store inbound receipt

	// bind JSON to receipt object - upon error, return bad request
	// unmarshaling JSON to struct, type checking for all fields
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The receipt is invalid"})
		return
	}

	// validate receipt
	if !validateReceipt(r) {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The receipt is invalid"})
		return
	}

	// process points
	points := processPoints(r)

	// generate ID
	id := uuid.New().String()

	// create a ReceiptPoints object and add to receipts map
	(*rs).ReceiptsMap[id] = ReceiptPoints{r, points}

	// return status created and receipt ID
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// Path: /receipts/{id}/points
// Method: GET
// Response: A JSON object containing the number of points awarded.
// Description: A simple Getter endpoint that looks up the receipt by the ID and returns an object specifying the points awarded.
func getPoints(c *gin.Context) {
	// get ID
	id := c.Param("id")
	// get receipt object with ID from receipts
	rp, present := (*rs).ReceiptsMap[id]
	if !present {
		c.JSON(http.StatusNotFound, gin.H{"description": "No receipt found for that id"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"points": rp.Points})
	}
}

// main function - start server
func main() {
	r := setupRouter()
	r.Run() // listen and serve on default port 8080 - otherwise port defined in env variable PORT
}
