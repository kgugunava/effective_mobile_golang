package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	api_models "github.com/kgugunava/effective_mobile_golang/internal/api/models"
)

type SubscriptionAPI struct {
	subscriptionService SubscriptionService
	logger             *slog.Logger
}

func NewSubscriptionAPI(service SubscriptionService, logger *slog.Logger) *SubscriptionAPI {
	return &SubscriptionAPI{
		subscriptionService: service,
		logger:              logger,
	}
}

func (api *SubscriptionAPI) SubscriptionCreatePost(c *gin.Context) {
	api.logger.Info("handling create subscription request", slog.String("method", "POST"), slog.String("path", "/subscriptions"))

	var newSubscription api_models.SubscriptionCreatePostRequest

	if err := c.ShouldBindJSON(&newSubscription); err != nil {
		api.logger.Error("failed to bind create subscription request",
			slog.String("method", "POST"),
			slog.Any("error", err),
		)
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	transferedNewSubscription, err := transferCreateRequestToServiceDomain(newSubscription)
	if err != nil {
		api.logger.Error("failed to map create request to domain",
			slog.String("method", "POST"),
			slog.Any("error", err),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_INPUT",
				Message: err.Error(),
			},
		})
		return
	}

	createdSubscription, err := api.subscriptionService.CreateSubscription(c.Request.Context(), transferedNewSubscription)
	if err != nil {
		api.logger.Error("failed to create subscription in service",
			slog.String("method", "POST"),
			slog.Any("error", err),
		)
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	api.logger.Info("subscription created successfully",
		slog.String("method", "POST"),
		slog.String("subscription_id", createdSubscription.SubscriptionID.String()),
	)
	c.JSON(201, transferServiceDomainToAPIModel(createdSubscription))
}

