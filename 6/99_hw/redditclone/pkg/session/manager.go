package session

import (
	"errors"
	"net/http"
	"redditclone/pkg/user"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrBadSign            = errors.New("bad sign method")
	ErrBadToken           = errors.New("bad token")
	ErrPayload            = errors.New("no payload")
	ErrExpTime            = errors.New("exp time >= now time")
    ErrJwtParse           = errors.New("jwtParse")
)

type SessionsManager struct {
	UserRepo *user.UserRepo
}

func (sm *SessionsManager) JwtParse(inToken string) (jwt.MapClaims, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, ErrBadSign
		}
		return tokenSecret, nil
	}
	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, ErrBadToken
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrPayload
	}
	return payload, nil
}

func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {
	// проверяем есть ли в базе
	inToken := r.Header.Get("Authorization")
	if inToken == "" {
		return nil, ErrNoAuth
	}
	inTokens := strings.Split(inToken, " ")
	if len(inTokens) < 2 {
		return nil, ErrNoAuth
	}
	payload, err := sm.JwtParse(inTokens[1])
	if err != nil {
		return nil, ErrJwtParse
	}
	username := payload["user"].(map[string]interface{})["username"].(string)
	ID := uint32(payload["user"].(map[string]interface{})["id"].(float64))
	u, err := sm.UserRepo.GetUserByUsername(username)
	if err != nil && u.ID != ID {
		return nil, ErrNoAuth
	}
	// проверяем время жизни
	ptime := payload["exp"].(string)
	t, err := time.Parse(time.RFC3339, ptime)
	if err != nil {
		return nil, ErrExpTime
	}
	if !time.Now().Before(t) {
		return nil, ErrExpTime
	}
	return &Session{ID: inTokens[1][:5], UserID: u.ID}, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, user *user.User) (*Session, error) {
	sess, err := NewSession(user)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
