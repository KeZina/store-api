package store

type StoreItem struct {
	Id    int    `json:"id" validate:"required, numeric"`
	Title string `json:"title" validate:"required"`
	Price int    `json:"price" validate:"required, numeric"`
}
