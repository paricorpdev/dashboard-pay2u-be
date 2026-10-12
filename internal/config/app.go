package config

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/delivery/http/controller"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/delivery/http/middleware"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/delivery/http/route"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	Redis    *redis.Storage
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	// Repository App
	redisRepository := repository.NewRedisRepository(config.Redis, config.Log)

	// Repository
	userRepository := repository.NewUserRepository(config.Log)
	permissionRepository := repository.NewPermissionRepository(config.Log)
	merchantRepository := repository.NewMerchantRepository(config.Log)
	roleRepository := repository.NewRoleRepository(config.Log)
	transactionRepository := repository.NewTransactionRepository(config.Log)

	// Use Case
	jwtUseCase := usecase.NewJwtUseCase(config.Config, config.Log, config.Validate, redisRepository)
	authUseCase := usecase.NewAuthUseCase(config.DB, config.Config, config.Log, config.Validate, jwtUseCase, userRepository, redisRepository)
	merchantUseCase := usecase.NewMerchantUseCase(config.DB, config.Log, config.Validate, merchantRepository)
	rbacUseCase := usecase.NewRBACUseCase(config.DB, config.Log, roleRepository, permissionRepository)
	analyticUseCase := usecase.NewAnalyticUseCase(config.DB, config.Log, transactionRepository, config.Config.GetString("app.timezone"))

	// Controller
	authController := controller.NewAuthController(authUseCase, config.Log)
	merchantController := controller.NewMerchantController(merchantUseCase, config.Log)
	rbacController := controller.NewRBACController(rbacUseCase, config.Log)
	analyticController := controller.NewAnalyticController(analyticUseCase, config.Log)

	// Middleware
	authMiddleware := middleware.NewAuth(config.Log, config.Config, jwtUseCase)
	rbacMiddleware := middleware.NewRBACMiddleware(config.DB, config.Log, permissionRepository)
	requestIDMiddleware := middleware.NewRequestIDMiddleware(config.Log)
	responseMiddleware := middleware.NewResponseMiddleware()

	routeConfig := route.RouteConfig{
		App:                 config.App,
		AuthMiddleware:      authMiddleware,
		RBACMiddleware:      rbacMiddleware,
		RequestIDMiddleware: requestIDMiddleware,
		ResponseMiddleware:  responseMiddleware,
		AuthController:      authController,
		MerchantController:  merchantController,
		RBACController:      rbacController,
		AnalyticController:  analyticController,
	}
	routeConfig.Setup()
}
