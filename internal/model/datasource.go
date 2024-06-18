package model

type DatasourceInOrganization struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
