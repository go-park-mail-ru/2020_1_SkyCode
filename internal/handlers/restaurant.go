package handlers

import (
	json "encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_models "github.com/2020_1_Skycode/internal/models"
	"github.com/gorilla/mux"
)

type RestaurantHandler struct {
	Restaurants _models.ResStorage
}

func (api *RestaurantHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `No such method`, 405)
		return
	}

	data, err := json.Marshal(api.Restaurants)

	if err != nil {
		fmt.Println(err)
		http.Error(w, `Server error`, 500)
		return
	}

	_, err = w.Write(data)

	if err != nil {
		fmt.Println(err)
		http.Error(w, `Server error`, 500)
	}
}

func (api *RestaurantHandler) GetRestaurantByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `No such method`, 405)
		return
	}

	vars := mux.Vars(r)

	restaurantId, err := strconv.Atoi(vars["restaurant_id"])

	if err != nil {
		http.Error(w, `Bad params`, 400)
		return
	}

	restaurant, err := _models.BaseResStorage.GetRestaurantByID(uint(restaurantId))

	if err != nil {
		fmt.Println(err)
		http.Error(w, `Not found`, 404)
		return
	}

	data, err := json.Marshal(restaurant)

	if err != nil {
		fmt.Println(err)
		http.Error(w, `Server error`, 500)
		return
	}

	_, err = w.Write(data)

	if err != nil {
		fmt.Println(err)
		http.Error(w, `Server error`, 500)
	}
}
