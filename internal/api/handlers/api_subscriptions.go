package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	api_models "github.com/kgugunava/effective_mobile_golang/internal/api/models"
)

type SubscriptionAPI struct {
	subscriptionService SubscriptionService
}

func NewSubscriptionAPI(service SubscriptionService) *SubscriptionAPI {
	return &SubscriptionAPI{
		subscriptionService: service,
	}
}

func (api *SubscriptionAPI) SubscriptionCreatePost(c *gin.Context) {
	var newSubscription api_models.SubscriptionCreatePostRequest

	if err := c.ShouldBindJSON(&newSubscription); err != nil {
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code: "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}


	transferedNewSubscription, err := transferCreateRequestToServiceDomain(newSubscription)
	if err != nil {
		fmt.Println("here")
	}
	createdSubscription, err := api.subscriptionService.CreateSubscription(c.Request.Context(), transferedNewSubscription)
	if err != nil {
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code: "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(201, transferServiceDomainToAPIModel(createdSubscription))

}

func (api *SubscriptionAPI) SubscriptionReadGet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code: "INVALID_ID",
				Message: "invalid subscription ID format",
			},
		})
		return
	}

	subscription, err := api.subscriptionService.GetSubscriptionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, api_models.ErrorResponse{
			Error: api_models.ErrorResponseError{
				Code: "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(200, api_models.SubscriptionReadGet200Response{
		Subscription: transferServiceDomainToAPIModel(&subscription),
	})
}