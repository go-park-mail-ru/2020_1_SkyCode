package main

import (
	"database/sql"
	"fmt"
	_ "github.com/2020_1_Skycode/docs"
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
	_usersDelivery "github.com/2020_1_Skycode/internal/users/delivery"
	_usersRepository "github.com/2020_1_Skycode/internal/users/repository"
	_usersUseCase "github.com/2020_1_Skycode/internal/users/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

type DatabaseInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

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

	connString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s sslmode=disable password=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.User,
		config.Database.Password)

	dbConn, err := sql.Open("postgres", connString)

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
	restUcase := _restUcase.NewRestaurantsUseCase(restRepo)

	userRepo := _usersRepository.NewUserRepository(dbConn)
	userUcase := _usersUseCase.NewUserUseCase(userRepo)

	sessionsRepo := _sessionsRepository.NewSessionRepository(dbConn)
	sessionsUcase := _sessionsUseCase.NewSessionUseCase(sessionsRepo)

	ordersRepo := _ordersRepository.NewOrdersRepository(dbConn)
	ordersUcase := _ordersUseCase.NewOrderUseCase(ordersRepo)

	mwareC := _middleware.NewMiddleWareController(e, sessionsUcase, userUcase)

	_ = _middleware.NewMiddleWareController(e, sessionsUcase, userUcase)
	_ = _sessionsDelivery.NewSessionHandler(e, sessionsUcase, userUcase, mwareC)
	_ = _usersDelivery.NewUserHandler(e, userUcase, sessionsUcase, mwareC)
	_ = _restDelivery.NewRestaurantHandler(e, restUcase)
	_ = _productDelivery.NewProductHandler(e, prodUcase, restUcase, mwareC)
	_ = _ordersDelivery.NewOrderHandler(e, ordersUcase, mwareC)

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(e.Run(":5000"))
}
