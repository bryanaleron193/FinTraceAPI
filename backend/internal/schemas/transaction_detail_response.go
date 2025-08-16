package schemas

import "github.com/google/uuid"

type Adjustment struct {
	AdjustmentName   string  `gorm:"column:adjustment_name"`
	AdjustmentAmount float64 `gorm:"column:adjustment_amount"`
}

type BorrowerDetail struct {
	DetailName string  `gorm:"column:detail_name"`
	Amount     float64 `gorm:"column:amount"`
	Quantity   int     `gorm:"column:quantity"`
	ItemTotal  float64 `gorm:"column:item_total"`
	IsPaid     bool    `gorm:"column:is_paid"`
	PaidAt     *string `gorm:"column:paid_at"`
}

type Borrower struct {
	BorrowerID   uuid.UUID        `gorm:"column:borrower_id"`
	BorrowerName string           `gorm:"column:borrower_name"`
	Details      []BorrowerDetail `gorm:"-"`
	Total        float64          `gorm:"column:total"`
}

type TransactionDetailResponse struct {
	TransactionHeaderID     uuid.UUID    `gorm:"column:transaction_header_id"`
	TransactionName         string       `gorm:"column:transaction_name"`
	TransactionCategoryName string       `gorm:"column:transaction_category_name"`
	TransactionDate         string       `gorm:"column:transaction_date"`
	LenderName              string       `gorm:"column:lender_name"`
	Borrowers               []Borrower   `gorm:"-"`
	SubTotal                float64      `gorm:"column:sub_total"`
	Adjustments             []Adjustment `json:"Adjustments"`
	TotalAdjustment         float64      `gorm:"column:total_adjustment"`
	GrandTotal              float64      `gorm:"column:grand_total"`
}
