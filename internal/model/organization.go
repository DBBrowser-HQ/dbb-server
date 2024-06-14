package model

const (
	OwnerRole    = "owner"
	AdminRole    = "admin"
	RedactorRole = "redactor"
	ReaderRole   = "reader"
)

type OrganizationForUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
