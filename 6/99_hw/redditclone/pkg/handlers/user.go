package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo *user.UserRepo
	Sessions *session.SessionsManager
}

var (
	tokenSecret = []byte("your-256-bit-secret")
	ErrGetParam = errors.New("interface to string")
)

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	_, err := session.SessionFromContext(r.Context())
	if err == nil {
		http.Redirect(w, r, "/api/posts/", 304)
		return
	}
	err = h.Tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	params, err := getJsonParams(r)
	if err != nil {
		http.Error(w, "400", http.StatusBadRequest)
		return
	}
	u, err := h.UserRepo.Authorize(params["username"], params["password"])
	if err == user.ErrNoUser {
		jsonError(w, map[string]interface{}{
			"message": user.ErrNoUser.Error(),
		}, http.StatusUnauthorized)
		return
	}
	if err == user.ErrBadPass {
		jsonError(w, map[string]interface{}{
			"message": user.ErrBadPass.Error(),
		}, http.StatusUnauthorized)
		return
	}

	token, err := h.Sessions.GetAccessToken(u)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(token)
	http.Redirect(w, r, "/", 302)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	params, err := getJsonParams(r)
	if err != nil {
		http.Error(w, "400", http.StatusBadRequest)
		return
	}
	fmt.Println(params["username"], params["password"])
	u, err := h.UserRepo.Register(params["username"], params["password"])
	if err == user.ErrAlreadyExist {
		jsonError(w, map[string]interface{}{
			"errors": []interface{}{
				map[string]interface{}{
					"location": "body",
					"param":    "username",
					"value":    params["username"],
					"msg":      user.ErrAlreadyExist.Error(),
				},
			},
		}, http.StatusUnprocessableEntity)
		return
	}
	h.Logger.Infof("registered user %v", u.ID)
	w.Header().Set("Content-Type", "application/json")

	token, err := h.Sessions.GetAccessToken(u)
	if err != nil {
		http.Error(w, "token jwt", http.StatusInternalServerError)
		return
	}
	w.Write(token)
	http.Redirect(w, r, "/", 302)
}

func jsonError(w http.ResponseWriter, err map[string]interface{}, status int) {
	errjs, err2 := json.Marshal(err)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err2 != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Write(errjs)
}

func getJsonParams(r *http.Request) (map[string]string, error) {
	decoder := json.NewDecoder(r.Body)
	var js map[string]interface{}
	err := decoder.Decode(&js)
	if err != nil {
		return nil, err
	}
	res := make(map[string]string, 4)
	for k, v := range js {
		val, ok := v.(string)
		if !ok {
			return nil, ErrGetParam
		}
		res[k] = val
	}
	return res, nil
}

func (h *UserHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["login"]

	ps, err := h.UserRepo.GetUserPosts(login)
	if err != nil {
		http.Error(w, "400", http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(ps)
	if err != nil {
		http.Error(w, `can't send as json`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
