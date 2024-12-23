package transaction

import (
	"context"
	"github.com/pkg/errors"
	"github.com/safayildirim/asset-management-service/internal/asset/entity"
	assetmock "github.com/safayildirim/asset-management-service/internal/asset/mock"
	transactionentity "github.com/safayildirim/asset-management-service/internal/transaction/entity"
	transactionmock "github.com/safayildirim/asset-management-service/internal/transaction/mock"
	"github.com/safayildirim/asset-management-service/internal/transaction/request"
	walletentity "github.com/safayildirim/asset-management-service/pkg/client/wallet/entity"
	walletmock "github.com/safayildirim/asset-management-service/pkg/client/wallet/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_ScheduleTransaction(t *testing.T) {
	tests := []struct {
		name                     string
		request                  *request.ScheduleTransactionRequest
		mockSourceWallet         bool
		mockSourceWalletError    error
		mockSourceWalletResponse *walletentity.Wallet
		mockDestWallet           bool
		mockDestWalletError      error
		mockDestWalletResponse   *walletentity.Wallet
		mockAsset                bool
		mockAssetsErr            error
		mockAssetsResponse       []*entity.Asset
		mockTransaction          bool
		mockTransactionErr       error
		mockTransactionResponse  *transactionentity.Transaction
		expectedResult           *transactionentity.Transaction
		expectedError            error
	}{
		{
			name: "when everything is ok then should return transaction",
			request: &request.ScheduleTransactionRequest{
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              10.0,
			},
			mockSourceWallet:         true,
			mockSourceWalletResponse: &walletentity.Wallet{ID: 1},
			mockDestWallet:           true,
			mockDestWalletResponse:   &walletentity.Wallet{ID: 2},
			mockAsset:                true,
			mockAssetsResponse: []*entity.Asset{
				{ID: 1, WalletID: 1, Name: "BTC", Amount: 20.0},
				{ID: 2, WalletID: 2, Name: "BTC", Amount: 0},
			},
			mockAssetsErr:   nil,
			mockTransaction: true,
			mockTransactionResponse: &transactionentity.Transaction{
				ID:                  1,
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              10.0,
				Status:              transactionentity.TransactionPending,
			},
			mockTransactionErr: nil,
			expectedResult: &transactionentity.Transaction{
				ID:                  1,
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              10.0,
				Status:              transactionentity.TransactionPending,
			},
			expectedError: nil,
		},
		{
			name: "when source wallet not found then should return error",
			request: &request.ScheduleTransactionRequest{
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              10.0,
			},
			mockSourceWallet:      true,
			mockSourceWalletError: errors.New("wallet not found"),
			expectedError:         errors.New("wallet not found"),
		},
		{
			name: "when destination wallet not found then should return error",
			request: &request.ScheduleTransactionRequest{
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              10.0,
			},
			mockSourceWallet:         true,
			mockSourceWalletResponse: &walletentity.Wallet{ID: 1},
			mockDestWallet:           true,
			mockDestWalletError:      errors.New("wallet not found"),
			expectedError:            errors.New("wallet not found"),
		},
		{
			name: "when asset not found then should return error",
			request: &request.ScheduleTransactionRequest{
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              10.0,
			},
			mockSourceWallet:         true,
			mockSourceWalletResponse: &walletentity.Wallet{ID: 1},
			mockDestWallet:           true,
			mockDestWalletResponse:   &walletentity.Wallet{ID: 2},
			mockAsset:                true,
			mockAssetsErr:            errors.New("asset not found"),
			expectedError:            errors.New("asset not found"),
		},
		{
			name: "when destination asset not found then should return error",
			request: &request.ScheduleTransactionRequest{
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              10.0,
			},
			mockSourceWallet:         true,
			mockSourceWalletResponse: &walletentity.Wallet{ID: 1},
			mockDestWallet:           true,
			mockDestWalletResponse:   &walletentity.Wallet{ID: 2},
			mockAsset:                true,
			mockAssetsResponse: []*entity.Asset{
				{ID: 1, WalletID: 1, Name: "BTC", Amount: 20.0},
			},
			expectedError: errors.New("asset not found"),
		},
		{
			name: "when balance is insufficient then should return error",
			request: &request.ScheduleTransactionRequest{
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              50.0,
			},
			mockSourceWallet:         true,
			mockSourceWalletResponse: &walletentity.Wallet{ID: 1},
			mockDestWallet:           true,
			mockDestWalletResponse:   &walletentity.Wallet{ID: 2},
			mockAsset:                true,
			mockAssetsResponse: []*entity.Asset{
				{ID: 1, WalletID: 1, Name: "BTC", Amount: 10.0},
				{ID: 2, WalletID: 2, Name: "BTC", Amount: 0},
			},
			mockTransactionErr: nil,
			expectedResult:     nil,
			expectedError:      ErrInsufficientBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAssetRepo := assetmock.NewMockAssetRepository(t)
			mockTransactionRepo := transactionmock.NewMockTransactionRepository(t)
			mockWalletClient := walletmock.NewMockWalletClient(t)
			s := NewService(mockAssetRepo, mockTransactionRepo, mockWalletClient)

			if tt.mockSourceWallet {
				mockWalletClient.EXPECT().GetWallet(mock.Anything, tt.request.SourceWalletID).
					Return(tt.mockSourceWalletResponse, tt.mockSourceWalletError).Once()
			}

			if tt.mockDestWallet {
				mockWalletClient.EXPECT().GetWallet(mock.Anything, tt.request.DestinationWalletID).
					Return(tt.mockDestWalletResponse, tt.mockDestWalletError).Once()
			}

			if tt.mockAsset {
				mockAssetRepo.EXPECT().GetAsset(mock.Anything, mock.Anything).
					Return(tt.mockAssetsResponse, tt.mockAssetsErr).Once()
			}

			if tt.mockTransaction {
				mockTransactionRepo.EXPECT().CreateTransaction(mock.Anything, mock.Anything,
					mock.Anything).Return(tt.mockTransactionResponse, tt.mockTransactionErr).Once()
			}

			result, err := s.ScheduleTransaction(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestService_GetTransactions(t *testing.T) {
	tests := []struct {
		name           string
		request        *request.GetTransactionsParams
		mockFilters    transactionentity.Filters
		mockService    bool
		mockReturn     []*transactionentity.Transaction
		mockError      error
		expectedResult []*transactionentity.Transaction
		expectedError  error
	}{
		{
			name:    "when no filters then should return all transactions",
			request: &request.GetTransactionsParams{},
			mockFilters: transactionentity.Filters{
				ID:                  nil,
				SourceWalletID:      nil,
				DestinationWalletID: nil,
				Status:              nil,
			},
			mockService: true,
			mockReturn: []*transactionentity.Transaction{
				{ID: 1, SourceWalletID: 1001, DestinationWalletID: 1002, Status: "completed"},
				{ID: 2, SourceWalletID: 1001, DestinationWalletID: 1003, Status: "pending"},
			},
			mockError: nil,
			expectedResult: []*transactionentity.Transaction{
				{ID: 1, SourceWalletID: 1001, DestinationWalletID: 1002, Status: "completed"},
				{ID: 2, SourceWalletID: 1001, DestinationWalletID: 1003, Status: "pending"},
			},
			expectedError: nil,
		},
		{
			name: "when filters then should return filtered transactions",
			request: &request.GetTransactionsParams{
				SourceWalletID: []uint{1001},
				Status:         []string{"pending"},
			},
			mockFilters: transactionentity.Filters{
				ID:                  nil,
				SourceWalletID:      []uint{1001},
				DestinationWalletID: nil,
				Status:              []string{"pending"},
			},
			mockService: true,
			mockReturn: []*transactionentity.Transaction{
				{ID: 2, SourceWalletID: 1001, DestinationWalletID: 1003, Status: "pending"},
			},
			mockError: nil,
			expectedResult: []*transactionentity.Transaction{
				{ID: 2, SourceWalletID: 1001, DestinationWalletID: 1003, Status: "pending"},
			},
			expectedError: nil,
		},
		{
			name: "when repository error then should return error",
			request: &request.GetTransactionsParams{
				SourceWalletID: []uint{1001},
			},
			mockFilters: transactionentity.Filters{
				ID:                  nil,
				SourceWalletID:      []uint{1001},
				DestinationWalletID: nil,
				Status:              nil,
			},
			mockService:    true,
			mockReturn:     nil,
			mockError:      errors.New("repository error"),
			expectedResult: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAssetRepo := assetmock.NewMockAssetRepository(t)
			mockTransactionRepo := transactionmock.NewMockTransactionRepository(t)
			mockWalletClient := walletmock.NewMockWalletClient(t)
			s := NewService(mockAssetRepo, mockTransactionRepo, mockWalletClient)

			if tt.mockService {
				mockTransactionRepo.EXPECT().GetTransactions(mock.Anything, tt.mockFilters).
					Return(tt.mockReturn, tt.mockError).Once()
			}

			result, err := s.GetTransactions(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestService_CancelTransaction(t *testing.T) {
	tests := []struct {
		name                       string
		transactionID              uint
		mockGetTransaction         bool
		mockGetTransactionError    error
		mockGetTransactionReturn   []*transactionentity.Transaction
		mockUpdateTransaction      bool
		mockUpdateTransactionError error
		expectedError              error
	}{
		{
			name:               "when transaction is pending then should cancel",
			transactionID:      1,
			mockGetTransaction: true,
			mockGetTransactionReturn: []*transactionentity.Transaction{
				{ID: 1, Status: transactionentity.TransactionPending},
			},
			mockUpdateTransaction: true,
			expectedError:         nil,
		},
		{
			name:                     "when transaction is not found then should return error",
			transactionID:            2,
			mockGetTransaction:       true,
			mockGetTransactionReturn: []*transactionentity.Transaction{},
			expectedError:            ErrTransactionNotFound,
		},
		{
			name:               "when transaction is not pending then should return error",
			transactionID:      3,
			mockGetTransaction: true,
			mockGetTransactionReturn: []*transactionentity.Transaction{
				{ID: 1, Status: transactionentity.TransactionCompleted},
			},
			expectedError: ErrTransactionCannotBeDeleted,
		},
		{
			name:                    "when repository error on get then should return error",
			transactionID:           4,
			mockGetTransaction:      true,
			mockGetTransactionError: errors.New("transaction not found"),
			expectedError:           errors.New("transaction not found"),
		},
		{
			name:               "when repository error on update then should return error",
			transactionID:      5,
			mockGetTransaction: true,
			mockGetTransactionReturn: []*transactionentity.Transaction{
				{ID: 1, Status: transactionentity.TransactionPending},
			},
			mockUpdateTransaction:      true,
			mockUpdateTransactionError: errors.New("update error"),
			expectedError:              errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAssetRepo := assetmock.NewMockAssetRepository(t)
			mockTransactionRepo := transactionmock.NewMockTransactionRepository(t)
			mockWalletClient := walletmock.NewMockWalletClient(t)
			s := NewService(mockAssetRepo, mockTransactionRepo, mockWalletClient)

			if tt.mockGetTransaction {
				mockTransactionRepo.EXPECT().
					GetTransactions(mock.Anything, transactionentity.Filters{ID: []uint{tt.transactionID}}).
					Return(tt.mockGetTransactionReturn, tt.mockGetTransactionError).
					Once()
			}

			if tt.mockUpdateTransaction {
				mockTransactionRepo.EXPECT().UpdateTransaction(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockUpdateTransactionError).Once()
			}

			err := s.CancelTransaction(context.Background(), tt.transactionID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
