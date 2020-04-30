package repository

import (
	"github.com/2020_1_Skycode/internal/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeoDataRepository_GetGeoPosByAddress(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{
						  "response": {
							"GeoObjectCollection": {
							  "metaDataProperty": {
								"GeocoderResponseMetaData": {
								  "boundedBy": {
									"Envelope": {
									  "lowerCorner": "-0.002497 -0.002496",
									  "upperCorner": "0.002497 0.002496"
									}
								  },
								  "request": "Россия, Москва, улица Большая Лубянка, 1с1",
								  "results": "10",
								  "found": "1"
								}
							  },
							  "featureMember": [
								{
								  "GeoObject": {
									"metaDataProperty": {
									  "GeocoderMetaData": {
										"precision": "exact",
										"text": "Россия, Москва, улица Большая Лубянка, 1с1",
										"kind": "house",
										"Address": {
										  "country_code": "RU",
										  "formatted": "Россия, Москва, улица Большая Лубянка, 1с1",
										  "postal_code": "107031",
										  "Components": [
											{
											  "kind": "country",
											  "name": "Россия"
											},
											{
											  "kind": "province",
											  "name": "Центральный федеральный округ"
											},
											{
											  "kind": "province",
											  "name": "Москва"
											},
											{
											  "kind": "locality",
											  "name": "Москва"
											},
											{
											  "kind": "street",
											  "name": "улица Большая Лубянка"
											},
											{
											  "kind": "house",
											  "name": "1с1"
											}
										  ]
										},
										"AddressDetails": {
										  "Country": {
											"AddressLine": "Россия, Москва, улица Большая Лубянка, 1с1",
											"CountryNameCode": "RU",
											"CountryName": "Россия",
											"AdministrativeArea": {
											  "AdministrativeAreaName": "Москва",
											  "Locality": {
												"LocalityName": "Москва",
												"Thoroughfare": {
												  "ThoroughfareName": "улица Большая Лубянка",
												  "Premise": {
													"PremiseNumber": "1с1",
													"PostalCode": {
													  "PostalCodeNumber": "107031"
													}
												  }
												}
											  }
											}
										  }
										}
									  }
									},
									"name": "улица Большая Лубянка, 1с1",
									"description": "Москва, Россия",
									"boundedBy": {
									  "Envelope": {
										"lowerCorner": "37.622127 55.758423",
										"upperCorner": "37.630337 55.763052"
									  }
									},
									"Point": {
									  "pos": "37.626232 55.760737"
									}
								  }
								}
							  ]
							}
						  }
						}`))
		}),
	)
	defer ts.Close()
	MapApi = ts.URL

	repo := NewGeoDataRepository("1234")

	exp := &models.GeoPos{
		Latitude:  37.626232,
		Longitude: 55.760737,
	}

	gp, err := repo.GetGeoPosByAddress("Pushkina dom Kolotushkina")
	require.NoError(t, err)
	require.EqualValues(t, exp, gp)
}
