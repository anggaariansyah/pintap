package utils

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"sync"
	"time"
)

type dbUtils struct {
	db *gorm.DB
}

var dbInstance *dbUtils
var dbOnce sync.Once

func GetDBConnection() *gorm.DB {
	dbOnce.Do(func() {
		log.Println("Initialize db connection...")
		connection := os.Getenv("USERNAME_DB") + ":" + os.Getenv("PASSWORD_DB") + "@tcp(" + os.Getenv("DATABASE_HOST") + ":" +
			os.Getenv("DATABASE_PORT") + ")/" + os.Getenv("DATABASE_NAME") + "?charset=utf8&parseTime=True&loc=Local"
		log.Println(connection)
		db, err := gorm.Open(os.Getenv("DATABASE_TYPE"), connection)

		if err != nil {
			log.Println(err)
		}

		db.DB().SetConnMaxLifetime(time.Second * 60)
		db.SingularTable(true)
		db.LogMode(false)

		if err != nil {
			log.Println(err)
		}

		dbInstance = &dbUtils{
			db: db,
		}
	})

	return dbInstance.db
}

