package model

type SignUser struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id           int    `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
}

type TokenClaimsExtension struct {
	UserId int `json:"userId"`
}
