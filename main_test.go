package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// TestMain runs before all other tests in the package.
// It's the perfect place to set up and tear down resources like a database connection.
func TestMain(m *testing.M) {
	// Call the initDB function from main.go to set up the database connection.
	initDB()

	// Exit code from the test run
	var code int
	// Use a defer function to ensure the database connection is closed after the tests run.
	defer func() {
		db.Close()
		os.Exit(code)
	}()
	// Run all the tests and store the exit code.
	code = m.Run()
}

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
			"name": "Test User From Test",
			"role": "Tester",
			"email": "test-user@example.com",
			"phone": "9876543210",
			"contacted": false
		}
	`)

	req, err := http.NewRequest("POST", "/customers", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
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
