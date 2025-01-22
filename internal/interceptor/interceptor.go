package interceptor

import (
	"context"
	"log"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"
	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/resiliency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Just a logging interceptor
func LogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.Println("[LOGGED BY UNARY SERVER INTERCEPTOR]", req)
		return handler(ctx, req)
	}

}

// Interceptor which modifies the metadata from the received request
func BasicUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// If the incoming request is from HelloRequest, modify the metadata
		switch request := req.(type) {
		case *hello.HelloRequest:
			request.Name = "[MODIFIED BY UNARY SERVER INTERCEPTOR - 1]" + request.Name
		}
		responseMetadata, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			responseMetadata = metadata.New(nil)
		}
		responseMetadata.Append("my-response-metadata-key-1", "my-response-metadata-value-1")
		responseMetadata.Append("my-response-metadata-key-2", "my-response-metadata-value-2")

		ctx = metadata.NewOutgoingContext(ctx, responseMetadata)
		grpc.SetHeader(ctx, responseMetadata)

		res, err := handler(ctx, req)
		if err != nil {
			return res, err
		}

		switch response := res.(type) {
		case *hello.HelloResponse:
			response.Greet = "[MODIFIED BY UNARY SERVER INTERCEPTOR - 2]" + response.Greet
		case *resiliency.ResiliencyResponse:
			response.DummyString = "[MODIFIED BY UNARY SERVER INTERCEPTOR - 2]" + response.DummyString
		}
		return res, nil
	}

}

func LogServerStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("[LOGGED BY SERVER STREAM INTERCEPTOR]", info)
		return handler(srv, ss)
	}
}

type InterceptedServerStream struct {
	grpc.ServerStream
}

// Interceptor updating response while sending response back
func BasicStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		//intercept stream
		interceptedServerStream := &InterceptedServerStream{
			ServerStream: ss,
		}

		responseMetadata, ok := metadata.FromOutgoingContext(interceptedServerStream.Context())
		if !ok {
			responseMetadata = metadata.New(nil)
		}

		responseMetadata.Append("my-response-metadata-key-1", "my-response-metadata-value-1")
		responseMetadata.Append("my-response-metadata-key-2", "my-response-metadata-value-2")

		interceptedServerStream.SetHeader(responseMetadata)
		return handler(srv, interceptedServerStream)
	}
}

// Interceptor updating request coming to the server
func (is *InterceptedServerStream) RecvMsg(msg interface{}) error {
	if err := is.ServerStream.RecvMsg(msg); err != nil {
		return err
	}
	switch request := msg.(type) {
	case *hello.HelloRequest:
		request.Name = "[MODIFIED BY SERVER INTERCEPTOR - 4]" + request.Name
	}
	return nil
}

// Intercepting response sent from Server
func (is InterceptedServerStream) SendMsg(msg interface{}) error {
	switch response := msg.(type) {
	case *hello.HelloResponse:
		response.Greet = "[MODIFIED BY SERVER INTERCEPTOR - 5]" + response.Greet
	case *resiliency.ResiliencyResponse:
		response.DummyString = "[MODIFIED BY SERVER INTERCEPTOR - 6]" + response.DummyString
	}
	return is.ServerStream.SendMsg(msg)
}
