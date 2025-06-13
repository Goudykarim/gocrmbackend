# Go CRM Backend API

This project is a backend server for a Customer Relationship Management (CRM) application, built entirely in Go. It provides a RESTful API for performing standard CRUD (Create, Read, Update, Delete) operations on customer data.

## Features

* **RESTful API:** A well-structured API for managing customers.
* **CRUD Operations:** Full support for creating, reading, updating, and deleting customers.
* **JSON Responses:** All data is returned in JSON format.
* **Method-Based Routing:** Utilizes `gorilla/mux` for clean and efficient routing.
* **Stand-out Feature:** Includes a batch update endpoint to modify multiple customers in a single request.

## Endpoints

| Method | Path                  | Description                               |
|--------|-----------------------|-------------------------------------------|
| `GET`    | `/`                   | Displays an HTML overview of the API.     |
| `GET`    | `/customers`          | Retrieves a list of all customers.        |
| `GET`    | `/customers/{id}`     | Retrieves a single customer by their ID.  |
| `POST`   | `/customers`          | Adds a new customer to the database.      |
| `PUT`    | `/customers/{id}`     | Updates an existing customer's data.      |
| `PUT`    | `/customers/batch`    | Updates multiple customers in one request.|
| `DELETE` | `/customers/{id}`     | Deletes a customer by their ID.           |

## Installation & Usage

### Prerequisites

* [Go](https://golang.org/doc/install) (version 1.18 or higher)
* A tool like [Postman](https://www.postman.com/downloads/) or [cURL](https://curl.se/) to test the API.

### 1. Clone the Repository

```bash
git clone <your-repository-url>
cd crm-backend
```

### 2. Install Dependencies

This project uses Go Modules for dependency management. The `gorilla/mux` package is required.

```bash
go mod tidy
```

### 3. Run the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

### 4. Run Unit Tests

```bash
go test
```

## Example `POST` Request Body

To add a new customer, send a `POST` request to `/customers` with a JSON body like this:

```json
{
    "name": "New Customer",
    "role": "Client",
    "email": "client@email.com",
    "phone": "123-456-7890",
    "contacted": false
}