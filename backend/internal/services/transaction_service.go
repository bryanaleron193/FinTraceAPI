package services

import (
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/schemas"
	"strings"

	"github.com/google/uuid"
)

func GetAllBorrowTransactions(auditedUserID uuid.UUID, input *schemas.BorrowTransactionRequest) ([]schemas.BorrowTransactionResponse, *schemas.Pagination, error) {
	borrowTransactions := []schemas.BorrowTransactionResponse{}
	pagination := new(schemas.Pagination)
	pagination.SetPagination(input.Limit, input.Page)

	input.TransactionName = strings.TrimSpace(input.TransactionName)
	input.LenderName = strings.TrimSpace(input.LenderName)

	baseQuery := database.DB.
		Table("tr_group_members AS a").
		Joins("JOIN ms_users AS b ON a.user_id = b.user_id AND b.is_deleted = FALSE").
		Joins(`JOIN tr_transaction_headers AS c ON b.user_id = c.lender_id AND c.is_deleted = FALSE`).
		Joins("JOIN ms_transaction_categories AS d ON c.transaction_category_id = d.transaction_category_id AND d.is_deleted = FALSE").
		Joins("JOIN tr_transaction_details AS e ON c.transaction_header_id = e.transaction_header_id AND e.is_deleted = FALSE AND e.borrower_id = ?", auditedUserID).
		Joins("JOIN tr_transaction_adjustments AS f ON c.transaction_header_id = f.transaction_header_id AND f.is_deleted = FALSE").
		Where("a.group_id = ? AND a.user_id <> ? AND a.is_deleted = FALSE", input.GroupID, auditedUserID)

	// Transaction Name filter
	if input.TransactionName != "" {
		baseQuery = baseQuery.Where("c.transaction_name ILIKE ?", input.TransactionName+"%")
	}

	// Lender filter
	if input.LenderID != uuid.Nil {
		baseQuery = baseQuery.Where("c.lender_id = ?", input.LenderID)
	} else if input.LenderName != "" {
		baseQuery = baseQuery.Where("b.name ILIKE ?", input.LenderName+"%")
	}

	// Transaction Category filter
	if input.TransactionCategoryID != uuid.Nil {
		baseQuery = baseQuery.Where("c.transaction_category_id = ?", input.TransactionCategoryID)
	}

	// IsPaid filter
	if input.IsPaid != nil {
		baseQuery = baseQuery.Where("e.is_paid = ?", *input.IsPaid)
	}

	// Count for pagination
	var total int64
	if err := baseQuery.
		Select("c.transaction_header_id").
		Group("c.transaction_header_id").
		Count(&total).Error; err != nil {
		return nil, nil, err
	}
	pagination.SetTotalPages(total)

	// Main data fetch
	if err := baseQuery.
		Select(`
			c.transaction_header_id,
			c.transaction_name,
			to_char(c.date_in, 'DD/MM/YYYY') AS transaction_date,
			d.transaction_category_name,
			b.name AS lender_name,
			SUM(e.amount * e.quantity) + (SUM(f.adjustment_amount) / NULLIF(COUNT(DISTINCT e.borrower_id), 0)) AS total_amount,
			CASE WHEN e.paid_at IS NULL THEN 'Not Yet Paid' ELSE 'Already Paid' END AS paid_status,
			CASE WHEN e.paid_at IS NULL THEN '-' ELSE to_char(e.paid_at, 'DD/MM/YYYY') END AS paid_at
		`).
		Group("c.transaction_header_id, c.transaction_name, c.date_in, d.transaction_category_name, b.name, e.paid_at").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Scan(&borrowTransactions).Error; err != nil {
		return nil, nil, err
	}

	return borrowTransactions, pagination, nil
}

