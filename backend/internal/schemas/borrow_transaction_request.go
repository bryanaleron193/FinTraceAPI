package schemas

import "github.com/google/uuid"

type BorrowTransactionRequest struct {
	GroupID               uuid.UUID `json:"group_id"`
	TransactionCategoryID uuid.UUID `json:"transaction_category_id"`
	TransactionName       string    `json:"transaction_name"`
	LenderID              uuid.UUID `json:"lender_id"`
	LenderName            string    `json:"lender_name"`
	IsPaid                *bool     `json:"is_paid"`
	Limit                 int       `json:"limit"`
	Page                  int       `json:"page"`
}
