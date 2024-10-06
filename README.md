# License Management API

This project is a License Management API built using **Go (Gin)** framework, **SQLite** as the database, and **slog** for logging. The API allows managing user licenses through various operations such as adding, freezing, renewing, and binding licenses.

## Database Schema

The database consists of two tables: `UserLicense` and `TransactionLogs`. Below is the schema:

![Database Schema](https://i.imgur.com/rUtTfGD.jpeg)

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
| GET    | `/get`            | Get details of a license       |
| GET    | `/all-licenses`           | Get details of all licenses    |
| POST   | `/add-license`            | Add a new license              |
| POST   | `/del-license`     | Delete a license               |
| POST   | `/freeze-license`     | Freeze a license               |
| POST   | `/unfreeze-license`   | Unfreeze a license             |
| POST   | `/renew-license`      | Renew a license                |
| POST   | `/bind-license`       | Bind a license to an HWID      |
| POST   | `/unbind-license`     | Unbind a license from HWID     |
| POST   | `/validate-license`   | Validate a license             |



## Technologies Used
- **Go (Gin)**: Fast, lightweight web framework for building the API.
- **SQLite**: Used as the database to store licenses and transaction logs.
- **slog**: Structured logger for handling logs.


## API Documentation
For more detailed information on the API endpoints and request/response formats, refer to the [API Documentation](https://comet-foundation-d6a.notion.site/API-Docs-117387c97789803db40cc889c53a62fb).


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
3. Rename `local.env.example` to `local.env` and configure the API key.
4. Ensure that the port is specified in `local.yaml`.
5. Run the API:
    ```bash
    go run cmd/api/main.go
    ```