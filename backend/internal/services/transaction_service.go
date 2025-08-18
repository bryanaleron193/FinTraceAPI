package services

import (
	"fmt"
	"math"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	transactionHeader := new(models.TrTransactionHeaders)
	if err := database.DB.
		Table("tr_transaction_headers").
		Where("transaction_header_id = ? AND is_deleted = FALSE", input.TransactionHeaderID).
		Find(transactionHeader).Error; err != nil {
		return nil, err
	}

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
		Order("u.name ASC").
		Scan(&borrowers).Error; err != nil {
		return nil, err
	}

	//special case for round down transactions
	if transactionHeader.IsRoundDown {
		var roundDownAmount float64 = 0.00

		for i, b := range borrowers {
			switch b.BorrowerName {
			case "Bryan Aleron":
				borrowers[i].FinalTotalOwed = roundDownTo5000(b.TotalOwed)
				roundDownAmount = borrowers[i].TotalOwed - borrowers[i].FinalTotalOwed

			case "Christina Alexandra":
				borrowers[i].FinalTotalOwed = borrowers[i].TotalOwed + roundDownAmount
			}
		}
	}

	return borrowers, nil
}

func roundDownTo5000(value float64) float64 {
	return math.Floor(value/5000) * 5000
}

func CreateTransaction(auditedUserID uuid.UUID, input *schemas.CreateTransactionRequest) error {
	timeNow := time.Now()
	transactionHeaderID := uuid.New()

	// start transaction
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Insert header
		newHeader := models.TrTransactionHeaders{
			BaseModel: models.BaseModel{
				UserIn:    auditedUserID,
				UserUp:    nil,
				DateIn:    timeNow,
				DateUp:    nil,
				IsDeleted: false,
			},
			TransactionHeaderID:   transactionHeaderID,
			TransactionCategoryID: input.Header.TransactionCategoryID,
			LenderID:              input.Header.LenderID,
			LenderName:            input.Header.LenderName,
			TransactionName:       input.Header.TransactionName,
			TransactionDate:       timeNow,
			IsRoundDown:           input.Header.IsRoundDown,
		}

		if err := tx.Create(&newHeader).Error; err != nil {
			return fmt.Errorf("error creating transaction header: %v", err)
		}

		// 2. Insert Details (bulk)
		var newDetails []models.TrTransactionDetails
		for _, d := range input.Details {
			newDetails = append(newDetails, models.TrTransactionDetails{
				BaseModel: models.BaseModel{
					UserIn:    auditedUserID,
					UserUp:    nil,
					DateIn:    timeNow,
					DateUp:    nil,
					IsDeleted: false,
				},
				TransactionDetailID: uuid.New(),
				TransactionHeaderID: transactionHeaderID,
				BorrowerID:          d.BorrowerID,
				BorrowerName:        d.BorrowerName,
				DetailName:          d.DetailName,
				Amount:              d.Amount,
				Quantity:            d.Quantity,
				IsPaid:              false,
				PaidAt:              nil, // default null
			})
		}

		if len(newDetails) > 0 {
			if err := tx.Create(&newDetails).Error; err != nil {
				return fmt.Errorf("error creating transaction details: %v", err)
			}
		}

		// 3. Insert Adjustments (bulk)
		var newAdjustments []models.TrTransactionAdjustments
		for _, a := range input.Adjustments {
			newAdjustments = append(newAdjustments, models.TrTransactionAdjustments{
				BaseModel: models.BaseModel{
					UserIn:    auditedUserID,
					UserUp:    nil,
					DateIn:    timeNow,
					DateUp:    nil,
					IsDeleted: false,
				},
				TransactionAdjustmentID: uuid.New(),
				TransactionHeaderID:     transactionHeaderID,
				AdjustmentName:          a.AdjustmentName,
				AdjustmentAmount:        a.AdjustmentAmount,
			})
		}

		if len(newAdjustments) > 0 {
			if err := tx.Create(&newAdjustments).Error; err != nil {
				return fmt.Errorf("error creating transaction adjustments: %v", err)
			}
		}

		return nil
	})
}

func UpdateTransaction(auditedUserID uuid.UUID, input *schemas.UpdateTransactionRequest) error {
	timeNow := time.Now()

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Update Header
		if err := UpdateTransactionHeader(tx, auditedUserID, input.Header, timeNow); err != nil {
			return err
		}

		// 2. Update Details
		if err := UpdateTransactionDetails(tx, auditedUserID, input.Details, input.Header.TransactionHeaderID, timeNow); err != nil {
			return err
		}

		// 3. Update Adjustments
		if err := UpdateTransactionAdjustments(tx, auditedUserID, input.Adjustments, input.Header.TransactionHeaderID, timeNow); err != nil {
			return err
		}

		return nil
	})
}

