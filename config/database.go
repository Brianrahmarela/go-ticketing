package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	// cek ada/tidak .env di disk (“titik awal” ketika program jln, biasanya root projek)
	//dijalankan saat container start (runtime), bukan saat build image.
	if _, err := os.Stat(".env"); err == nil {
		//jika ada .env (development) → isikan ke environment proses → barulah os.Getenv bisa mengambilnya.
		if err := godotenv.Load(); err != nil {
			//jika .env tdk ada di image, docker inject value env ke container utk dibaca go
			log.Println("Warning: gagal load .env, pakai environment dari sistem")
		}
	} else {
		log.Println("No .env file found, relying on system env")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Truncate(time.Second)
		},
	})

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	return db
}
