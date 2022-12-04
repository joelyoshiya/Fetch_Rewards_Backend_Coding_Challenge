package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// define structs
type item struct {
	ShortDescription string  `json:"shortDescription"` // from client
	Price            float64 `json:"price"`            // from client
}

type items struct {
	Items []item `json:"items"` // from client
}

type receipt struct {
	ID           string  `json:"id"`           // from service
	Points       int     `json:"points"`       // from service
	Retailer     string  `json:"retailer"`     // from client
	PurchaseDate string  `json:"purchaseDate"` // from client
	PurchaseTime string  `json:"purchaseTime"` // from client
	Total        float64 `json:"total"`        // from client
	Items        []item  `json:"items"`        // from client
}

type receipts struct {
	Receipts []receipt `json:"receipts"`
}

// setup router
func setupRouter() *gin.Engine {
	r := gin.Default()
	// define routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// process receipts
	r.POST("/receipts/process", processReceipts)
	// get points
	r.GET("/receipts/:id/points", getPoints)

	return r
}

func main() {
	r := setupRouter()
	// run server
	r.Run(":8080")
}

// Path: /receipts/process
// Method: POST
// Payload: Receipt JSON
// Response: JSON containing an id for the receipt.
// Description:
// Takes in a JSON receipt (see example in the example directory) and returns a JSON object with an ID generated by your code.
// The ID returned is the ID that should be passed into /receipts/{id}/points to get the number of points the receipt was awarded.
// How many points should be earned are defined by the rules below.
// Reminder: Data does not need to survive an application restart. This is to allow you to use in-memory solutions to track any data generated by this endpoint.
func processReceipts(c *gin.Context) {
	// process receipts
}

// Path: /receipts/{id}/points
// Method: GET
// Response: A JSON object containing the number of points awarded.
// A simple Getter endpoint that looks up the receipt by the ID and returns an object specifying the points awarded.
func getPoints(c *gin.Context) {
	// get ID
	// get points

}
