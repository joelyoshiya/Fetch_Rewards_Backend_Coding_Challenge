package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var body1 = []byte(`{
	"retailer": "Target",
	"purchaseDate": "2022-01-01",
	"purchaseTime": "13:01",
	"items": [
	  {
		"shortDescription": "Mountain Dew 12PK",
		"price": "6.49"
	  },{
		"shortDescription": "Emils Cheese Pizza",
		"price": "12.25"
	  },{
		"shortDescription": "Knorr Creamy Chicken",
		"price": "1.26"
	  },{
		"shortDescription": "Doritos Nacho Cheese",
		"price": "3.35"
	  },{
		"shortDescription": "Klarbrunn 12PK 12 FL OZ",
		"price": "12.00"
	  }
	],
	"total": "35.35"
  }`)

var body2 = []byte(`{
	"retailer": "M&M Corner Market",
	"purchaseDate": "2022-03-20",
	"purchaseTime": "14:33",
	"items": [
	  {
		"shortDescription": "Gatorade",
		"price": "2.25"
	  },{
		"shortDescription": "Gatorade",
		"price": "2.25"
	  },{
		"shortDescription": "Gatorade",
		"price": "2.25"
	  },{
		"shortDescription": "Gatorade",
		"price": "2.25"
	  }
	],
	"total": "9.00"
  }`)

var body1_pts = 25
var body2_pts = 109

func TestPingRoute(t *testing.T) {
	// set up router, recorder, and request
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestProcessReceipt(t *testing.T) {
	// set up router, recorder, and request
	router := setupRouter()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(body1))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// assert response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"id"`)
	assert.Contains(t, w.Body.String(), `"points"`)
}

func TestGetPoints(t *testing.T) {
	// set up router, recorder, and request
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(body2))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"id"`)
	assert.Contains(t, w.Body.String(), `"points"`)

	// grab id from response
	var resp map[string]interface{} // referring to: https://bitfieldconsulting.com/golang/map-string-interface
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}
	id := resp["id"].(string)

	// use id to query for points
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/receipts/"+id+"/points", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// get points from body:
	var resp2 map[string]interface{} // referring to: https://bitfieldconsulting.com/golang/map-string-interface
	err2 := json.Unmarshal(w.Body.Bytes(), &resp2)
	if err2 != nil {
		t.Fatal(err)
	}
	points := resp2["points"].(float64)

	// check if points valid
	assert.Equal(t, body2_pts, int(points))
}
