package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBClient *gorm.DB

func init() {
	loadDBClient(true)
}
func loadDBClient(fillDB bool) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("db_host"),
		os.Getenv("db_user"),
		os.Getenv("db_password"),
		os.Getenv("db_name"),
		os.Getenv("db_port"),
	)

	log.Println(dsn)

	for {
		log.Println("Connecting to DB...")
		db, err := gorm.Open(postgres.Open(dsn))
		if err != nil {
			log.Println("Connection failled, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		DBClient = db
		break
	}
	log.Println("Connexion to DB successful")
	return nil
}

func GetDBClient() *gorm.DB {
	db, err := DBClient.DB()
	if err != nil {
		return nil
	}
	if db.Ping() != nil {
		loadDBClient(false)
		return DBClient
	} else {
		return DBClient
	}
}
