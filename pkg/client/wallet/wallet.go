package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/safayildirim/asset-management-service/pkg/client/wallet/entity"
	"io"
	"net/http"
	"time"
)

type Client interface {
	GetWallet(ctx context.Context, id uint) (*entity.Wallet, error)
}

type client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string) Client {
	return &client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
	}
}

func (c client) GetWallet(ctx context.Context, id uint) (*entity.Wallet, error) {
	// Construct the request URL
	url := fmt.Sprintf("%s/wallets/%d", c.baseURL, id)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrWalletNotFound
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	// Parse the response body
	var wallet entity.Wallet
	if err = json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &wallet, nil
}
