package handlers

import (
	json "encoding/json"
	"net/http"
	"strconv"

	_models "github.com/2020_1_Skycode/internal/models"
	"github.com/gorilla/mux"
)

type RestaurantHandler struct {
	Restaurants _models.ResStorage
}

func (api *RestaurantHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(api.Restaurants)
}

func (api *RestaurantHandler) GetRestaurantByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	restaurantId, err := strconv.Atoi(vars["restaurant_id"])

	if err != nil {
		HttpResponseBody(w, "Bad params", 400)
		return
	}

	restaurant, err := _models.BaseResStorage.GetRestaurantByID(uint(restaurantId))

	if err != nil {
		HttpResponseBody(w, "Not found", 404)
		return
	}

	_ = json.NewEncoder(w).Encode(restaurant)
}
