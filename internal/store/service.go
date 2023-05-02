package store

import (
	"encoding/json"
	"net/http"
	"proj/helper"

	"github.com/go-playground/validator/v10"
)

type StoreService struct {
	StoreRepo StoreRepository
}

func (service StoreService) PurchaseStoreItem(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int)
	var storeItem StoreItem

	err := json.NewDecoder(r.Body).Decode(&storeItem)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, "bad request")

		return
	}

	validate := validator.New()
	err = validate.Struct(storeItem)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, "bad request")

		return
	}

	err = service.StoreRepo.PurchaseStoreItem(userId, storeItem)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, "bad request")

		return
	}

	helper.SendJSON(w, http.StatusOK, nil)
}

func (service StoreService) GetAvailableStoreItems(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int)

	storeItems, err := service.StoreRepo.GetAvailableStoreItems(userId)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	helper.SendJSON(w, http.StatusOK, storeItems)
}
