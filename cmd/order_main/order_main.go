package main

import (
	"database/sql"
	"fmt"
	_notificationsRepo "github.com/2020_1_Skycode/internal/notifications/repository"
	protobuf_order "github.com/2020_1_Skycode/internal/orders/delivery/protobuf"
	_orderRepo "github.com/2020_1_Skycode/internal/orders/repository"
	_restRepo "github.com/2020_1_Skycode/internal/restaurants/repository"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"

	_ "github.com/lib/pq"
)

func main() {
	config, err := tools.LoadConf()

	if err != nil {
		log.Fatal(err)
	}

	connString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s sslmode=disable password=%s",
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

	restRepo := _restRepo.NewRestaurantRepository(dbConn)
	orderRepo := _orderRepo.NewOrdersRepository(dbConn, restRepo)
	notificationsRepo := _notificationsRepo.NewNotificationsRepository(dbConn)

	orderManager := protobuf_order.NewOrderProtoManager(orderRepo, notificationsRepo)

	port := ":5003"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Error("Cant listen port: ", err)
	}

	server := grpc.NewServer()

	protobuf_order.RegisterOrderWorkerServer(server, orderManager)

	logrus.Info("Starting server on port", port)
	if err := server.Serve(lis); err != nil {
		logrus.Error(err)
	}
}
