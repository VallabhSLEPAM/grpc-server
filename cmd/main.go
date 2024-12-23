package main

import (
	"database/sql"
	"log"

	dbmigration "github.com/VallabhSLEPAM/grpc-server/db"
	grpc "github.com/VallabhSLEPAM/grpc-server/internal/adapters/grpc"
	app "github.com/VallabhSLEPAM/grpc-server/internal/application"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	sqlDB, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/grpc?sslmode=disable")

	if err != nil {
		log.Fatalln("Unable to connect to DB: ", err)
	}

	dbmigration.Migrate(sqlDB)

	helloService := &app.HelloService{}
	bankService := &app.BankService{}

	grpcAdapter := grpc.NewGRPCAdapter(helloService, bankService, 9090)

	grpcAdapter.Run()

}
