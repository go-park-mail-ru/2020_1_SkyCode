package usecase

import (
	mock_geodata "github.com/2020_1_Skycode/internal/geodata/mocks"
	"github.com/2020_1_Skycode/internal/models"
	mock_orders "github.com/2020_1_Skycode/internal/orders/mocks"
	mock_prodtags "github.com/2020_1_Skycode/internal/product_tags/mocks"
	mock_restpoints "github.com/2020_1_Skycode/internal/restaurant_points/mocks"
	mock_restaurants "github.com/2020_1_Skycode/internal/restaurants/mocks"
	mock_resttags "github.com/2020_1_Skycode/internal/restaurants_tags/mocks"
	mock_reviews "github.com/2020_1_Skycode/internal/reviews/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRestaurantUseCase_GetRestaurants(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	resultRests := []*models.Restaurant{
		{ID: 1, Name: "test1", Rating: 4.2, Image: "./default.jpg"},
	}

	expectTotal := uint64(1)

	mockRestaurantsRepo.EXPECT().GetAll(uint64(1), uint64(1), uint64(0)).Return(resultRests, expectTotal, nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	rests, total, err := restUCase.GetRestaurants(uint64(1), uint64(1), uint64(0))

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, resultRests, rests)
	require.EqualValues(t, expectTotal, total)
}

func TestRestaurantUseCase_GetRestaurantByID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	resultRest := &models.Restaurant{
		ID: uint64(1), ManagerID: uint64(1), Name: "test1",
		Description: "smthng", Rating: 4.2, Image: "./default.jpg",
	}
	var elemID uint64 = 1

	mockRestaurantsRepo.EXPECT().GetByID(elemID).Return(resultRest, nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	rest, err := restUCase.GetRestaurantByID(elemID)

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, resultRest, rest)
}

func TestRestaurantUseCase_CreateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testRest := &models.Restaurant{
		ManagerID:   1,
		Name:        "test",
		Description: "smthng",
		Rating:      4.2,
		Image:       "./default.jpg",
	}

	mockRestaurantsRepo.EXPECT().InsertInto(testRest).Return(nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	err := restUCase.CreateRestaurant(testRest)
	require.NoError(t, err)
}

func TestRestaurantUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	restID := uint64(1)

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	mockRestaurantsRepo.EXPECT().Delete(restID).Return(nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	err := restUCase.Delete(restID)
	require.NoError(t, err)
}

func TestRestaurantUseCase_UpdateImage(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testRest := &models.Restaurant{
		ID:    uint64(1),
		Image: "./default.jpg",
	}

	mockRestaurantsRepo.EXPECT().UpdateImage(testRest).Return(nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	err := restUCase.UpdateImage(testRest.ID, testRest.Image)
	require.NoError(t, err)
}

func TestRestaurantUseCase_UpdateRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testRest := &models.Restaurant{
		ID:          1,
		Name:        "test",
		Description: "smthng",
	}

	mockRestaurantsRepo.EXPECT().Update(testRest).Return(nil)
	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	err := restUCase.UpdateRestaurant(testRest)
	require.NoError(t, err)
}

func TestRestaurantUseCase_AddPoint(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)

	testRest := &models.Restaurant{
		ID:          1,
		Name:        "test",
		Description: "smthng",
	}

	testGP := &models.GeoPos{
		Longitude: 55.753227,
		Latitude:  37.619030,
	}

	testPoint := &models.RestaurantPoint{
		Address:       "Pushkina dom Kolotushkina",
		MapPoint:      testGP,
		ServiceRadius: 5,
		RestID:        testRest.ID,
	}

	mockRestaurantsRepo.EXPECT().GetByID(testRest.ID).Return(testRest, nil)
	mockGeoDataRepo.EXPECT().GetGeoPosByAddress(testPoint.Address).Return(testGP, nil)
	mockRestPointsRepo.EXPECT().InsertInto(testPoint).Return(nil)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	err := restUCase.AddPoint(testPoint)
	require.NoError(t, err)
}

func TestRestaurantUseCase_AddReview(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testRest := &models.Restaurant{
		ID:          1,
		Name:        "test",
		Description: "smthng",
	}

	testReview := &models.Review{
		RestID:       testRest.ID,
		Author:       &models.User{ID: 1},
		Text:         "Review",
		CreationDate: time.Now(),
		Rate:         5,
	}

	mockRestaurantsRepo.EXPECT().GetByID(testRest.ID).Return(testRest, nil)
	mockReviewsRepo.EXPECT().GetRestaurantReviewByUser(testRest.ID, testReview.Author.ID).
		Return(nil, nil)
	mockReviewsRepo.EXPECT().CreateReview(testReview).Return(nil)

	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	err := restUCase.AddReview(testReview)
	require.NoError(t, err)
}

