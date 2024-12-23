package transaction

import (
	"context"
	"github.com/safayildirim/asset-management-service/internal/transaction/entity"
	"gorm.io/gorm"
)

type Repository interface {
	CreateTransaction(ctx context.Context, tx *gorm.DB, entity *entity.Transaction) (*entity.Transaction, error)
	GetTransactions(ctx context.Context, filters entity.Filters) ([]*entity.Transaction, error)
	DeleteTransaction(ctx context.Context, tx *gorm.DB, id uint) error
	UpdateTransaction(ctx context.Context, tx *gorm.DB, item *entity.Transaction) error
	InTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateTransaction(ctx context.Context, tx *gorm.DB,
	entity *entity.Transaction) (*entity.Transaction, error) {
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

func (r *repository) GetTransactions(ctx context.Context, filters entity.Filters) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction

	query := r.db.WithContext(ctx).Model(&entity.Transaction{})

	if len(filters.ID) > 0 {
		query = query.Where("id IN ?", filters.ID)
	}
	if len(filters.SourceWalletID) > 0 {
		query = query.Where("source_wallet_id IN ?", filters.SourceWalletID)
	}
	if len(filters.DestinationWalletID) > 0 {
		query = query.Where("destination_wallet_id IN ?", filters.DestinationWalletID)
	}
	if len(filters.Status) > 0 {
		query = query.Where("status IN ?", filters.Status)
	}
	if !filters.ScheduledStart.IsZero() {
		query = query.Where("scheduled_at >= ?", filters.ScheduledStart)
	}
	if !filters.ScheduledEnd.IsZero() {
		query = query.Where("scheduled_at <= ?", filters.ScheduledEnd)
	}

	err := query.Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *repository) DeleteTransaction(ctx context.Context, tx *gorm.DB, id uint) error {
	db := tx
	if db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).Delete(&entity.Transaction{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateTransaction(ctx context.Context, tx *gorm.DB, item *entity.Transaction) error {
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

func (r *repository) InTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.db.WithContext(ctx).Begin() // Start a transaction
	if tx.Error != nil {
		return tx.Error
	}

	// Execute the transactional logic
	if err := fn(tx); err != nil {
		tx.Rollback() // Rollback on error
		return err
	}

	// Commit if everything is successful
	return tx.Commit().Error
}
