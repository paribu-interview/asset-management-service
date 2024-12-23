package asset

import (
	"context"
	"github.com/pkg/errors"
	"github.com/safayildirim/asset-management-service/internal/asset/entity"
	"github.com/safayildirim/asset-management-service/internal/asset/request"
	"github.com/safayildirim/asset-management-service/pkg/client/wallet"
	"gorm.io/gorm"
)

type Service interface {
	CreateAsset(ctx context.Context, tx *gorm.DB, request *request.CreateAssetRequest) (*entity.Asset, error)
	GetAssets(ctx context.Context, request *request.GetAssetsParams) ([]*entity.Asset, error)
	Deposit(ctx context.Context, tx *gorm.DB, request *request.CreateDepositRequest) (*entity.Asset, error)
	Withdraw(ctx context.Context, tx *gorm.DB, request *request.CreateWithdrawRequest) (*entity.Asset, error)
}

type service struct {
	assetRepository Repository
	walletClient    wallet.Client
}

func NewService(assetRepository Repository, walletClient wallet.Client) Service {
	return &service{assetRepository: assetRepository, walletClient: walletClient}
}

func (s *service) CreateAsset(ctx context.Context, tx *gorm.DB, request *request.CreateAssetRequest) (*entity.Asset,
	error) {
	item := entity.Asset{
		WalletID: request.WalletID,
		Name:     request.Name,
		Amount:   request.Amount,
	}
	return s.assetRepository.CreateAsset(ctx, tx, &item)
}

func (s *service) GetAssets(ctx context.Context, request *request.GetAssetsParams) ([]*entity.Asset, error) {
	filters := entity.Filters{
		ID:       request.ID,
		Name:     request.Name,
		WalletID: request.WalletID,
	}
	return s.assetRepository.GetAsset(ctx, filters)
}

// Deposit adds the specified amount of an asset to the wallet.
//
// Parameters:
// - ctx: Context for managing request lifecycle and cancellation.
// - tx: Optional database transaction for atomic operations.
// - request: Request object containing the deposit details, including:
//   - WalletID: The ID of the wallet to deposit into.
//   - Name: The name of the asset being deposited.
//   - Amount: The amount to deposit.
//
// Returns:
// - The updated asset entity after the deposit.
// - An error if any validation or persistence step fails.
//
// Errors:
// - Returns an error if the wallet does not exist, or if asset retrieval or update fails.
func (s *service) Deposit(ctx context.Context, tx *gorm.DB, request *request.CreateDepositRequest) (*entity.Asset,
	error) {
	// Verify that the wallet exists using the wallet client
	w, err := s.walletClient.GetWallet(ctx, request.WalletID)
	if err != nil {
		return nil, err
	}

	var assetEntity *entity.Asset

	// Fetch existing assets for the specified wallet and asset name
	assets, err := s.assetRepository.GetAsset(ctx, entity.Filters{
		Name:     []string{request.Name},
		WalletID: []uint{request.WalletID},
	})
	if err != nil {
		return nil, err
	}

	// Check if the asset exists for the wallet
	if len(assets) == 0 {
		assetEntity, err = s.assetRepository.CreateAsset(ctx, tx, &entity.Asset{
			WalletID: w.ID,
			Name:     request.Name,
		})
		if err != nil {
			return nil, err
		}
	} else {
		assetEntity = assets[0]
	}

	// Increase the asset amount by the specified deposit value
	assetEntity.Amount += request.Amount

	// Update the asset in the repository
	err = s.assetRepository.UpdateAsset(ctx, tx, assetEntity)
	if err != nil {
		return nil, err
	}

	// Return the updated asset
	return assetEntity, nil
}

// Withdraw deducts the specified amount of an asset from the wallet.
//
// Parameters:
// - ctx: Context for managing request lifecycle and cancellation.
// - tx: Optional database transaction for atomic operations.
// - request: Request object containing the withdrawal details, including:
//   - WalletID: The ID of the wallet to withdraw from.
//   - Name: The name of the asset being withdrawn.
//   - Amount: The amount to withdraw.
//
// Returns:
// - The updated asset entity after the withdrawal.
// - An error if any validation or persistence step fails.
//
// Errors:
// - Returns an error if the wallet does not exist, if asset retrieval or update fails, or if the balance is insufficient.
func (s *service) Withdraw(ctx context.Context, tx *gorm.DB, request *request.CreateWithdrawRequest) (*entity.Asset,
	error) {
	// Verify that the wallet exists using the wallet client
	w, err := s.walletClient.GetWallet(ctx, request.WalletID)
	if err != nil {
		return nil, err
	}

	var assetEntity *entity.Asset

	// Fetch the asset associated with the specified wallet and asset name
	assets, err := s.assetRepository.GetAsset(ctx, entity.Filters{
		Name:     []string{request.Name},
		WalletID: []uint{request.WalletID},
	})
	if err != nil {
		return nil, err
	}

	// Check if the asset exists for the wallet
	if len(assets) == 0 {
		// Create a new asset if it does not exist (with a zero balance)
		assetEntity, err = s.assetRepository.CreateAsset(ctx, tx, &entity.Asset{
			WalletID: w.ID,
			Name:     request.Name,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Use the existing asset
		assetEntity = assets[0]
	}

	// Validate if the wallet has sufficient balance for the withdrawal
	if assetEntity.Amount < request.Amount {
		return nil, errors.New("amount is not enough to withdraw")
	}

	// Deduct the specified amount from the asset's balance
	assetEntity.Amount -= request.Amount

	// Update the asset in the repository
	err = s.assetRepository.UpdateAsset(ctx, tx, assetEntity)
	if err != nil {
		return nil, err
	}

	// Return the updated asset
	return assetEntity, nil
}
