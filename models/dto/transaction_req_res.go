package dto

import "time"

type TrxReport struct {
	TrxId            string `json:"transaction_id"`
	AccountNumber    string `json:"account_number"`
	Kredit           int    `json:"kredit"`
	Debet            int    `json:"debet"`
	Type             string `json:"type"`
	AccountNumberXId string `json:"account_numberx"` //account number yang menjadi pengirim, penerima,atau account sendiri(khusus kasus topup)
	Balance          int    `json:"balance"`
	CreatedAt        string `json:"created_at"`
}

type TrxReq struct {
	ID              string    `json:"id"`
	SenderAccount   string    `json:"sender_account"`
	ReceiverAccount string    `json:"receiver_account"`
	Type            string    `json:"type"`
	Amount          int       `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
}

type TrxData struct {
	SenderAccount   string
	ReceiverAccount string
	Type            string
	Amount          int
	CreatedAt       string
}
