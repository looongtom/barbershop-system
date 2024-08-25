package service

import (
	"DoAn/entity"
	"DoAn/pkg/criteria"
	"DoAn/pkg/criteria/api"
	"DoAn/pkg/criteria/common"
	"context"
	"database/sql"
	"errors"
	"github.com/go-kit/kit/log"
	"strconv"
	"time"
)

type CriteriaStruct struct {
	repository criteria.CriteriaRepository
	logger     log.Logger
}

func NewService(repo criteria.CriteriaRepository, logger log.Logger) criteria.CriteriaService {
	return &CriteriaStruct{
		repository: repo,
		logger:     logger,
	}
}

func (c CriteriaStruct) GetCriteria(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) GetListCriteria(ctx context.Context) (interface{}, error) {
	listCriteria, err := c.repository.GetListCriteria(ctx)
	if err != nil {
		return nil, err
	}
	return listCriteria, nil
}

func (c CriteriaStruct) CreateOrUpdateCriteria(ctx context.Context, criteria api.CreateOrUpdateCriteria) (interface{}, error) {
	ts := time.Now()
	if &criteria.ID != nil && criteria.ID != 0 {
		checkCriteria, err := c.repository.GetCriteria(ctx, strconv.Itoa(criteria.ID))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		updatedCriteria := entity.Criteria{
			ID:         criteria.ID,
			Name:       criteria.Name,
			CategoryId: criteria.CategoryId,
			CreatedAt:  checkCriteria.CreatedAt,
		}

		checkCategory, err := c.repository.GetCategory(ctx, strconv.Itoa(criteria.CategoryId))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("error cannot find category")
		}

		updatedCriteria.CategoryId = checkCategory.ID

		if criteria.Img != nil {
			urlImage, err := common.UploadImageToCloud(criteria.Img)
			if err != nil {
				c.logger.Log("error while uploading image")
				return nil, err
			}
			updatedCriteria.Img = urlImage
		}
		if criteria.Name == "" || &criteria.Name == nil {
			updatedCriteria.Name = checkCriteria.Name
		}
		if criteria.CategoryId == 0 || &criteria.CategoryId == nil {
			updatedCriteria.CategoryId = checkCriteria.CategoryId
		}
		updatedCriteria.UpdatedAt = ts.Unix()
		resp, err := c.repository.UpdateCriteria(ctx, updatedCriteria)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	categoryValue := 0
	checkCategory, err := c.repository.GetCategory(ctx, strconv.Itoa(criteria.CategoryId))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("error cannot find category")
	}

	urlImage, err := common.UploadImageToCloud(criteria.Img)
	if err != nil {
		return nil, err
	}

	categoryValue = checkCategory.ID

	resp, err := c.repository.CreateCriteria(ctx, entity.Criteria{
		Name:       criteria.Name,
		Img:        urlImage,
		CategoryId: categoryValue,
		CreatedAt:  ts.Unix(),
		UpdatedAt:  ts.Unix(),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c CriteriaStruct) DeleteCriteria(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) FindCriteria(ctx context.Context, findCriteriaReq api.FindCriteria) (interface{}, error) {
	list, err := c.repository.FindCriteria(ctx, findCriteriaReq.Name, findCriteriaReq.CategoryId)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c CriteriaStruct) GetCategory(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c CriteriaStruct) GetListCategory(ctx context.Context) (interface{}, error) {
	listCategory, err := c.repository.GetListCategory(ctx)
	if err != nil {
		return nil, err
	}
	return listCategory, nil
}

func (c CriteriaStruct) CreateOrUpdateCategory(ctx context.Context, category api.CreateOrUpdateCategory) (interface{}, error) {
	if &category.ID != nil && category.ID != 0 {
		existestCategory, err := c.repository.GetCategory(ctx, strconv.Itoa(category.ID))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if &existestCategory != nil {
			if category.Name == "" || &category.Name == nil {
				category.Name = existestCategory.Name
			}
			resp, err := c.repository.UpdateCategory(ctx, entity.Category{
				ID:        category.ID,
				Name:      category.Name,
				CreatedAt: existestCategory.CreatedAt,
				UpdatedAt: time.Now().Unix(),
			})
			if err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
	resp, err := c.repository.CreateCategory(ctx, category.Name)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c CriteriaStruct) DeleteCategory(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