func TestRestaurantUseCase_GetPoints(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testRest := &models.Restaurant{
		ID:          1,
		Name:        "test",
		Description: "smthng",
		Points:      nil,
	}

	testGP := &models.GeoPos{
		Longitude: 55.753227,
		Latitude:  37.619030,
	}

	testPoints := []*models.RestaurantPoint{
		{
			Address:       "Pushkina dom Kolotushkina",
			MapPoint:      testGP,
			ServiceRadius: 5,
			RestID:        testRest.ID,
		},
	}

	expectTotal := uint64(1)

	mockRestaurantsRepo.EXPECT().GetByID(testRest.ID).Return(testRest, nil)
	mockRestPointsRepo.EXPECT().GetPointsByRestID(testRest.ID, uint64(1), uint64(1)).
		Return(testPoints, expectTotal, nil)

	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	result, total, err := restUCase.GetPoints(testRest.ID, uint64(1), uint64(1))
	require.NoError(t, err)
	require.EqualValues(t, testPoints, result)
	require.EqualValues(t, expectTotal, total)
}

func TestRestaurantUseCase_GetRestaurantsInServiceRadius(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testGP := &models.GeoPos{
		Longitude: 55.753227,
		Latitude:  37.619030,
	}

	testPoint := &models.RestaurantPoint{
		Address:       "Pushkina dom Kolotushkina",
		MapPoint:      testGP,
		ServiceRadius: 5,
		RestID:        1,
	}

	testRests := []*models.Restaurant{
		{
			ID:          1,
			Name:        "test",
			Description: "smthng",
			Points:      []*models.RestaurantPoint{testPoint},
		},
	}

	expectTotal := uint64(1)

	mockGeoDataRepo.EXPECT().GetGeoPosByAddress(testPoint.Address).Return(testGP, nil)
	mockRestaurantsRepo.EXPECT().GetAllInServiceRadius(testGP, uint64(1), uint64(1), uint64(0)).Return(testRests, expectTotal, nil)
	mockRestPointsRepo.EXPECT().GetCloserPointByRestID(testPoint.RestID, testGP).Return(testPoint, nil)

	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	result, total, err := restUCase.GetRestaurantsInServiceRadius(testPoint.Address, uint64(1), uint64(1), uint64(0))
	require.NoError(t, err)
	require.EqualValues(t, testRests, result)
	require.EqualValues(t, expectTotal, total)
}

func TestRestaurantUseCase_GetReviews(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)
	mockRestTagsRepo := mock_resttags.NewMockRepository(ctrl)
	mockProdTagsRepo := mock_prodtags.NewMockRepository(ctrl)
	mockOrdersRepo := mock_orders.NewMockRepository(ctrl)

	testRest := &models.Restaurant{
		ID:          1,
		Name:        "test",
		Description: "smthng",
	}

	testReview := &models.Review{
		RestID:       testRest.ID,
		Author:       &models.User{ID: 1},
		Text:         "Review",
		CreationDate: time.Now(),
		Rate:         5,
	}

	mockRestaurantsRepo.EXPECT().GetByID(testRest.ID).Return(testRest, nil)
	mockReviewsRepo.EXPECT().GetReviewsByRestID(testRest.ID, uint64(1), uint64(1)).
		Return([]*models.Review{testReview}, nil)
	mockReviewsRepo.EXPECT().GetReviewsCountByRestID(testRest.ID).Return(uint64(1), nil)
	mockReviewsRepo.EXPECT().GetRestaurantReviewByUser(testRest.ID, testReview.Author.ID).
		Return(testReview, nil)

	restUCase := NewRestaurantsUseCase(mockRestaurantsRepo, mockRestPointsRepo, mockReviewsRepo, mockGeoDataRepo,
		mockRestTagsRepo, mockProdTagsRepo, mockOrdersRepo)

	result, current, total, err := restUCase.GetReviews(testRest.ID, testReview.Author.ID, uint64(1), uint64(1))
	require.NoError(t, err)
	require.EqualValues(t, []*models.Review{testReview}, result)
	require.EqualValues(t, testReview, current)
	require.EqualValues(t, uint64(1), total)
}
