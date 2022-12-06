// referring to: https://github.com/olliefr/docker-gs-ping/blob/main/main_test.go
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

var body_valid_1_docker = []byte(`{
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

var body_valid_2_docker = []byte(`{
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

// body with no purchaseDate entry - should fail
var body_bad_empty_date_docker = []byte(`{
	"retailer": "M&M Corner Market",
	"purchaseDate": "",
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
	  }
	],
	"total": "9.00"
	  }`)

var body_bad_empty_items_arr_docker = []byte(`{
	"retailer": "M&M Corner Market",
	"purchaseDate": "2022-03-20",
	"purchaseTime": "14:33",
	"items": []
	"total": "9.00"
	  }`)

var body_bad_empty_items_elts_docker = []byte(`{
	"retailer": "M&M Corner Market",
	"purchaseDate": "2022-03-20",
	"purchaseTime": "14:33",
	"items": [
	  {
		"shortDescription": "",
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

var body_bad_negative_total_docker = []byte(`{
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
	  }
	],
	"total": "-7.75"
	  }`)

var body_bad_negative_price_docker = []byte(`{
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
		"price": "-2.25"
	  }
	],
	"total": "7.75"
	  }`)

var body_bad_empty_docker = []byte(``)

// expected points for body1 and body2
var body_valid_1_pts_docker = 25
var body_valid_2_pts_docker = 109

// post and process receipt to server running on docker container
// also test if valid points returned
func TestProcessReceipt_1_Docker(t *testing.T) {
	// variables used
	var id string
	var pts int64
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Post(fmt.Sprintf("http://localhost:%s/receipts/process", resource.GetPort("8080/tcp")), "application/json", bytes.NewBuffer(body_valid_1_docker))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("received non-201 response: %d", resp.StatusCode)
		}
		var body []byte
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		id, err = jsonparser.GetString(body, "id")
		if err != nil {
			return err
		}
		resp, err = http.Get(fmt.Sprintf("http://localhost:%s/receipts/%s/points", resource.GetPort("8080/tcp"), id))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
		}
		var body2 []byte
		body2, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		pts, err = jsonparser.GetInt(body2, "points")
		if err != nil {
			return err
		}
		return nil
	}))
	require.Equal(t, body_valid_1_pts_docker, int(pts))

}

// post and process receipt to server running on docker container
// also test if valid points returned
func TestProcessReceipt_2_Docker(t *testing.T) {
	// variables used
	var id string
	var pts int64
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Post(fmt.Sprintf("http://localhost:%s/receipts/process", resource.GetPort("8080/tcp")), "application/json", bytes.NewBuffer(body_valid_2_docker))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("received non-201 response: %d", resp.StatusCode)
		}
		var body []byte
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		id, err = jsonparser.GetString(body, "id")
		if err != nil {
			return err
		}
		resp, err = http.Get(fmt.Sprintf("http://localhost:%s/receipts/%s/points", resource.GetPort("8080/tcp"), id))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
		}
		var body2 []byte
		body2, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		pts, err = jsonparser.GetInt(body2, "points")
		if err != nil {
			return err
		}
		return nil
	}))
	require.Equal(t, body_valid_2_pts_docker, int(pts))

}

// Bad Input - Process Receipt

func TestProcessReceipt_Bad_Date_Docker(t *testing.T) {
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Post(fmt.Sprintf("http://localhost:%s/receipts/process", resource.GetPort("8080/tcp")), "application/json", bytes.NewBuffer(body_bad_empty_date_docker))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("received non-400 response: %d", resp.StatusCode)
		}
		return nil
	}))
}

func TestProcessReceipt_Bad_Items_Empty_Arr_Docker(t *testing.T) {
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Post(fmt.Sprintf("http://localhost:%s/receipts/process", resource.GetPort("8080/tcp")), "application/json", bytes.NewBuffer(body_bad_empty_items_arr_docker))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("received non-400 response: %d", resp.StatusCode)
		}
		return nil
	}))
}

func TestProcessReceipt_Bad_Items_Empty_Item_Docker(t *testing.T) {
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Post(fmt.Sprintf("http://localhost:%s/receipts/process", resource.GetPort("8080/tcp")), "application/json", bytes.NewBuffer(body_bad_empty_items_elts_docker))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("received non-400 response: %d", resp.StatusCode)
		}
		return nil
	}))
}

func TestProcessReceipt_Bad_Negative_Total_Docker(t *testing.T) {
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Post(fmt.Sprintf("http://localhost:%s/receipts/process", resource.GetPort("8080/tcp")), "application/json", bytes.NewBuffer(body_bad_negative_total_docker))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("received non-400 response: %d", resp.StatusCode)
		}
		return nil
	}))
}

func TestProcessReceipt_Bad_Negative_Item_Price_Docker(t *testing.T) {
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Post(fmt.Sprintf("http://localhost:%s/receipts/process", resource.GetPort("8080/tcp")), "application/json", bytes.NewBuffer(body_bad_negative_price_docker))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("received non-400 response: %d", resp.StatusCode)
		}
		return nil
	}))
}

func TestProcessReceipt_Bad_Empty_Body_Docker(t *testing.T) {
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Post(fmt.Sprintf("http://localhost:%s/receipts/process", resource.GetPort("8080/tcp")), "application/json", bytes.NewBuffer(body_bad_empty_docker))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("received non-400 response: %d", resp.StatusCode)
		}
		return nil
	}))
}

// Bad Input - Get Points
func TestGetPoints_Bad_ID_Docker(t *testing.T) {
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Get(fmt.Sprintf("http://localhost:%s/receipts/%s/points/", resource.GetPort("8080/tcp"), "bad_id"))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("received non-400 response: %d", resp.StatusCode)
		}
		return nil
	}))
}

func TestGetPoints_Bad_Empty_ID_Docker(t *testing.T) {
	// setup docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	resource, err := pool.Run("receipt-processor-service", "latest", nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}()
	// wait for docker to start
	require.NoError(t, pool.Retry(func() error {
		var err error
		var resp *http.Response
		resp, err = http.Get(fmt.Sprintf("http://localhost:%s/receipts/%s/points/", resource.GetPort("8080/tcp"), ""))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("received non-400 response: %d", resp.StatusCode)
		}
		return nil
	}))
}
