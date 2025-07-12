package schemas

// CreateItemSchemaIn represents the input data for the creation of item
type CreateItemSchemaIn struct {
	Name string `json:"name" binding:"required,min=1,max=100" example:"Sample Item"`
}

// CreateItemSchemaIn represents the input data for the creation of item
type UpdateItemSchemaIn struct {
	Name string `json:"name" binding:"required,min=1,max=100" example:"Sample Item"`
}
