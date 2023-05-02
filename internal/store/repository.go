package store

import (
	"database/sql"
)

type StoreRepository struct {
	DB *sql.DB
}

func (repo StoreRepository) PurchaseStoreItem(userId int, storeItem StoreItem) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE users
		SET currency=users.currency - $1
		WHERE id=$2
	`, storeItem.Price, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO user_store_items (user_id, store_item_id)
		VALUES ($1, $2)
	`, userId, storeItem.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repo StoreRepository) GetAvailableStoreItems(userId int) ([]StoreItem, error) {
	var storeItems []StoreItem

	query, err := repo.DB.Prepare(`
		SELECT i.id, i.title, i.price FROM store_items i 
		WHERE i.id NOT IN (
			SELECT store_item_id FROM user_store_items WHERE user_id=$1
		)
	`)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	rows, err := query.Query(userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var storeItem StoreItem

		if err := rows.Scan(&storeItem.Id, &storeItem.Title, &storeItem.Price); err != nil {
			return nil, err
		}

		storeItems = append(storeItems, storeItem)
	}

	return storeItems, nil
}
