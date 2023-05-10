package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"proj/helper"
	"strconv"

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

func (service StoreService) DeleteUserStoreItem(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int)
	storeItemId, err := strconv.Atoi(r.URL.Query().Get("storeItemId"))
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, "bad request")

		return
	}

	err = service.StoreRepo.DeleteUserStoreItem(userId, storeItemId)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, "internal server error")
	}

	helper.SendJSON(w, http.StatusOK, nil)
}

func (service StoreService) GetUserStoreItems(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int)

	storeItems, err := service.StoreRepo.GetUserStoreItems(userId)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	validate := validator.New()

	for _, storeItem := range storeItems {
		err = validate.Struct(storeItem)
		if err != nil {
			helper.SendError(w, http.StatusInternalServerError, "internal server error")

			return
		}
	}

	helper.SendJSON(w, http.StatusOK, storeItems)
}

func (service StoreService) GetAvailableStoreItems(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int)

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, "bad request")

		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, "bad request")

		return
	}

	storeItems, err := service.StoreRepo.GetAvailableStoreItems(userId, page, limit)
	if err != nil {
		fmt.Println(err)
		helper.SendError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	validate := validator.New()

	for _, storeItem := range storeItems {
		err = validate.Struct(storeItem)
		if err != nil {
			helper.SendError(w, http.StatusInternalServerError, "internal server error")

			return
		}
	}

	helper.SendJSON(w, http.StatusOK, storeItems)
}
