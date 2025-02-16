package models

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type ReceivedCoinsItem struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentCoinsItem struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []ReceivedCoinsItem `json:"received"`
	Sent     []SentCoinsItem     `json:"sent"`
}

type MerchItem struct {
	ID    uint
	Name  string
	Price uint
}
