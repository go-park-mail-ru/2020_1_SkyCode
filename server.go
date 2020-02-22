package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type SessionHandler struct {
	store *SessionStore
}

type RestaurantsHandler struct {
}

func (api *SessionHandler) Handle(w http.ResponseWriter, r *http.Request) {

}

func DefaultHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Напиши для моего пути обработчик"))
}

func main() {
	router := mux.NewRouter()

	apiSession := &SessionHandler{
		store: NewSessionStore(),
	}

	router.HandleFunc("/session", apiSession.Handle)
	router.HandleFunc("/user", DefaultHandle)
	router.HandleFunc("/profile", DefaultHandle)
	router.HandleFunc("/restaurants", DefaultHandle)
	router.HandleFunc("/restaurants/{restaurant_id:[0-9]+}", DefaultHandle)

	http.ListenAndServe(":8080", router)
}
