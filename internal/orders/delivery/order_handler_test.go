package delivery

import (
	"encoding/json"
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	mock_orders "github.com/2020_1_Skycode/internal/orders/mocks"
	mock_sessions "github.com/2020_1_Skycode/internal/sessions/mocks"
	"github.com/2020_1_Skycode/internal/tools"
	_csrfManager "github.com/2020_1_Skycode/internal/tools/CSRFManager"
	_rValidator "github.com/2020_1_Skycode/internal/tools/requestValidator"
	mock_users "github.com/2020_1_Skycode/internal/users/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOrderHandler_Checkout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUC := mock_orders.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	orderReq := orderRequest{
		UserID:    1,
		RestID:    1,
		Address:   "Pushkina dom Kolotushkina",
		Comment:   "Faster",
		Phone:     "89765433221",
		PersonNum: 3,
		Products: []*models.OrderProduct{
			{ProductID: 1, Count: 3},
		},
		Price: 300,
	}

	order := &models.Order{
		UserID:    orderReq.UserID,
		RestID:    orderReq.RestID,
		Address:   orderReq.Address,
		Phone:     orderReq.Phone,
		Comment:   orderReq.Comment,
		PersonNum: orderReq.PersonNum,
		Products:  nil,
		Price:     orderReq.Price,
	}

	orderProd := orderReq.Products

	reqJson, err := json.Marshal(orderReq)
	require.NoError(t, err)

	sessRes := &models.Session{ID: 1, UserId: 1}
	userRes := &models.User{ID: sessRes.UserId, Role: "User"}

	expectRes := &tools.Message{Message: "success"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(sessRes.UserId).Return(userRes, nil)
	mockOrderUC.EXPECT().CheckoutOrder(order, orderProd).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewOrderHandler(privateGroup, publicGroup, mockOrderUC, reqValidator, mwareC)

	target := "/api/v1/orders/checkout"
	req, err := http.NewRequest("POST", target, strings.NewReader(string(reqJson)))
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

	require.EqualValues(t, expectRes, result)
}

func TestOrderHandler_GetUserOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUC := mock_orders.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	orderReq := orderRequest{
		UserID:    1,
		RestID:    1,
		Address:   "Pushkina dom Kolotushkina",
		Comment:   "Faster",
		Phone:     "89765433221",
		PersonNum: 3,
		Price:     300,
	}

	order := &models.Order{
		ID:        1,
		UserID:    orderReq.UserID,
		RestID:    orderReq.RestID,
		Address:   orderReq.Address,
		Phone:     orderReq.Phone,
		Comment:   orderReq.Comment,
		PersonNum: orderReq.PersonNum,
		Price:     orderReq.Price,
	}

	expectRes := &tools.Body{"order": order}

	sessRes := &models.Session{ID: 1, UserId: 1}
	userRes := &models.User{ID: sessRes.UserId, Role: "User"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(sessRes.UserId).Return(userRes, nil)
	mockOrderUC.EXPECT().GetOrderByID(order.ID, order.UserID).Return(order, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewOrderHandler(privateGroup, publicGroup, mockOrderUC, reqValidator, mwareC)

	target := "/api/v1/orders/1"
	req, err := http.NewRequest("GET", target, nil)
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

	result, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	expectResp, err := json.Marshal(expectRes)
	require.NoError(t, err)

	require.EqualValues(t, expectResp, result)
}

func TestOrderHandler_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUC := mock_orders.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	orderReq := orderRequest{
		UserID:    1,
		RestID:    1,
		Address:   "Pushkina dom Kolotushkina",
		Comment:   "Faster",
		Phone:     "89765433221",
		PersonNum: 3,
		Price:     300,
	}

	orders := []*models.Order{
		{
			ID:        1,
			UserID:    orderReq.UserID,
			RestID:    orderReq.RestID,
			Address:   orderReq.Address,
			Phone:     orderReq.Phone,
			Comment:   orderReq.Comment,
			PersonNum: orderReq.PersonNum,
			Price:     orderReq.Price,
		},
	}

	total := uint64(1)

	expectRes := &tools.Body{
		"orders": orders,
		"total":  total,
	}

	sessRes := &models.Session{ID: 1, UserId: 1}
	userRes := &models.User{ID: sessRes.UserId, Role: "User"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(sessRes.UserId).Return(userRes, nil)
	mockOrderUC.EXPECT().GetAllUserOrders(sessRes.UserId, uint64(1), uint64(1)).Return(orders, total, nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewOrderHandler(privateGroup, publicGroup, mockOrderUC, reqValidator, mwareC)

	target := "/api/v1/orders"
	req, err := http.NewRequest("GET", target, nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("count", "1")
	req.URL.RawQuery = q.Encode()

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

	result, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	expectResp, err := json.Marshal(expectRes)
	require.NoError(t, err)

	require.EqualValues(t, expectResp, result)
}

func TestOrderHandler_DeleteOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUC := mock_orders.NewMockUseCase(ctrl)
	mockSessUC := mock_sessions.NewMockUseCase(ctrl)
	mockUserUC := mock_users.NewMockUseCase(ctrl)

	sessRes := &models.Session{ID: 1, UserId: 1}
	userRes := &models.User{ID: 1, Role: "User"}

	expectRes := &tools.Message{Message: "success"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(sessRes.UserId).Return(userRes, nil)
	mockOrderUC.EXPECT().DeleteOrder(sessRes.UserId, uint64(1)).Return(nil)

	g := gin.New()
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)

	csrfManager := _csrfManager.NewCSRFManager()
	mwareC := _middleware.NewMiddleWareController(g, mockSessUC, mockUserUC, csrfManager)

	publicGroup := g.Group("/api/v1")
	privateGroup := g.Group("/api/v1")
	reqValidator := _rValidator.NewRequestValidator()

	_ = NewOrderHandler(privateGroup, publicGroup, mockOrderUC, reqValidator, mwareC)

	target := "/api/v1/orders/1"
	req, err := http.NewRequest("DELETE", target, nil)
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

	require.EqualValues(t, expectRes, result)
}
