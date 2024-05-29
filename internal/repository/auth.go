package repository

import (
	"dbb-server/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(login, passwordHash string) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (login, password_hash) VALUES ($1, $2) RETURNING id`, RefreshSessionsTable)

	var id int
	err := r.db.Get(&id, query, login, passwordHash)
	return id, err
}

func (r *AuthPostgres) GetUserByCredentials(login, password string) (model.User, error) {
	query := fmt.Sprintf(`SELECT id, login, password_hash FROM %s WHERE login=$1 AND password_hash=$2`, UsersTable)

	var user model.User
	err := r.db.Get(&user, query, login, password)
	return user, err
}

func (r *AuthPostgres) SaveSessionData(refreshToken, jti string, userId int) error {
	query := fmt.Sprintf(`INSERT INTO %s (refresh_token, jti, user_id) VALUES ($1, $2, $3)`, RefreshSessionsTable)

	_, err := r.db.Exec(query, refreshToken, jti, userId)
	return err
}

func (r *AuthPostgres) GetSessionData(userId int) (model.RefreshSession, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id=$1`, RefreshSessionsTable)

	var refreshSession model.RefreshSession
	err := r.db.Get(&refreshSession, query, userId)
	return refreshSession, err
}

func (r *AuthPostgres) DeleteSessionForUser(userId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id=$1`, RefreshSessionsTable)
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *AuthPostgres) UpdateSessionData(refreshToken, jti string, userId int) error {
	query := fmt.Sprintf(`UPDATE %s SET refresh_token=$1, jti = $2 WHERE user_id=$3`, RefreshSessionsTable)

	_, err := r.db.Exec(query, refreshToken, jti, userId)
	return err
}
