package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	api_models "github.com/kgugunava/effective_mobile_golang/internal/api/models"
	service_domain "github.com/kgugunava/effective_mobile_golang/internal/domain"
)

func transferStringMonthYearToDate(s string) (time.Time, error) {
	if t, err := time.Parse("01-2006", s); err == nil {
		return t, err
	}

	if t, err := time.Parse("2006-01", s); err == nil {
		return t, err
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s, expected MM-YYYY or YYYY-MM", s)
}

func transferDatetoString(d time.Time) string {
	s := d.Format("01-2006")
	return s
}

func transferCreateRequestToServiceDomain(req api_models.SubscriptionCreatePostRequest) (*service_domain.Subscription, error) {
	var startDate time.Time
	var err error
	if startDate, err = transferStringMonthYearToDate(req.StartDate); err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}
	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := transferStringMonthYearToDate(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		endDate = &parsed
	}
	return &service_domain.Subscription{
		SubscriptionID: uuid.UUID{},
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func transferServiceDomainToAPIModel(s *service_domain.Subscription) api_models.Subscription {
	resp := api_models.Subscription{
		SubscriptionID: s.SubscriptionID,
		ServiceName: s.ServiceName,
		Price: s.Price,
		UserID: s.UserID,
		StartDate: transferDatetoString(s.StartDate),
	}
	if s.EndDate != nil {
		str := s.EndDate.Format("01-2006")
		resp.EndDate = str
	}

	return resp
}