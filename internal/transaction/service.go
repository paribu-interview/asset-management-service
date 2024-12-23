package transaction

import (
	"context"
	"github.com/safayildirim/asset-management-service/internal/asset"
	"github.com/safayildirim/asset-management-service/internal/asset/entity"
	transactionentity "github.com/safayildirim/asset-management-service/internal/transaction/entity"
	"github.com/safayildirim/asset-management-service/internal/transaction/request"
	"github.com/safayildirim/asset-management-service/pkg/client/wallet"
)

type Service interface {
	ScheduleTransaction(ctx context.Context,
		request *request.ScheduleTransactionRequest) (*transactionentity.Transaction, error)
	GetTransactions(ctx context.Context,
		request *request.GetTransactionsParams) ([]*transactionentity.Transaction, error)
	CancelTransaction(ctx context.Context, id uint) error
}

type service struct {
	assetRepository       asset.Repository
	transactionRepository Repository
	walletClient          wallet.Client
}

func NewService(assetRepository asset.Repository, transactionRepository Repository,
	walletClient wallet.Client) Service {
	return &service{assetRepository: assetRepository, transactionRepository: transactionRepository,
		walletClient: walletClient}
}

// ScheduleTransaction schedules a transaction between two wallets for a specific asset.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - request: A request object containing details for the scheduled transaction, including:
//   - SourceWalletID: The ID of the wallet sending the asset.
//   - DestinationWalletID: The ID of the wallet receiving the asset.
//   - AssetName: The name of the asset to be transferred.
//   - Amount: The amount of the asset to transfer.
//   - ScheduledAt: The scheduled time for the transaction.
//
// Returns:
//   - A pointer to the newly created transaction entity.
//   - An error if any validation or persistence step fails.
//
// Errors:
//   - ErrAssetNotFound: If the asset is not found for either the source or destination wallet.
//   - ErrInsufficientBalance: If the source wallet does not have enough balance for the transaction.
//   - Any other error encountered during wallet or asset retrieval, or transaction persistence.
func (s *service) ScheduleTransaction(ctx context.Context,
	request *request.ScheduleTransactionRequest) (*transactionentity.Transaction, error) {
	// Validate that the source wallet exists by fetching it from the wallet client
	_, err := s.walletClient.GetWallet(ctx, request.SourceWalletID)
	if err != nil {
		return nil, err
	}

	// Validate that the destination wallet exists by fetching it from the wallet client
	_, err = s.walletClient.GetWallet(ctx, request.DestinationWalletID)
	if err != nil {
		return nil, err
	}

	// Fetch the assets for both source and destination wallets with the specified asset name
	assets, err := s.assetRepository.GetAsset(ctx, entity.Filters{
		Name:     []string{request.AssetName},
		WalletID: []uint{request.SourceWalletID, request.DestinationWalletID},
	})
	if err != nil {
		return nil, err
	}

	// Ensure that assets exist for both wallets
	if len(assets) < 2 {
		return nil, ErrAssetNotFound
	}

	// Find the asset associated with the source wallet
	var sourceAsset *entity.Asset
	for _, a := range assets {
		if a.WalletID == request.SourceWalletID {
			sourceAsset = a
		}
	}

	// Check if the source wallet has sufficient balance for the transaction
	if sourceAsset.Amount < request.Amount {
		return nil, ErrInsufficientBalance
	}

	// Create a transaction object with the provided details and set its status to pending
	transaction := &transactionentity.Transaction{
		SourceWalletID:      request.SourceWalletID,
		DestinationWalletID: request.DestinationWalletID,
		Amount:              request.Amount,
		AssetName:           request.AssetName,
		Status:              transactionentity.TransactionPending,
		ScheduledAt:         request.ScheduledAt,
	}

	// Create a transaction object with the provided details and set its status to pending
	return s.transactionRepository.CreateTransaction(ctx, nil, transaction)
}

// GetTransactions retrieves a list of transactions based on the provided filters.
//
// This method constructs a filter object from the request parameters and delegates
// the query to the transaction repository.
//
// Parameters:
//   - ctx: Context for managing request lifecycle and cancellation.
//   - request: Request object containing the filter criteria for fetching transactions, including:
//   - ID: A list of transaction IDs to filter by.
//   - SourceWalletID: A list of source wallet IDs to filter by.
//   - DestinationWalletID: A list of destination wallet IDs to filter by.
//   - Status: A list of transaction statuses to filter by.
//
// Returns:
//   - A slice of transactions that match the filter criteria.
//   - An error if the repository query fails.
func (s *service) GetTransactions(ctx context.Context,
	request *request.GetTransactionsParams) ([]*transactionentity.Transaction, error) {
	filters := transactionentity.Filters{
		ID:                  request.ID,
		SourceWalletID:      request.SourceWalletID,
		DestinationWalletID: request.DestinationWalletID,
		Status:              request.Status,
	}
	return s.transactionRepository.GetTransactions(ctx, filters)
}

// CancelTransaction cancels a transaction with the given ID.
//
// Parameters:
//   - ctx: The context for managing request lifecycle and cancellation.
//   - id: The ID of the transaction to be cancelled.
//
// Returns:
//   - An error if the transaction cannot be fetched, validated, or updated.
//
// Errors:
//   - ErrTransactionNotFound: If the transaction with the given ID does not exist.
//   - ErrTransactionCannotBeDeleted: If the transaction is not in a "Pending" state.
func (s *service) CancelTransaction(ctx context.Context, id uint) error {
	// Fetch the transaction by ID from the transaction repository
	transaction, err := s.transactionRepository.GetTransactions(ctx, transactionentity.Filters{ID: []uint{id}})
	if err != nil {
		return err
	}

	// Check if the transaction exists
	if len(transaction) == 0 {
		return ErrTransactionNotFound
	}

	// Ensure that the transaction is in a pending state before cancellation
	if transaction[0].Status != transactionentity.TransactionPending {
		return ErrTransactionCannotBeDeleted
	}

	// Update the transaction status to "Cancelled"
	transaction[0].Status = transactionentity.TransactionCancelled

	// Persist the updated transaction to the database
	return s.transactionRepository.UpdateTransaction(ctx, nil, transaction[0])
}
