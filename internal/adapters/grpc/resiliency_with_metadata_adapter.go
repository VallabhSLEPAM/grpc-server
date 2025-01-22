package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/resiliency"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func readRequestMetadata(ctx context.Context) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println("Request metadata:")
		for k, v := range md {
			log.Printf("%v:%v", k, v)
		}
	} else {
		log.Println("Request metadata not found")
	}
}

func sendResponseMetadata() metadata.MD {
	md := map[string]string{
		"grpc-server-time":     fmt.Sprint(time.Now().Format("12:40:00")),
		"grpc-server-location": "Nasik, Maharashtra, IN",
		"grpc-response-uuid":   uuid.New().String(),
	}
	return metadata.New(md)
}

func (grpcAdapter *GRPCAdapter) UnaryResiliencyWithMetadata(ctx context.Context, req *resiliency.ResiliencyRequest) (*resiliency.ResiliencyResponse, error) {
	log.Print("UnaryResiliencyWithMetadata called: max,min,statuscode", req.MaxDelaySecond, req.MinDelaySecond, req.StatusCodes)
	str, stc := grpcAdapter.resiliencyService.GenerateResiliency(int(req.MinDelaySecond), int(req.MaxDelaySecond), req.StatusCodes)

	readRequestMetadata(ctx)

	if errStatus := generateErrStatus(stc); errStatus != nil {
		return nil, errStatus
	}

	grpc.SendHeader(ctx, sendResponseMetadata())
	return &resiliency.ResiliencyResponse{
		DummyString: str,
	}, nil
}

func (grpcAdapter *GRPCAdapter) ServerResiliencyWithMetadata(req *resiliency.ResiliencyRequest, stream grpc.ServerStreamingServer[resiliency.ResiliencyResponse]) error {
	log.Println("ServerResiliencyWithMetadata called")

	ctx := stream.Context()

	readRequestMetadata(ctx)

	if err := stream.SendHeader(sendResponseMetadata()); err != nil {
		log.Println("Error while sending response metadata:", err)
	}
	for {
		select {
		case <-ctx.Done():
			log.Println("Client cancelled request")
			return nil
		default:
			str, stc := grpcAdapter.resiliencyService.GenerateResiliency(int(req.MinDelaySecond), int(req.MaxDelaySecond), req.StatusCodes)

			if errStatus := generateErrStatus(stc); errStatus != nil {
				return errStatus
			}

			stream.Send(&resiliency.ResiliencyResponse{
				DummyString: str,
			})
		}
	}

}

func (grpcAdapter *GRPCAdapter) ClientResiliencyWithMetadata(stream grpc.ClientStreamingServer[resiliency.ResiliencyRequest, resiliency.ResiliencyResponse]) error {
	log.Println("ClientResiliencyWithMetadata called")

	i := 0
	for {

		req, err := stream.Recv()

		if err == io.EOF {
			res := resiliency.ResiliencyResponse{
				DummyString: fmt.Sprintf("Received %v requests from client:", i),
			}
			if err := stream.SendHeader(sendResponseMetadata()); err != nil {
				log.Println("Error while sending response metadata:", err)
			}
			return stream.SendAndClose(&res)
		}
		ctx := stream.Context()
		readRequestMetadata(ctx)
		if req != nil {
			_, stc := grpcAdapter.resiliencyService.GenerateResiliency(int(req.MinDelaySecond), int(req.MaxDelaySecond), req.StatusCodes)

			if errStatus := generateErrStatus(stc); errStatus != nil {
				return errStatus
			}
		}
		i = i + 1
	}
}

func (grpcAdapter *GRPCAdapter) BiDirectionalResiliencyWithMetadata(bidirectionalStream grpc.BidiStreamingServer[resiliency.ResiliencyRequest, resiliency.ResiliencyResponse]) error {
	log.Println("BiDirectionalResiliencyWithMetadata called")

	ctx := bidirectionalStream.Context()

	if err := bidirectionalStream.SendHeader(sendResponseMetadata()); err != nil {
		log.Println("Error while sending response metadata:", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Client cancelled request")
			return nil
		default:
			req, err := bidirectionalStream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				log.Fatalln("Error while reading from client: ", err)
			}

			readRequestMetadata(ctx)

			str, stc := grpcAdapter.resiliencyService.GenerateResiliency(int(req.MinDelaySecond), int(req.MaxDelaySecond), req.StatusCodes)

			if errStatus := generateErrStatus(stc); errStatus != nil {
				return errStatus
			}
			err = bidirectionalStream.Send(&resiliency.ResiliencyResponse{
				DummyString: str,
			})
			if err != nil {
				log.Fatalln("Error while sending response to client: ", err)
			}

		}
	}
}
