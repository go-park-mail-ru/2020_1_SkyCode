package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/require"

	_handlers "github.com/2020_1_Skycode/internal/handlers"
	_models "github.com/2020_1_Skycode/internal/models"
)

func TestGetRestaurants(t *testing.T) {
	t.Parallel()

	h := _handlers.RestaurantHandler{
		Restaurants: _models.BaseResStorage,
	}

	expectedRestaurants := h.Restaurants

	r := httptest.NewRequest("GET", "/restaurants", nil)
	w := httptest.NewRecorder()

	h.GetRestaurants(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	var result _models.ResStorage
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, expectedRestaurants.Restaurants, result.Restaurants)
}

func TestGetRestaurantByID(t *testing.T) {
	t.Parallel()

	testID := uint(2)

	h := _handlers.RestaurantHandler{
		Restaurants: _models.BaseResStorage,
	}

	expectedRest, err := _models.BaseResStorage.GetRestaurantByID(testID)

	if err != nil {
		t.Error("Error with get by ID BasResStorage")
		return
	}

	idstr := strconv.Itoa(int(testID))
	r := httptest.NewRequest("GET", "/restaurants/", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"restaurant_id": idstr,
	})
	h.GetRestaurantByID(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	var result _models.Restaurant
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, *expectedRest, result)
}

func TestGetRestaurantByID404(t *testing.T) {
	t.Parallel()

	testID := uint(10)

	h := _handlers.RestaurantHandler{
		Restaurants: _models.BaseResStorage,
	}

	idstr := strconv.Itoa(int(testID))
	r := httptest.NewRequest("GET", "/restaurants/", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"restaurant_id": idstr,
	})
	h.GetRestaurantByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("Status is not 404")
	}
}
