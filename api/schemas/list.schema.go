package schemas

type CreateList struct {
	Title		string `json:"title" binding:"required"`
	Description	string `json:"description" binding:"required"`
}

type UpdateList struct {
	Title		*string	`json:"title"`
	Description	*string	`json:"description"`
}