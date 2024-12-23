package asset

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/safayildirim/asset-management-service/internal/asset/entity"
	assetmock "github.com/safayildirim/asset-management-service/internal/asset/mock"
	walletpkg "github.com/safayildirim/asset-management-service/pkg/client/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_GetAssets(t *testing.T) {
	e := echo.New()

	// Table-driven test cases
	tests := []struct {
		name                 string
		query                map[string]string
		mockService          bool
		mockReturnData       []*entity.Asset
		mockReturnErr        error
		expectedStatus       int
		expectErr            bool
		expectedErrorMessage string
	}{
		{
			name:        "when name and wallet_id query parameters are provided then should return assets",
			query:       map[string]string{"name": "BTC", "wallet_id": "1"},
			mockService: true,
			mockReturnData: []*entity.Asset{
				{ID: 1, Name: "BTC", WalletID: 1, Amount: 10},
				{ID: 2, Name: "ETH", WalletID: 1, Amount: 5},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "when no query parameters are provided then should return all assets",
			query:          nil,
			mockService:    true,
			mockReturnData: []*entity.Asset{},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:                 "when invalid query parameter is provided then should return bad request",
			query:                map[string]string{"invalid": "BTC"},
			mockService:          false,
			mockReturnData:       nil,
			mockReturnErr:        nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "schema: invalid path",
		},
		{
			name:                 "when service returns error then should return internal server error",
			query:                map[string]string{"name": "BTC"},
			mockService:          true,
			mockReturnData:       nil,
			mockReturnErr:        errors.New("service error"),
			expectedStatus:       http.StatusInternalServerError,
			expectErr:            true,
			expectedErrorMessage: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock service
			mockService := assetmock.NewMockAssetService(t)
			handler := NewHandler(mockService)

			// Mock service behavior based on test case
			if tt.mockService {
				mockService.EXPECT().GetAssets(mock.Anything, mock.Anything).Return(tt.mockReturnData,
					tt.mockReturnErr).Once()
			}

			// Create request and recorder
			req := httptest.NewRequest(http.MethodGet, "/assets", nil)
			q := req.URL.Query()
			for key, value := range tt.query {
				q.Add(key, value)
			}
			rec := httptest.NewRecorder()
			req.URL.RawQuery = q.Encode()
			ctx := e.NewContext(req, rec)

			err := handler.GetAssets(ctx)
			if tt.expectErr {
				assert.Error(t, err)
				httpErr := err.(*echo.HTTPError)
				assert.Contains(t, httpErr.Message, tt.expectedErrorMessage)
			} else {
				assert.NoError(t, err)

			}
		})
	}
}

