package store

import (
	"fmt"
	"net/http"
	"proj/helper"
	"proj/internal/user"
)

type StoreService struct {
	StoreRepo StoreRepository
	UserRepo  user.UserRepository
}

func (service StoreService) GetAvailableStoreItems(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int)

	storeItems, err := service.StoreRepo.GetAvailableStoreItems(userId)
	if err != nil {
		fmt.Println(err.Error())
		helper.SendError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	helper.SendJSON(w, http.StatusOK, storeItems)
}
