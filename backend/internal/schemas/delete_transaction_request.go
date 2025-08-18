package schemas

import "github.com/google/uuid"

type DeleteTransactionRequest struct {
	TransactionHeaderID uuid.UUID `json:"transaction_header_id"`
}
