package controllers

import (
	"net/http"

	"github.com/dharlequin/go-auth-service/api/responses"
)

//Home welcomes us in our API
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome Auth API")
}