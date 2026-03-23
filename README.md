# pismo-account

Simple Transaction Routine service, which manages accounts and transactions associated with those accounts.

## What you need

- Docker
- Docker Compose
- Make

## How to run it

1. Make sure you have a `PROJECT_ROOT/config/config.json` file (a `config` folder with a `config.json` file inside it).
2. In a terminal, inside PROJECT_ROOT(pismo-account) folder, run:

```bash
make up
```

This will:
- starts MySQL
- build the pismo-account app
- start the pismo-account app once the mysql is up and running


To see logs:

```bash
make logs
```

To stop everything:

```bash
make down
```

## Project Structure
The project is organized as follows:

- `cmd/`  
  Entry point for the application (main Go program).

- `config/`  
  Configuration files, such as `config.json` for app and database settings.

- `dto/`  
  Data Transfer Objects used for communication between api and logic layer.

- `internal/`  
  Core business logic, divided into domains like `account` and `transaction`, which hold the actual logic.

- `storage/`  
  Data access layer, acts as an interface between the logic layer and the actual storage.

- `api/`  
  Defines the API endpoints. Provides the handlers for the http routes and responsible for parsing and validating the request and also to call the underlying logic layer.

- `Makefile`  
  Script for building, running, and managing the app with Docker and Docker Compose.


This layered structure separates configuration, business logic, API endpoints, data access, and infrastructure for maintainability and clarity.

`Server` -> `API Layer` -> `Logic Layer` -> `Storage Layer`

## API Contract

For detailed information on API endpoints and request/response formats, see [API_CONTRACT.md](API_CONTRACT.md)

## Product Requirements

https://docs.google.com/document/d/1ibohkkWR0WzgX_f-Cd3HH4f2UBnrWvvQ8f4O_C-4v40/edit?tab=t.0 