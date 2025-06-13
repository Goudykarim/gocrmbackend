package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// TestGetCustomersHandler tests the happy path of submitting a well-formed GET /customers request
func TestGetCustomersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/customers", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getCustomers)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("getCustomers returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type does not match: got %v want %v",
			ctype, "application/json")
	}
}

// TestAddCustomerHandler tests the happy path of submitting a well-formed POST /customers request
func TestAddCustomerHandler(t *testing.T) {
	// The request body includes a valid JSON payload
	requestBody := strings.NewReader(`
		{
			"name": "Test User",
			"role": "Tester",
			"email": "test@example.com",
			"phone": "1234567890",
			"contacted": true
		}
	`)

	req, err := http.NewRequest("POST", "/customers", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	// The handler for adding a customer is tested directly
	handler := http.HandlerFunc(addCustomer)
	handler.ServeHTTP(rr, req)

	// We expect a 201 Created status for a successful addition
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("addCustomer returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

// TestDeleteCustomerHandler tests the unhappy path of deleting a user that doesn't exist
func TestDeleteCustomerHandler(t *testing.T) {
	// The request now uses an integer ID that is unlikely to exist
	req, err := http.NewRequest("DELETE", "/customers/9999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	// To test a handler with URL parameters like {id}, we need a router
	router := mux.NewRouter()
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	router.ServeHTTP(rr, req)

	// We expect a 404 Not Found because the customer ID 9999 doesn't exist
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("deleteCustomer returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// TestGetCustomerHandler tests the unhappy path of getting a user that doesn't exist
func TestGetCustomerHandler(t *testing.T) {
	// The request uses an integer ID that is unlikely to exist
	req, err := http.NewRequest("GET", "/customers/9999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	// A router is needed here as well to handle the {id} parameter
	router := mux.NewRouter()
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.ServeHTTP(rr, req)

	// We expect a 404 Not Found because the customer does not exist
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("getCustomer returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
