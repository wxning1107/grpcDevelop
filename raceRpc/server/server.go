package main

import (
	"context"
	"github.com/EDDYCJY/go-grpc-example/pkg/gtls"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"google.golang.org/grpc"
	"grpcClient/raceRpc"
	"io"
	"log"
	"net"
)

const (
	port = ":1107"

	SERVICE_NAME              = "simple_zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://127.0.0.1:9411/api/v1/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
)

type HelloService struct {
}

func (p *HelloService) Hello(ctx context.Context, args *raceRpc.String) (*raceRpc.String, error) {
	reply := &raceRpc.String{Value: "pubsubGrpc" + args.GetValue()}

	return reply, nil
}

func (p *HelloService) Channel(stream raceRpc.HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &raceRpc.String{Value: "hello " + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func main() {
	collector, err := zipkin.NewHTTPCollector(ZIPKIN_HTTP_ENDPOINT)
	if err != nil {
		log.Fatalf("zipkin.NewHTTPCollector err: %v", err)
	}

	recorder := zipkin.NewRecorder(collector, true, ZIPKIN_RECORDER_HOST_PORT, SERVICE_NAME)

	tracer, err := zipkin.NewTracer(recorder, zipkin.ClientServerSameSpan(false))
	if err != nil {
		log.Fatalf("zipkin.NewTracer err: %v", err)
	}

	tlsServer := gtls.Server{
		CaFile:   "../config/ca.pem",
		CertFile: "../config/server/server.pem",
		KeyFile:  "../config/server/server.key",
	}

	creds, err := tlsServer.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}

	//opts := []grpc.ServerOption{
	//	grpc.Creds(creds),
	//	grpc_middleware.WithUnaryServerChain(
	//		otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
	//	),
	//}

	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc_middleware.WithUnaryServerChain(
		otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
	))
	// RegisterHelloServiceServer把HelloService的方法反射给grpcServer
	raceRpc.RegisterHelloServiceServer(grpcServer, new(HelloService))

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("could not listen on %s: %s", port, err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}

}
