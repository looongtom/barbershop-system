package db

import (
	"DoAn/pkg/criteria"
	entity2 "DoAn/pkg/criteria/entity"
	"DoAn/pkg/servicing/entity"
	"context"
	"database/sql"
	"github.com/go-kit/kit/log"
	"time"
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepository(db *sql.DB, logger log.Logger) (criteria.CriteriaRepository, error) {
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}

func (r repo) CreateCategory(ctx context.Context, category string) (*entity.Category, error) {
	ts := time.Now()
	query := `INSERT INTO category(name,created_at,updated_at) VALUES($1,$2,$3) RETURNING id,name,created_at,updated_at`
	var resp entity.Category
	err := r.db.QueryRowContext(ctx, query, category, ts.Unix(), ts.Unix()).Scan(&resp.ID, &resp.Name, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) CreateCriteria(ctx context.Context, criteria entity2.Criteria) (*entity2.Criteria, error) {
	ts := time.Now()
	query := `INSERT INTO criteria(name,img,category_id,created_at,updated_at) VALUES($1,$2,$3,$4,$5) RETURNING id,name,img,category_id,created_at,updated_at`
	var resp entity2.Criteria
	err := r.db.QueryRowContext(ctx, query, criteria.Name, criteria.Img, criteria.CategoryId, ts.Unix(), ts.Unix()).
		Scan(&resp.ID, &resp.Name, &resp.Img, &resp.CategoryId, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) UpdateCriteria(ctx context.Context, criteria entity2.Criteria) (*entity2.Criteria, error) {
	ts := time.Now()
	query := `UPDATE criteria SET name=$1,img=$2,category_id=$3,updated_at=$4 WHERE id=$5 RETURNING id,name,img,category_id,created_at,updated_at`
	var resp entity2.Criteria
	err := r.db.QueryRowContext(ctx, query, criteria.Name, criteria.Img, criteria.CategoryId, ts.Unix(), criteria.ID).
		Scan(&resp.ID, &resp.Name, &resp.Img, &resp.CategoryId, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) UpdateCategory(ctx context.Context, category entity.Category) (*entity.Category, error) {
	ts := time.Now()
	query := `UPDATE category SET name=$1,updated_at=$2 WHERE id=$3 RETURNING id,name,created_at,updated_at`
	var resp entity.Category
	err := r.db.QueryRowContext(ctx, query, category.Name, ts.Unix(), category.ID).Scan(&resp.ID, &resp.Name, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) GetCriteria(ctx context.Context, id string) (*entity2.Criteria, error) {
	query := `SELECT id,name,img,category_id,created_at,updated_at FROM criteria WHERE id=$1`
	var resp entity2.Criteria
	err := r.db.QueryRowContext(ctx, query, id).Scan(&resp.ID, &resp.Name, &resp.Img, &resp.CategoryId, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) GetCategory(ctx context.Context, id string) (*entity.Category, error) {
	query := `SELECT id,name,created_at,updated_at FROM category WHERE id=$1`
	var resp entity.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&resp.ID, &resp.Name, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) GetListCategory(ctx context.Context) ([]entity.Category, error) {
	query := `SELECT id,name,created_at,updated_at FROM category`
	var list []entity.Category
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var resp entity.Category
		err = rows.Scan(&resp.ID, &resp.Name, &resp.CreatedAt, &resp.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, resp)
	}
	return list, nil
}

func (r repo) FindCriteria(ctx context.Context, name string, category int) ([]entity2.Criteria, error) {
	query := `SELECT id,name,img,category_id,created_at,updated_at FROM criteria WHERE name=$1 OR category_id=$2`
	var list []entity2.Criteria
	rows, err := r.db.QueryContext(ctx, query, name, category)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var resp entity2.Criteria
		err = rows.Scan(&resp.ID, &resp.Name, &resp.Img, &resp.CategoryId, &resp.CreatedAt, &resp.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, resp)
	}
	return list, nil
}

func (r repo) GetListCriteria(ctx context.Context) ([]entity2.Criteria, error) {
	query := `SELECT id,name,img,category_id,created_at,updated_at FROM criteria`
	var list []entity2.Criteria
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var resp entity2.Criteria
		err = rows.Scan(&resp.ID, &resp.Name, &resp.Img, &resp.CategoryId, &resp.CreatedAt, &resp.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, resp)
	}
	return list, nil
}
