package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

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
)

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	_, err := session.SessionFromContext(r.Context())
	if err == nil {
		http.Redirect(w, r, "/api/posts/", 302)
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
		jsonError(w, map[string]interface{}{
			"message": user.ErrNoUser.Error(),
		})
		return
	}
	if err == user.ErrBadPass {
		jsonError(w, map[string]interface{}{
			"message": user.ErrBadPass.Error(),
		})
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
	u, err := h.UserRepo.Register(r.FormValue("username"), r.FormValue("password"))
	if err == user.ErrAlreadyExist {
		jsonError(w, map[string]interface{}{
			"location": "body",
			"param":    "username",
			"value":    r.FormValue("username"),
			"msg":      user.ErrAlreadyExist.Error(),
		})
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

func jsonError(w http.ResponseWriter, err map[string]interface{}) {
	errjs, err2 := json.Marshal(err)
	w.Header().Set("Content-Type", "application/json")
	if err2 != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Write(errjs)
}
