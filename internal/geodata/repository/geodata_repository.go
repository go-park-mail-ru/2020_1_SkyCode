package repository

import (
	"encoding/json"
	"github.com/2020_1_Skycode/internal/geodata"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GeoDataRepository struct {
	apikey string
}

func NewGeoDataRepository(key string) geodata.Repository {
	return &GeoDataRepository{
		apikey: key,
	}
}

func (gr *GeoDataRepository) GetGeoPosByAddress(addr string) (*models.GeoPos, error) {
	gp := &models.GeoPos{}

	client := &http.Client{
		Timeout: time.Second * 4,
	}

	req, err := http.NewRequest("GET", "https://geocode-maps.yandex.ru/1.x/", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("apikey", gr.apikey)
	q.Add("geocode", addr)
	q.Add("format", "json")
	q.Add("results", "1")

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, tools.ApiResponseStatusNotOK
	}

	type YandexResponce struct {
		Resp struct {
			GeoObjColl struct {
				FeatMem []*struct {
					GeoObj struct {
						Point struct {
							Pos string `json:"pos"`
						} `json:"Point"`
						MetaDataProp struct {
							GeocoderMeta struct {
								Kind string `json:"kind"`
							} `json:"GeocoderMetaData"`
						} `json:"metaDataProperty"`
					} `json:"GeoObject"`
				} `json:"featureMember"`
			} `json:"GeoObjectCollection"`
		} `json:"response"`
	}

	p := &YandexResponce{}

	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	if len(p.Resp.GeoObjColl.FeatMem) == 0 {
		return nil, tools.ApiAnswerEmptyResult
	}

	if len(p.Resp.GeoObjColl.FeatMem) > 1 {
		return nil, tools.ApiMultiAnswerError
	}

	if p.Resp.GeoObjColl.FeatMem[0].GeoObj.Point.Pos == "" {
		return nil, tools.ApiAnswerEmptyResult
	}

	if p.Resp.GeoObjColl.FeatMem[0].GeoObj.MetaDataProp.GeocoderMeta.Kind != "house" {
		return nil, tools.ApiNotHouseAnswerError
	}

	data := strings.Split(p.Resp.GeoObjColl.FeatMem[0].GeoObj.Point.Pos, " ")
	gp.Latitude, err = strconv.ParseFloat(data[0], 64)
	if err != nil {
		return nil, err
	}
	gp.Longitude, err = strconv.ParseFloat(data[1], 64)
	if err != nil {
		return nil, err
	}

	return gp, nil
}