func (api *SubscriptionAPI) SubscriptionReadGet(c *gin.Context) {
	idStr := c.Param("id")
	api.logger.Info("handling get subscription request",
		slog.String("method", "GET"),
		slog.String("path", fmt.Sprintf("/subscriptions/%s", idStr)),
		slog.String("subscription_id", idStr),
	)

	id, err := uuid.Parse(idStr)
	if err != nil {
		api.logger.Warn("invalid subscription ID format",
			slog.String("method", "GET"),
			slog.String("subscription_id", idStr),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_ID",
				Message: "invalid subscription ID format",
			},
		})
		return
	}

	subscription, err := api.subscriptionService.GetSubscriptionByID(c.Request.Context(), id)
	if err != nil {
		api.logger.Error("failed to get subscription",
			slog.String("method", "GET"),
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	api.logger.Debug("subscription retrieved successfully",
		slog.String("method", "GET"),
		slog.String("subscription_id", id.String()),
	)
	c.JSON(200, api_models.SubscriptionReadGet200Response{
		Subscription: transferServiceDomainToAPIModel(&subscription),
	})
}

func (api *SubscriptionAPI) SubscriptionUpdatePatch(c *gin.Context) {
	idStr := c.Param("id")
	api.logger.Info("handling patch subscription request",
		slog.String("method", "PATCH"),
		slog.String("path", fmt.Sprintf("/subscriptions/%s", idStr)),
		slog.String("subscription_id", idStr),
	)

	id, err := uuid.Parse(idStr)
	if err != nil {
		api.logger.Warn("invalid subscription ID format in patch",
			slog.String("method", "PATCH"),
			slog.String("subscription_id", idStr),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_ID",
				Message: "invalid subscription ID format",
			},
		})
		return
	}

	var newSubscription api_models.SubscriptionUpdatePutRequest

	if err := c.ShouldBindJSON(&newSubscription); err != nil {
		api.logger.Error("failed to bind patch request",
			slog.String("method", "PATCH"),
			slog.String("subscription_id", idStr),
			slog.Any("error", err),
		)
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	transferedNewSubscription, err := transferUpdatePutRequestToServiceDomain(newSubscription, id)
	if err != nil {
		api.logger.Error("failed to map patch request to domain",
			slog.String("method", "PATCH"),
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_INPUT",
				Message: err.Error(),
			},
		})
		return
	}

	updatedSubscription, err := api.subscriptionService.UpdateSubscriptionPatch(c.Request.Context(), id, &transferedNewSubscription)
	if err != nil {
		api.logger.Error("failed to patch subscription",
			slog.String("method", "PATCH"),
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	api.logger.Info("subscription patched successfully",
		slog.String("method", "PATCH"),
		slog.String("subscription_id", id.String()),
	)
	c.JSON(200, api_models.SubscriptionUpdatePut200Response{
		Subscription: transferServiceDomainToAPIModel(updatedSubscription),
	})
}

func (api *SubscriptionAPI) SubscriptionUpdatePut(c *gin.Context) {
	idStr := c.Param("id")
	api.logger.Info("handling put subscription request",
		slog.String("method", "PUT"),
		slog.String("path", fmt.Sprintf("/subscriptions/%s", idStr)),
		slog.String("subscription_id", idStr),
	)

	id, err := uuid.Parse(idStr)
	if err != nil {
		api.logger.Warn("invalid subscription ID format in put",
			slog.String("method", "PUT"),
			slog.String("subscription_id", idStr),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_ID",
				Message: "invalid subscription ID format",
			},
		})
		return
	}

	var newSubscription api_models.SubscriptionUpdatePutRequest

	if err := c.ShouldBindJSON(&newSubscription); err != nil {
		api.logger.Error("failed to bind put request",
			slog.String("method", "PUT"),
			slog.String("subscription_id", idStr),
			slog.Any("error", err),
		)
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	transferedNewSubscription, err := transferUpdatePutRequestToServiceDomain(newSubscription, id)
	if err != nil {
		api.logger.Error("failed to map put request to domain",
			slog.String("method", "PUT"),
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_INPUT",
				Message: err.Error(),
			},
		})
		return
	}

	updatedSubscription, err := api.subscriptionService.UpdateSubscriptionPut(c.Request.Context(), id, &transferedNewSubscription)
	if err != nil {
		api.logger.Error("failed to put subscription",
			slog.String("method", "PUT"),
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	api.logger.Info("subscription updated (PUT) successfully",
		slog.String("method", "PUT"),
		slog.String("subscription_id", id.String()),
	)
	c.JSON(200, api_models.SubscriptionUpdatePut200Response{
		Subscription: transferServiceDomainToAPIModel(updatedSubscription),
	})
}

func (api *SubscriptionAPI) SubscriptionDelete(c *gin.Context) {
	idStr := c.Param("id")
	api.logger.Info("handling delete subscription request",
		slog.String("method", "DELETE"),
		slog.String("path", fmt.Sprintf("/subscriptions/%s", idStr)),
		slog.String("subscription_id", idStr),
	)

	id, err := uuid.Parse(idStr)
	if err != nil {
		api.logger.Warn("invalid subscription ID format in delete",
			slog.String("method", "DELETE"),
			slog.String("subscription_id", idStr),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_ID",
				Message: "invalid subscription ID format",
			},
		})
		return
	}

	if err := api.subscriptionService.DeleteSubscriptionByID(c.Request.Context(), id); err != nil {
		api.logger.Error("failed to delete subscription",
			slog.String("method", "DELETE"),
			slog.String("subscription_id", id.String()),
			slog.Any("error", err),
		)
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	api.logger.Info("subscription deleted successfully",
		slog.String("method", "DELETE"),
		slog.String("subscription_id", id.String()),
	)
	c.Status(http.StatusNoContent)
}

func (api *SubscriptionAPI) SubscriptionListGet(c *gin.Context) {
	serviceNameStr := c.Query("service_name")
	userIDStr := c.Query("user_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		api.logger.Warn("invalid user ID format in list request",
			slog.String("method", "GET"),
			slog.String("user_id", userIDStr),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_ID",
				Message: "invalid user ID format",
			},
		})
		return
	}

	var startDate, endDate time.Time

	if startDate, err = time.Parse("01-2006", startDateStr); err != nil {
		api.logger.Warn("invalid start date format in list request",
			slog.String("method", "GET"),
			slog.String("start_date", startDateStr),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_START_DATE",
				Message: "invalid start date format",
			},
		})
		return
	}

	if endDate, err = time.Parse("01-2006", endDateStr); err != nil {
		api.logger.Warn("invalid end date format in list request",
			slog.String("method", "GET"),
			slog.String("end_date", endDateStr),
		)
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INVALID_END_DATE",
				Message: "invalid end date format",
			},
		})
		return
	}

	domainSubscriptionsList, err := api.subscriptionService.ListSubscriptions(c.Request.Context(), serviceNameStr, userID, startDate, endDate)

	if err != nil {
		api.logger.Error("failed to get subscriptions list",
			slog.String("method", "GET"),
			slog.String("service_name", serviceNameStr),
			slog.String("user_id", userIDStr),
			slog.String("start_date", startDateStr),
			slog.String("end_date", endDateStr),
			slog.Any("error", err),
		)
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	subscriptions := transferServiceDomainListToAPIModelList(domainSubscriptionsList)

	c.JSON(200, api_models.SubscriptionListGetResponse200{
		Subscriptions: subscriptions,
	})

}