func TestHandler_CreateAsset(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name                 string
		body                 string
		mockService          bool
		mockReturn           *entity.Asset
		mockError            error
		expectedStatus       int
		expectErr            bool
		expectedErrorMessage string
	}{
		{
			name:        "when valid request body is provided then should create asset",
			body:        `{"name":"BTC","amount":10, "wallet_id":1}`,
			mockService: true,
			mockReturn: &entity.Asset{
				ID:        1,
				CreatedAt: time.Time{},
				UpdatedAt: null.Time{},
				WalletID:  1,
				Name:      "BTC",
				Amount:    10,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:                 "when invalid request body is provided then should return bad request",
			body:                 `{"name":123}`,
			mockReturn:           nil,
			mockError:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "Unmarshal type error",
		},
		{
			name:                 "when empty request body is provided then should return bad request",
			body:                 `{"name":""}`,
			mockReturn:           nil,
			mockError:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "name: cannot be blank",
		},
		{
			name:                 "when asset already exists then should return conflict",
			body:                 `{"name":"BTC","amount":10, "wallet_id":1}`,
			mockService:          true,
			mockReturn:           nil,
			mockError:            ErrDuplicateAsset,
			expectedStatus:       http.StatusConflict,
			expectErr:            true,
			expectedErrorMessage: "asset already exist",
		},
		{
			name:                 "when service returns error then should return internal server error",
			body:                 `{"name":"BTC","amount":10, "wallet_id":1}`,
			mockService:          true,
			mockReturn:           nil,
			mockError:            errors.New("internal server error"),
			expectedStatus:       http.StatusInternalServerError,
			expectErr:            true,
			expectedErrorMessage: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := assetmock.NewMockAssetService(t)
			handler := NewHandler(mockService)

			if tt.mockService {
				mockService.EXPECT().CreateAsset(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockError).Once()
			}

			// Create request and response recorder
			req := httptest.NewRequest(http.MethodPost, "/assets", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Call handler
			err := handler.CreateAsset(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandler_Deposit(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name                 string
		body                 string
		mockService          bool
		mockReturn           *entity.Asset
		mockError            error
		expectedStatus       int
		expectErr            bool
		expectedErrorMessage string
	}{
		{
			name:           "when valid request body is provided then should deposit asset",
			body:           `{"wallet_id":1,"name":"BTC","amount":10}`,
			mockService:    true,
			mockReturn:     &entity.Asset{ID: 1, WalletID: 1, Name: "BTC", Amount: 10},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:                 "when invalid request body is provided then should return bad request",
			body:                 `{"wallet_id":"invalid"}`, // Invalid type for wallet_id
			mockReturn:           nil,
			mockError:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "cannot unmarshal",
		},
		{
			name:                 "when empty request body is provided then should return bad request",
			body:                 `{"wallet_id":0,"name":"","amount":0}`, // Missing required fields
			mockReturn:           nil,
			mockError:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "validation error",
		},
		{
			name:                 "when wallet not found then should return bad request",
			body:                 `{"wallet_id":1,"name":"BTC","amount":10}`,
			mockService:          true,
			mockReturn:           nil,
			mockError:            walletpkg.ErrWalletNotFound,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "wallet not found",
		},
		{
			name:                 "when service returns error then should return internal server error",
			body:                 `{"wallet_id":1,"name":"BTC","amount":10}`,
			mockService:          true,
			mockReturn:           nil,
			mockError:            errors.New("database error"),
			expectedStatus:       http.StatusInternalServerError,
			expectErr:            true,
			expectedErrorMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := assetmock.NewMockAssetService(t)
			handler := NewHandler(mockService)

			if tt.mockService {
				mockService.EXPECT().Deposit(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockError).Once()
			}

			// Create request and response recorder
			req := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Call handler
			err := handler.Deposit(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandler_Withdraw(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name                 string
		body                 string
		mockService          bool
		mockReturn           *entity.Asset
		mockError            error
		expectedStatus       int
		expectErr            bool
		expectedErrorMessage string
	}{
		{
			name:           "when valid request body is provided then should withdraw asset",
			body:           `{"wallet_id":1,"name":"BTC","amount":10}`,
			mockService:    true,
			mockReturn:     &entity.Asset{ID: 1, WalletID: 1, Name: "BTC", Amount: 10},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:                 "when invalid request body is provided then should return bad request",
			body:                 `{"wallet_id":"invalid"}`, // Invalid type for wallet_id
			mockReturn:           nil,
			mockError:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "cannot unmarshal",
		},
		{
			name:                 "when empty request body is provided then should return bad request",
			body:                 `{"wallet_id":0,"name":"","amount":0}`, // Missing required fields
			mockReturn:           nil,
			mockError:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "validation error",
		},
		{
			name:                 "when wallet not found then should return bad request",
			body:                 `{"wallet_id":1,"name":"BTC","amount":10}`,
			mockService:          true,
			mockReturn:           nil,
			mockError:            walletpkg.ErrWalletNotFound,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "wallet not found",
		},
		{
			name:                 "when service returns error then should return internal server error",
			body:                 `{"wallet_id":1,"name":"BTC","amount":10}`,
			mockService:          true,
			mockReturn:           nil,
			mockError:            errors.New("database error"),
			expectedStatus:       http.StatusInternalServerError,
			expectErr:            true,
			expectedErrorMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := assetmock.NewMockAssetService(t)
			handler := NewHandler(mockService)

			if tt.mockService {
				mockService.EXPECT().Withdraw(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockError).Once()
			}

			// Create request and response recorder
			req := httptest.NewRequest(http.MethodPost, "/withdraw", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Call handler
			err := handler.Withdraw(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
