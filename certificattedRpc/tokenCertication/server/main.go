package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpcClient/certificattedRpc"
	"io"
	"log"
	"net"
)

var (
	port = ":1107"
)

type myGrpcServer struct {
}

func (p *myGrpcServer) Hello(ctx context.Context, args *certificattedRpc.String) (*certificattedRpc.String, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing credentials")
	}

	var (
		appId  string
		appKey string
	)

	if val, ok := md["user"]; ok {
		appId = val[0]
	}
	if val, ok := md["password"]; ok {
		appKey = val[0]
	}

	if appId != "wxning" || appKey != "gopher" {

		return nil, status.Errorf(codes.Unauthenticated, "invalid token: appId=%s, appKey=%s", appId, appKey)
	}

	reply := &certificattedRpc.String{Value: "hello" + args.GetValue()}

	return reply, nil
}

func (p *myGrpcServer) Channel(stream certificattedRpc.HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &certificattedRpc.String{Value: "hello " + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func main() {
	//grpcServer := grpc.NewServer(grpc.UnaryInterceptor(filter))
	grpcServer := grpc.NewServer()
	certificattedRpc.RegisterHelloServiceServer(grpcServer, new(myGrpcServer))

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("could not listen on %s: %s", port, err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}

}
