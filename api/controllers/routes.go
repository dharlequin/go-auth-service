package controllers

import (
	"github.com/dharlequin/go-auth-service/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/sign-in", middlewares.SetMiddlewareJSON(s.SignIn)).Methods("GET")
	s.Router.HandleFunc("/register", middlewares.SetMiddlewareJSON(s.RegisterNewUser)).Methods("POST")
	s.Router.HandleFunc("/auth", middlewares.SetMiddlewareJSON(s.ValidateSessionID)).Methods("GET")
	s.Router.HandleFunc("/logout", middlewares.SetMiddlewareJSON(s.Logout)).Methods("GET")
}
