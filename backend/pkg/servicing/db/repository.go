package db

import (
	"DoAn/pkg/servicing"
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

func NewRepository(db *sql.DB, logger log.Logger) (servicing.ServicingRepository, error) {
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}

func (r repo) GetService(ctx context.Context, id string) (*entity.Servicing, error) {
	var service entity.Servicing
	err := r.db.QueryRow("SELECT * FROM servicing WHERE id = $1", id).Scan(&service.ID, &service.Name, &service.Price, &service.Description, &service.Url, &service.CategoryID, &service.CreatedAt, &service.UpdatedAt)
	if err != nil {
		r.logger.Log("error while fetching data")
		return nil, err
	}
	return &service, nil
}
func (r repo) GetCategory(ctx context.Context, id string) (*entity.Servicing, error) {
	var service entity.Servicing
	err := r.db.QueryRow("SELECT * FROM category WHERE id = $1", id).Scan(&service.ID, &service.Name, &service.CreatedAt, &service.UpdatedAt)
	if err != nil {
		r.logger.Log("error while fetching data")
		return nil, err
	}
	return &service, nil
}
func (r repo) GetListCategory(ctx context.Context) ([]entity.Category, error) {
	var listCate []entity.Category
	rows, err := r.db.Query("SELECT * FROM category")
	if err != nil {
		r.logger.Log("error while fetching data")
		return nil, err
	}
	for rows.Next() {
		var cate entity.Category
		err := rows.Scan(&cate.ID, &cate.Name, &cate.CreatedAt, &cate.UpdatedAt)
		if err != nil {
			r.logger.Log("error while fetching data")
			return nil, err
		}
		listCate = append(listCate, cate)
	}
	return listCate, nil
}

func (r repo) GetListService(ctx context.Context) ([]entity.Servicing, error) {
	var listService []entity.Servicing
	rows, err := r.db.Query("SELECT * FROM servicing")
	if err != nil {
		r.logger.Log("error while fetching data")
		return nil, err
	}
	for rows.Next() {
		var cate entity.Servicing
		err := rows.Scan(&cate.ID, &cate.Name, &cate.Price, &cate.Description, &cate.Url, &cate.CategoryID, &cate.CreatedAt, &cate.UpdatedAt)
		if err != nil {
			r.logger.Log("error while fetching data")
			return nil, err
		}
		listService = append(listService, cate)
	}
	return listService, nil
}

func (r repo) UpdateService(ctx context.Context, service entity.Servicing) (*entity.Servicing, error) {
	query := `UPDATE servicing SET name = $1, price = $2, description = $3, url = $4, category_id = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.Exec(query, service.Name, service.Price, service.Description, service.Url, service.CategoryID, time.Now().Unix(), service.ID)
	if err != nil {
		r.logger.Log("error while updating data")
		return nil, err
	}
	return &service, nil
}
func (r repo) UpdateCategory(ctx context.Context, category entity.Category) (*entity.Category, error) {
	query := `UPDATE category SET name = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, category.Name, time.Now().Unix(), category.ID)
	if err != nil {
		r.logger.Log("error while updating data")
		return nil, err
	}
	return &category, nil
}
func (r repo) CreateService(ctx context.Context, service entity.Servicing) error {
	createTable := `CREATE TABLE IF NOT EXISTS servicing(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price int NOT NULL,
    description VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    category_id int NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL);
    `
	_, err := r.db.Exec(createTable)
	if err != nil {
		r.logger.Log("error while creating table")
		return err
	}
	_, err = r.db.Exec("INSERT INTO servicing(name, price, description, url, category_id, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7)", service.Name, service.Price, service.Description, service.Url, service.CategoryID, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		r.logger.Log("error while inserting data")
		return err
	}
	return nil
}

func (r repo) CreateCategory(ctx context.Context, category string) (*entity.Category, error) {
	createTable := `CREATE TABLE IF NOT EXISTS category(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL);
    `
	_, err := r.db.Exec(createTable)
	if err != nil {
		r.logger.Log("error while creating table")
		return nil, err
	}
	newCategory := entity.Category{
		Name:      category,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	_, err = r.db.Exec("INSERT INTO category(name, created_at, updated_at) VALUES($1, $2, $3)", category, newCategory.CreatedAt, newCategory.UpdatedAt)
	if err != nil {
		r.logger.Log("error while inserting data")
		return nil, err
	}
	return &newCategory, nil
}
