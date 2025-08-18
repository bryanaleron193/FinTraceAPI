package schemas

import "github.com/google/uuid"

type TransactionHeaderRequest struct {
	TransactionHeaderID   uuid.UUID `json:"transaction_header_id"`
	TransactionCategoryID uuid.UUID `json:"transaction_category_id"`
	LenderID              uuid.UUID `json:"lender_id"`
	LenderName            string    `json:"lender_name"`
	TransactionName       string    `json:"transaction_name"`
	IsRoundDown           bool      `json:"is_round_down"`
}
