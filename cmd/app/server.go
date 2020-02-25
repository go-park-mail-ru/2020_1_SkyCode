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

	apiRestaurants := &_handlers.RestaurantHandler{
		Restaurants: _models.BaseResStorage,
	}

	router.HandleFunc("/session", apiSession.SessionHandle).Methods("DELETE", "POST", "OPTIONS")
	router.HandleFunc("/user", apiSession.UserHandle)
	router.HandleFunc("/profile", apiSession.GetUserProfile)
	router.HandleFunc("/restaurants", apiRestaurants.GetRestaurants)
	router.HandleFunc("/restaurants/{restaurant_id:[0-9]+}", apiRestaurants.GetRestaurantByID)

	fmt.Println("Server started")

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println(err)
	}
}
