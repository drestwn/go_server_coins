package handlers

import (
    "encoding/json"
    "net/http"
    "go_server/internal/tools"
    log "github.com/sirupsen/logrus"
)

type CreateCoinBalanceRequest struct {
    Username string `json:"username"`
    Coins    int64  `json:"coins"` // Changed to int64 to match tools.CoinDetails
}

type CreateCoinBalanceResponse struct {
    Message string `json:"message"`
    Code    int    `json:"code"`
}

func CreateCoinBalance(w http.ResponseWriter, r *http.Request) {
    // Ensure the request method is POST
    if r.Method != http.MethodPost {
 http.Error(w, "Method not allowed (only post)", http.StatusMethodNotAllowed)
        return
    }

    // Decode JSON request body
    var req CreateCoinBalanceRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        log.Errorf("Failed to decode request body: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest) // Fixed typo
        return
    }

    // Validate input
    if req.Username == "" || req.Coins < 0 { // Changed " " to "" for proper empty check
        log.Errorf("Invalid input: username=%s, coins=%d", req.Username, req.Coins)
        http.Error(w, "Username cannot be empty and coins cannot be negative", http.StatusBadRequest)
        return
    }

    // Initialize database
    var database *tools.DatabaseInterface
    database, err = tools.NewDatabase()
    if err != nil {
        log.Errorf("Failed to initialize database: %v", err)
        InternalErrorHandler(w) // Removed "api." since it's in the same package
        return
    }

    // Check if user already exists
    if existingDetails := (*database).GetUserCoins(req.Username); existingDetails != nil {
        log.Errorf("User %s already exists", req.Username)
        http.Error(w, "User already exists", http.StatusConflict)
        return
    }

    // Create new coin details
    newDetails := tools.CoinDetails{
        Username: req.Username,
        Coins:    req.Coins,
    }

    // Add to database
    err = (*database).CreateUserCoins(req.Username, newDetails) // Fixed typo "databse" to "database"
    if err != nil {
        log.Errorf("Failed to create coin data: %v", err)
        InternalErrorHandler(w) // Removed "api."
        return
    }

    // Prepare response
    response := CreateCoinBalanceResponse{
        Message: "Coin data created successfully", // Added comma
        Code:    http.StatusCreated,
    }

    // Send JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Errorf("Failed to encode response: %v", err)
        InternalErrorHandler(w) // Removed "api."
        return
    }
}