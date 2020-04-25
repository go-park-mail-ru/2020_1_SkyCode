package main

import (
	"database/sql"
	"fmt"
	_ "github.com/2020_1_Skycode/docs"
	_geodataRepository "github.com/2020_1_Skycode/internal/geodata/repository"
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	_ordersDelivery "github.com/2020_1_Skycode/internal/orders/delivery"
	_ordersRepository "github.com/2020_1_Skycode/internal/orders/repository"
	_ordersUseCase "github.com/2020_1_Skycode/internal/orders/usecase"
	_productDelivery "github.com/2020_1_Skycode/internal/products/delivery"
	_productRepo "github.com/2020_1_Skycode/internal/products/repository"
	_productUseCase "github.com/2020_1_Skycode/internal/products/usecase"
	_restPointsDelivery "github.com/2020_1_Skycode/internal/restaurant_points/delivery"
	_restPointsRepository "github.com/2020_1_Skycode/internal/restaurant_points/repository"
	_restPointsUseCase "github.com/2020_1_Skycode/internal/restaurant_points/usecase"
	_restDelivery "github.com/2020_1_Skycode/internal/restaurants/delivery"
	_restRepo "github.com/2020_1_Skycode/internal/restaurants/repository"
	_restUcase "github.com/2020_1_Skycode/internal/restaurants/usecase"
	_reviewsDelivery "github.com/2020_1_Skycode/internal/reviews/delivery"
	_reviewsRepository "github.com/2020_1_Skycode/internal/reviews/repository"
	_reviewsUseCase "github.com/2020_1_Skycode/internal/reviews/usecase"
	_sessionsDelivery "github.com/2020_1_Skycode/internal/sessions/delivery"
	_sessionsRepository "github.com/2020_1_Skycode/internal/sessions/repository"
	_sessionsUseCase "github.com/2020_1_Skycode/internal/sessions/usecase"
	"github.com/2020_1_Skycode/internal/tools"
	_csrfManager "github.com/2020_1_Skycode/internal/tools/CSRFManager"
	_rValidator "github.com/2020_1_Skycode/internal/tools/requestValidator"
	_usersDelivery "github.com/2020_1_Skycode/internal/users/delivery"
	_usersRepository "github.com/2020_1_Skycode/internal/users/repository"
	_usersUseCase "github.com/2020_1_Skycode/internal/users/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Swagger SkyDelivery API
// @version 1.0
// @description This is a SkyDelivery server for Technopark Project.

// @contact.name API Support

// @host localhost:5000
// @BasePath /api/v1
func main() {
	config, err := tools.LoadConf("./configs/config.json")

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

	geoCoderKey := config.ApiKeys.YandexGeoCoder

	prodRepo := _productRepo.NewProductRepository(dbConn)
	prodUcase := _productUseCase.NewProductUseCase(prodRepo)

	reviewRepo := _reviewsRepository.NewReviewsRepository(dbConn)
	reviewUcase := _reviewsUseCase.NewReviewsUseCase(reviewRepo)

	geoDataRepo := _geodataRepository.NewGeoDataRepository(geoCoderKey)

	restPointsRepo := _restPointsRepository.NewRestPosintsRepository(dbConn)
	restPointsUCase := _restPointsUseCase.NewRestPointsUseCase(restPointsRepo)

	restRepo := _restRepo.NewRestaurantRepository(dbConn)
	restUcase := _restUcase.NewRestaurantsUseCase(restRepo, restPointsRepo, reviewRepo, geoDataRepo)

	userRepo := _usersRepository.NewUserRepository(dbConn)
	userUcase := _usersUseCase.NewUserUseCase(userRepo)

	sessionsRepo := _sessionsRepository.NewSessionRepository(dbConn)
	sessionsUcase := _sessionsUseCase.NewSessionUseCase(sessionsRepo)

	ordersRepo := _ordersRepository.NewOrdersRepository(dbConn, restRepo)
	ordersUcase := _ordersUseCase.NewOrderUseCase(ordersRepo)

	csrfManager := _csrfManager.NewCSRFManager()

	reqValidator := _rValidator.NewRequestValidator()

	mwareC := _middleware.NewMiddleWareController(e, sessionsUcase, userUcase, csrfManager)

	publicGroup := e.Group("/api/v1")
	privateGroup := e.Group("/api/v1")

	privateGroup.Use(mwareC.CSRFControl())

	_ = _sessionsDelivery.NewSessionHandler(privateGroup, publicGroup, sessionsUcase, userUcase, reqValidator, csrfManager, mwareC)
	_ = _usersDelivery.NewUserHandler(privateGroup, publicGroup, userUcase, sessionsUcase, reqValidator, mwareC)
	_ = _restDelivery.NewRestaurantHandler(privateGroup, publicGroup, reqValidator, restUcase, mwareC)
	_ = _productDelivery.NewProductHandler(privateGroup, publicGroup, prodUcase, reqValidator, restUcase, mwareC)
	_ = _ordersDelivery.NewOrderHandler(privateGroup, publicGroup, ordersUcase, reqValidator, mwareC)
	_ = _reviewsDelivery.NewReviewsHandler(privateGroup, publicGroup, reviewUcase, reqValidator, mwareC)
	_ = _restPointsDelivery.NewRestPointsHandler(privateGroup, publicGroup, restPointsUCase, mwareC)

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(e.Run(":5000"))
}
