# go_server

A simple, read-only RESTful API built in Go that allows authenticated users to retrieve their coin balance via the `/account/coins` endpoint. The project uses the Chi router for routing, Logrus for logging, and Gorilla Schema for query parameter parsing, with a modular structure for handlers, middleware, and database interactions.

## Table of Contents
- [Overview](#overview)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Dependencies](#dependencies)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Overview
`go_server` is a minimal Go-based API designed for read-only access to user coin balances. Key features include:
- **Authentication**: A middleware layer validates the `username` query parameter and `Authorization` header token.
- **Endpoint**: A single GET endpoint (`/account/coins`) returns the user’s coin balance as JSON.
- **Modular Design**: Organized into packages for API types (`api`), handlers (`handlers`), middleware (`middleware`), and database utilities (`tools`).

The API is currently read-only, with potential for future enhancements.

## Project Structure



go_server/
├── .gitignore               # Ignores .vscode/, binaries, etc.
├── go.mod                   # Module definition and dependencies
├── main.go                  # Entry point, sets up Chi router
├── api/
│   └── api.go               # Defines request/response structs and error handlers
├── internal/
│   ├── handlers/
│   │   ├── handlers.go      # Defines routes
│   │   └── get_coin_balance.go  # Handles /account/coins endpoint
│   ├── middleware/
│   │   └── middleware.go    # Authorization middleware
│   └── tools/
│       └── tools.go         # Database interface and data structures



### Package Details
- **`api`**: Defines structs (`CoinBalanceParams`, `CoinBalanceResponse`) and error handlers (`InternalErrorHandler`, `RequestErrorHandler`).
- **`handlers`**: Configures routes (`Handler` in `handlers.go`) and handles the `/account/coins` endpoint (`GetCoinBalance` in `get_coin_balance.go`).
- **`middleware`**: Implements `Authorization` middleware to validate username and token.
- **`tools`**: Provides database interfaces (`DatabaseInterface`, `CoinDetails`, `LoginDetails`) for data access.

## Prerequisites
- **Go**: Version 1.24.2 or later.
- **VSCode** (optional): Recommended with the Go extension for auto-imports and development.
- **Git**: For cloning the repository.

## Installation
### Clone the Repository
```bash
git clone <repository-url>
cd go_server
```

### Install Dependencies
```bash
go mod tidy
```

#### Build the Project
```bash

go build
```

#### Run the Server
```bash

./go_server
```

The server runs on http://localhost:8080.

#### Usage

```bash 
curl "http://localhost:8080/account/coins?username=john" -H "Authorization: some-token"
```