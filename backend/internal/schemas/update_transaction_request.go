package schemas

type UpdateTransactionRequest struct {
	Header      TransactionHeaderRequest       `json:"header"`
	Details     []TransactionDetailRequest     `json:"details"`
	Adjustments []TransactionAdjustmentRequest `json:"adjustments"`
}
