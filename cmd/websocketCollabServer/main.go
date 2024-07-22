package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/masterkusok/websocketCollab/internal/api"
	"github.com/masterkusok/websocketCollab/internal/businnesLogic"
	"github.com/masterkusok/websocketCollab/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func getDsn() string {
	host := os.Getenv("SQL_HOST")
	port := os.Getenv("SQL_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		host, user, password, db, port)
}

func main() {
	// this connection can be used to connect to postgres inside docker container.
	// db, err := gorm.Open(postgres.Open(getDsn()), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	// this connection uses sqlite
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&businnesLogic.Document{})
	if err != nil {
		log.Fatal(err)
	}
	handler := api.NewHandler(repositories.NewDocumentRepository(db))
	app := fiber.New()
	api.CreateRouting(app, handler)

	log.Fatal(app.Listen(":8080"))
}
