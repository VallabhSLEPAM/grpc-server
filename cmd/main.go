package main

import (
	"database/sql"
	"log"

	dbmigrations "github.com/VallabhSLEPAM/grpc-server/db/migrations"
	grpc "github.com/VallabhSLEPAM/grpc-server/internal/adapters/grpc"
	app "github.com/VallabhSLEPAM/grpc-server/internal/application"
)

func main() {

	sqlDB, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/grpc?sslmode=disable")

	if err != nil {
		log.Fatalln("Unable to connect to DB: ", err)
	}

	dbmigrations.Migrate(sqlDB)

	helloService := &app.HelloService{}
	bankService := &app.BankService{}

	grpcAdapter := grpc.NewGRPCAdapter(helloService, bankService, 9090)

	grpcAdapter.Run()

}
