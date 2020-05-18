package main

import (
	"context"
	"google.golang.org/grpc"
	"grpcClient/certificattedRpc"
	"grpcClient/certificattedRpc/tokenCertication/token"
	"io"
	"log"
	"net"
)

var (
	port = ":1107"
)

type myGrpcServer struct {
	auth *token.Authentication
}

func (p *myGrpcServer) Hello(ctx context.Context, args *certificattedRpc.String) (*certificattedRpc.String, error) {
	//md, ok := metadata.FromIncomingContext(ctx)
	//if !ok {
	//	return nil, fmt.Errorf("missing credentials")
	//}
	//
	//var (
	//	appId  string
	//	appKey string
	//)
	//
	//if val, ok := md["user"]; ok {
	//	appId = val[0]
	//}
	//if val, ok := md["password"]; ok {
	//	appKey = val[0]
	//}
	//
	//if appId != "wxning" || appKey != "gopher" {
	//
	//	return nil, status.Errorf(codes.Unauthenticated, "invalid token: appId=%s, appKey=%s", appId, appKey)
	//}

	if err := p.auth.Auth(ctx); err != nil {
		return nil, err
	}

	reply := &certificattedRpc.String{Value: "hello" + args.GetValue()}

	return reply, nil
}

func (p *myGrpcServer) Channel(stream certificattedRpc.HelloService_ChannelServer) error {
	//md, ok := metadata.FromIncomingContext(stream.Context())
	//if !ok {
	//	return fmt.Errorf("missing credentials")
	//}
	//
	//var (
	//	appId  string
	//	appKey string
	//)
	//
	//if val, ok := md["user"]; ok {
	//	appId = val[0]
	//}
	//if val, ok := md["password"]; ok {
	//	appKey = val[0]
	//}
	//
	//if appId != "wxning" || appKey != "gopher" {
	//	return status.Errorf(codes.Unauthenticated, "invalid token: appId=%s, appKey=%s", appId, appKey)
	//}

	if err := p.auth.Auth(stream.Context()); err != nil {
		return err
	}

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
