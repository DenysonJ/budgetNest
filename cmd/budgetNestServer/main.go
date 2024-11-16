package main

import (
	"budgetNest/database"
	"budgetNest/internal/helpers"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	migrate := flag.Bool("migrate", false, "Run database migrations")
	migrateDown := flag.Bool("migrate-down", false, "Run database migrations down")
	server := flag.Bool("server", false, "Run the server")

	flag.Parse()

	if *migrate && *migrateDown {
		fmt.Println("Please specify only one option: -migrate or -migrate-down")
		flag.Usage()
		return
	}

	err := godotenv.Load()
	helpers.CheckFatal(err, "Error loading .env file")

	if *migrateDown {
		database.RunMigrations(false)
		return
	}

	if *migrate {
		database.RunMigrations(true)
		return
	}

	if *server {
		router := gin.Default()
		err := router.Run("localhost:8080")
		helpers.CheckFatal(err, "Error running server")
		return
	}

	fmt.Println("Please specify an option: -migrate or -server")
	flag.Usage()
}
