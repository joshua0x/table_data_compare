package db

import (
	"github.com/joshua0x/table_data_compare/config"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

//	db, cErr := gorm.Open(mysql.Open(dsn), &gorm.Config{}) ,
var hostaDb *gorm.DB
var hostbDb *gorm.DB


func InitDB(cfg config.DbConfig) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:              time.Second * 3,   // Slow SQL threshold
			LogLevel:                   logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,          // Disable color
		},
	)

	var err error
	hostaDb,err = gorm.Open(mysql.Open(cfg.HostA), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}


	hostbDb,err = gorm.Open(mysql.Open(cfg.HostB), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

}

func GetDb() (*gorm.DB,*gorm.DB){
	return hostaDb,hostbDb
}
