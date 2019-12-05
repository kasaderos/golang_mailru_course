package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
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
	elems, err := h.PostsRepo.GetAll()
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(elems)
	if err != nil {
		http.Error(w, `can't send as json`, http.StatusInternalServerError)
		return
	}
	/*if h.PostsRepo.Changed {
		w.WriteHeader(http.StatusOK)
		h.PostsRepo.Changed = false
	} else {
		w.WriteHeader(http.StatusNotModified)
	}*/
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// (type url) params {"category":"music","type":"link","title":"youtube","url":"http://youtube.com"}
// response {"score":1,"views":0,"type":"link","title":"youtube","url":"http://youtube.com","author":{"username":"alisher","id":"5dde28b549c115e4af02238b"},"category":"music","votes":[{"user":"5dde28b549c115e4af02238b","vote":1}],"comments":[],"created":"2019-12-01T18:39:32.297Z","upvotePercentage":100,"id":"5de408e4584517a1f7461866"}

// (type text) params {"category":"music","type":"text","title":"You","text":"youadf"}
// response {"score":1,"views":0,"type":"text","title":"You","author":{"username":"alisher","id":"5dde28b549c115e4af02238b"},"category":"music","text":"youadf","votes":[{"user":"5dde28b549c115e4af02238b","vote":1}],"comments":[],"created":"2019-12-01T18:45:15.454Z","upvotePercentage":100,"id":"5de40a3b5845176559461867"}

func (h *PostsHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	params, err := getJsonParams(r)
	sess, err := session.SessionFromContext(r.Context())

	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	p := &items.Post{
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
	h.PostsRepo.Add(p)
	data, err := json.Marshal(p)
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
	p, err := h.PostsRepo.GetByID(uint32(id))
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	p.Views++
	data, err := json.Marshal(p)
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

	err = h.PostsRepo.Delete(uint32(id))
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

	ps, err := h.PostsRepo.GetUserPosts(login)
	if err != nil {
		http.Error(w, "400", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(ps)
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
	p, err := h.PostsRepo.GetPost(uint32(p_id))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
	}
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	params, err := getJsonParams(r)
	if _, ok := params["comment"]; !ok {
		http.Error(w, "400 no param", http.StatusInternalServerError)
		return
	}
	u := h.UserRepo.GetUserById(sess.UserID)
	p.AddComment(params["comment"], u.ID, u.Login)
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
	p, err := h.PostsRepo.GetPost(uint32(p_id))
	p.DeleteComment(uint32(c_id))
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
	elems := h.PostsRepo.GetCategoryPosts(category)
	data, err := json.Marshal(elems)
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
	p, err := h.PostsRepo.GetPost(uint32(p_id))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
	}
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
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
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
