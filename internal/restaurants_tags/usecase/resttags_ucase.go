package usecase

import (
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/restaurants_tags"
	"github.com/2020_1_Skycode/internal/tools"
	"os"
	"path/filepath"
)

type RestTagsUseCase struct {
	restTagsRepo restaurants_tags.Repository
}

func NewRestTagsUCase(rtr restaurants_tags.Repository) restaurants_tags.UseCase {
	return &RestTagsUseCase{
		restTagsRepo: rtr,
	}
}

func (rtUC *RestTagsUseCase) CreateTag(tag *models.RestTag) error {
	if err := rtUC.restTagsRepo.InsertInto(tag); err != nil {
		return err
	}

	return nil
}

func (rtUC *RestTagsUseCase) GetTagByID(id uint64) (*models.RestTag, error) {
	tag, err := rtUC.restTagsRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.RestTagNotFound
		}

		return nil, err
	}

	return tag, nil
}

func (rtUC *RestTagsUseCase) GetAllTags() ([]*models.RestTag, error) {
	tags, err := rtUC.restTagsRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (rtUC *RestTagsUseCase) UpdateTag(tag *models.RestTag) error {
	t, err := rtUC.restTagsRepo.GetByID(tag.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tools.RestTagNotFound
		}

		return err
	}

	if tag.Image == "" {
		tag.Image = t.Image
	}

	if tag.Name == "" {
		tag.Name = t.Name
		if t.Image != "" {
			rootDir, _ := os.Getwd()
			if err := os.Remove(filepath.Join(rootDir, tools.RestTagsImagesPath, t.Image)); err != nil {
				return err
			}
		}
	}

	if err := rtUC.restTagsRepo.Update(tag); err != nil {
		return err
	}

	return nil
}

func (rtUC *RestTagsUseCase) DeleteTag(id uint64) error {
	t, err := rtUC.restTagsRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return tools.RestTagNotFound
		}

		return err
	}

	if t.Image != "" {
		rootDir, _ := os.Getwd()
		if err := os.Remove(filepath.Join(rootDir, tools.RestTagsImagesPath, t.Image)); err != nil {
			return err
		}
	}

	if err := rtUC.restTagsRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
