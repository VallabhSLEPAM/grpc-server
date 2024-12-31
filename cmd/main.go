package main

import (
	"database/sql"
	"log"

	dbmigration "github.com/VallabhSLEPAM/grpc-server/db"
	"github.com/VallabhSLEPAM/grpc-server/internal/adapters/database"
	grpc "github.com/VallabhSLEPAM/grpc-server/internal/adapters/grpc"
	app "github.com/VallabhSLEPAM/grpc-server/internal/application"
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

	grpcAdapter := grpc.NewGRPCAdapter(helloService, bankService, 9090)

	grpcAdapter.Run()

}
