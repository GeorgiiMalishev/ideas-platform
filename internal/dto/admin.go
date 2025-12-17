package dto

type RegisterAdminRequest struct {
	Login          string `json:"login" binding:"required"`
	Password       string `json:"password" binding:"required"`
	CoffeeShopName string `json:"coffee_shop_name" binding:"required"`
	Address        string `json:"address" binding:"required"`
}

type AdminLoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
