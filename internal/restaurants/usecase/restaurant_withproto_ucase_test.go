package usecase

import (
	mock_geodata "github.com/2020_1_Skycode/internal/geodata/mocks"
	"github.com/2020_1_Skycode/internal/models"
	mock_restpoints "github.com/2020_1_Skycode/internal/restaurant_points/mocks"
	mock_restaurants "github.com/2020_1_Skycode/internal/restaurants/mocks"
	mock_reviews "github.com/2020_1_Skycode/internal/reviews/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestRestaurantWithProtoUseCase_GetRestaurants(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)

	resultRests := []*models.Restaurant{
		{ID: 1, Name: "test1", Rating: 4.2, Image: "./default.jpg"},
	}

	expectTotal := uint64(1)

	mockRestaurantsRepo.EXPECT().GetAll(uint64(1), uint64(1)).Return(resultRests, expectTotal, nil)

	c := &grpc.ClientConn{}
	restUCase := NewRestaurantsWithProtoUseCase(mockRestaurantsRepo, mockRestPointsRepo,
		mockReviewsRepo, mockGeoDataRepo, c)

	rests, total, err := restUCase.GetRestaurants(uint64(1), uint64(1))

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, resultRests, rests)
	require.EqualValues(t, expectTotal, total)
}

func TestRestaurantWithProtoUseCase_GetRestaurantByID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)

	resultRest := &models.Restaurant{
		ID: uint64(1), ManagerID: uint64(1), Name: "test1",
		Description: "smthng", Rating: 4.2, Image: "./default.jpg",
	}
	var elemID uint64 = 1

	mockRestaurantsRepo.EXPECT().GetByID(elemID).Return(resultRest, nil)
	c := &grpc.ClientConn{}
	restUCase := NewRestaurantsWithProtoUseCase(mockRestaurantsRepo, mockRestPointsRepo,
		mockReviewsRepo, mockGeoDataRepo, c)
	rest, err := restUCase.GetRestaurantByID(elemID)

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
		return
	}

	require.EqualValues(t, resultRest, rest)
}

func TestRestaurantWithProtoUseCase_AddReview(t *testing.T) {
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

	c := &grpc.ClientConn{}
	restUCase := NewRestaurantsWithProtoUseCase(mockRestaurantsRepo, mockRestPointsRepo,
		mockReviewsRepo, mockGeoDataRepo, c)

	err := restUCase.AddReview(testReview)
	require.NoError(t, err)
}

func TestRestaurantWithProtoUseCase_GetPoints(t *testing.T) {
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

	c := &grpc.ClientConn{}
	restUCase := NewRestaurantsWithProtoUseCase(mockRestaurantsRepo, mockRestPointsRepo,
		mockReviewsRepo, mockGeoDataRepo, c)

	result, total, err := restUCase.GetPoints(testRest.ID, uint64(1), uint64(1))
	require.NoError(t, err)
	require.EqualValues(t, testPoints, result)
	require.EqualValues(t, expectTotal, total)
}

func TestRestaurantWithProtoUseCase_GetRestaurantsInServiceRadius(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRestaurantsRepo := mock_restaurants.NewMockRepository(ctrl)
	mockGeoDataRepo := mock_geodata.NewMockRepository(ctrl)
	mockReviewsRepo := mock_reviews.NewMockRepository(ctrl)
	mockRestPointsRepo := mock_restpoints.NewMockRepository(ctrl)

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
	mockRestaurantsRepo.EXPECT().GetAllInServiceRadius(testGP, uint64(1), uint64(1)).Return(testRests, expectTotal, nil)
	mockRestPointsRepo.EXPECT().GetCloserPointByRestID(testPoint.RestID, testGP).Return(testPoint, nil)

	c := &grpc.ClientConn{}
	restUCase := NewRestaurantsWithProtoUseCase(mockRestaurantsRepo, mockRestPointsRepo,
		mockReviewsRepo, mockGeoDataRepo, c)

	result, total, err := restUCase.GetRestaurantsInServiceRadius(testPoint.Address, uint64(1), uint64(1))
	require.NoError(t, err)
	require.EqualValues(t, testRests, result)
	require.EqualValues(t, expectTotal, total)
}

func TestRestaurantWithProtoUseCase_GetReviews(t *testing.T) {
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

	c := &grpc.ClientConn{}
	restUCase := NewRestaurantsWithProtoUseCase(mockRestaurantsRepo, mockRestPointsRepo,
		mockReviewsRepo, mockGeoDataRepo, c)

	result, current, total, err := restUCase.GetReviews(testRest.ID, testReview.Author.ID, uint64(1), uint64(1))
	require.NoError(t, err)
	require.EqualValues(t, []*models.Review{testReview}, result)
	require.EqualValues(t, testReview, current)
	require.EqualValues(t, uint64(1), total)
}
