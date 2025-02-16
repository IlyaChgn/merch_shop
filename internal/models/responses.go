package models

type ErrResponse struct {
	Errors string `json:"errors"`
}

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
