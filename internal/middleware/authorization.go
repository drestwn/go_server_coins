package middleware

import (
	"errors"
	"net/http"
	"go_server/api"            // Updated from github.com/avukadin/goapi/api
	"go_server/internal/tools" // Updated from github.com/avukadin/goapi/internal/tools
	log "github.com/sirupsen/logrus"
)

var UnAuthorizedError = errors.New("Invalid username or Token")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		var username string = r.URL.Query().Get("username")
		var token = r.Header.Get("Authorization")
log.Infof("Received request - Username: %s, Token: %s", username, token)

		var err error
log.Infof("Received request - Username: %s, Token: %s", username, token)
		if username == "" || token == "" {
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		var database *tools.DatabaseInterface
		database, err = tools.NewDatabase()
		if err != nil {
			api.InternalErrorHandler(w)
			return
		}

		var loginDetails *tools.LoginDetails
		loginDetails = (*database).GetUserLoginDetails(username)
		if loginDetails == nil || (token != (*loginDetails).AuthToken) {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
