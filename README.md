# Bookstore Users API

A Go-based microservice for managing users within the "Bookstore" project ecosystem. It handles user registration, profile management, and authentication-related queries using a MySQL database.

## Technology Stack

*   **Framework:** [Gin](https://github.com/gin-gonic/gin) for high-performance HTTP routing.
*   **Authentication:** Handled via the custom `bookstore-oauth-go` SDK, which performs local, stateless JWT validation.
*   **API Documentation:** [Swagger UI](https://swagger.io/tools/swagger-ui/) for interactive API exploration.
*   **Configuration:** [GoDotEnv](https://github.com/joho/godotenv) for loading environment variables during local development.
*   **Primary Datastore:** MySQL.

## Architectural Notes

*   **Clean Architecture:** This service is built following the principles of Clean Architecture, separating core business logic from external concerns.
*   **Delegated JWT Authentication:** This service does not handle JWT validation logic directly. Instead, it relies on the **`bookstore-oauth-go`** library, which was refactored to perform local, stateless validation without network calls. This keeps the `users-api` clean of authentication specifics and promotes reusable, consistent security across the microservice ecosystem.
*   **Mutating Validation:** The `User.Validate()` method performs data sanitization **before** validation by trimming whitespace and converting email to lowercase.

## API Documentation (Swagger)

This service provides an interactive API documentation powered by Swagger UI. It allows you to explore all available endpoints, view their models, and execute requests directly from your browser.

The Swagger UI is only available when the `APP_ENV` is set to `development`. Once the application is running in development mode, you can access it at:

**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

*(Note: The port `8080` should match the `GIN_PORT` configured in your `.env` file.)*

To regenerate Swagger docs after modifying controller annotations:

```bash
swag init --parseDependency --parseInternal
```

## Prerequisites

- Go (1.23 or newer)
- MySQL
- An up-to-date version of the `bookstore-oauth-go` library.

## Configuration

This project uses [GoDotEnv](https://github.com/joho/godotenv) to load configuration for local development. Create a `.env` file in the project root and use the following template.

```env
# Application Settings
GIN_PORT=:8080
CTX_TIMEOUT=2s
APP_ENV=development # Use 'development' to enable Swagger, or 'production'

# ----- JWT Authentication -----
# This secret is used by the bookstore-oauth-go SDK to verify JWT signatures.
# It must be the same key used by the OAuth API to sign tokens.
SECRET_KEY=your_super_secret_signing_key_must_be_long_and_secure

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

The API is now split into public and protected endpoints. For detailed information on models and error responses, please refer to the [Swagger documentation](#api-documentation-swagger).

### Public Endpoints
These endpoints do not require authentication.
-   **`POST /users`**: Create a new user.
-   **`POST /users/login`**: Log in a user.

### Protected Endpoints
These endpoints are protected by a JWT middleware and require a valid `Authorization: Bearer <token>` header. They operate on the user identified by the `user_id` within the token.

-   **`GET /users`**: Get the details of the currently authenticated user.
-   **`PUT /users`**: Update the entire profile of the currently authenticated user.
-   **`PATCH /users`**: Partially update the profile of the currently authenticated user.
-   **`DELETE /users`**: Delete the currently authenticated user.
-   **`GET /internal/users/search`**: Search for users by status. (This endpoint is also protected).