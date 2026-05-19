# Bookstore Users API

A Go-based microservice for managing users within the "Bookstore" project ecosystem. It handles user registration, profile management, and authentication-related queries using a MySQL database.

## Technology Stack

*   **Framework:** [Gin](https://github.com/gin-gonic/gin) for high-performance HTTP routing.
*   **API Documentation:** [Swagger UI](https://swagger.io/tools/swagger-ui/) for interactive API exploration and testing.
*   **Configuration:** [GoDotEnv](https://github.com/joho/godotenv) for loading environment variables from a `.env` file during local development.
*   **Primary Datastore:** MySQL.

## Architectural Notes

*   **Clean Architecture:** This service is built following the principles of Clean Architecture, separating the core business logic from external concerns like the database and web framework.
*   **Mutating Validation:** The `User.Validate()` method performs data sanitization **before** validation. It intentionally **mutates** the `User` object it's called on by trimming whitespace and converting email to lowercase. This ensures data consistency before persistence.

## API Documentation (Swagger)

This service provides an interactive API documentation powered by Swagger UI. It allows you to explore all available endpoints, view their models, and execute requests directly from your browser.

Once the application is running, you can access the Swagger interface at:

**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

*(Note: The port `8080` should match the `GIN_PORT` configured in your `.env` file.)*

To regenerate Swagger docs after modifying controller annotations:

```bash
swag init --parseDependency --parseInternal
```

## API Behavior

### Response Payload Control (`X-Public` Header)

Several endpoints support a mechanism to control the level of detail in the JSON response via the `X-Public` HTTP header.

*   **`X-Public: true`**: Returns a **public** representation of the user object, with a limited set of non-sensitive fields.
*   **Header Absent or Not `true`**: Returns a **private** representation with a more complete set of user fields (excluding the password).

## Prerequisites

- Go (1.18 or newer)
- MySQL
- A running instance of the **Bookstore OAuth microservice** for validating access tokens.

## Configuration

This project uses [GoDotEnv](https://github.com/joho/godotenv) to load configuration for local development. Create a `.env` file in the project root and use the following template.

```env
# Application Settings
GIN_PORT=:8080
CTX_TIMEOUT=2s

# Logger
LEVEL=info
OUTPUT_PATHS=stdout

# ----- Dependencies -----
# URL for the external OAuth microservice. Required for startup.
OAUTH_API_BASE_URL=http://localhost:8081
OAUTH_TIMEOUT=150ms

# ----- Database Connection -----
DB_USER=root
DB_PASSWORD=your_secret_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=users_db
```

## Getting Started

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd <repository-directory>
    ```

2.  **Install/Update dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Set up the database:**
    Ensure you have a MySQL database created and that its name matches `DB_NAME` in your `.env` file. Apply the required schema migrations.

4.  **Run the application:**
    ```bash
    go run main.go
    ```

## Running Tests

To run the full suite of tests:

```bash
go test ./...
```

To view test coverage:

```bash
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## API Endpoints Overview

Below is a high-level overview. For detailed information on request bodies, parameters, and response models, please refer to the [Swagger documentation](#api-documentation-swagger).

-   **`POST /users`**: Create a new user.
-   **`POST /users/login`**: Log in a user.
-   **`GET /users/:user_id`**: Get user details.
    -   **Authentication**: Protected. Requires a valid `Bearer` token.
-   **`PUT /users/:user_id`**: Update an entire user's information.
    -   **Note**: If a user with the specified ID does not exist, the API will return a `404 Not Found` error.
-   **`PATCH /users/:user_id`**: Partially update a user's information.
    -   **Note**: If a user with the specified ID does not exist, the API will return a `404 Not Found` error.
-   **`DELETE /users/:user_id`**: Delete a user.
    -   **Note**: If a user with the specified ID does not exist, the API will return a `404 Not Found` error.
-   **`GET /internal/users/search`**: Search for users by status.