package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	userStorage "github.com/francotraversa/Sliceflow/internal/database/user_utils"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseInstance struct {
	DB *gorm.DB
}

var DBInstance = DatabaseInstance{}

func (database DatabaseInstance) NewDataBase() {
	var db *gorm.DB
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	ssl := os.Getenv("POSTGRES_SSL")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, pass, dbName, port, ssl,
	)
	for {
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Println("Waiting for DB...", err)
		time.Sleep(2 * time.Second)
	}
	log.Println("DB CONNECTED")

	db.AutoMigrate(&types.User{})
	if err := userStorage.EnsureHardcodedUser(db); err != nil {
		log.Fatalf("seed harcoded user error %s", err)
	}
	db.AutoMigrate(&types.StockItem{})
	DBInstance.DB = db
}

func (database DatabaseInstance) Instance() *gorm.DB {
	return DBInstance.DB
}

func OverrideDatabaseInstance(db *gorm.DB) {
	DBInstance.DB = db
}
