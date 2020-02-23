package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type RestaurantsHandler struct {
}

func DefaultHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Напиши для моего пути обработчик"))
}

func main() {
	router := mux.NewRouter()

	apiSession := &SessionHandler{
		userStore: NewUserStore(),
		sessions:  make(map[string]uint, 10),
	}

	router.HandleFunc("/session", apiSession.SessionHandle)
	router.HandleFunc("/user", apiSession.UserHandle)
	router.HandleFunc("/profile", DefaultHandle)
	router.HandleFunc("/restaurants", DefaultHandle)
	router.HandleFunc("/restaurants/{restaurant_id:[0-9]+}", DefaultHandle)

	fmt.Println("Server started")

	http.ListenAndServe(":8080", router)
}
