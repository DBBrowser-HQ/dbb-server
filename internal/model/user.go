package model

import "github.com/golang-jwt/jwt"

type SignUser struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id           int    `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
}

type AccessTokenClaimsExtension struct {
	UserId int `json:"userId"`
}

type RefreshTokenClaimsExtension struct {
	UserId int    `json:"userId"`
	JTI    string `json:"jti"`
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	AccessTokenClaimsExtension
}

type RefreshTokenClaims struct {
	jwt.StandardClaims
	RefreshTokenClaimsExtension
}

type RefreshSession struct {
	Id           int    `db:"id"`
	RefreshToken string `db:"refresh_token"`
	JTI          string `db:"jti"`
	UserId       int    `db:"user_id"`
}

type UserWithoutPassword struct {
	Id    int    `json:"id" db:"id"`
	Login string `json:"login" db:"login"`
}

type UserInOrganization struct {
	Id    int    `json:"id" db:"id"`
	Login string `json:"login" db:"login"`
	Role  string `json:"role" db:"role"`
}
