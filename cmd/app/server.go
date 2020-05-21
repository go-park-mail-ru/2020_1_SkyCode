package main

import (
	"database/sql"
	"fmt"
	_ "github.com/2020_1_Skycode/docs"
	_chatsDelivery "github.com/2020_1_Skycode/internal/chats/delivery"
	_chatsRepository "github.com/2020_1_Skycode/internal/chats/repository"
	_chatsUseCase "github.com/2020_1_Skycode/internal/chats/usecase"
	_geodataDelivery "github.com/2020_1_Skycode/internal/geodata/delivery"
	_geodataRepository "github.com/2020_1_Skycode/internal/geodata/repository"
	_geodataUseCase "github.com/2020_1_Skycode/internal/geodata/usecase"
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	_notificationsDelivery "github.com/2020_1_Skycode/internal/notifications/delivery"
	_notificationsRepository "github.com/2020_1_Skycode/internal/notifications/repository"
	_notificationsUseCase "github.com/2020_1_Skycode/internal/notifications/usecase"
	_ordersDelivery "github.com/2020_1_Skycode/internal/orders/delivery"
	_orderRepo "github.com/2020_1_Skycode/internal/orders/repository"
	_ordersUseCase "github.com/2020_1_Skycode/internal/orders/usecase"
	_prodTagsRepository "github.com/2020_1_Skycode/internal/product_tags/repository"
	_productDelivery "github.com/2020_1_Skycode/internal/products/delivery"
	_productRepo "github.com/2020_1_Skycode/internal/products/repository"
	_productUseCase "github.com/2020_1_Skycode/internal/products/usecase"
	_restPointsDelivery "github.com/2020_1_Skycode/internal/restaurant_points/delivery"
	_restPointsRepository "github.com/2020_1_Skycode/internal/restaurant_points/repository"
	_restPointsUseCase "github.com/2020_1_Skycode/internal/restaurant_points/usecase"
	_restDelivery "github.com/2020_1_Skycode/internal/restaurants/delivery"
	_restRepo "github.com/2020_1_Skycode/internal/restaurants/repository"
	_restUcase "github.com/2020_1_Skycode/internal/restaurants/usecase"
	_restTagsDelivery "github.com/2020_1_Skycode/internal/restaurants_tags/delivery"
	_restTagsRepository "github.com/2020_1_Skycode/internal/restaurants_tags/repository"
	_restTagsUseCase "github.com/2020_1_Skycode/internal/restaurants_tags/usecase"
	_reviewsDelivery "github.com/2020_1_Skycode/internal/reviews/delivery"
	_reviewsRepository "github.com/2020_1_Skycode/internal/reviews/repository"
	_reviewsUseCase "github.com/2020_1_Skycode/internal/reviews/usecase"
	_sessionsDelivery "github.com/2020_1_Skycode/internal/sessions/delivery"
	_sessionsUseCase "github.com/2020_1_Skycode/internal/sessions/usecase"
	"github.com/2020_1_Skycode/internal/tools"
	_csrfManager "github.com/2020_1_Skycode/internal/tools/CSRFManager"
	"github.com/2020_1_Skycode/internal/tools/notificationsWS"
	_rValidator "github.com/2020_1_Skycode/internal/tools/requestValidator"
	_usersDelivery "github.com/2020_1_Skycode/internal/users/delivery"
	_usersRepository "github.com/2020_1_Skycode/internal/users/repository"
	_usersUseCase "github.com/2020_1_Skycode/internal/users/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"log"
	"os"
)

// @title Swagger SkyDelivery API
// @version 1.0
// @description This is a SkyDelivery server for Technopark Project.

// @contact.name API Support

