# License Management API Overview

This License Management API is developed using the **Go (Gin)** framework, with **SQLite** as the database and **slog** for logging. It facilitates various operations for managing user licenses, including adding, freezing, renewing, and binding licenses.

## Database Schema

The API utilizes two primary tables in its SQLite database:

### UserLicense Table
- **id**: Integer (Primary Key)
- **license**: Varchar
- **UserId**: Varchar
- **createdAt**: Timestamp
- **updatedAt**: Timestamp
- **expiresAt**: Timestamp
- **hwid**: Varchar
- **status**: Varchar

### TransactionLogs Table
- **id**: Integer (Primary Key)
- **timestamp**: Datetime
- **description**: Text

## API Endpoints

### License Management Endpoints

| Method | Endpoint               | Description                     |
|--------|-----------------------|---------------------------------|
| GET    | `/ping`               | Health check for the API       |
| GET    | `/license`            | Get details of a license       |
| GET    | `/licenses`           | Get details of all licenses    |
| POST   | `/license`            | Add a new license              |
| POST   | `/license/delete`     | Delete a license               |
| POST   | `/license/freeze`     | Freeze a license               |
| POST   | `/license/unfreeze`   | Unfreeze a license             |
| POST   | `/license/renew`      | Renew a license                |
| POST   | `/license/bind`       | Bind a license to an HWID      |
| POST   | `/license/unbind`     | Unbind a license from HWID     |
| POST   | `/license/validate`   | Validate a license             |

For detailed information on the API endpoints and request/response formats, refer to the [API Documentation](https://comet-foundation-d6a.notion.site/API-Docs-117387c97789803db40cc889c53a62fb).

## Technologies Used
- **Go (Gin)**: A fast and lightweight web framework suitable for building APIs.
- **SQLite**: A compact database solution ideal for handling small data efficiently.
- **slog**: A structured logger that aids in managing logs effectively.

## Setup and Installation

To set up the License Management API, follow these steps:

1. Clone the repository:
    ```bash
    git clone https://github.com/dzhisl/license-manager
    ```
2. Install dependencies:
    ```bash
    go mod download
    ```
3. Rename `local.env.example` to `local.env` to configure environment variables, including the API key.
4. Ensure that the port is specified in `local.yaml`.
5. Run the API:
    ```bash
    go run cmd/api/main.go
    ```
