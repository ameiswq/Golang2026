package app

import (
	"net/http"

	"assignment-02/internal/handler"
	"assignment-02/internal/middleware"
	"assignment-02/internal/usecase"
)

func newRouter(u usecase.UserUsecase) http.Handler {
	mux := http.NewServeMux()
	h := handler.NewUserHandler(u)
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetUsers(w, r)
		case http.MethodPost:
			h.CreateUser(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetUserByID(w, r)
		case http.MethodPut, http.MethodPatch:
			h.UpdateUser(w, r)
		case http.MethodDelete:
			h.DeleteUser(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	var root http.Handler = mux
	root = middleware.APIKey("secret")(root)
	root = middleware.Logging(root)

	return root
}