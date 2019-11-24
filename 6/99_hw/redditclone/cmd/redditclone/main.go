package main

import (
	"html/template"
	"net/http"

	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	templates := template.Must(template.ParseFiles("./template/index.html"))

	sm := session.NewSessionsMem()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	userRepo := user.NewUserRepo()

	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", userHandler.Index).Methods("GET")
	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	//r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
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
