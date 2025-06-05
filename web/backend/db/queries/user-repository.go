package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/krup4/goPetProject/web/backend/db"
)

func GetUserByLogin(login string) (db.User, error) {
	var user db.User
	err := db.Pool.QueryRow(context.Background(), `
		SELECT (id, login, password, name) FROM users WHERE login = $1
	`, login).Scan(&user)

	return user, err
}

func CreateNewUser(login string, password []byte, name string) (pgconn.CommandTag, error) {
	commandTag, err := db.Pool.Exec(context.Background(), `
		INSERT INTO users (login, password, name) 
		VALUES ($1, $2, $3) ON CONFLICT DO NOTHING
	`, login, password, name)

	return commandTag, err
}
