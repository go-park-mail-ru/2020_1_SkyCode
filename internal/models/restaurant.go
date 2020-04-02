package models

import "errors"

type Restaurant struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rating      float32   `json:"rating"`
	Image       string    `json:"image"`
	Products    []Product `json:"products"`
}

type ResStorage struct {
	Restaurants map[uint]*Restaurant
	nextID      uint
}

func (storage *ResStorage) GetRestaurantByID(id uint) (*Restaurant, error) {
	result := storage.Restaurants[id]

	if result == nil {
		return result, errors.New("no such element")
	}

	return result, nil
}

var ResArray = map[uint]*Restaurant{
	1: &Restaurant{
		Name:        "KFC",
		Description: "My lovely KFC",
		Rating:      3.14,
		Image:       "/kfc.jpg",
		Products: []Product{
			{
				Name:  "Potato free",
				Price: 150,
				Image: "/images/Vestern-Gurme-MakKombo.png",
			},
			{
				Name:  "Hamburger",
				Price: 230,
				Image: "/images/Fish-Roll-MakKombo.png",
			},
			{
				Name:  "Cola",
				Price: 70,
				Image: "/images/Vestern-Gurme-MakKombo.png",
			},
		},
	},
	2: &Restaurant{
		Name:        "Mac",
		Description: "My lovely Mac",
		Rating:      3.14,
		Image:       "/mac.jpg",
		Products: []Product{
			{
				Name:  "Potato free",
				Price: 150,
				Image: "/mcdonalds/Vestern-Gurme-MakKombo.png",
			},
			{
				Name:  "Hamburger",
				Price: 230,
				Image: "/mcdonalds/Fish-Roll-MakKombo.png",
			},
			{
				Name:  "Cola",
				Price: 70,
				Image: "/mcdonalds/Vestern-Gurme-MakKombo.png",
			},
		},
	},
	3: &Restaurant{
		Name:        "BK",
		Description: "My lovely BK (no)",
		Rating:      3.14,
		Image:       "/bk.jpg",
		Products: []Product{
			{
				Name:  "Potato free",
				Price: 150,
				Image: "/potatoFree.jpg",
			},
			{
				Name:  "Hamburger",
				Price: 230,
				Image: "/hamburger.jpg",
			},
			{
				Name:  "Cola",
				Price: 70,
				Image: "/cola.jpg",
			},
		},
	},
}

var BaseResStorage = ResStorage{
	Restaurants: ResArray,
	nextID:      4,
}
