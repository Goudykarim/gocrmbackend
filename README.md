# **Go CRM Backend API**

This project is a complete, production-ready backend server for a Customer Relationship Management (CRM) application, built entirely in Go. It features a RESTful API for all standard CRUD operations, backed by a PostgreSQL database, and is fully containerized with Docker for seamless deployment to cloud platforms like Heroku.

**Live API URL:** https://karimgoudycrmbackend-0997c5859f5b.herokuapp.com/

## **Features**

- **RESTful API:** A well-structured API for creating, reading, updating, and deleting customer data.
- **PostgreSQL Database:** Uses a robust, persistent SQL database for data storage.
- **Dockerized Environment:** The application is fully containerized, ensuring consistency between local development and production deployment.
- **Cloud Deployment:** Deployed and fully functional on Heroku.
- **Configuration Management:** Securely configured using environment variables for database credentials and port management.
- **Unit & Integration Testing:** Includes a suite of tests to validate the application's logic and database connectivity.
- **Stand-Out Feature: Batch Updates:** Includes a `/customers/batch` endpoint to efficiently update multiple customers in a single API call.

## **Technologies Used**

- **Language:** Go
- **API Framework:** `gorilla/mux`
- **Database:** PostgreSQL
- **Database Driver:** `lib/pq`
- **Containerization:** Docker
- **Deployment:** Heroku

## **API Endpoints**

| **Method** | **Path** | **Description** |
| --- | --- | --- |
| `GET` | `/` | Displays a simple HTML overview of the API. |
| `GET` | `/customers` | Retrieves a list of all customers. |
| `GET` | `/customers/{id}` | Retrieves a single customer by their ID. |
| `POST` | `/customers` | Adds a new customer to the database. |
| `PUT` | `/customers/{id}` | Updates an existing customer's data. |
| `PUT` | `/customers/batch` | (Batch Operation) Updates multiple customers. |
| `DELETE` | `/customers/{id}` | Deletes a customer by their ID. |

## **Local Development Setup**

### **Prerequisites**

- [Go](https://golang.org/doc/install) (version 1.24 or higher)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli)
- A tool like [Postman](https://www.postman.com/downloads/) to test the API.

### **Instructions**

1. **Clone the Repository:**
    
    ```
    git clone https://github.com/Goudykarim/gocrmbackend.git
    cd gocrmbackend
    
    ```
    
2. **Install Dependencies:**
    
    ```
    go mod tidy
    
    ```
    
3. **Set Up Local Database:**
    - Start your local PostgreSQL service.
    - Use `psql` to create a new database named `crm`.
    - Connect to the `crm` database (`\c crm`) and run the `CREATE TABLE` script found in the project guide to create the `customers` table.
4. Run the Server:
    
    The application reads the database credentials from an environment variable. For local development, it falls back to a default.
    
    ```
    go run main.go
    
    ```
    
    The server will start on `http://localhost:8080`.
    
5. Run Tests:
    
    To run the unit and integration tests, ensure your local database is running.
    
    ```
    go test
    
    ```
    

## **Deployment**

This application is deployed on Heroku using Docker. The deployment process is managed via the `Dockerfile` and `heroku.yml` in the repository.

The live application is configured with a Heroku Postgres add-on, and the `DATABASE_URL` is securely passed to the running container as the `CRM_DB_CONNECTION_STRING`.

## **Example `POST` Request Body**

To add a new customer, send a `POST` request to `/customers` with a JSON body like this:

```
{
    "name": "New Customer",
    "role": "Client",
    "email": "new.client@example.com",
    "phone": "555-555-5555",
    "contacted": false
}

```