func UpdateTransactionHeader(tx *gorm.DB, auditedUserID uuid.UUID, inputHeader schemas.TransactionHeaderRequest, timeNow time.Time) error {
	var header models.TrTransactionHeaders
	if err := tx.Where("transaction_header_id = ? AND is_deleted = false", inputHeader.TransactionHeaderID).
		First(&header).Error; err != nil {
		return fmt.Errorf("transaction header not found: %v", err)
	}

	headerUpdates, changed := diffHeader(header, inputHeader, auditedUserID, timeNow)
	if !changed {
		return nil // no update needed
	}

	// insert history
	header.DateUp = &timeNow
	if err := InsertHistoryTransactionHeader(&header); err != nil {
		return err
	}

	// update
	if err := tx.Model(&header).Updates(headerUpdates).Error; err != nil {
		return fmt.Errorf("error updating header: %v", err)
	}
	return nil
}

func UpdateTransactionDetails(tx *gorm.DB, auditedUserID uuid.UUID, detailsInput []schemas.TransactionDetailRequest, headerID uuid.UUID, timeNow time.Time) error {
	var existingDetails []models.TrTransactionDetails
	if err := tx.Where("transaction_header_id = ? AND is_deleted = false", headerID).Find(&existingDetails).Error; err != nil {
		return err
	}
	detailsMap := make(map[uuid.UUID]models.TrTransactionDetails)
	for _, d := range existingDetails {
		detailsMap[d.TransactionDetailID] = d
	}

	var detailsHistory []models.TrTransactionDetails
	var detailUpdates []struct {
		ID      uuid.UUID
		Updates map[string]interface{}
	}

	for _, d := range detailsInput {
		existing, ok := detailsMap[d.TransactionDetailID]
		if !ok {
			return fmt.Errorf("transaction detail not found: %v", d.TransactionDetailID)
		}

		updates, changed := diffDetail(existing, d, auditedUserID, timeNow)
		if changed {
			existing.DateUp = &timeNow
			detailsHistory = append(detailsHistory, existing)
			detailUpdates = append(detailUpdates, struct {
				ID      uuid.UUID
				Updates map[string]interface{}
			}{existing.TransactionDetailID, updates})
		}
	}

	if len(detailsHistory) > 0 {
		if err := InsertHistoryTransactionDetails(detailsHistory); err != nil {
			return err
		}
	}
	for _, u := range detailUpdates {
		if err := tx.Model(&models.TrTransactionDetails{}).Where("transaction_detail_id = ?", u.ID).Updates(u.Updates).Error; err != nil {
			return fmt.Errorf("error updating detail %v: %v", u.ID, err)
		}
	}
	return nil
}

func UpdateTransactionAdjustments(tx *gorm.DB, auditedUserID uuid.UUID, adjustmentsInput []schemas.TransactionAdjustmentRequest, headerID uuid.UUID, timeNow time.Time) error {
	var existingAdjustments []models.TrTransactionAdjustments
	if err := tx.Where("transaction_header_id = ? AND is_deleted = false", headerID).Find(&existingAdjustments).Error; err != nil {
		return err
	}
	adjustmentsMap := make(map[uuid.UUID]models.TrTransactionAdjustments)
	for _, a := range existingAdjustments {
		adjustmentsMap[a.TransactionAdjustmentID] = a
	}

	var adjustmentsHistory []models.TrTransactionAdjustments
	var adjustmentUpdates []struct {
		ID      uuid.UUID
		Updates map[string]interface{}
	}

	for _, a := range adjustmentsInput {
		existing, ok := adjustmentsMap[a.TransactionAdjustmentID]
		if !ok {
			return fmt.Errorf("transaction adjustment not found: %v", a.TransactionAdjustmentID)
		}

		updates, changed := diffAdjustment(existing, a, auditedUserID, timeNow)
		if changed {
			existing.DateUp = &timeNow
			adjustmentsHistory = append(adjustmentsHistory, existing)
			adjustmentUpdates = append(adjustmentUpdates, struct {
				ID      uuid.UUID
				Updates map[string]interface{}
			}{existing.TransactionAdjustmentID, updates})
		}
	}

	if len(adjustmentsHistory) > 0 {
		if err := InsertHistoryTransactionAdjustments(adjustmentsHistory); err != nil {
			return err
		}
	}
	for _, u := range adjustmentUpdates {
		if err := tx.Model(&models.TrTransactionAdjustments{}).Where("transaction_adjustment_id = ?", u.ID).Updates(u.Updates).Error; err != nil {
			return fmt.Errorf("error updating adjustment %v: %v", u.ID, err)
		}
	}

	return nil
}

func diffHeader(existing models.TrTransactionHeaders, input schemas.TransactionHeaderRequest, auditedUserID uuid.UUID, timeNow time.Time) (map[string]interface{}, bool) {
	updates := map[string]interface{}{}
	changed := false

	if existing.TransactionCategoryID != input.TransactionCategoryID {
		updates["transaction_category_id"] = input.TransactionCategoryID
		changed = true
	}
	if existing.LenderID != input.LenderID {
		updates["lender_id"] = input.LenderID
		changed = true
	}
	if existing.LenderName != input.LenderName {
		updates["lender_name"] = input.LenderName
		changed = true
	}
	if existing.TransactionName != input.TransactionName {
		updates["transaction_name"] = input.TransactionName
		changed = true
	}
	if existing.IsRoundDown != input.IsRoundDown {
		updates["is_round_down"] = input.IsRoundDown
		changed = true
	}

	if changed {
		updates["user_up"] = auditedUserID
		updates["date_up"] = timeNow
	}
	return updates, changed
}