func GetAllLendTransactions(auditedUserID uuid.UUID, input *schemas.LendTransactionRequest) ([]schemas.LendTransactionResponse, *schemas.Pagination, error) {
	lendTransactions := []schemas.LendTransactionResponse{}
	pagination := new(schemas.Pagination)
	pagination.SetPagination(input.Limit, input.Page)

	input.TransactionName = strings.TrimSpace(input.TransactionName)

	baseQuery := database.DB.
		Table("tr_group_members AS a").
		Joins("JOIN ms_users AS b ON a.user_id = b.user_id").
		Joins(`JOIN tr_transaction_headers AS c ON b.user_id = c.lender_id AND c.is_deleted = FALSE`).
		Joins("JOIN ms_transaction_categories AS d ON c.transaction_category_id = d.transaction_category_id AND d.is_deleted = FALSE").
		Joins("JOIN tr_transaction_details AS e ON c.transaction_header_id = e.transaction_header_id AND e.is_deleted = FALSE").
		Joins("JOIN tr_transaction_adjustments AS f ON c.transaction_header_id = f.transaction_header_id AND f.is_deleted = FALSE").
		Where("a.group_id = ? AND a.user_id = ? AND a.is_deleted = FALSE", input.GroupID, auditedUserID)

	// Transaction Name filter
	if input.TransactionName != "" {
		baseQuery = baseQuery.Where("c.transaction_name ILIKE ?", input.TransactionName+"%")
	}

	// Transaction Category filter
	if input.TransactionCategoryID != uuid.Nil {
		baseQuery = baseQuery.Where("c.transaction_category_id = ?", input.TransactionCategoryID)
	}

	// IsPaid filter
	if input.IsPaid != nil {
		baseQuery = baseQuery.Where("e.is_paid = ?", *input.IsPaid)
	}

	// Count for pagination
	var total int64
	if err := baseQuery.
		Select("c.transaction_header_id").
		Group("c.transaction_header_id").
		Count(&total).Error; err != nil {
		return nil, nil, err
	}
	pagination.SetTotalPages(total)

	// Main data fetch
	if err := baseQuery.
		Select(`
			c.transaction_header_id,
			c.transaction_name,
			to_char(c.date_in, 'DD/MM/YYYY') AS transaction_date,
			d.transaction_category_name,
			SUM(e.amount * e.quantity) + SUM(f.adjustment_amount) AS total_amount,
			CASE WHEN e.paid_at IS NULL THEN 'Not Yet Paid' ELSE 'Already Paid' END AS paid_status,
			CASE WHEN e.paid_at IS NULL THEN '-' ELSE to_char(e.paid_at, 'DD/MM/YYYY') END AS paid_at
		`).
		Group("c.transaction_header_id, c.transaction_name, c.date_in, d.transaction_category_name, e.paid_at").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Scan(&lendTransactions).Error; err != nil {
		return nil, nil, err
	}

	return lendTransactions, pagination, nil
}

