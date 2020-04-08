package delivery

import (
	"encoding/json"
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	mock_products "github.com/2020_1_Skycode/internal/products/mocks"
	mock_restaurants "github.com/2020_1_Skycode/internal/restaurants/mocks"
	mock_sessions "github.com/2020_1_Skycode/internal/sessions/mocks"
	"github.com/2020_1_Skycode/internal/tools"
	_csrfManager "github.com/2020_1_Skycode/internal/tools/CSRFManager"
	_rValidator "github.com/2020_1_Skycode/internal/tools/requestValidator"
	mock_users "github.com/2020_1_Skycode/internal/users/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

//func TestProductHandler_CreateProduct(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockProdUC := mock_products.NewMockUseCase(ctrl)
//	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
//	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
//	mockUserUC := mock_users.NewMockUseCase(ctrl)
//
//	reqProd := &models.Product{
//		Name:   "test1",
//		Price: 	2.50,
//	}
//
//	type productRequest struct {
//		Name  string  `json:"name, omitempty" binding:"required" validate:"min=2"`
//		Price float32 `json:"price, omitempty" binding:"required"`
//	}
//
//	restID := uint64(1)
//	userID := uint64(1)
//	prodReq := &productRequest{
//		Name:        reqProd.Name,
//		Price: 		 reqProd.Price,
//	}
//
//	restRes := &models.Restaurant{ID: restID}
//	sessRes := &models.Session{UserId: userID}
//	userRes := &models.User{Role: "Admin"}
//
//	j, err := json.Marshal(prodReq)
//	require.NoError(t, err)
//
//	expectResult := &tools.Message{ Message:"Product has been created" }
//
//	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
//	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
//	mockRestUC.EXPECT().GetRestaurantByID(restID).Return(restRes, nil)
//	mockProdUC.EXPECT().CreateProduct(reqProd).Return(nil)
//
//	g := gin.New()
//
//	csrfManager := _csrfManager.NewCSRFManager()
//	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)
//
//	publicGroup := g.Group("/api/v1")
//	privateGroup := g.Group("/api/v1")
//	reqValidator := _rValidator.NewRequestValidator()
//
//	_ = NewProductHandler(privateGroup, publicGroup, mockProdUC, reqValidator, mockRestUC, mwareC)
//
//	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/product"
//	req, err:= http.NewRequest("POST", target, strings.NewReader(string(j)))
//	require.NoError(t, err)
//	req.Header.Set("Content-Type", "application/json")
//	req.AddCookie(&http.Cookie{
//		Name:"SkyDelivery",
//		Value:"1234",
//	})
//	w := httptest.NewRecorder()
//
//	g.ServeHTTP(w, req)
//
//	if w.Code != http.StatusOK {
//		t.Error("Status is not ok")
//		return
//	}
//
//	var result *tools.Message
//	_ = json.NewDecoder(w.Result().Body).Decode(&result)
//
//	require.EqualValues(t, expectResult, result)
//}

func TestProductHandler_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProdUC := mock_products.NewMockUseCase(ctrl)
	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	prodID := uint64(1)
	userID := uint64(1)

	resProd := &models.Product{
		ID:    1,
		Name:  "test1",
		Price: 2.50,
		Image: "",
	}

	resResp := &tools.Message{Message: "success"}

	sessRes := &models.Session{UserId: userID}
	userRes := &models.User{Role: "Admin"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockProdUC.EXPECT().GetProductByID(resProd.ID).Return(resProd, nil)
	mockProdUC.EXPECT().DeleteProduct(prodID).Return(nil)

	g := gin.New()

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewProductHandler(privateGroup, publicGroup, mockProdUC, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/products/" + strconv.Itoa(int(prodID)) + "/delete"
	req, err := http.NewRequest("DELETE", target, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, resResp, result)
}

func TestProductHandler_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProdUC := mock_products.NewMockUseCase(ctrl)
	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	resProd := &models.Product{
		ID:    1,
		Name:  "test1",
		Price: 2.50,
	}

	userID := uint64(1)

	sessRes := &models.Session{UserId: userID}
	userRes := &models.User{Role: "User"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockProdUC.EXPECT().GetProductByID(resProd.ID).Return(resProd, nil)

	g := gin.New()

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewProductHandler(privateGroup, publicGroup, mockProdUC, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/products/" + strconv.Itoa(int(resProd.ID))
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	var result *models.Product
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, resProd, result)
}

func TestProductHandler_GetProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProdUC := mock_products.NewMockUseCase(ctrl)
	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	resProd := []*models.Product{
		{
			ID:     1,
			Name:   "test1",
			Price:  2.50,
			RestId: 1,
		},
		{
			ID:     2,
			Name:   "test2",
			Price:  2.50,
			RestId: 1,
		},
	}

	restID := uint64(1)
	userID := uint64(1)

	sessRes := &models.Session{UserId: userID}
	userRes := &models.User{Role: "User"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockProdUC.EXPECT().GetProductsByRestaurantID(restID).Return(resProd, nil)

	g := gin.New()

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewProductHandler(privateGroup, publicGroup, mockProdUC, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/restaurants/" + strconv.Itoa(int(restID)) + "/product"
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	var result []*models.Product
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, resProd, result)
}

//func TestProductHandler_UpdateImage(t *testing.T) {
//
//}

func TestProductHandler_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProdUC := mock_products.NewMockUseCase(ctrl)
	mockRestUC := mock_restaurants.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	prodID := uint64(1)

	type productRequest struct {
		Name  string  `json:"name, omitempty" binding:"required" validate:"min=2"`
		Price float32 `json:"price, omitempty" binding:"required"`
	}

	reqProd := &models.Product{
		ID:    prodID,
		Name:  "test1",
		Price: 2.50,
	}

	reqJson := &productRequest{
		Name:  reqProd.Name,
		Price: reqProd.Price,
	}

	j, err := json.Marshal(reqJson)
	require.NoError(t, err)

	userID := uint64(1)

	resResp := &tools.Message{Message: "Product has been updated"}

	sessRes := &models.Session{UserId: userID}
	userRes := &models.User{Role: "Admin"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(userID).Return(userRes, nil)
	mockProdUC.EXPECT().UpdateProduct(reqProd).Return(nil)

	g := gin.New()

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewProductHandler(privateGroup, publicGroup, mockProdUC, reqValidator, mockRestUC, mwareC)

	target := "/api/v1/products/" + strconv.Itoa(int(reqProd.ID)) + "/update"
	req, err := http.NewRequest("PUT", target, strings.NewReader(string(j)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "SkyDelivery",
		Value: "1234",
	})
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	var result *tools.Message
	_ = json.NewDecoder(w.Result().Body).Decode(&result)

	require.EqualValues(t, resResp, result)
}
