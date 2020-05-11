package session

import (
	"errors"
	"net/http"
	"time"

	uuid "github.com/google/uuid"
)

var Error400 = errors.New("bad request")
var ErrorNotSet = errors.New("cookie is not set")

type SessionManager struct{}

const sessionKey = "session"

func (m *SessionManager) ReadCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	c, error := r.Cookie(sessionKey)
	if error != nil {
		if error == http.ErrNoCookie {
			error = ErrorNotSet
			return "", error
		}
		error = Error400
	}
	str := c.Value
	return str, nil
}

func (m *SessionManager) SetCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:    sessionKey,
		Value:   token,
		Expires: time.Now().Add(60 * time.Second),
	})
}

func (m *SessionManager) NewToken() string {
	return uuid.New().String()
}
