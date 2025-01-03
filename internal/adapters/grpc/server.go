package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/bank"
	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"
	port "github.com/VallabhSLEPAM/grpc-server/internal/ports.go"
	"google.golang.org/grpc"
)

type GRPCAdapter struct {
	helloService port.HelloServicePort
	bankService  port.BankServicePort
	grpcPort     int
	server       *grpc.Server
	hello.HelloServiceServer
	bank.BankServiceServer
}

func NewGRPCAdapter(helloService port.HelloServicePort, bankService port.BankServicePort, grpcPort int) *GRPCAdapter {
	return &GRPCAdapter{
		helloService: helloService,
		bankService:  bankService,
		grpcPort:     grpcPort,
	}
}

func (adapter *GRPCAdapter) Run() {

	var err error
	// Create a listener for TCP connections
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", adapter.grpcPort))

	if err != nil {
		log.Fatalf("Failed to listen on port :%d : %v\n", adapter.grpcPort, err)
	}

	log.Println("Server listening on port ", adapter.grpcPort)

	// Create a gRPC server
	grpcServiceRegistrar := grpc.NewServer()
	adapter.server = grpcServiceRegistrar

	// Associate the gRPC server with gRPC service registrar and pass it the struct which will implement the rpc methods
	hello.RegisterHelloServiceServer(grpcServiceRegistrar, adapter)
	bank.RegisterBankServiceServer(grpcServiceRegistrar, adapter)

	// Now the service registrar will start serving the request taking the TCP listener as input
	if err = grpcServiceRegistrar.Serve(listen); err != nil {
		log.Fatalf("Failed to server gRPC on port :%v : %v\n", adapter.grpcPort, err)
	}

}

func (adapter *GRPCAdapter) Stop() {
	adapter.server.Stop()
}
