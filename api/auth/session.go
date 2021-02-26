package auth

import (
	"errors"
	"math/rand"
	"time"

	"github.com/dharlequin/go-auth-service/api/models"
)

//SESSIONS ...
var SESSIONS = make(map[string]*models.User)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

//CreateSession creates new Sessions for logged-in User and stores his data
func CreateSession(user *models.User) string {
	sessionID := generateSessionID(40)
	SESSIONS[sessionID] = user

	return sessionID
}

func generateSessionID(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

//GetUserDataFromSessions checks if there is data for specified Session ID and returns User data
func GetUserDataFromSessions(sessionID string) (*models.User, error) {
	if SESSIONS[sessionID] == nil {
		return nil, errors.New("No data found for session ID: " + sessionID)
	}

	return SESSIONS[sessionID], nil
}
