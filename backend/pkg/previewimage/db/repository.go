package db

import (
	"DoAn/entity"
	"DoAn/pkg/previewimage"
	"context"
	"database/sql"
	"github.com/go-kit/kit/log"
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepository(db *sql.DB, logger log.Logger) (previewimage.PreviewImageRepository, error) {
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}

func (r repo) CreatePreviewImage(ctx context.Context, previewImage entity.PreviewImage) (*entity.PreviewImage, error) {
	query := `INSERT INTO preview_image (account_id, image_url, created_at) VALUES ($1, $2, $3) RETURNING id,account_id, image_url, created_at`

	var resp entity.PreviewImage
	err := r.db.QueryRowContext(ctx, query, previewImage.AccountId, previewImage.Url, previewImage.CreatedAt).Scan(&resp.ID, &resp.AccountId, &resp.Url, &resp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) GetPreviewImage(ctx context.Context, id int) (*entity.PreviewImage, error) {
	query := `SELECT id, account_id, image_url, created_at FROM preview_image WHERE id=$1`
	var resp entity.PreviewImage
	err := r.db.QueryRowContext(ctx, query, id).Scan(&resp.ID, &resp.AccountId, &resp.Url, &resp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) GetListPreviewImageByAccountId(ctx context.Context, accountId int) ([]entity.PreviewImage, error) {
	query := `SELECT id, account_id, image_url, created_at FROM preview_image WHERE account_id=$1`
	var list []entity.PreviewImage
	rows, err := r.db.QueryContext(ctx, query, accountId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var resp entity.PreviewImage
		err = rows.Scan(&resp.ID, &resp.AccountId, &resp.Url, &resp.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, resp)
	}
	return list, nil
}
