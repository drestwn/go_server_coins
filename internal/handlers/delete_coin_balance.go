package handlers

import (
    "encoding/json"
    "net/http"
    "go_server/internal/tools"
    log "github.com/sirupsen/logrus"
)

type DeleteCoinBalanceResponse struct {
    Message string `json:"message"`
    Code    int    `json:"code"`
}

func DeleteCoinBalance(w http.ResponseWriter, r *http.Request) {
    // Get username from query parameter
    username := r.URL.Query().Get("username")
    if username == "" {
        log.Error("Username query parameter is missing")
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }

    // Initialize database
    var database *tools.DatabaseInterface
    database, err := tools.NewDatabase()
    if err != nil {
        log.Errorf("Failed to initialize database: %v", err)
        InternalErrorHandler(w)
        return
    }

    // Check if user exists in coin_data
    if existingDetails := (*database).GetUserCoins(username); existingDetails == nil {
        log.Errorf("User %s not found in coin_data", username)
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Delete from coin_data
    err = (*database).DeleteUserCoins(username) // Assuming this method exists or will be added
    if err != nil {
        log.Errorf("Failed to delete coin data for %s: %v", username, err)
        InternalErrorHandler(w)
        return
    }

    // Optionally delete from login_data (synchronize with coin_data)
    err = (*database).DeleteUserLoginDetails(username)
    if err != nil {
        log.Errorf("Failed to delete login data for %s: %v", username, err)
        InternalErrorHandler(w)
        return
    }

    // Prepare response
    response := DeleteCoinBalanceResponse{
        Message: "Coin data and login details deleted successfully",
        Code:    http.StatusOK,
    }

    // Send JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Errorf("Failed to encode response: %v", err)
        InternalErrorHandler(w)
        return
    }
}