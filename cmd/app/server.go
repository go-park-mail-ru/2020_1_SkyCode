package main

import (
	"fmt"
	_handlers "github.com/2020_1_Skycode/internal/handlers"
	mw "github.com/2020_1_Skycode/internal/middlewares"
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

	mwController := &mw.MWController{}

	router.Use(mwController.CORS)
	router.Use(mwController.AccessLogging)

	router.HandleFunc("/session", apiSession.SessionHandle).Methods("DELETE", "POST", "OPTIONS")
	router.HandleFunc("/user", apiSession.UserHandle).Methods("POST", "OPTIONS")
	router.HandleFunc("/profile", apiSession.GetUserProfile).Methods("GET", "PUT", "OPTIONS")
	router.HandleFunc("/restaurants", apiRestaurants.GetRestaurants).Methods("GET")
	router.HandleFunc("/restaurants/{restaurant_id:[0-9]+}", apiRestaurants.GetRestaurantByID).Methods("GET")

	fmt.Println("Server started")

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println(err)
	}
}
