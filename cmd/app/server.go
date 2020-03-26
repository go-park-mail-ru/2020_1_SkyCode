package main

import (
	"database/sql"
	"fmt"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/users/delivery"
	"github.com/2020_1_Skycode/internal/users/repository"
	"github.com/2020_1_Skycode/internal/users/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	path, err := os.Getwd()

	fmt.Println(path)

	config, err := tools.LoadConf("../../configs/config.json")

	if err != nil {
		log.Fatal(err)
	}

	dataSourceName := "host=%s dbname=%s dbname=%s sslmode=disable"

	dbConn, err := sql.Open("postgres", fmt.Sprintf(dataSourceName,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	))

	if err != nil {
		log.Fatal(err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	e := gin.New()

	userRepo := repository.NewUserRepository(dbConn)
	userUcase := usecase.NewUserUseCase(userRepo)
	_ = delivery.NewUserHandler(e, userUcase)

	log.Fatal(e.Run())
}
