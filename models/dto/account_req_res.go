package dto

import "time"

// untuk kebutuhan register customer/user dan account
type AccountReq struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Password  string    `json:"password"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// untuk respon body,profil dan kebutuhan compare email password untuk jwt
type AccountData struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"*"`
	Address       string `json:"address"`
	Phone         string `json:"phone"`
	AccountNumber string `json:"account_number"`
	Balance       int    `json:"balance"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
