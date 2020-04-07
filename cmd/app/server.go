package main

import (
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	_ordersDelivery "github.com/2020_1_Skycode/internal/orders/delivery"
	_ordersRepository "github.com/2020_1_Skycode/internal/orders/repository"
	_ordersUseCase "github.com/2020_1_Skycode/internal/orders/usecase"
	_productDelivery "github.com/2020_1_Skycode/internal/products/delivery"
	_productRepo "github.com/2020_1_Skycode/internal/products/repository"
	_productUseCase "github.com/2020_1_Skycode/internal/products/usecase"
	_restDelivery "github.com/2020_1_Skycode/internal/restaurants/delivery"
	_restRepo "github.com/2020_1_Skycode/internal/restaurants/repository"
	_restUcase "github.com/2020_1_Skycode/internal/restaurants/usecase"
	_sessionsDelivery "github.com/2020_1_Skycode/internal/sessions/delivery"
	_sessionsRepository "github.com/2020_1_Skycode/internal/sessions/repository"
	_sessionsUseCase "github.com/2020_1_Skycode/internal/sessions/usecase"
	"github.com/2020_1_Skycode/internal/tools"
	_csrfManager "github.com/2020_1_Skycode/internal/tools/CSRFManager"
	_usersDelivery "github.com/2020_1_Skycode/internal/users/delivery"
	_usersRepository "github.com/2020_1_Skycode/internal/users/repository"
	_usersUseCase "github.com/2020_1_Skycode/internal/users/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"

	_ "github.com/2020_1_Skycode/docs"
)

// @title Swagger SkyDelivery API
// @version 1.0
// @description This is a SkyDelivery server for Technopark Project.

// @contact.name API Support

// @host localhost:5000
// @BasePath /api/v1
func main() {
	config, err := tools.LoadConf("../../configs/config.json")

	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := pgx.Connect(pgx.ConnConfig{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Database: config.Database.Name,
		User:     config.Database.User,
		Password: config.Database.Password,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	e := gin.New()

	prodRepo := _productRepo.NewProductRepository(dbConn)
	prodUcase := _productUseCase.NewProductUseCase(prodRepo)

	restRepo := _restRepo.NewRestaurantRepository(dbConn)
	restUcase := _restUcase.NewRestaurantsUseCase(restRepo, prodRepo)

	userRepo := _usersRepository.NewUserRepository(dbConn)
	userUcase := _usersUseCase.NewUserUseCase(userRepo)

	sessionsRepo := _sessionsRepository.NewSessionRepository(dbConn)
	sessionsUcase := _sessionsUseCase.NewSessionUseCase(sessionsRepo)

	ordersRepo := _ordersRepository.NewOrdersRepository(dbConn)
	ordersUcase := _ordersUseCase.NewOrderUseCase(ordersRepo)

	privateGroup := e.Group("/api/v1")
	publicGroup := e.Group("/api/v1")

	csrfManager := _csrfManager.NewCSRFManager()

	mwareC := _middleware.NewMiddleWareController(privateGroup, publicGroup, sessionsUcase, userUcase, csrfManager)

	_ = _sessionsDelivery.NewSessionHandler(privateGroup, publicGroup, sessionsUcase, userUcase, csrfManager, mwareC)
	_ = _usersDelivery.NewUserHandler(privateGroup, publicGroup, userUcase, sessionsUcase, mwareC)
	_ = _restDelivery.NewRestaurantHandler(privateGroup, publicGroup, restUcase)
	_ = _productDelivery.NewProductHandler(privateGroup, publicGroup, prodUcase, restUcase, mwareC)
	_ = _ordersDelivery.NewOrderHandler(privateGroup, publicGroup, ordersUcase, mwareC)

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(e.Run(":5000"))
}
