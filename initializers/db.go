package initializers

import (
	"os"

	"github.com/RIKUGHI/go-pos-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnnectToDb() {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}

	DB = db
}

func SyncDb() {
	DB.AutoMigrate(&models.User{})
}
