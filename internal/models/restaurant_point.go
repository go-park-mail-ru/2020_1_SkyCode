package models

type RestaurantPoint struct {
	ID            uint64  `json:"id"`
	Address       string  `json:"address"`
	MapPoint      *GeoPos `json:"map_point"`
	ServiceRadius float64 `json:"service_radius"`
	RestID        uint64  `json:"rest_id"`
}

type GeoPos struct {
	Latitude  float64 `json:"latitude"`  //широта
	Longitude float64 `json:"longitude"` //долгота
}
