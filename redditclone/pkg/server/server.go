package server

import (
	"net/http"
	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"
	"regexp"

	"github.com/gorilla/mux"
)

func handlersRouter(userHandler *handlers.UserHandler, itemHandler *handlers.ItemsHandler) http.Handler {
	staticHandler := http.StripPrefix(
		"/static",
		http.FileServer(http.Dir("./web/static")),
	)

	r := mux.NewRouter()
	r.Handle("/static/js/{file}", staticHandler)
	r.Handle("/static/css/{file}", staticHandler)
	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	}))

	r.HandleFunc("/api/register", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	r.HandleFunc("/api/posts/", itemHandler.List).Methods("GET")
	r.HandleFunc("/api/posts", itemHandler.Add).Methods("POST")
	r.HandleFunc("/api/post/{id}", itemHandler.Read).Methods("GET")

	r.HandleFunc("/api/post/{id}", itemHandler.Delete).Methods("DELETE")
	r.HandleFunc("/api/post/{id}", itemHandler.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{id}/upvote", itemHandler.Upvote).Methods("GET")
	r.HandleFunc("/api/post/{id}/downvote", itemHandler.Downvote).Methods("GET")
	r.HandleFunc("/api/post/{id}/unvote", itemHandler.Unvote).Methods("GET")
	r.HandleFunc("/api/posts/{catName}", itemHandler.Category).Methods("GET")
	r.HandleFunc("/api/user/{username}", itemHandler.UserItems).Methods("GET")
	r.HandleFunc("/api/post/{postID}/{comID}", itemHandler.DeleteComment).Methods("DELETE")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	mux := routerWrapper(userHandler, r)

	return mux
}

func routerWrapper(userHandler *handlers.UserHandler, r *mux.Router) http.Handler {
	mux := middleware.Auth(userHandler.Sessions, r, userHandler.UserRepo, []middleware.AuthURL{
		{
			URL:    regexp.MustCompile("/api/posts/?$"),
			Method: "POST",
		},
		{
			URL:    regexp.MustCompile("/api/post/[0-9a-zA-Z-]+/?$"),
			Method: "POST",
		},
		{
			URL:    regexp.MustCompile("/api/post/[0-9a-zA-Z-]+/?$"),
			Method: "DELETE",
		},
		{
			URL:    regexp.MustCompile("/api/post/[0-9a-zA-Z-]+/[0-9a-zA-Z-]+/?$"),
			Method: "DELETE",
		},
		{
			URL:    regexp.MustCompile("/api/post/[0-9a-zA-Z-]+/upvote/?$"),
			Method: "GET",
		},
		{
			URL:    regexp.MustCompile("/api/post/[0-9a-zA-Z-]+/downvote/?$"),
			Method: "GET",
		},
		{
			URL:    regexp.MustCompile("/api/post/[0-9a-zA-Z-]+/unvote/?$"),
			Method: "GET",
		},
	})

	mux = middleware.AccessLog(userHandler.Logger, mux)
	mux = middleware.Panic(mux)

	return mux
}

func Run(userHandler *handlers.UserHandler, itemHandler *handlers.ItemsHandler, addr string) error {
	mux := handlersRouter(userHandler, itemHandler)

	ErrHTTP := http.ListenAndServe(addr, mux)

	return ErrHTTP
}
