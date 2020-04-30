package delivery

import (
	"encoding/json"
	mock_geodata "github.com/2020_1_Skycode/internal/geodata/mocks"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeoDataHandler_CheckAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	returnGeoData := &models.GeoPos{
		Latitude:  34.675,
		Longitude: 55.567,
	}

	testAddress := "Pushkina dom Kolotushkina"

	geodataUcase := mock_geodata.NewMockUseCase(ctrl)

	geodataUcase.EXPECT().CheckGeoPos(testAddress).Return(returnGeoData, nil)

	expectResult := tools.Body{"geopos": returnGeoData}

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")

	_ = NewGeoDataHandler(privateGroup, publicGroup, geodataUcase)

	target := "/api/v1/check_address"
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("address", testAddress)
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	result, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	expectResp, err := json.Marshal(expectResult)
	require.NoError(t, err)

	require.EqualValues(t, expectResp, result)
}