func GetTransactionDetailById(input *schemas.TransactionDetailRequest) (*schemas.TransactionDetailResponse, error) {
	transaction := new(schemas.TransactionDetailResponse)

	// Fetch transaction header + sub_total + adjustments + grand_total
	if err := database.DB.
		Table("tr_transaction_headers AS a").
		Select(`
			a.transaction_header_id,
			a.transaction_name,
			b.transaction_category_name,
			c.name,
			to_char(a.date_in, 'DD/MM/YYYY') AS transaction_date,
			COALESCE(SUM(d.amount * d.quantity), 0) AS sub_total,
			COALESCE(SUM(f.adjustment_amount), 0) AS total_adjustment,
			COALESCE(SUM(d.amount * d.quantity), 0) + COALESCE(SUM(f.adjustment_amount), 0) AS grand_total
		`).
		Joins("JOIN ms_transaction_categories AS b ON a.transaction_category_id = b.transaction_category_id AND b.is_deleted = FALSE").
		Joins("JOIN ms_users AS c ON a.lender_id = c.user_id").
		Joins("LEFT JOIN tr_transaction_details AS d ON a.transaction_header_id = d.transaction_header_id AND d.is_deleted = FALSE").
		Joins("LEFT JOIN tr_transaction_adjustments AS f ON a.transaction_header_id = f.transaction_header_id AND f.is_deleted = FALSE").
		Where("a.transaction_header_id = ? AND a.is_deleted = FALSE", input.TransactionHeaderID).
		Group("a.transaction_header_id, a.transaction_name, b.transaction_category_name, c.name, a.date_in").
		Scan(transaction).Error; err != nil {
		return nil, err
	}

	// Fetch borrowers
	var borrowers []schemas.Borrower
	if err := database.DB.
		Table("tr_transaction_details AS d").
		Select(`
			d.borrower_id,
			u.name,
			COALESCE(SUM(d.amount * d.quantity),0) AS total
		`).
		Joins("JOIN ms_users AS u ON d.borrower_id = u.user_id").
		Where("d.transaction_header_id = ? AND d.is_deleted = FALSE", input.TransactionHeaderID).
		Group("d.borrower_id, u.name").
		Scan(&borrowers).Error; err != nil {
		return nil, err
	}

	// Fetch details per borrower
	for i := range borrowers {
		var details []schemas.BorrowerDetail
		if err := database.DB.
			Table("tr_transaction_details AS d").
			Select(`
				d.detail_name,
				d.amount,
				d.quantity,
				(d.amount * d.quantity) AS item_total,
				CASE WHEN d.paid_at IS NULL THEN FALSE ELSE TRUE END AS is_paid,
				CASE WHEN d.paid_at IS NULL THEN NULL ELSE to_char(d.paid_at, 'DD/MM/YYYY') END AS paid_at
			`).
			Where("d.transaction_header_id = ? AND d.borrower_id = ? AND d.is_deleted = FALSE",
				input.TransactionHeaderID, borrowers[i].BorrowerID).
			Scan(&details).Error; err != nil {
			return nil, err
		}
		borrowers[i].Details = details
	}

	transaction.Borrowers = borrowers

	// Fetch adjustments
	var adjustments []schemas.Adjustment
	if err := database.DB.
		Table("tr_transaction_adjustments").
		Select("adjustment_name, adjustment_amount").
		Where("transaction_header_id = ? AND is_deleted = FALSE", input.TransactionHeaderID).
		Scan(&adjustments).Error; err != nil {
		return nil, err
	}

	transaction.Adjustments = adjustments

	return transaction, nil
}

func GetBorrowersTotalByTransaction(input *schemas.TransactionDetailRequest) ([]schemas.BorrowerTotalResponse, error) {
	var borrowers []schemas.BorrowerTotalResponse

	// Count total distinct borrowers first for splitting adjustments
	var borrowerCount int64
	if err := database.DB.
		Table("tr_transaction_details").
		Where("transaction_header_id = ? AND is_deleted = FALSE", input.TransactionHeaderID).
		Distinct("borrower_id").
		Count(&borrowerCount).Error; err != nil {
		return nil, err
	}

	// Main query: total per borrower + their share of adjustments
	if err := database.DB.
		Table("tr_transaction_details AS d").
		Select(`
			d.borrower_id,
			u.name AS borrower_name,
			COALESCE(SUM(d.amount * d.quantity), 0) 
			+ COALESCE(SUM(f.adjustment_amount) / ?, 0) AS total_owed
		`, borrowerCount).
		Joins("JOIN ms_users AS u ON d.borrower_id = u.user_id").
		Joins("LEFT JOIN tr_transaction_adjustments AS f ON d.transaction_header_id = f.transaction_header_id AND f.is_deleted = FALSE").
		Where("d.transaction_header_id = ? AND d.is_deleted = FALSE", input.TransactionHeaderID).
		Group("d.borrower_id, u.name").
		Scan(&borrowers).Error; err != nil {
		return nil, err
	}

	return borrowers, nil
}
