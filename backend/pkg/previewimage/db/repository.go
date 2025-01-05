package db

import (
	"context"
	"database/sql"

	previewimage "DoAn"
	"DoAn/entity"

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

func (r repo) UploadImages(ctx context.Context, previewImage entity.PreviewImage) (*entity.PreviewImage, error) {
	query := `INSERT INTO preview_image (account_id,generated_img, self_img,shape_img,color_img, created_at) VALUES ($1, $2, $3,$4,$5,$6) RETURNING id,account_id,generated_img,self_img,shape_img,color_img, created_at`
	var resp entity.PreviewImage
	err := r.db.QueryRowContext(ctx, query, previewImage.AccountId, previewImage.GeneratedImg, previewImage.SelfImg, previewImage.ShapeImg, previewImage.ColorImg, previewImage.CreatedAt).Scan(&resp.ID, &resp.AccountId, &resp.GeneratedImg, &resp.SelfImg, &resp.ShapeImg, &resp.ColorImg, &resp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) GetPreviewImage(ctx context.Context, id int) (*entity.PreviewImage, error) {
	query := `SELECT id, account_id,generated_img, self_img,shape_img,color_img, created_at FROM preview_image WHERE id=$1`
	var resp entity.PreviewImage
	err := r.db.QueryRowContext(ctx, query, id).Scan(&resp.ID, &resp.AccountId, &resp.GeneratedImg, &resp.SelfImg, &resp.ShapeImg, &resp.ColorImg, &resp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r repo) GetListPreviewImageByAccountId(ctx context.Context, accountId int) ([]entity.PreviewImage, error) {
	query := `SELECT id, account_id,generated_img, self_img,shape_img,color_img, created_at FROM preview_image WHERE account_id=$1 ORDER BY created_at DESC`
	var list []entity.PreviewImage
	rows, err := r.db.QueryContext(ctx, query, accountId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var resp entity.PreviewImage
		err = rows.Scan(&resp.ID, &resp.AccountId, &resp.GeneratedImg, &resp.SelfImg, &resp.ShapeImg, &resp.ColorImg, &resp.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, resp)
	}
	return list, nil
}
