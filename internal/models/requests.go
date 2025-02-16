package models

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
