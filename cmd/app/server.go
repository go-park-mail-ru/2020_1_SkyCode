package main

import (
	"fmt"
	_handlers "github.com/2020_1_Skycode/internal/handlers"
	_models "github.com/2020_1_Skycode/internal/models"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	apiSession := &_handlers.SessionHandler{
		UserStore: _models.NewUserStore(),
		Sessions:  make(map[string]uint, 10),
	}

	router.HandleFunc("/session", apiSession.SessionHandle)
	router.HandleFunc("/user", apiSession.UserHandle)
	router.HandleFunc("/profile", _handlers.DefaultHandle)
	router.HandleFunc("/restaurants", _handlers.DefaultHandle)
	router.HandleFunc("/restaurants/{restaurant_id:[0-9]+}", _handlers.DefaultHandle)

	fmt.Println("Server started")

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println(err)
	}
}
