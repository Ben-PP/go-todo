package schemas

type CreateUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"is_admin"`
}

type UpdateUser struct {
	Username string `json:"username" binding:"required"`
	IsAdmin  *bool   `json:"is_admin" binding:"required"`
}