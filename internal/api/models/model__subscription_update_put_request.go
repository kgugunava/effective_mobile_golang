package models

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionUpdatePutRequest struct {
	ServiceName *string `json:"service_name"`
	Price *int `json:"price"`
	UserID *uuid.UUID `json:"user_id"`
	StartDate *time.Time `json:"start_date"`
	EndDate *time.Time `json:"end_date"`
}