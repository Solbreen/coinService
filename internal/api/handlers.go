package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"coinService/internal/database"
	"coinService/internal/models"
)

type MerchService struct {
	db *database.DB
}

func NewMerchService(db *database.DB) *MerchService {
	return &MerchService{db: db}
}

func (s *MerchService) HandleInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	info, err := s.db.GetUserInfo(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func (s *MerchService) HandleSendCoin(w http.ResponseWriter, r *http.Request) {
	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)
	if err := s.db.SendCoin(userID, req.ToUser, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *MerchService) HandleBuy(w http.ResponseWriter, r *http.Request) {
	item := strings.TrimPrefix(r.URL.Path, "/api/buy/")
	userID := r.Context().Value("userID").(int)

	if err := s.db.BuyItem(userID, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
