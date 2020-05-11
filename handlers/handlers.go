package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/grandima/user-count/session"
)

type StorageInterface interface {
	Get(key string) (int, error)
	Set(key string, value string) error
	SetIfNotExistsToZero(key string) error
	Increment(key string) (int, error)
	Exists(key string) (int, error)
}

type SessionManagerInterface interface {
	ReadCookie(http.ResponseWriter, *http.Request) (string, error)
	SetCookie(http.ResponseWriter, string)
	NewToken() string
}

type Handler struct {
	Storage        StorageInterface
	SessionManager SessionManagerInterface
}

const userCountKey = "userCountKey"

func (handler *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var count int
	token, error := handler.SessionManager.ReadCookie(w, r)
	if error != nil {
		if error == session.ErrorNotSet {
			token = handler.SessionManager.NewToken()
			error = handler.Storage.Set(token, token)
			if error != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(error)
				return
			}

			handler.SessionManager.SetCookie(w, token)

			error = handler.Storage.SetIfNotExistsToZero(userCountKey)
			if error != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(error)
				return
			}

			count, error := handler.Storage.Increment(userCountKey)
			log.Print("New user count: ", count)
			if error != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(error)
				return
			}
		} else {
			log.Print(error)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		isTokenExist, error := handler.Storage.Exists(token)
		if isTokenExist == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad token!"))
			log.Print(error)
			log.Print("session token:" + token + " does not exist")
			return
		}
	}

	count, error = handler.Storage.Get(userCountKey)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(error)
		return
	}
	s := strconv.Itoa(count)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}
