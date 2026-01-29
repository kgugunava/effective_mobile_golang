package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kgugunava/effective_mobile_golang/internal/api/handlers"
)

type Route struct {
	// Name is the name of this Route.
	Name		string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method		string
	// Pattern is the pattern of the URI.
	Pattern	 	string
	// HandlerFunc is the handler function of this route.
	HandlerFunc	gin.HandlerFunc
}

func NewRouter(apiHandler handlers.SubscriptionAPI) *gin.Engine {
	return NewRouterWithGinEngine(gin.Default(), apiHandler)
}

func NewRouterWithGinEngine(router *gin.Engine, apiHandler handlers.SubscriptionAPI) *gin.Engine {
	for _, route := range getRoutes(apiHandler) {
		if route.HandlerFunc == nil {
			route.HandlerFunc = DefaultHandleFunc
		}
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

func DefaultHandleFunc(c *gin.Context) {
	c.String(http.StatusNotImplemented, "501 not implemented")
}

func getRoutes(apiHandler handlers.SubscriptionAPI) []Route {
	return []Route{ 
		{
			"SubscriptionCreatePost",
			http.MethodPost,
			"/create",
			apiHandler.SubscriptionCreatePost,
		},
		{
			"SubscriptionReadGet",
			http.MethodGet,
			"/read/:id",
			apiHandler.SubscriptionReadGet,
		},
		{
			"SubscriptionUpdatePut",
			http.MethodPut,
			"/update_put/:id",
			apiHandler.SubscriptionUpdatePut,
		},
		{
			"SubscriptionUpdatePatch",
			http.MethodPatch,
			"/update_patch/:id",
			apiHandler.SubscriptionUpdatePatch,
		},
		{
			"SubscriptionDelete",
			http.MethodDelete,
			"/delete/:id",
			apiHandler.SubscriptionDelete,
		},
		{
			"SubscriptionsListGet",
			http.MethodGet,
			"/subscriptions_list/",
			apiHandler.SubscriptionListGet,
		},
	}
}
