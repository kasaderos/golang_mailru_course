package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"redditclone/pkg/user"
	"time"
)

type Session struct {
	ID     string
	UserID uint32
}

var (
	SessionKey  sessKey = "sessionKey"
	tokenSecret         = []byte("your-256-bit-secret")
	ErrNoAuth           = errors.New("No session found")
	ErrSignedString     = errors.New("signed string")
	ErrJsonMarshal      = errors.New("can't marshal")
    ErrAccessToken      = errors.New("get access token error")
)

func genAccessToken(u *user.User) ([]byte, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": u.Login,
			"id":       u.ID,
		},
		"iat": time.Now(),
		"exp": time.Now().Add(time.Hour * 24),
	})
	tokenString, err := token.SignedString(tokenSecret)
	if err != nil {
		return nil, ErrSignedString
	}
	tokenjs, err := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		return nil, ErrJsonMarshal
	}
	return tokenjs, nil
}

func NewSession(u *user.User) (*Session, error) {
	token, err := genAccessToken(u)
	if err != nil {
		return nil, ErrAccessToken
	}
	return &Session{UserID: u.ID, ID: string(token[:5])}, nil
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}
