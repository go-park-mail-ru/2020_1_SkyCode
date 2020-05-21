package restaurants_tags

import "github.com/2020_1_Skycode/internal/models"

type UseCase interface {
	CreateTag(tag *models.RestTag) error
	GetTagByID(id uint64) (*models.RestTag, error)
	GetAllTags() ([]*models.RestTag, error)
	UpdateTag(tag *models.RestTag) error
	DeleteTag(id uint64) error
}
