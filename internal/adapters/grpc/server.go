package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/bank"
	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/resiliency"
	port "github.com/VallabhSLEPAM/grpc-server/internal/ports.go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCAdapter struct {
	helloService      port.HelloServicePort
	bankService       port.BankServicePort
	resiliencyService port.ResiliencyServicePort
	grpcPort          int
	server            *grpc.Server
	hello.HelloServiceServer
	bank.BankServiceServer
	resiliency.ResiliencyServiceServer
	resiliency.ResiliencyServiceWithMetadataServer
}

func NewGRPCAdapter(helloService port.HelloServicePort, bankService port.BankServicePort, resiliencyService port.ResiliencyServicePort, grpcPort int) *GRPCAdapter {
	return &GRPCAdapter{
		helloService:      helloService,
		bankService:       bankService,
		resiliencyService: resiliencyService,
		// ResiliencyServiceWithMetadataServer: res,
		grpcPort: grpcPort,
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

	creds, err := credentials.NewClientTLSFromFile("ssl/server.crt", "ssl/server.pem")
	if err != nil {
		log.Fatalln("Can't create server credentials: ", err)
	}

	// Create a gRPC server
	grpcServiceRegistrar := grpc.NewServer(
		grpc.Creds(creds),
	// grpc.ChainUnaryInterceptor(
	// 	interceptor.BasicUnaryServerInterceptor(),
	// 	interceptor.LogUnaryServerInterceptor(),
	// ),
	// grpc.ChainStreamInterceptor(
	// 	interceptor.BasicStreamServerInterceptor(),
	// 	interceptor.LogServerStreamInterceptor(),
	// ),
	)
	adapter.server = grpcServiceRegistrar

	// Associate the gRPC server with gRPC service registrar and pass it the struct which will implement the rpc methods
	hello.RegisterHelloServiceServer(grpcServiceRegistrar, adapter)
	bank.RegisterBankServiceServer(grpcServiceRegistrar, adapter)
	resiliency.RegisterResiliencyServiceServer(grpcServiceRegistrar, adapter)
	resiliency.RegisterResiliencyServiceWithMetadataServer(grpcServiceRegistrar, adapter)
	// Now the service registrar will start serving the request taking the TCP listener as input
	if err = grpcServiceRegistrar.Serve(listen); err != nil {
		log.Fatalf("Failed to server gRPC on port :%v : %v\n", adapter.grpcPort, err)
	}

}

func (adapter *GRPCAdapter) Stop() {
	adapter.server.Stop()
}
