package schemas

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}