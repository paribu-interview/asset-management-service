package transaction

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/safayildirim/asset-management-service/internal/transaction/entity"
	transactionmock "github.com/safayildirim/asset-management-service/internal/transaction/mock"
	walletpkg "github.com/safayildirim/asset-management-service/pkg/client/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_GetTransactions(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name                 string
		query                map[string]string
		mockService          bool
		mockReturnData       []*entity.Transaction
		mockReturnErr        error
		expectedStatus       int
		expectErr            bool
		expectedErrorMessage string
	}{
		{
			name: "when valid query parameters are provided then should return transactions",
			query: map[string]string{
				"id":                    "1",
				"source_wallet_id":      "1001",
				"destination_wallet_id": "1003",
				"status":                "pending",
			},
			mockService: true,
			mockReturnData: []*entity.Transaction{
				{ID: 1, SourceWalletID: 1001, DestinationWalletID: 1003, Status: "pending"},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "when no query parameters are provided then should return all transactions",
			query:          nil,
			mockService:    true,
			mockReturnData: []*entity.Transaction{},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "when invalid query parameter is provided then should return bad request",
			query: map[string]string{
				"id": "not-integer",
			},
			mockService:          false,
			mockReturnData:       nil,
			mockReturnErr:        nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "schema: error",
		},
		{
			name: "when service returns error then should return internal server error",
			query: map[string]string{
				"id":                    "1",
				"source_wallet_id":      "1001",
				"destination_wallet_id": "1003",
				"status":                "pending",
			},
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
			mockService := transactionmock.NewMockTransactionService(t)
			handler := NewHandler(mockService)

			if tt.mockService {
				mockService.EXPECT().GetTransactions(mock.Anything, mock.Anything).Return(tt.mockReturnData,
					tt.mockReturnErr).Once()
			}

			req := httptest.NewRequest(http.MethodGet, "/assets", nil)
			q := req.URL.Query()
			for key, value := range tt.query {
				q.Add(key, value)
			}
			rec := httptest.NewRecorder()
			req.URL.RawQuery = q.Encode()
			ctx := e.NewContext(req, rec)

			err := handler.GetTransactions(ctx)
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

func TestHandler_DeleteTransaction(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name                 string
		transactionID        string
		mockService          bool
		mockError            error
		expectedStatus       int
		expectErr            bool
		expectedErrorMessage string
	}{
		{
			name:                 "when valid transaction ID is provided then should delete transaction",
			transactionID:        "1",
			mockService:          true,
			mockError:            nil,
			expectedStatus:       http.StatusNoContent,
			expectedErrorMessage: "",
		},
		{
			name:                 "when invalid transaction ID is provided then should return bad request",
			transactionID:        "invalid",
			mockError:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "invalid syntax",
		},
		{
			name:                 "when transaction not found then should return not found",
			transactionID:        "2",
			mockService:          true,
			mockError:            errors.New("transaction not found"),
			expectedStatus:       http.StatusInternalServerError,
			expectErr:            true,
			expectedErrorMessage: "transaction not found",
		},
		{
			name:                 "when service returns error then should return internal server error",
			transactionID:        "3",
			mockService:          true,
			mockError:            errors.New("internal server error"),
			expectedStatus:       http.StatusInternalServerError,
			expectErr:            true,
			expectedErrorMessage: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := transactionmock.NewMockTransactionService(t)
			handler := NewHandler(mockService)

			if tt.mockService {
				mockService.EXPECT().CancelTransaction(mock.Anything, mock.Anything).Return(tt.mockError).Once()
			}

			req := httptest.NewRequest(http.MethodDelete, "/transactions/:id", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/transactions/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.transactionID)

			err := handler.DeleteTransaction(ctx)

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

func TestHandler_ScheduleTransaction(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name                 string
		body                 string
		mockReturn           *entity.Transaction
		mockService          bool
		mockError            error
		expectedStatus       int
		expectedResult       *entity.Transaction
		expectErr            bool
		expectedErrorMessage string
	}{
		{
			name:        "when valid request body is provided then should schedule transaction",
			body:        `{"source_wallet_id":1,"destination_wallet_id":2, "asset_name":"BTC","amount":10,"scheduled_at":"2024-01-01T12:00:00Z"}`,
			mockService: true,
			mockReturn: &entity.Transaction{
				ID:                  1,
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              10,
				Status:              "pending",
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedResult: &entity.Transaction{
				ID:                  1,
				SourceWalletID:      1,
				DestinationWalletID: 2,
				AssetName:           "BTC",
				Amount:              10,
				Status:              "pending",
			},
		},
		{
			name:                 "when invalid request body is provided then should return bad request",
			body:                 `{"amount":"invalid"}`,
			mockReturn:           nil,
			mockError:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "Unmarshal type error",
		},
		{
			name:                 "when asset name is empty then should return bad request",
			body:                 `{"source_wallet_id":1,"destination_wallet_id":2, "asset_name":"","amount":10,"scheduled_at":"2024-01-01T12:00:00Z"}`,
			mockReturn:           nil,
			mockError:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "asset_name: cannot be blank",
		},
		{
			name:                 "when wallet not found then should return bad request",
			body:                 `{"source_wallet_id":1,"destination_wallet_id":2, "asset_name":"BTC","amount":10,"scheduled_at":"2024-01-01T12:00:00Z"}`,
			mockService:          true,
			mockReturn:           nil,
			mockError:            walletpkg.ErrWalletNotFound,
			expectedStatus:       http.StatusBadRequest,
			expectErr:            true,
			expectedErrorMessage: "wallet not found",
		},
		{
			name:                 "when internal server error then should return internal server error",
			body:                 `{"source_wallet_id":1,"destination_wallet_id":2, "asset_name":"BTC","amount":10,"scheduled_at":"2024-01-01T12:00:00Z"}`,
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
			mockService := transactionmock.NewMockTransactionService(t)
			handler := NewHandler(mockService)

			if tt.mockService {
				mockService.EXPECT().ScheduleTransaction(mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockError).Once()
			}

			req := httptest.NewRequest(http.MethodPost, "/transactions/schedule", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			err := handler.ScheduleTransaction(ctx)

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
