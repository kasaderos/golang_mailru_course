package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"redditclone/pkg/items"

	"go.uber.org/zap"
)

type PostsHandler struct {
	Tmpl      *template.Template
	PostsRepo *items.PostsRepo
	Logger    *zap.SugaredLogger
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// (type url) params {"category":"music","type":"link","title":"youtube","url":"http://youtube.com"}
// response {"score":1,"views":0,"type":"link","title":"youtube","url":"http://youtube.com","author":{"username":"alisher","id":"5dde28b549c115e4af02238b"},"category":"music","votes":[{"user":"5dde28b549c115e4af02238b","vote":1}],"comments":[],"created":"2019-12-01T18:39:32.297Z","upvotePercentage":100,"id":"5de408e4584517a1f7461866"}

// (type text) params {"category":"music","type":"text","title":"You","text":"youadf"}
// response {"score":1,"views":0,"type":"text","title":"You","author":{"username":"alisher","id":"5dde28b549c115e4af02238b"},"category":"music","text":"youadf","votes":[{"user":"5dde28b549c115e4af02238b","vote":1}],"comments":[],"created":"2019-12-01T18:45:15.454Z","upvotePercentage":100,"id":"5de40a3b5845176559461867"}
/*
func (h *PostsHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	//if r.Context().Value("sessionKey")
	params, err := getJsonParams(r)
	if _, ok := params["url"]; ok {
		//url
		p := &items.Post{
			Score:  1,
			Views:  0,
			Type:   "link",
			Title:  params["title"],
			Url:    params["url"],
			Author: items.Author{

			},
		}
	} else {
		//text
	}
	if err != nil {
		http.Error(w, "400", http.StatusBadRequest)
		return
	}
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

/*
func (h *ItemsHandler) AddForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) Add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	item := new(items.Item)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(item, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}

	sess, _ := session.SessionFromContext(r.Context())
	item.CreatedBy = sess.UserID

	lastID, err := h.ItemsRepo.Add(item)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Insert with id LastInsertId: %v", lastID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *ItemsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByID(uint32(id))
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, `no item`, http.StatusNotFound)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", item)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}

	r.ParseForm()
	item := new(items.Item)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(item, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}
	item.ID = uint32(id)

	ok, err := h.ItemsRepo.Update(item)
	if err != nil {
		http.Error(w, `db error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("update: %v %v", item, ok)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	ok, err := h.ItemsRepo.Delete(uint32(id))
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]bool{
		"success": ok,
	})
	w.Write(respJSON)
}
*/
