package scheduler

import (
	"context"
	"github.com/safayildirim/asset-management-service/internal/asset"
	"github.com/safayildirim/asset-management-service/internal/asset/request"
	"github.com/safayildirim/asset-management-service/internal/transaction"
	"github.com/safayildirim/asset-management-service/internal/transaction/entity"
	"github.com/safayildirim/asset-management-service/pkg/config"
	"github.com/safayildirim/asset-management-service/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Scheduler struct {
	cfg                   config.SchedulerConfig
	assetService          asset.Service
	transactionRepository transaction.Repository
}

// NewScheduler initializes a new Scheduler instance.
//
// Parameters:
// - cfg: Configuration for the scheduler, including the interval between runs.
// - assetService: Service to handle asset-related operations such as deposits and withdrawals.
// - transactionRepository: Repository to handle transaction-related database operations.
//
// Returns:
// - A pointer to a newly created Scheduler instance.
func NewScheduler(cfg config.SchedulerConfig, assetService asset.Service,
	transactionRepository transaction.Repository) *Scheduler {
	return &Scheduler{cfg: cfg, assetService: assetService, transactionRepository: transactionRepository}
}

// Start runs the scheduler in an infinite loop to process pending transactions.
//
// Parameters:
// - ctx: Context for managing request lifecycle and cancellation.
//
// Notes:
// - This function runs indefinitely, using the interval from the configuration for pauses between runs.
// - Logs errors and skips failed transactions, allowing the scheduler to continue processing others.
func (s *Scheduler) Start(ctx context.Context) {
	log.Logger.Info("scheduler started")

	for {
		// Fetch pending transactions scheduled to run before the current time
		transactions, err := s.transactionRepository.GetTransactions(ctx, entity.Filters{
			Status:       []string{string(entity.TransactionPending)},
			ScheduledEnd: time.Now(),
		})
		if err != nil {
			log.Logger.Error("failed to get transactions", zap.Error(err))

			return
		}

		// Log if no pending transactions are found
		if len(transactions) == 0 {
			log.Logger.Info("no pending transactions")
		}

		// Process each transaction
		for _, t := range transactions {
			// Run the transaction processing in a database transaction
			err = s.transactionRepository.InTransaction(ctx, func(tx *gorm.DB) error {
				// Withdraw the specified amount from the source wallet
				_, err = s.assetService.Withdraw(ctx, tx, &request.CreateWithdrawRequest{
					WalletID: t.SourceWalletID,
					Name:     t.AssetName,
					Amount:   t.Amount,
				})
				if err != nil {
					return err
				}

				// Deposit the specified amount to the destination wallet
				_, err = s.assetService.Deposit(ctx, tx, &request.CreateDepositRequest{
					WalletID: t.DestinationWalletID,
					Name:     t.AssetName,
					Amount:   t.Amount,
				})
				if err != nil {
					return err
				}

				// Update the transaction status to "Completed"
				t.Status = entity.TransactionCompleted
				err = s.transactionRepository.UpdateTransaction(ctx, tx, t)
				if err != nil {
					return err
				}

				return nil // Commit the transaction if all operations succeed
			})
			if err != nil {
				// Log the error for the failed transaction and continue with the next one
				log.Logger.Error("transaction failed", zap.Error(err))
				continue
			}

			// Log the successful completion of the transaction
			log.Logger.Info("transaction completed", zap.Uint("id", t.ID))
		}

		// Pause the scheduler for the configured interval before running again
		time.Sleep(time.Duration(s.cfg.Interval) * time.Second) // Run every 10 seconds
	}
}
