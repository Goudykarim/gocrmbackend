package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Customer struct represents a single customer
type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

var db *sql.DB

// initDB initializes the database connection.
func initDB() {
	connStr := os.Getenv("CRM_DB_CONNECTION_STRING")
	if connStr == "" {
		connStr = "postgres://karimabdelaziz@localhost/crm?sslmode=disable"
		log.Println("Warning: CRM_DB_CONNECTION_STRING not set. Using fallback for local development.")
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to database.")
}

// getCustomers retrieves all customers from the database
func getCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, role, email, phone, contacted FROM customers ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.Phone, &c.Contacted); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// getCustomer retrieves a single customer
func getCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var c Customer
	err := db.QueryRow("SELECT id, name, role, email, phone, contacted FROM customers WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.Phone, &c.Contacted)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// addCustomer creates a new customer
func addCustomer(w http.ResponseWriter, r *http.Request) {
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO customers(name, role, email, phone, contacted) VALUES($1, $2, $3, $4, $5) RETURNING id",
		c.Name, c.Role, c.Email, c.Phone, c.Contacted).Scan(&c.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

// updateCustomer updates an existing customer
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := db.Exec(
		"UPDATE customers SET name=$1, role=$2, email=$3, phone=$4, contacted=$5 WHERE id=$6",
		c.Name, c.Role, c.Email, c.Phone, c.Contacted, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	c.ID, _ = strconv.Atoi(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// updateCustomersBatch handles PUT requests to /customers/batch for bulk updates
func updateCustomersBatch(w http.ResponseWriter, r *http.Request) {
	var updates []Customer
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	updatedCount := 0
	for _, u := range updates {
		res, err := tx.Exec("UPDATE customers SET name=$1, role=$2, email=$3, phone=$4, contacted=$5 WHERE id=$6",
			u.Name, u.Role, u.Email, u.Phone, u.Contacted, u.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to update customer: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected > 0 {
			updatedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result":            "Batch update completed",
		"customers_updated": updatedCount,
	})
}

// deleteCustomer deletes a customer
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	res, err := db.Exec("DELETE FROM customers WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": fmt.Sprintf("Customer %s deleted", id)})
}

// homePage handles requests to the root URL "/"
func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<h1>Welcome to the CRM Backend API (PostgreSQL Version)</h1>
		<p>This API allows you to manage customer data stored in a PostgreSQL database.</p>
	`)
}

func main() {
	initDB()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	// IMPORTANT: More specific routes must be registered before general ones.
	router.HandleFunc("/customers/batch", updateCustomersBatch).Methods("PUT")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
