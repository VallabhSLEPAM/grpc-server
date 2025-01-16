package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	dbmigration "github.com/VallabhSLEPAM/grpc-server/db"
	"github.com/VallabhSLEPAM/grpc-server/internal/adapters/database"
	grpc "github.com/VallabhSLEPAM/grpc-server/internal/adapters/grpc"
	app "github.com/VallabhSLEPAM/grpc-server/internal/application"
	"github.com/VallabhSLEPAM/grpc-server/internal/application/domain/bank"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Create SQL connection
	sqlDB, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/grpc?sslmode=disable")
	if err != nil {
		log.Fatalln("Unable to connect to DB: ", err)
	}

	dbmigration.Migrate(sqlDB)
	helloService := &app.HelloService{}

	// We create a DB adapter here which will return us the gorm.DB object for doing DB operations
	dbAdapter, err := database.NewDatabaseAdapter(sqlDB)
	if err != nil {
		log.Fatalln("Error creating DB adapter:", err)
	}

	bankService := app.NewBankService(dbAdapter)
	go generateExchangeRates(*bankService, "USD", "INR", 5*time.Second)

	resiliencyService := &app.ResiliencyService{}

	grpcAdapter := grpc.NewGRPCAdapter(helloService, bankService, resiliencyService, 9090)

	grpcAdapter.Run()

}

func generateExchangeRates(bs app.BankService, fromCurrency, toCurrency string, duration time.Duration) {

	ticker := time.NewTicker(duration)

	for range ticker.C {

		now := time.Now()
		validFrom := now.Truncate(time.Second).Add(3 * time.Second)
		validTo := validFrom.Add(duration).Add(-1 * time.Millisecond)

		dummyRate := bank.ExchangeRate{
			FromCurrency:       fromCurrency,
			ToCurrency:         toCurrency,
			ValidFromTimeStamp: validFrom,
			ValidToTimeStamp:   validTo,
			Rate:               80 + float64(rand.Intn(300)),
		}

		bs.CreateExchangeRate(dummyRate)
	}

}
