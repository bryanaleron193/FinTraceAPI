package schemas

import "github.com/google/uuid"

type BorrowTransactionResponse struct {
	TransactionHeaderID     uuid.UUID `json:"group_id"`
	TransactionName         string    `json:"transaction_name"`
	TransactionCategoryName uuid.UUID `json:"transaction_category_name"`
	LenderName              string    `json:"lender_name"`
	TransactionDate         string    `json:"transaction_date"`
	TotalAmount             float64   `json:"total_amount"`
	PaidStatus              string    `json:"paid_status"`
	PaidAt                  string    `json:"paid_at"`
}
