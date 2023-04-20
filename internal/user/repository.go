package user

import (
	"database/sql"
	"fmt"
)

type UserRepository struct {
	DB *sql.DB
}

func (repo UserRepository) CheckIfUserExists(credentials Credentials) (bool, error) {
	var isExists bool

	query, err := repo.DB.Prepare("SELECT EXISTS(SELECT 1 FROM users WHERE name=$1 AND password=$2)")
	if err != nil {
		fmt.Println(query, err.Error())
		return false, err
	}

	defer query.Close()

	result := query.QueryRow(credentials.Name, credentials.Password)

	err = result.Scan(&isExists)
	if err != nil {
		return false, err
	}

	return isExists, nil
}
