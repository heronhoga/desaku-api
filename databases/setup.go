package databases

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	errLoad := godotenv.Load()
	if errLoad != nil {
		panic("Failed to load ENV")
	}

	dbname := os.Getenv("DATABASE_NAME")
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_URL")
	port := os.Getenv("DATABASE_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	database, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection Opened to Database")

	DB = database
}
