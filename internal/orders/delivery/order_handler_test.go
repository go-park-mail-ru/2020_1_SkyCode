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
	"github.com/stretchr/testify/require"
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
		Address:   orderReq.Address,
		Phone:     orderReq.Phone,
		Comment:   orderReq.Comment,
		PersonNum: orderReq.PersonNum,
		Products:  orderReq.Products,
		Price:     orderReq.Price,
	}

	reqJson, err := json.Marshal(orderReq)
	require.NoError(t, err)

	sessRes := &models.Session{UserId: 1}
	userRes := &models.User{Role: "User"}

	expectRes := &tools.Message{Message: "success"}

	mockSessUC.EXPECT().GetSession("1234").Return(sessRes, nil)
	mockUserUC.EXPECT().GetUserById(sessRes.UserId).Return(userRes, nil)
	mockOrderUC.EXPECT().CheckoutOrder(order).Return(nil)

	g := gin.New()

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
