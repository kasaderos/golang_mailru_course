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

type sessKey string

var (
	SessionKey  sessKey = "sessionKey"
	tokenSecret         = []byte("your-256-bit-secret")
	ErrNoAuth           = errors.New("No session found")
)

func getAccessToken(u *user.User) ([]byte, error) {
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
		return nil, fmt.Errorf("signed string")
	}
	tokenjs, err := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		return nil, fmt.Errorf("can't marshal")
	}
	return tokenjs, nil
}

func NewSession(u *user.User) (*Session, error) {
	token, err := getAccessToken(u)
	if err != nil {
		return nil, fmt.Errorf("getAccessToken")
	}
	return &Session{UserID: u.ID, ID: string(token[:5])}, nil
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		fmt.Println(ErrNoAuth.Error())
		return nil, ErrNoAuth
	}
	return sess, nil
}