func diffDetail(existing models.TrTransactionDetails, input schemas.TransactionDetailRequest, auditedUserID uuid.UUID, timeNow time.Time) (map[string]interface{}, bool) {
	updates := map[string]interface{}{}
	changed := false

	if existing.BorrowerID != input.BorrowerID {
		updates["borrower_id"] = input.BorrowerID
		changed = true
	}
	if existing.BorrowerName != input.BorrowerName {
		updates["borrower_name"] = input.BorrowerName
		changed = true
	}
	if existing.DetailName != input.DetailName {
		updates["detail_name"] = input.DetailName
		changed = true
	}
	if existing.Amount != input.Amount {
		updates["amount"] = input.Amount
		changed = true
	}
	if existing.Quantity != input.Quantity {
		updates["quantity"] = input.Quantity
		changed = true
	}
	if existing.IsPaid != input.IsPaid {
		updates["is_paid"] = input.IsPaid
		if input.IsPaid {
			updates["paid_at"] = timeNow
		} else {
			updates["paid_at"] = nil
		}
		changed = true
	}

	if changed {
		updates["user_up"] = auditedUserID
		updates["date_up"] = timeNow
	}
	return updates, changed
}

func diffAdjustment(existing models.TrTransactionAdjustments, input schemas.TransactionAdjustmentRequest, auditedUserID uuid.UUID, timeNow time.Time) (map[string]interface{}, bool) {
	updates := map[string]interface{}{}
	changed := false

	if existing.AdjustmentName != input.AdjustmentName {
		updates["adjustment_name"] = input.AdjustmentName
		changed = true
	}
	if existing.AdjustmentAmount != input.AdjustmentAmount {
		updates["adjustment_amount"] = input.AdjustmentAmount
		changed = true
	}

	if changed {
		updates["user_up"] = auditedUserID
		updates["date_up"] = timeNow
	}
	return updates, changed
}

func DeleteTransaction(auditedUserID uuid.UUID, input *schemas.DeleteTransactionRequest) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// --- 1. Get header ---
		var header models.TrTransactionHeaders
		if err := tx.Where("transaction_header_id = ? AND is_deleted = false", input.TransactionHeaderID).
			First(&header).Error; err != nil {
			return fmt.Errorf("transaction header not found: %v", err)
		}

		// --- 2. Get details ---
		var details []models.TrTransactionDetails
		if err := tx.Where("transaction_header_id = ? AND is_deleted = false", input.TransactionHeaderID).
			Find(&details).Error; err != nil {
			return fmt.Errorf("failed to get transaction details: %v", err)
		}

		// --- 3. Get adjustments ---
		var adjustments []models.TrTransactionAdjustments
		if err := tx.Where("transaction_header_id = ? AND is_deleted = false", input.TransactionHeaderID).
			Find(&adjustments).Error; err != nil {
			return fmt.Errorf("failed to get transaction adjustments: %v", err)
		}

		// --- 4. Insert into history ---
		if err := InsertHistoryTransactionHeader(&header); err != nil {
			return err
		}
		if err := InsertHistoryTransactionDetails(details); err != nil {
			return err
		}
		if err := InsertHistoryTransactionAdjustments(adjustments); err != nil {
			return err
		}

		// use time.Now() for DateUp
		timeNow := time.Now()

		// --- 5. Soft delete ---
		if err := tx.Model(&models.TrTransactionHeaders{}).
			Where("transaction_header_id = ?", input.TransactionHeaderID).
			Updates(map[string]interface{}{
				"user_up":    auditedUserID,
				"date_up":    timeNow,
				"is_deleted": true,
			}).Error; err != nil {
			return fmt.Errorf("failed to soft delete header: %v", err)
		}

		if err := tx.Model(&models.TrTransactionDetails{}).
			Where("transaction_header_id = ?", input.TransactionHeaderID).
			Updates(map[string]interface{}{
				"user_up":    auditedUserID,
				"date_up":    timeNow,
				"is_deleted": true,
			}).Error; err != nil {
			return fmt.Errorf("failed to soft delete details: %v", err)
		}

		if err := tx.Model(&models.TrTransactionAdjustments{}).
			Where("transaction_header_id = ?", input.TransactionHeaderID).
			Updates(map[string]interface{}{
				"user_up":    auditedUserID,
				"date_up":    timeNow,
				"is_deleted": true,
			}).Error; err != nil {
			return fmt.Errorf("failed to soft delete adjustments: %v", err)
		}

		return nil
	})
}
