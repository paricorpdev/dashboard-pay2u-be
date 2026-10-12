package route

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/delivery/http/controller"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type RouteConfig struct {
	App                 *fiber.App
	AuthMiddleware      fiber.Handler
	RBACMiddleware      fiber.Handler
	RequestIDMiddleware fiber.Handler
	ResponseMiddleware  fiber.Handler
	AuthController      *controller.AuthController
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
	// Auth middleware: verify JWT
	c.App.Use(c.AuthMiddleware)
	// RBAC middleware: derive resource+action from path, check DB
	c.App.Use(c.RBACMiddleware)

	v1 := c.App.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.Post("/logout", c.AuthController.Logout)
		}

		users := v1.Group("/users")
		{
			users.Get("/current", c.AuthController.Current)
		}
	}
}
