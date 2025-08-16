package schemas

import "github.com/google/uuid"

type TransactionDetailRequest struct {
	TransactionHeaderID uuid.UUID `json:"transaction_header_id"`
}
