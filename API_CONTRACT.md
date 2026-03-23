# API Contracts

## Create Account

### Purpose: 
Create an account for a user (cardholder).

### Path: 
POST /account

### Request: 
| Field Name  | Data Type | Required | Example          | Description                                               |
|-------------|-----------|----------|------------------|-----------------------------------------------------------|
| msg_id      | string    | Yes      | "some_unique_id" | Unique identifier for the request                         |
| document_id | string    | Yes      | "12345678900"    | Government-issued identifier for the cardholder           |
| currency    | string    | Yes      | "USD"            | 3 digit ISO currency code for the account                 |

Example:
```json
{
  "msg_id": "some_unique_id",
  "document_id": "12345678900",
  "currency": "USD"
}
```

### Response:
#### Success : 200 
| Field Name  | Data Type | Required | Example | Description                          |
|-------------|-----------|----------|---------|--------------------------------------|
| account_id  | int64     | Yes      | 123     | Unique identifier for the new account|

Example:
```json
{
  "account_id": 1
}
```

#### Errors (4xx/5xx)
| Field Name | Data Type | Required | Example            | Description                                 |
|------------|-----------|----------|--------------------|---------------------------------------------|
| code       | string    | Yes      | "INVALID_ARGUMENT" | Error code, describes type of error         |
| message    | string    | Yes      | "invalid request"  | Human-readable message explaining the error |

Example:
```json
{
  "code": "INVALID_ARGUMENT",
  "message": "invalid request"
}
```

## Get Account

### Purpose:
Retrieve account information by account ID.

### Path:
GET /account/{account_id}

### Request Params:
| Param Name         | Data Type | Required | Example           | Description                                  |
|--------------------|-----------|----------|-------------------|----------------------------------------------|
| account_id (path)  | int64     | Yes      | 123               | Unique identifier of the account to retrieve |
| msgID (query)      | string    | Yes      | "some_unique_id"  | Unique identifier for request                |

Example Request:
```
GET /account/1?msgID=some_unique_id
```

### Response:
#### Success : 200
| Field Name   | Data Type | Required | Example         | Description                                      |
|--------------|-----------|----------|-----------------|--------------------------------------------------|
| account_id   | int64     | Yes      | 123             | Account's unique ID                              |
| document_id  | string    | Yes      | "12345678900"   | Cardholder's document/identifier                 |
| currency     | string    | Yes      | "USD"           | ISO currency code for the account                |
| status       | string    | Yes      | "ACTIVE"        | Account status (e.g. ACTIVE, BLOCKED, etc.)      |

Example:
```json
{
  "account_id": 1,
  "document_id": "12345678900",
  "currency": "USD",
  "status": "ACTIVE"
}
```

#### Errors (4xx/5xx)
| Field Name | Data Type | Required | Example             | Description                               |
|------------|-----------|----------|---------------------|-------------------------------------------|
| code       | string    | Yes      | "NOT_FOUND"         | Error code, such as NOT_FOUND             |
| message    | string    | Yes      | "account not found" | Human-readable error details              |

Example:
```json
{
  "code": "NOT_FOUND",
  "message": "account not found"
}
```

## Create Transaction

Record a transaction for a given account.

### Path:
POST /transactions

### Request Body:
| Field Name        | Data Type | Required | Example               | Description                                              |
|-------------------|-----------|----------|-----------------------|----------------------------------------------------------|
| msg_id            | string    | Yes      | "4e01cff74f2585c..."  | Unique identifier for this request                       |
| reference_id      | string    | Yes      | "4e01cff74f2585c..."  | client-generated reference for tracking(idempotency key) |
| account_id        | int64     | Yes      | 1                     | ID of account associated with this transaction           |
| operation_type_id | int       | Yes      | 4                     | Type of operation (e.g. payment, withdrawal)             |
| amount            | float64   | Yes      | 12.5                  | Amount for the transaction (positive value required)     |

Note : 
```
If a POST is retried with the same reference_id, the API will return the existing transaction (idempotency).
```

Example Request:
```json
{
  "msg_id": "4e01cff74f2585cbcdb38f66fb74533",
  "reference_id": "1e01cff74f2585cbcdb38f66fb74522",
  "account_id": 1,
  "operation_type_id": 4,
  "amount": 12.5
}
```

### Response:
#### Success : 201
| Field Name      | Data Type | Required | Example | Description                        |
|-----------------|-----------|----------|---------|------------------------------------|
| transaction_id  | int64     | Yes      | 101     | Unique identifier for the transaction|

Example:
```json
{
  "transaction_id": 101
}
```

#### Errors (4xx/5xx)
| Field Name | Data Type | Required | Example            | Description                          |
|------------|-----------|----------|--------------------|--------------------------------------|
| code       | string    | Yes      | "NOT_FOUND"        | Error code, such as NOT_FOUND        |
| message    | string    | Yes      | "account not found"| Human-readable error details         |

Example:
```json
{
  "code": "NOT_FOUND",
  "message": "account not found"
}
```

Another Example (idempotency):
```json
{
  "transaction_id": 101
}
```

## Error Codes for 4xx/5xx

| Code            | Message                  | HTTP Code | Description                                                               |
|-----------------|--------------------------|-----------|---------------------------------------------------------------------------|
| NOT_FOUND       | account not found        | 404       | The specified account was not found.                                      |
| BAD_REQUEST     | invalid input            | 400       | The request contains invalid or missing parameters.                       |
| UNPROCESSABLE   | unprocessable entity     | 422       | The server understands the request but cannot process it,                 |
| CONFLICT        | conflict                 | 409       | Resource already exists (e.g., duplicate reference_id for idempotency).   |
| INTERNAL_ERROR  | internal server error    | 500       | An unexpected error occurred on the server.                               |
