package route

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/delivery/http/controller"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type RouteConfig struct {
	App                   *fiber.App
	AuthMiddleware        fiber.Handler
	RBACMiddleware        fiber.Handler
	RequestIDMiddleware   fiber.Handler
	ResponseMiddleware    fiber.Handler
	AuthController        *controller.AuthController
	MerchantController    *controller.MerchantController
	RBACController        *controller.RBACController
	AnalyticController    *controller.AnalyticController
	TransactionController *controller.TransactionController
}

func (c *RouteConfig) Setup() {
	config := cors.ConfigDefault
	config.AllowOrigins = strings.Join([]string{
		"http://localhost",
		"http://127.0.0.1",
	}, ",")
	config.AllowCredentials = true
	config.ExposeHeaders = "Content-Disposition"
	config.AllowMethods = strings.Join([]string{
		fiber.MethodGet,
		fiber.MethodHead,
		fiber.MethodPut,
		fiber.MethodPatch,
		fiber.MethodPost,
		fiber.MethodDelete,
		fiber.MethodOptions,
	}, ",")
	config.AllowHeaders = strings.Join([]string{
		fiber.HeaderOrigin,
		fiber.HeaderContentType,
		fiber.HeaderAccept,
		fiber.HeaderContentLength,
		fiber.HeaderAcceptEncoding,
		fiber.HeaderAuthorization,
	}, ",")
	c.App.Use(cors.New(config))

	// default route
	c.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"app":    c.App.Config().AppName,
			"status": "OK",
		})
	})

	c.App.Use(c.RequestIDMiddleware)
	c.App.Use(c.ResponseMiddleware)

	c.SetupGuestRoute()
	c.SetupRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	v1 := c.App.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.Post("/login", c.AuthController.Login)
			auth.Post("/refresh", c.AuthController.Refresh)
		}
	}
}

func (c *RouteConfig) SetupRoute() {
	// Protected routes: require valid JWT token (AuthMiddleware)
	v1 := c.App.Group("/api/v1", c.AuthMiddleware)
	{
		auth := v1.Group("/auth")
		{
			auth.Post("/logout", c.AuthController.Logout)
		}

		users := v1.Group("/users")
		{
			users.Get("/current", c.AuthController.Current)
		}

		// Merchant options: authenticated only, no RBAC — used for dropdown filters
		v1.Get("/merchants/options", c.MerchantController.GetOptions)

		// RBAC routes: require valid token and permission check (RBACMiddleware)
		rbac := v1.Group("", c.RBACMiddleware)
		{
			merchants := rbac.Group("/merchants")
			{
				merchants.Post("", c.MerchantController.Create)
				merchants.Get("", c.MerchantController.FindAll)
				merchants.Get("/:id", c.MerchantController.FindByID)
				merchants.Put("/:id", c.MerchantController.Update)
				merchants.Delete("/:id", c.MerchantController.Delete)
				merchants.Patch("/:id/regenerate-secret", c.MerchantController.RegenerateSecret)
			}

			roles := rbac.Group("/roles")
			{
				roles.Get("", c.RBACController.FindAllRoles)
				roles.Get("/:id", c.RBACController.FindRoleByID)
			}

			permissions := rbac.Group("/permissions")
			{
				permissions.Get("", c.RBACController.FindAllPermissions)
			}

			analytics := rbac.Group("/analytics")
			{
				analytics.Get("/summary", c.AnalyticController.GetSummary)
				analytics.Get("/trend", c.AnalyticController.GetDailyTrend)
				analytics.Get("/payment-methods", c.AnalyticController.GetPaymentMethodDistribution)
				analytics.Get("/status-distribution", c.AnalyticController.GetStatusDistribution)
			}

			transactions := rbac.Group("/transactions")
			{
				transactions.Get("", c.TransactionController.FindAll)
				transactions.Get("/:id", c.TransactionController.FindByID)
			}
		}
	}
}
