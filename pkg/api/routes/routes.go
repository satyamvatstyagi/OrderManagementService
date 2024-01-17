package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/api/controller"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/config"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/repository"
	"github.com/satyamvatstyagi/OrderManagementService/pkg/app/usecase"
	"go.elastic.co/apm/module/apmgin/v2"
)

func Setup() {

	cfg := config.Config{}

	// Initialize the database
	db := cfg.InitDb()

	// Initialize the repositories
	OrderRepository := repository.NewOrderRepository(db)

	// Initialize the usecases
	OrderUsecase := usecase.NewOrderUsecase(OrderRepository)

	// Initialize the controller
	OrderController := &controller.OrderController{OrderUsecase: OrderUsecase}

	ginMode := os.Getenv("GIN_MODE")
	gin.SetMode(ginMode)
	router := gin.Default()
	router.Use(apmgin.Middleware(router))

	// Setup the routes
	setupOrderRoutes(OrderController, router)

	err := router.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}

func setupOrderRoutes(c *controller.OrderController, router *gin.Engine) {
	username := os.Getenv("BASIC_AUTH_USER")
	password := os.Getenv("BASIC_AUTH_PASSWORD")
	orderService := router.Group("/order", gin.BasicAuth(gin.Accounts{username: password}))
	{
		orderService.POST("/", c.CreateOrder)
		orderService.GET("/:username", c.GetOrderByOrderUserName)
	}
}
