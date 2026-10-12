package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	errFailedConnection = "failed to connect database"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("database.postgresql.username")
	password := viper.GetString("database.postgresql.password")
	host := viper.GetString("database.postgresql.host")
	port := viper.GetInt("database.postgresql.port")
	database := viper.GetString("database.postgresql.name")
	idleConnection := viper.GetInt("database.postgresql.pool.idle")
	maxConnection := viper.GetInt("database.postgresql.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.postgresql.pool.lifetime")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, username, password, database, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger: logger.New(&gormLogger{Log: log}, logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatalf("%s: %+v", errFailedConnection, err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("%s: %+v", errFailedConnection, err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Duration(maxLifeTimeConnection) * time.Second)

	if err = connection.Ping(); err != nil {
		log.Fatalf("%s: %+v", errFailedConnection, err)
	}

	return db
}

type gormLogger struct {
	Log *logrus.Logger
}

func (l *gormLogger) Printf(format string, args ...interface{}) {
	l.Log.Tracef(format, args...)
}
