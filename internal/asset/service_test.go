package asset

import (
	"context"
	"github.com/pkg/errors"
	"github.com/safayildirim/asset-management-service/internal/asset/entity"
	assetmock "github.com/safayildirim/asset-management-service/internal/asset/mock"
	"github.com/safayildirim/asset-management-service/internal/asset/request"
	walletentity "github.com/safayildirim/asset-management-service/pkg/client/wallet/entity"
	walletmock "github.com/safayildirim/asset-management-service/pkg/client/wallet/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_CreateAsset(t *testing.T) {
	tests := []struct {
		name           string
		request        *request.CreateAssetRequest
		mockRepo       bool
		mockReturn     *entity.Asset
		mockError      error
		expectedResult *entity.Asset
		expectedError  error
	}{
		{
			name: "when request is valid then should create asset",
			request: &request.CreateAssetRequest{
				WalletID: 1,
				Name:     "BTC",
				Amount:   10.0,
			},
			mockRepo: true,
			mockReturn: &entity.Asset{
				ID:       1,
				WalletID: 1,
				Name:     "BTC",
				Amount:   10.0,
			},
			mockError:      nil,
			expectedResult: &entity.Asset{ID: 1, WalletID: 1, Name: "BTC", Amount: 10.0},
			expectedError:  nil,
		},
		{
			name:     "when repository returns error then should return error",
			mockRepo: true,
			request: &request.CreateAssetRequest{
				WalletID: 1,
				Name:     "BTC",
				Amount:   10.0,
			},
			mockReturn:     nil,
			mockError:      errors.New("repository error"),
			expectedResult: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := assetmock.NewMockAssetRepository(t)
			mockWalletClient := walletmock.NewMockWalletClient(t)
			s := NewService(mockRepository, mockWalletClient)
			if tt.mockRepo {
				mockRepository.EXPECT().CreateAsset(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockError).Once()
			}

			// Call the service method
			result, err := s.CreateAsset(context.Background(), nil, tt.request)

			// Assertions
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

func TestService_GetAssets(t *testing.T) {
	tests := []struct {
		name           string
		request        *request.GetAssetsParams
		mockFilters    entity.Filters
		mockRepo       bool
		mockReturn     []*entity.Asset
		mockError      error
		expectedResult []*entity.Asset
		expectedError  error
	}{
		{
			name:     "when request is valid then should return assets",
			mockRepo: true,
			request: &request.GetAssetsParams{
				ID:       []uint{1, 2},
				Name:     []string{"BTC", "ETH"},
				WalletID: []uint{1001},
			},
			mockFilters: entity.Filters{
				ID:       []uint{1, 2},
				Name:     []string{"BTC", "ETH"},
				WalletID: []uint{1001},
			},
			mockReturn: []*entity.Asset{
				{ID: 1, WalletID: 1001, Name: "BTC", Amount: 10.0},
				{ID: 2, WalletID: 1001, Name: "ETH", Amount: 5.0},
			},
			mockError: nil,
			expectedResult: []*entity.Asset{
				{ID: 1, WalletID: 1001, Name: "BTC", Amount: 10.0},
				{ID: 2, WalletID: 1001, Name: "ETH", Amount: 5.0},
			},
			expectedError: nil,
		},
		{
			name:     "when no assets found then should return empty list",
			mockRepo: true,
			request: &request.GetAssetsParams{
				ID:       []uint{3},
				Name:     []string{"LTC"},
				WalletID: []uint{2001},
			},
			mockFilters: entity.Filters{
				ID:       []uint{3},
				Name:     []string{"LTC"},
				WalletID: []uint{2001},
			},
			mockReturn:     []*entity.Asset{},
			mockError:      nil,
			expectedResult: []*entity.Asset{},
			expectedError:  nil,
		},
		{
			name:     "when repository error then should return error",
			mockRepo: true,
			request: &request.GetAssetsParams{
				ID:       []uint{4},
				Name:     []string{"DOGE"},
				WalletID: []uint{3001},
			},
			mockFilters: entity.Filters{
				ID:       []uint{4},
				Name:     []string{"DOGE"},
				WalletID: []uint{3001},
			},
			mockReturn:     nil,
			mockError:      errors.New("repository error"),
			expectedResult: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := assetmock.NewMockAssetRepository(t)
			mockWalletClient := walletmock.NewMockWalletClient(t)
			s := NewService(mockRepository, mockWalletClient)
			if tt.mockRepo {
				mockRepository.EXPECT().GetAsset(mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockError).Once()
			}

			// Call the service method
			result, err := s.GetAssets(context.Background(), tt.request)

			// Assertions
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

func TestService_Deposit(t *testing.T) {
	tests := []struct {
		name               string
		request            *request.CreateDepositRequest
		mockWallet         *walletentity.Wallet
		mockWalletErr      error
		mockAsset          bool
		mockAssetsResponse []*entity.Asset
		mockAssetsErr      error
		mockCreate         *entity.Asset
		mockCreateErr      error
		mockUpdate         bool
		mockUpdateErr      error
		expectedResult     *entity.Asset
		expectedError      error
	}{
		{
			name: "when request is valid then should deposit amount",
			request: &request.CreateDepositRequest{
				WalletID: 1,
				Name:     "BTC",
				Amount:   10.0,
			},
			mockWallet:         &walletentity.Wallet{ID: 1},
			mockWalletErr:      nil,
			mockAsset:          true,
			mockAssetsResponse: []*entity.Asset{{ID: 1, WalletID: 1, Name: "BTC", Amount: 5.0}},
			mockAssetsErr:      nil,
			mockCreate:         nil,
			mockCreateErr:      nil,
			mockUpdate:         true,
			mockUpdateErr:      nil,
			expectedResult:     &entity.Asset{ID: 1, WalletID: 1, Name: "BTC", Amount: 15.0},
			expectedError:      nil,
		},
		{
			name: "when new asset created then should deposit amount",
			request: &request.CreateDepositRequest{
				WalletID: 2,
				Name:     "ETH",
				Amount:   20.0,
			},
			mockWallet:         &walletentity.Wallet{ID: 2},
			mockWalletErr:      nil,
			mockAsset:          true,
			mockAssetsResponse: []*entity.Asset{},
			mockAssetsErr:      nil,
			mockCreate:         &entity.Asset{ID: 2, WalletID: 2, Name: "ETH", Amount: 0},
			mockCreateErr:      nil,
			mockUpdate:         true,
			mockUpdateErr:      nil,
			expectedResult:     &entity.Asset{ID: 2, WalletID: 2, Name: "ETH", Amount: 20.0},
			expectedError:      nil,
		},
		{
			name: "when wallet not found then should return error",
			request: &request.CreateDepositRequest{
				WalletID: 3,
				Name:     "LTC",
				Amount:   15.0,
			},
			mockWallet:     nil,
			mockWalletErr:  errors.New("wallet not found"),
			mockAssetsErr:  nil,
			mockCreate:     nil,
			mockCreateErr:  nil,
			expectedResult: nil,
			expectedError:  errors.New("wallet not found"),
		},
		{
			name: "when repository error then should return error",
			request: &request.CreateDepositRequest{
				WalletID: 1,
				Name:     "BTC",
				Amount:   10.0,
			},
			mockWallet:     &walletentity.Wallet{ID: 1},
			mockWalletErr:  nil,
			mockAsset:      true,
			mockAssetsErr:  errors.New("repository error"),
			mockCreate:     nil,
			mockCreateErr:  nil,
			mockUpdateErr:  nil,
			expectedResult: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := assetmock.NewMockAssetRepository(t)
			mockWalletClient := walletmock.NewMockWalletClient(t)
			s := NewService(mockRepository, mockWalletClient)

			mockWalletClient.EXPECT().GetWallet(mock.Anything, tt.request.WalletID).
				Return(tt.mockWallet, tt.mockWalletErr).Once()

			if tt.mockAsset {
				mockRepository.EXPECT().GetAsset(mock.Anything, mock.Anything).
					Return(tt.mockAssetsResponse, tt.mockAssetsErr).Once()
			}
			if tt.mockCreate != nil || tt.mockCreateErr != nil {
				mockRepository.EXPECT().CreateAsset(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockCreate, tt.mockCreateErr).Once()
			}
			if tt.mockUpdate {
				mockRepository.EXPECT().UpdateAsset(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockUpdateErr).Once()
			}

			// Call the service method
			result, err := s.Deposit(context.Background(), nil, tt.request)

			// Assertions
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

func TestService_Withdraw(t *testing.T) {
	tests := []struct {
		name               string
		request            *request.CreateWithdrawRequest
		mockWallet         *walletentity.Wallet
		mockWalletErr      error
		mockAsset          bool
		mockAssetsResponse []*entity.Asset
		mockAssetsErr      error
		mockCreate         *entity.Asset
		mockCreateErr      error
		mockUpdate         bool
		mockUpdateErr      error
		expectedResult     *entity.Asset
		expectedError      error
	}{
		{
			name: "when balance is enough then should withdraw amount",
			request: &request.CreateWithdrawRequest{
				WalletID: 1,
				Name:     "BTC",
				Amount:   5.0,
			},
			mockWallet:         &walletentity.Wallet{ID: 1},
			mockWalletErr:      nil,
			mockAsset:          true,
			mockAssetsResponse: []*entity.Asset{{ID: 1, WalletID: 1, Name: "BTC", Amount: 10.0}},
			mockAssetsErr:      nil,
			mockUpdate:         true,
			mockUpdateErr:      nil,
			expectedResult:     &entity.Asset{ID: 1, WalletID: 1, Name: "BTC", Amount: 5.0},
			expectedError:      nil,
		},
		{
			name: "when balance is not enough then should return error",
			request: &request.CreateWithdrawRequest{
				WalletID: 1,
				Name:     "BTC",
				Amount:   15.0,
			},
			mockWallet:         &walletentity.Wallet{ID: 1},
			mockWalletErr:      nil,
			mockAsset:          true,
			mockAssetsResponse: []*entity.Asset{{ID: 1, WalletID: 1, Name: "BTC", Amount: 10.0}},
			mockAssetsErr:      nil,
			mockUpdateErr:      nil,
			expectedResult:     nil,
			expectedError:      errors.New("amount is not enough to withdraw"),
		},
		{
			name: "when wallet not found then should return error",
			request: &request.CreateWithdrawRequest{
				WalletID: 2,
				Name:     "ETH",
				Amount:   10.0,
			},
			mockWallet:         nil,
			mockWalletErr:      errors.New("wallet not found"),
			mockAssetsResponse: nil,
			mockAssetsErr:      nil,
			mockUpdateErr:      nil,
			expectedResult:     nil,
			expectedError:      errors.New("wallet not found"),
		},
		{
			name: "when repository error then should return error",
			request: &request.CreateWithdrawRequest{
				WalletID: 1,
				Name:     "BTC",
				Amount:   5.0,
			},
			mockWallet:         &walletentity.Wallet{ID: 1},
			mockWalletErr:      nil,
			mockAsset:          true,
			mockAssetsResponse: nil,
			mockAssetsErr:      errors.New("repository error"),
			mockUpdateErr:      nil,
			expectedResult:     nil,
			expectedError:      errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := assetmock.NewMockAssetRepository(t)
			mockWalletClient := walletmock.NewMockWalletClient(t)
			s := NewService(mockRepository, mockWalletClient)

			mockWalletClient.EXPECT().GetWallet(mock.Anything, tt.request.WalletID).
				Return(tt.mockWallet, tt.mockWalletErr).Once()

			if tt.mockAsset {
				mockRepository.EXPECT().GetAsset(mock.Anything, mock.Anything).
					Return(tt.mockAssetsResponse, tt.mockAssetsErr).Once()
			}
			if tt.mockCreate != nil || tt.mockCreateErr != nil {
				mockRepository.EXPECT().CreateAsset(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockCreate, tt.mockCreateErr).Once()
			}
			if tt.mockUpdate {
				mockRepository.EXPECT().UpdateAsset(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockUpdateErr).Once()
			}

			// Call the service method
			result, err := s.Withdraw(context.Background(), nil, tt.request)

			// Assertions
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
