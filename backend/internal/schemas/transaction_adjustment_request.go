package schemas

import "github.com/google/uuid"

type TransactionAdjustmentRequest struct {
	TransactionAdjustmentID uuid.UUID `json:"transaction_adjustment_id"`
	TransactionHeaderID     uuid.UUID `json:"transaction_header_id"`
	AdjustmentName          string    `json:"adjustment_name"`
	AdjustmentAmount        float64   `json:"adjustment_amount"`
}
