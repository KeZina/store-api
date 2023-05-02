package user

import (
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func (repo UserRepository) GetUserById(userId int) (User, error) {
	var user User

	query, err := repo.DB.Prepare(`
		SELECT id, name, currency FROM users
		WHERE id=$1
	`)
	if err != nil {
		return User{}, err
	}

	defer query.Close()

	row := query.QueryRow(userId)

	err = row.Scan(&user.Id, &user.Name, &user.Currency)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (repo UserRepository) CheckUserCredentials(credentials Credentials) (int, error) {
	var userId int

	query, err := repo.DB.Prepare(`
		SELECT id FROM users WHERE name=$1 AND password=$2
	`)
	if err != nil {
		return 0, err
	}

	defer query.Close()

	row := query.QueryRow(credentials.Name, credentials.Password)

	err = row.Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
