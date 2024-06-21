package model

type DatasourceInOrganization struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Datasource struct {
	Id   int    `json:"id" db:"id"`
	Host string `json:"host" db:"host"`
	Port int    `json:"port" db:"port"`
	Name string `json:"name" db:"name"`
}

type DatasourceUser struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
