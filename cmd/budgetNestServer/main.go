package main

import (
	"budgetNest/database"
	"budgetNest/internal/helpers"
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"sync"
)

type DatabaseSession struct {
	db *sql.DB
}

var once sync.Once
var session *DatabaseSession

func CreateDatabaseSession() *DatabaseSession {
	once.Do(func() {
		config := mysql.Config{
			User:                 os.Getenv("DB_USER"),
			Passwd:               os.Getenv("DB_PASS"),
			DBName:               os.Getenv("DB_NAME"),
			AllowNativePasswords: false,
		}
		db, err := sql.Open("mysql", config.FormatDSN())
		helpers.CheckFatal(err, "Error connecting to database")
		session = &DatabaseSession{db: db}
	})
	return session
}

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

	db := CreateDatabaseSession().db

	if *migrateDown {
		database.RunMigrations(db, false)
		return
	}

	if *migrate {
		database.RunMigrations(db, true)
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
