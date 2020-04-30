package main

import (
	"database/sql"
	"fmt"
	_geodataRepository "github.com/2020_1_Skycode/internal/geodata/repository"
	_productRepo "github.com/2020_1_Skycode/internal/products/repository"
	_restPointsRepository "github.com/2020_1_Skycode/internal/restaurant_points/repository"
	protobuf_admin_rest "github.com/2020_1_Skycode/internal/restaurants/delivery/protobuf"
	_restRepo "github.com/2020_1_Skycode/internal/restaurants/repository"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/tools/protobuf/adminwork"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"

	_ "github.com/lib/pq"
)

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

	geoCoderKey := config.ApiKeys.YandexGeoCoder

	restRepo := _restRepo.NewRestaurantRepository(dbConn)
	prodRepo := _productRepo.NewProductRepository(dbConn)
	restPointRepo := _restPointsRepository.NewRestPosintsRepository(dbConn)
	geoDataRepo := _geodataRepository.NewGeoDataRepository(geoCoderKey)

	adminManager := protobuf_admin_rest.NewAdminRestaurantManager(restRepo, prodRepo, restPointRepo, geoDataRepo)

	port := ":5002"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Error("Cant listen port: ", err)
	}

	server := grpc.NewServer()

	adminwork.RegisterRestaurantAdminWorkerServer(server, adminManager)

	logrus.Info("Starting server on port", port)
	server.Serve(lis)
}
