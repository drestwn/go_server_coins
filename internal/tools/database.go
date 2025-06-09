package tools

import (
	log "github.com/sirupsen/logrus"
)

type LoginDetails struct {
	AuthToken string
	Username  string
}

type CoinDetails struct {
	Coins    int64
	Username string
}

type DatabaseInterface interface {
GetUserLoginDetails(username string) *LoginDetails
    GetUserCoins(username string) *CoinDetails
    CreateUserCoins(username string, details CoinDetails) error
    CreateUserLoginDetails(username string, details LoginDetails) error
    SetupDatabase() error
	DeleteUserCoins(username string) error         // New method
    DeleteUserLoginDetails(username string) error // New method
}

func NewDatabase() (*DatabaseInterface, error) {
	var database DatabaseInterface = &mockDB{}

	var err error = database.SetupDatabase()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &database, nil
}
