package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"redditclone/pkg/user"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type SessionsManager struct {
	UserRepo *user.UserRepo
}

func (sm *SessionsManager) JwtParse(inToken string) (jwt.MapClaims, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return tokenSecret, nil
	}
	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("bad token")
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("no payload")
	}
	return payload, nil
}

func (sm *SessionsManager) GetAccessToken(u *user.User) ([]byte, error) {
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

func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {
	// проверяем есть ли в базе
	inToken := r.Header.Get("Authorization")
	if inToken == "" {
		return nil, ErrNoAuth
	}
	fmt.Println("token", inToken)
	payload, err := sm.JwtParse(inToken)
	if err != nil {
		return nil, fmt.Errorf("jwtParse")
	}
	username := payload["user"].(map[string]interface{})["username"].(string)
	ID := payload["user"].(map[string]interface{})["ID"].(uint32)
	u, ok := sm.UserRepo.GetData(username)
	if !ok && u.ID != ID {
		return nil, ErrNoAuth
	}
	// проверяем время жизни
	ptime := payload["iat"].(time.Time)
	if ptime.Second() >= time.Now().Second() {
		return nil, fmt.Errorf("exp time >= now time")
	}
	return &Session{token: []byte(inToken)}, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, user *user.User) (*Session, error) {
	sess, err := NewSession(user)
	if err != nil {
		return nil, err
	}
	w.Header().Set("Authorization", string(sess.token))
	return sess, nil
}
