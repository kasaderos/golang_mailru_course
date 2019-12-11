package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"

	"redditclone/pkg/items"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostsHandler struct {
	Tmpl      *template.Template
	PostsRepo *items.PostsRepo
	Logger    *zap.SugaredLogger
	UserRepo  *user.UserRepo
}

func (h *PostsHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	h.PostsRepo.Mu.RLock()
	elems, err := h.PostsRepo.GetAll()
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(elems)
	h.PostsRepo.Mu.RUnlock()
	if err != nil {
		http.Error(w, `can't send as json`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *PostsHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	params, err := getJsonParams(r)
	sess, err := session.SessionFromContext(r.Context())

	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	p := &items.Post{
		Mu:    &sync.RWMutex{},
		Score: 1,
		Views: 0,
		Title: params["title"],
		Author: items.Author{
			Username: h.UserRepo.GetUserById(sess.UserID).Login,
			Id:       sess.UserID,
		},
		Category: params["category"],
		Votes: []*items.Vote{
			&items.Vote{
				User: sess.UserID,
				Vote: 1,
			},
		},
		Created:          time.Now().Format(time.RFC3339),
		UpvotePercentage: 100,
	}
	if _, ok := params["url"]; ok {
		p.Type = "link"
		p.Url = params["url"]
	} else if _, ok := params["text"]; ok {
		p.Type = "text"
		p.Text = params["text"]
	}
	h.PostsRepo.Mu.Lock()
	h.PostsRepo.Add(p)
	p.Mu.RLock()
	h.PostsRepo.Mu.Unlock()
	data, err := json.Marshal(p)
	p.Mu.RUnlock()
	if err != nil {
		http.Error(w, `can't send as json`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *PostsHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["POST_ID"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	h.PostsRepo.Mu.RLock()
	p, err := h.PostsRepo.GetByID(uint32(id))
	p.Mu.Lock()
	h.PostsRepo.Mu.RUnlock()
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	p.Views++
	data, err := json.Marshal(p)
	p.Mu.Unlock()
	if err != nil {
		http.Error(w, `can't send as json`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["POST_ID"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	h.PostsRepo.Mu.Lock()
	err = h.PostsRepo.Delete(uint32(id))
	h.PostsRepo.Mu.Unlock()
	if err != nil {
		http.Error(w, "not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(map[string]interface{}{
		"message": "success",
	})
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *PostsHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["LOGIN"]
	h.PostsRepo.Mu.RLock()
	ps, err := h.PostsRepo.GetUserPosts(login)

	if err != nil {
		http.Error(w, "400", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(ps)
	h.PostsRepo.Mu.RUnlock()
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ps_id := vars["POST_ID"]
	p_id, err := strconv.Atoi(ps_id)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	params, err := getJsonParams(r)
	if _, ok := params["comment"]; !ok {
		http.Error(w, "400 no param", http.StatusInternalServerError)
		return
	}
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	h.PostsRepo.Mu.RLock()
	p, err := h.PostsRepo.GetPost(uint32(p_id))
	h.PostsRepo.Mu.RUnlock()
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
	}
	h.UserRepo.Mu.Lock()
	u := h.UserRepo.GetUserById(sess.UserID)
	p.Mu.Lock()
	p.AddComment(params["comment"], u.ID, u.Login)
	h.UserRepo.Mu.Unlock()
	p.Mu.Unlock()
	data, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ps_id := vars["POST_ID"]
	cm_id := vars["COMMENT_ID"]
	p_id, err := strconv.Atoi(ps_id)
	c_id, err2 := strconv.Atoi(cm_id)
	if err != nil || err2 != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	h.PostsRepo.Mu.RLock()
	p, err := h.PostsRepo.GetPost(uint32(p_id))
	p.Mu.Lock()
	h.PostsRepo.Mu.Unlock()
	p.DeleteComment(uint32(c_id))
	p.Mu.Unlock()
	data, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *PostsHandler) GetCategoryPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["CATEGORY_NAME"]
	h.PostsRepo.Mu.RLock()
	elems := h.PostsRepo.GetCategoryPosts(category)
	data, err := json.Marshal(elems)
	h.PostsRepo.Mu.RUnlock()
	if err != nil {
		http.Error(w, `can't send as json`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func (h *PostsHandler) DownUpVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ps_id := vars["POST_ID"]
	choice := vars["CHOICE"]
	p_id, err := strconv.Atoi(ps_id)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	h.PostsRepo.Mu.RLock()
	p, err := h.PostsRepo.GetPost(uint32(p_id))
	p.Mu.Lock()
	h.PostsRepo.Mu.RUnlock()
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
	}
	switch choice {
	case "unvote":
		err = p.DeleteVote(sess.UserID)
	case "upvote":
		err = p.AddVote(sess.UserID, 1)
	case "downvote":
		err = p.AddVote(sess.UserID, -1)
	}
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	data, err := json.Marshal(p)
	p.Mu.Unlock()
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
