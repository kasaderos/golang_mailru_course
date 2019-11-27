package handlers

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
	"time"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo *user.UserRepo
	Sessions *session.SessionsManager
}

type Error struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}
type JsonError struct {
	Errors []Error `json:"errors"`
}

var (
	tokenSecret = []byte("your-256-bit-secret")
)

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	_, err := session.SessionFromContext(r.Context())
	if err == nil {
		fmt.Println("INDEX: err nil")
		http.Redirect(w, r, "/", 302)
		return
	}
	err = h.Tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserRepo.Authorize(r.FormValue("username"), r.FormValue("password"))
	if err == user.ErrNoUser {
		http.Error(w, `no user`, http.StatusBadRequest)
		return
	}
	if err == user.ErrBadPass {
		http.Error(w, `bad pass`, http.StatusBadRequest)
		return
	}

	token, err := getAccessToken(u)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(token)
	http.Redirect(w, r, "/", 302)
}

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

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserRepo.Register(r.FormValue("username"), r.FormValue("password"))
	if err == user.ErrAlreadyExist {
		//w.Header().Set("Content-Type, application/json")
		errjs, err2 := json.Marshal(&JsonError{
			Errors: []Error{
				Error{
					Location: "body",
					Param:    "username",
					Value:    r.FormValue("username"),
					Msg:      user.ErrAlreadyExist.Error(),
				},
			},
		})
		if err2 != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
		w.Write(errjs)
	}

	h.Logger.Infof("registered user %v", u.ID)

	w.Header().Set("Content-Type", "application/json")

	token, err := getAccessToken(u)
	if err != nil {
		http.Error(w, "token jwt", http.StatusInternalServerError)
		return
	}
	w.Write(token)
	http.Redirect(w, r, "/", 302)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Sessions.DestroyCurrent(w, r)
	http.Redirect(w, r, "/", 302)
}

/*
func jsonError(w io.Writer, status int, msg string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})
	w.Write(resp)
}*/
