package usecase

import (
	"github.com/2020_1_Skycode/internal/models"
	mock_reviews "github.com/2020_1_Skycode/internal/reviews/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestReviewsUseCase_DeleteReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mock_reviews.NewMockRepository(ctrl)

	testReview := &models.Review{
		ID:     1,
		RestID: 1,
		Text:   "Good",
		Author: &models.User{
			ID:        1,
			FirstName: "abc",
			LastName:  "asd",
		},
		CreationDate: time.Now(),
		Rate:         5,
	}

	mockReviewRepo.EXPECT().GetReviewByID(testReview.ID).Return(testReview, nil)
	mockReviewRepo.EXPECT().DeleteReview(testReview.ID).Return(nil)

	reviewUCase := NewReviewsUseCase(mockReviewRepo)

	err := reviewUCase.DeleteReview(testReview.ID, testReview.Author)
	require.NoError(t, err)
}

func TestReviewsUseCase_GetReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mock_reviews.NewMockRepository(ctrl)

	testReview := &models.Review{
		ID:     1,
		RestID: 1,
		Text:   "Good",
		Author: &models.User{
			ID:        1,
			FirstName: "abc",
			LastName:  "asd",
		},
		CreationDate: time.Now(),
		Rate:         5,
	}

	mockReviewRepo.EXPECT().GetReviewByID(testReview.ID).Return(testReview, nil)

	reviewUCase := NewReviewsUseCase(mockReviewRepo)

	result, err := reviewUCase.GetReview(testReview.ID)
	require.NoError(t, err)

	require.EqualValues(t, testReview, result)
}

func TestReviewsUseCase_GetUserReviews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mock_reviews.NewMockRepository(ctrl)

	testID := uint64(1)

	testReviews := []*models.Review{
		{
			ID:     1,
			RestID: 1,
			Text:   "Good",
			Author: &models.User{
				ID:        testID,
				FirstName: "abc",
				LastName:  "asd",
			},
			CreationDate: time.Now(),
			Rate:         5,
		},
	}

	mockReviewRepo.EXPECT().GetReviewsByUserID(testID, uint64(1), uint64(1)).Return(testReviews, nil)
	mockReviewRepo.EXPECT().GetReviewsCountByUserID(testID).Return(uint64(1), nil)

	reviewUCase := NewReviewsUseCase(mockReviewRepo)

	result, total, err := reviewUCase.GetUserReviews(testID, 1, 1)
	require.NoError(t, err)

	require.EqualValues(t, testReviews, result)
	require.EqualValues(t, uint64(1), total)
}

func TestReviewsUseCase_UpdateReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := mock_reviews.NewMockRepository(ctrl)

	testReview := &models.Review{
		ID:     1,
		RestID: 1,
		Text:   "Good",
		Author: &models.User{
			ID:        1,
			FirstName: "abc",
			LastName:  "asd",
		},
		CreationDate: time.Now(),
		Rate:         5,
	}

	mockReviewRepo.EXPECT().GetReviewByID(testReview.ID).Return(testReview, nil)
	mockReviewRepo.EXPECT().UpdateReview(testReview).Return(nil)

	reviewUCase := NewReviewsUseCase(mockReviewRepo)

	err := reviewUCase.UpdateReview(testReview, testReview.Author)
	require.NoError(t, err)
}
