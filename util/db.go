package util

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var Database *gorm.DB

func InitDB() {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", dbinfo)
	if err != nil {
		log.WithField("Database", dbinfo).Errorln("Error:", err)
		db = nil
	} else {
		log.WithField("Database", dbinfo).Infof("Connected with database")
	}

	Database = db
	db.LogMode(true)
}

func GetDB(w http.ResponseWriter) *gorm.DB {
	if Database == nil || Database.DB().Ping() != nil {
		if !reconnectDatabase() {
			log.WithField("Database", Database).Errorln("Database session failed")
			if w != nil {
				SendNoConnectionWithDb(w)
			}
		}
	}
	return Database
}

func reconnectDatabase() bool {
	InitDB()
	if Database.DB() == nil {
		return false
	}
	return true
}
