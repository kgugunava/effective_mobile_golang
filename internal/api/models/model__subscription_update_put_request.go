package models

import (
	"github.com/google/uuid"
)

type SubscriptionUpdatePutRequest struct {
	ServiceName string `json:"service_name"`
	Price int `json:"price"`
	UserID uuid.UUID `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}