// @host localhost:5000
// @BasePath /api/v1
func main() {
	config, err := tools.LoadConf()

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

	grpcSessionConn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer grpcSessionConn.Close()

	grpcAdminConn, err := grpc.Dial("localhost:5002", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer grpcAdminConn.Close()

	grpcOrderConn, err := grpc.Dial("localhost:5003", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer grpcOrderConn.Close()

	f, err := os.OpenFile("skydelivery.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	logrus.SetOutput(f)

	e := gin.New()

	geoCoderKey := config.ApiKeys.YandexGeoCoder

	prodRepo := _productRepo.NewProductRepository(dbConn)

	prodTagsRepo := _prodTagsRepository.NewProductTagsRepository(dbConn)

	prodUcase := _productUseCase.NewProductWithProtoUseCase(prodRepo, prodTagsRepo, grpcAdminConn)

	reviewRepo := _reviewsRepository.NewReviewsRepository(dbConn)
	reviewUcase := _reviewsUseCase.NewReviewsUseCase(reviewRepo)

	geoDataRepo := _geodataRepository.NewGeoDataRepository(geoCoderKey)
	geoDataUcase := _geodataUseCase.NewGeoDataUseCase(geoDataRepo)

	restPointsRepo := _restPointsRepository.NewRestPosintsRepository(dbConn)
	restPointsUCase := _restPointsUseCase.NewRestPointsWithProtoUseCase(restPointsRepo, grpcAdminConn)

	restTagsRepo := _restTagsRepository.NewRestTagRepository(dbConn)
	restTagsUcase := _restTagsUseCase.NewRestTagsUCase(restTagsRepo)

	restRepo := _restRepo.NewRestaurantRepository(dbConn)

	ordersRepo := _orderRepo.NewOrdersRepository(dbConn, restRepo)

	restUcase := _restUcase.NewRestaurantsWithProtoUseCase(restRepo, restPointsRepo, reviewRepo,
		geoDataRepo, restTagsRepo, prodTagsRepo, ordersRepo, grpcAdminConn)

	userRepo := _usersRepository.NewUserRepository(dbConn)
	userUcase := _usersUseCase.NewUserUseCase(userRepo)

	sessionsUcase := _sessionsUseCase.NewSessionProtoUseCase(grpcSessionConn)

	ordersUcase := _ordersUseCase.NewOrderProtoUseCase(grpcOrderConn)

	chatsRepo := _chatsRepository.NewChatsRepository(dbConn)
	chatsUcase := _chatsUseCase.NewChatUseCase(chatsRepo)

	notificationsRepo := _notificationsRepository.NewNotificationsRepository(dbConn)
	notificationsUCase := _notificationsUseCase.NewNotificationsUseCase(notificationsRepo)
	notificationsServer := notificationsWS.NewNotificationServer()

	csrfManager := _csrfManager.NewCSRFManager()

	reqValidator := _rValidator.NewRequestValidator()

	_ = _middleware.NewMetricsController(e)
	mwareC := _middleware.NewMiddleWareController(e, sessionsUcase, userUcase, csrfManager)

	publicGroup := e.Group("/api/v1")
	privateGroup := e.Group("/api/v1")

	privateGroup.Use(mwareC.CSRFControl())

	_ = _sessionsDelivery.NewSessionHandler(privateGroup, publicGroup, sessionsUcase, userUcase, reqValidator,
		csrfManager, mwareC)
	_ = _usersDelivery.NewUserHandler(privateGroup, publicGroup, userUcase, sessionsUcase, reqValidator,
		csrfManager, mwareC)
	_ = _restDelivery.NewRestaurantHandler(privateGroup, publicGroup, reqValidator, restUcase, mwareC)
	_ = _productDelivery.NewProductHandler(privateGroup, publicGroup, prodUcase, reqValidator, restUcase, mwareC)
	_ = _ordersDelivery.NewOrderHandler(privateGroup, publicGroup, ordersUcase, reqValidator, mwareC,
		notificationsServer)
	_ = _reviewsDelivery.NewReviewsHandler(privateGroup, publicGroup, reviewUcase, reqValidator, mwareC)
	_ = _restPointsDelivery.NewRestPointsHandler(privateGroup, publicGroup, restPointsUCase, mwareC)
	_ = _geodataDelivery.NewGeoDataHandler(privateGroup, publicGroup, geoDataUcase)
	_ = _chatsDelivery.NewChatsHandler(privateGroup, publicGroup, chatsUcase, mwareC)
	_ = _restTagsDelivery.NewRestTagHandler(privateGroup, publicGroup, restTagsUcase, reqValidator, mwareC)
	_ = _notificationsDelivery.NewNotificationsHandler(privateGroup, publicGroup, notificationsUCase, mwareC,
		notificationsServer)

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(e.Run(":5000"))
}
