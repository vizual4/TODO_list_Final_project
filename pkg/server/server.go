package server

import (
	"TODO_List/pkg/api"
	"log"
	"net/http"
	"os"
	"time"
)

func Server(l *log.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	api.Init(mux)

	env := os.Getenv("TODO_LIST")
	if len(env) == 0 {
		env = "7540"
	}

	var server = http.Server{
		Addr:         ":" + env,
		Handler:      mux,
		ErrorLog:     l,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	return &server
}
