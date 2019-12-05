package main

import (
	"html/template"
	"net/http"
	"redditclone/pkg/handlers"
	"redditclone/pkg/items"
	"redditclone/pkg/middleware"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	templates := template.Must(template.ParseFiles("./template/index.html"))

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	userRepo := user.NewUserRepo()
	postsRepo := items.NewPostRepo()
	sm := &session.SessionsManager{UserRepo: userRepo}
	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	handlers := &handlers.PostsHandler{
		Tmpl:      templates,
		Logger:    logger,
		PostsRepo: postsRepo,
		UserRepo:  userRepo,
	}
	dir := "./template/static/"
	r := mux.NewRouter()
	r.HandleFunc("/", userHandler.Index).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/posts/", handlers.GetPosts).Methods("GET")
	r.HandleFunc("/api/posts", handlers.AddPost).Methods("POST")
	r.HandleFunc("/api/post/{POST_ID}", handlers.GetPost).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", handlers.DeletePost).Methods("DELETE")
	r.HandleFunc("/api/user/{LOGIN}", handlers.GetUserPosts).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", handlers.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", handlers.DeleteComment).Methods("DELETE")
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", handlers.GetCategoryPosts).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/{CHOICE:upvote|downvote|unvote}", handlers.DownUpVote).Methods("GET")
	mux := middleware.Auth(sm, r)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, mux)
}
