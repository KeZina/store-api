package store

import (
	"database/sql"
)

type StoreRepository struct {
	DB *sql.DB
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
