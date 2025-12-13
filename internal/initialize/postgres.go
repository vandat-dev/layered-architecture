package initialize

import (
	"app/global"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Postgres() {
	/*
		refer https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL for details
	*/
	global.Logger.Info("Start connecting to postgres")
	p := global.Config.Postgres

	// PostgreSQL DSN without timezone - will use server default (UTC+0)
	dsn := "host=%s user=%s password=%s dbname=%s port=%v sslmode=%s"
	var s = fmt.Sprintf(dsn, p.Host, p.UserName, p.Password, p.DBName, p.Port, p.SSLMode)

	db, err := gorm.Open(postgres.Open(s), &gorm.Config{
		SkipDefaultTransaction: true, // turn on transaction
	})
	checkErrPanic(err, "Initialize PostgreSQL database failed")
	global.Postgres = db
	global.Logger.Info("PostgreSQL connect successfully!")
	setPostgresPool()
	// Migration is now handled manually, not automatically
	// migratePostgresTables()
}

func setPostgresPool() {
	p := global.Config.Postgres
	sqlDB, err := global.Postgres.DB()
	checkErrPanic(err, "Set Pool PostgreSQL database failed")
	sqlDB.SetMaxIdleConns(p.MaxIdleConn)
	sqlDB.SetMaxOpenConns(p.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(p.ConnMaxLifeTime) * time.Second)
}

func checkErrPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString)
		panic(err)
	}
}
