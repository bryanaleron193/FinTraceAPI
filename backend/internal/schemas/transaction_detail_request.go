package schemas

import "github.com/google/uuid"

type TransactionDetailRequest struct {
	TransactionDetailID uuid.UUID `json:"transaction_detail_id"`
	TransactionHeaderID uuid.UUID `json:"transaction_header_id"`
	BorrowerID          uuid.UUID `json:"borrower_id"`
	BorrowerName        string    `json:"borrower_name"`
	DetailName          string    `json:"detail_name"`
	Amount              float64   `json:"amount"`
	Quantity            int       `json:"quantity"`
	IsPaid              bool      `json:"is_paid"`
}
