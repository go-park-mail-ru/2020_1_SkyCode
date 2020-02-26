package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/require"

	json2 "encoding/json"

	_handlers "github.com/2020_1_Skycode/internal/handlers"
	_models "github.com/2020_1_Skycode/internal/models"
)

func TestGetRestaurants(t *testing.T) {
	t.Parallel()

	h := _handlers.RestaurantHandler{
		Restaurants: _models.BaseResStorage,
	}

	expectedRestaurants, err := json2.Marshal(_models.BaseResStorage)

	if err != nil {
		t.Error("Error with BasResStorage")
		return
	}

	r := httptest.NewRequest("GET", "/restaurants", nil)
	w := httptest.NewRecorder()

	h.GetRestaurants(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	result, errRead := ioutil.ReadAll(w.Result().Body)
	if errRead != nil {
		t.Error("Error with read response body")
		return
	}

	require.EqualValues(t, expectedRestaurants, result)
}

func TestGetRestaurantByID(t *testing.T) {
	t.Parallel()

	testID := uint(2)

	h := _handlers.RestaurantHandler{
		Restaurants: _models.BaseResStorage,
	}

	restaurant, err := _models.BaseResStorage.GetRestaurantByID(testID)

	if err != nil {
		t.Error("Error with get by ID BasResStorage")
		return
	}

	expected, errMarsh := json2.Marshal(restaurant)

	if errMarsh != nil {
		t.Error("Error with marshal json")
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

	result, errRead := ioutil.ReadAll(w.Result().Body)
	if errRead != nil {
		t.Error("Error with read response body")
		return
	}

	require.EqualValues(t, expected, result)
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
