package wallet

import (
	"context"
	"encoding/json"
	"github.com/safayildirim/asset-management-service/pkg/client/wallet/entity"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWallet(t *testing.T) {
	tests := []struct {
		name           string
		walletID       uint
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedResult *entity.Wallet
		expectedError  string
	}{
		{
			name:     "when wallet is found then should return wallet",
			walletID: 1,
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/wallets/1", r.URL.Path)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(entity.Wallet{ID: 1, Address: "Test Wallet"})
			},
			expectedResult: &entity.Wallet{ID: 1, Address: "Test Wallet"},
			expectedError:  "",
		},
		{
			name:     "when wallet is not found then should return not found error",
			walletID: 2,
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/wallets/2", r.URL.Path)
				w.WriteHeader(http.StatusNotFound)
			},
			expectedResult: nil,
			expectedError:  "wallet not found",
		},
		{
			name:     "when server returns an error then should return error",
			walletID: 3,
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/wallets/3", r.URL.Path)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error"}`))
			},
			expectedResult: nil,
			expectedError:  "unexpected status code: 500",
		},
		{
			name:     "when server returns invalid JSON then should return error",
			walletID: 4,
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/wallets/4", r.URL.Path)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{invalid-json}`))
			},
			expectedResult: nil,
			expectedError:  "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Initialize the client
			c := NewClient(server.URL)

			// Call the GetWallet function
			result, err := c.GetWallet(context.Background(), tt.walletID)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
