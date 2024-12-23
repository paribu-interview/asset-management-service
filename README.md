# Asset Management Service

The Asset Management Service (AMS) is a microservice designed for managing assets within wallets. It provides
functionality for creating, depositing, withdrawing, and scheduling asset transactions.

## Features

- Create and manage assets in wallets.
- Deposit and withdraw assets.
- Schedule transactions between wallets.

## Requirements

- Go 1.19 or later
- PostgreSQL
- Docker and Docker Compose (optional for containerized deployment)

## Setup and Installation

### Clone the Repository

```bash
git clone https://github.com/safayildirim/asset-management-service.git
cd asset-management-service
```

## Run Locally

### Set Up the Database

1. Create a new PostgreSQL database.
2. Update the database connection details in the `local.env` file.

   ```env 
   PG_USER=your_db_user
   PG_PASSWORD=your_db_password
   PG_NAME=your_db_name
   PG_HOST=your_db_host
   PG_PORT=your_db_port
   ```

3. Set APP_ENV to `local`.

    ```env
    export APP_ENV=local
    ```
4. Install the dependencies.
   ```bash
   go mod tidy
   ```
5. Run the server.
   ```bash
   go run cmd/main.go
   ```

## Run with Docker

1. Build and run the Docker container.
   ```bash
   docker-compose up --build
   ```

2. The application will start on `http://localhost:8081/api`.

## Migration

Migration is automatically handled by the GORM library. The database schema is created when the application starts.

## API Endpoints

The following endpoints are available:

- `POST /api/assets`: Create a new asset.
- `GET /api/assetss`: Retrieve all assets.
- `POST /api/assets/deposit`: Deposit assets into a wallet.
- `POST /api/assets/withdraw`: Withdraw assets from a wallet.
- `POST /api/transactions/schedule`: Schedule a transaction between wallets.
- `GET /api/transactions`: Retrieve all transactions.
- `DELETE /api/transactions/{id}`: Cancel a scheduled transaction.

### Create a new asset:

- Request:

   ```http
  POST /api/assets
  Content-Type: application/json
   ```
- Request Body:
  ```json
  {
    "wallet_id": 1,
    "name": "BTC",
    "amount": 0
  }
  ```

- Response Body:

   ```json
   {
    "id": 1,
    "created_at": "2022-01-01T00:00:00Z",
    "updated_at": "2022-01-01T00:00:00Z",
    "wallet_id": 1,
    "name": "BTC",
    "amount": 0
   }
   ```
- Response
    - 201 Created: Asset created successfully.
    - 400 Bad Request: Invalid input.
    - 409 Conflict: Asset already exists.

### Retrieve all assets:

- Request:

   ```http
   GET /api/assets
   ```
- Query Parameters:
    - `id`: Filter assets by ID.
    - `wallet_id`: Filter assets by wallet ID.
    - `name`: Filter assets by name.


- Response Body:

   ```json
   {
    "data": [
        {
            "id": 1,
            "created_at": "2022-01-01T00:00:00Z",
            "updated_at": "2022-01-01T00:00:00Z",
            "wallet_id": 1,
            "name": "BTC",
            "amount": 0
        }
    ]
   }
   ```
- Response
    - 200 OK: Assets retrieved successfully.
    - 400 Bad Request: Invalid input.
    - 500 Internal Server Error: Server error.

### Deposit assets into a wallet:

- Request:

  ```http
  POST /api/assets/deposit
  ```

- Request Body:
  ```json
  {
    "wallet_id": 1,
    "name": "BTC",
    "amount": 10
  }
  ```
- Response Body:

    ```json
    {
        "id": 1,
        "created_at": "2022-01-01T00:00:00Z",
        "updated_at": "2022-01-01T00:00:00Z",
        "wallet_id": 1,
        "name": "BTC",
        "amount": 10
    }
    ```
- Response
    - 200 OK: Deposit successful.
    - 400 Bad Request: Invalid input.
    - 404 Not Found: Asset not found.

### Withdraw assets from a wallet:

- Request:

  ```http
  POST /api/assets/withdraw
  ```
- Request Body:
  ```json
  {
    "wallet_id": 1,
    "name": "BTC",
    "amount": 5
  }
  ```
- Response Body:

    ```json
    {
        "id": 1,
        "created_at": "2022-01-01T00:00:00Z",
        "updated_at": "2022-01-01T00:00:00Z",
        "wallet_id": 1,
        "name": "BTC",
        "amount": 5
    }
    ```
- Response
    - 200 OK: Withdrawal successful.
    - 400 Bad Request: Invalid input.
    - 404 Not Found: Asset not found.
    - 409 Conflict: Insufficient balance.
    - 500 Internal Server Error: Server error.

### Schedule a transaction between wallets:

- Request:

  ```http
  POST /api/transactions/schedule
  ```
- Request Body:
  ```json
  {
    "source_wallet_id": 1,
    "destination_wallet_id": 2,
    "asset_name": "BTC",
    "amount": 5,
    "scheduled_at": "2022-01-01T00:00:00Z"
  }
  ```
- Response Body:

    ```json
    {
        "id": 1,
        "created_at": "2022-01-01T00:00:00Z",
        "updated_at": "2022-01-01T00:00:00Z",
        "source_wallet_id": 1,
        "destination_wallet_id": 2,
        "asset_name": "BTC",
        "amount": 5,
        "status": "pending",
        "scheduled_at": "2022-01-01T00:00:00Z"
    }
    ```
- Response
    - 200 OK: Transaction scheduled successfully.
    - 400 Bad Request: Invalid input.
    - 404 Not Found: Asset not found.
    - 409 Conflict: Insufficient balance.
    - 500 Internal Server Error: Server error.

### Retrieve all transactions:

- Request:

  ```http
  GET /api/transactions
  ```
- Query Parameters:
    - `id`: Filter transactions by ID.
    - `source_wallet_id`: Filter transactions by source wallet ID.
    - `destination_wallet_id`: Filter transactions by destination wallet ID.
    - `status`: Filter transactions by status.

- Response Body:

    ```json
    {
        "data": [
            {
                "id": 1,
                "created_at": "2022-01-01T00:00:00Z",
                "updated_at": "2022-01-01T00:00:00Z",
                "source_wallet_id": 1,
                "destination_wallet_id": 2,
                "asset_name": "BTC",
                "amount": 5,
                "status": "pending",
                "scheduled_at": "2022-01-01T00:00:00Z"
            }
        ]
    }
    ```

- Response
    - 200 OK: Transactions retrieved successfully.
    - 400 Bad Request: Invalid input.
    - 500 Internal Server Error: Server error.

### Cancel a scheduled transaction:

- Request:

  ```http
  DELETE /api/transactions/1
  ```
- Response Body:

    ```json
    {
        "id": 1,
        "created_at": "2022-01-01T00:00:00Z",
        "updated_at": "2022-01-01T00:00:00Z",
        "source_wallet_id": 1,
        "destination_wallet_id": 2,
        "asset_name": "BTC",
        "amount": 5,
        "status": "cancelled",
        "scheduled_at": "2022-01-01T00:00:00Z"
    }
    ```
- Response
    - 200 OK: Transaction cancelled successfully.
    - 400 Bad Request: Invalid input.
    - 404 Not Found: Transaction not found.
    - 500 Internal Server Error: Server error.

## Testing

Run the tests using the following command:

```bash
go test ./...
```
