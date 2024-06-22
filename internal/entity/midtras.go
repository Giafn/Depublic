package entity

type MidtransTransactionDetails struct {
	OrderID  string `json:"order_id"`
	GrossAmt int    `json:"gross_amount"`
}
