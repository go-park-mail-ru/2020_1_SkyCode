package main

import (
	"database/sql"
	"fmt"
	protobuf_session "github.com/2020_1_Skycode/internal/sessions/delivery/protobuf"
	_sessionRepository "github.com/2020_1_Skycode/internal/sessions/repository"
	"github.com/2020_1_Skycode/internal/tools"
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

	sessionRepo := _sessionRepository.NewSessionRepository(dbConn)
	sessionManager := protobuf_session.NewSessionManager(sessionRepo)

	port := ":5001"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Error("Cant listen port: ", err)
	}

	server := grpc.NewServer()

	protobuf_session.RegisterSessionWorkerServer(server, sessionManager)

	logrus.Info("Startin server on port", port)
	server.Serve(lis)
}
