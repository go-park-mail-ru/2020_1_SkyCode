package main

import (
	_middleware "github.com/2020_1_Skycode/internal/middlewares"
	_sessionsDelivery "github.com/2020_1_Skycode/internal/sessions/delivery"
	_sessionsRepository "github.com/2020_1_Skycode/internal/sessions/repository"
	_sessionsUseCase "github.com/2020_1_Skycode/internal/sessions/usecase"
	"github.com/2020_1_Skycode/internal/tools"
	_usersDelivery "github.com/2020_1_Skycode/internal/users/delivery"
	_usersRepository "github.com/2020_1_Skycode/internal/users/repository"
	_usersUseCase "github.com/2020_1_Skycode/internal/users/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"log"
)

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
		Password: "",
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

	userRepo := _usersRepository.NewUserRepository(dbConn)
	userUcase := _usersUseCase.NewUserUseCase(userRepo)

	sessionsRepo := _sessionsRepository.NewSessionRepository(dbConn)
	sessionsUcase := _sessionsUseCase.NewSessionUseCase(sessionsRepo)

	mwareC := _middleware.NewMiddleWareController(e, sessionsUcase)


	_ = _middleware.NewMiddleWareController(e, sessionsUcase)
	_ = _sessionsDelivery.NewSessionHandler(e, sessionsUcase, userUcase, mwareC)
	_ = _usersDelivery.NewUserHandler(e, userUcase)

	log.Fatal(e.Run())
}
