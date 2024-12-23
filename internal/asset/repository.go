package asset

import (
	"context"
	"github.com/safayildirim/asset-management-service/internal/asset/entity"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Deposit(ctx context.Context, tx *gorm.DB, entity *entity.Asset) (*entity.Asset, error)
	Withdraw(ctx context.Context, tx *gorm.DB, entity *entity.Asset) (*entity.Asset, error)
	GetAsset(ctx context.Context, filters entity.Filters) ([]*entity.Asset, error)
	CreateAsset(ctx context.Context, tx *gorm.DB, item *entity.Asset) (*entity.Asset, error)
	UpdateAsset(ctx context.Context, tx *gorm.DB, item *entity.Asset) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Deposit(ctx context.Context, tx *gorm.DB, entity *entity.Asset) (*entity.Asset, error) {
	db := tx
	if db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *repository) Withdraw(ctx context.Context, tx *gorm.DB, entity *entity.Asset) (*entity.Asset, error) {
	db := tx
	if db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *repository) CreateAsset(ctx context.Context, tx *gorm.DB, item *entity.Asset) (*entity.Asset, error) {
	db := tx
	if db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Create(item).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, ErrDuplicateAsset
		}

		return nil, err
	}

	return item, nil
}

func (r *repository) GetAsset(ctx context.Context, filters entity.Filters) ([]*entity.Asset, error) {
	var assets []*entity.Asset

	query := r.db.Model(&entity.Asset{})

	if len(filters.ID) > 0 {
		query = query.Where("id IN ?", filters.ID)
	}
	if len(filters.Name) > 0 {
		query = query.Where("name IN ?", filters.Name)
	}
	if len(filters.WalletID) > 0 {
		query = query.Where("wallet_id IN ?", filters.WalletID)
	}

	err := query.Find(&assets).Error
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *repository) UpdateAsset(ctx context.Context, tx *gorm.DB, item *entity.Asset) error {
	db := tx
	if db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Save(item).Error
	if err != nil {
		return err
	}

	return nil
}
