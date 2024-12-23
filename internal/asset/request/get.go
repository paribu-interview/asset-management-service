package request

type GetAssetsParams struct {
	ID       []uint   `json:"id" schema:"id"`
	WalletID []uint   `json:"wallet_id" schema:"wallet_id"`
	Name     []string `json:"name" schema:"name"`
}
