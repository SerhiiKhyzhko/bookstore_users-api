# Bookstore Users API

A Go-based microservice for managing users within the "Bookstore" project ecosystem. It handles user registration, profile management, and identity verification.

## Technology Stack

*   **Framework:** [Gin](https://github.com/gin-gonic/gin) for high-performance HTTP routing and handling.
*   **HTTP Client:** [Resty](https://github.com/go-resty/resty) for expressive and simple client-side HTTP requests to other services.
*   **Configuration:** [GoDotEnv](https://github.com/joho/godotenv) for loading environment variables from a `.env` file during local development.

## Architectural Notes

This service is built following the principles of **Clean Architecture**, separating the core business logic from external concerns.

*   **Extensible Datasources:** The `datasources` directory is structured to support multiple database backends. The current MySQL implementation can be found in `datasources/mysql/`, and other databases (e.g., PostgreSQL) can be added in parallel without affecting existing code.

*   **Mutating Validation:** The `User.Validate()` method performs data sanitization **before** validation. It intentionally **mutates** the `User` object it's called on: it trims whitespace from fields and converts the email to lowercase. This ensures data consistency before persistence.

## API Behavior

### Response Payload Control (`X-Public` Header)

Several endpoints support a mechanism to control the level of detail in the JSON response via the `X-Public` HTTP header.

*   **`X-Public: true`**: When this header is present and set to `true`, the API returns a **public** representation of the user object, containing a limited set of non-sensitive fields.
*   **Header Absent or Not `true`**: If the header is not provided or has any other value, the API returns a **private** representation with a more complete set of user fields.

This allows the same endpoint to serve different consumers (e.g., a public-facing frontend vs. an internal admin panel) securely.

## Prerequisites

- Go (1.18 or newer)
- MySQL
- A running instance of the **Bookstore OAuth microservice** that this service can call.

## Configuration

This project uses [GoDotEnv](https://github.com/joho/godotenv) to load configuration for local development. Simply create a `.env` file in the project root and use the following template. The application will automatically load these variables on startup.

```env
# Application Settings
GIN_PORT=:8080               # The port for the Gin server to listen on (the leading colon is important)
CTX_TIMEOUT=2s               # Optional: Default timeout for requests. Defaults to 2s if not set.

# Logger
LEVEL=info                   # e.g., debug, info, warn, error
OUTPUT_PATHS=stdout          # Can be stdout, stderr, or a file path

# ----- Dependencies -----
# URL for the external OAuth microservice. Required for startup.
OAUTH_API_BASE_URL=http://localhost:8081
OAUTH_TIMEOUT=150ms          # Optional: Timeout for OAuth calls. Defaults to 150ms.

# ----- Database Connection -----
DB_USER=root
DB_PASSWORD=your_secret_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=users_db             # The name of the database to use
```

## Getting Started

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd <repository-directory>
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Set up the database:**
    Ensure you have a MySQL database created and that its name matches the `DB_NAME` in your `.env` file. Apply the required schema migrations.

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

## API Endpoints

This is a high-level overview. For a detailed and interactive specification, please refer to the upcoming Swagger documentation.

The `X-Public` header can be used on endpoints that return a user object to control response detail.

- `POST /users`: Create a new user.
- `POST /users/login`: Log in a user.
- `GET /users/:user_id`: Get user details.
  - **Authentication**: This endpoint is protected. To use it, you must first **obtain an `access_token` from the corresponding OAuth API**. Then, include this token in the `Authorization` header as a `Bearer` token.
    ```
    Authorization: Bearer <your_access_token>
    ```
- `PUT /users/:user_id`: Update a user's information.
- `PATCH /users/:user_id`: Partially update a user's information.
- `DELETE /users/:user_id`: Delete a user.
- `GET /internal/users/search`: Search for users by status.