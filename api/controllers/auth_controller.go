package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dharlequin/go-auth-service/api/auth"
	"github.com/dharlequin/go-auth-service/api/models"
	"github.com/dharlequin/go-auth-service/api/responses"
	"github.com/dharlequin/go-auth-service/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

const cookieName = "session_id"

//Login checks user credentianls and creates Session ID
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	sessionID, err := server.signIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	cookie := http.Cookie{Name: cookieName, Value: sessionID}
	http.SetCookie(w, &cookie)

	responses.JSON(w, http.StatusOK, sessionID)
}

func (server *Server) signIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateSession(&user), nil
}

//SignIn is a placeholder route for API Gateway
func (server *Server) SignIn(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Use login endpoint and provide credentials")
}

//RegisterNewUser creates new User and saves to DB
func (server *Server) RegisterNewUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

//ValidateSessionID authorises user by checking Session ID in cookies
func (server *Server) ValidateSessionID(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	if cookie == nil || cookie.Value == "" {
		err = errors.New("User not logged in")

		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	user, err := auth.GetUserDataFromSessions(cookie.Value)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	w.Header().Set("X-UserId", strconv.FormatUint(uint64(user.ID), 10))
	w.Header().Set("X-User", user.Nickname)
	w.Header().Set("X-Email", user.Email)

	responses.JSON(w, http.StatusOK, user)
}

//Logout removes http Cookie
func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: cookieName, Value: ""}
	http.SetCookie(w, &cookie)

	responses.JSON(w, http.StatusOK, "User logged out")
}
