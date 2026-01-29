package models

import (
	"github.com/google/uuid"
)

type SubscriptionListGetRequest struct {
	ServiceName string `json:"service_name"`
	UserID uuid.UUID `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}