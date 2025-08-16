package services

import (
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
)

func GetAllTransactionCategories() ([]schemas.TransactionCategoryResponse, error) {
	transactionCategories := []schemas.TransactionCategoryResponse{}

	query := database.DB.
		Model(&models.MsTransactionCategories{}).
		Select("transaction_category_id, transaction_category_name").
		Where("is_deleted = ?", false).
		Scan(&transactionCategories)

	if err := query.Error; err != nil {
		return nil, err
	}

	return transactionCategories, nil
}
