// Package db connects the Database
package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SyafaHadyan/worku/internal/infra/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(env *env.Env) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.DBUsername,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Println("database connection failed")
	}

	log.Println("database connected")

	Migrate(db)

	return db
}
