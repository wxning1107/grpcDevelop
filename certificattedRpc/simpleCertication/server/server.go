package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpcClient/certificattedRpc"
	"grpcClient/certificattedRpc/tokenCertication/token"
	"io"
	"io/ioutil"

	"log"
	"net"
)

var (
	port      = ":1107"
	serverCrt = "/Users/mac/go/src/grpcDemo/certificattedRpc/simpleCertication/certification/server.crt"
	serverKey = "/Users/mac/go/src/grpcDemo/certificattedRpc/simpleCertication/certification/server.key"
	caCrt     = "/Users/mac/go/src/grpcDemo/certificattedRpc/simpleCertication/certification/ca.crt"
)

type HelloService struct {
}

func (p *HelloService) Hello(ctx context.Context, args *certificattedRpc.String) (*certificattedRpc.String, error) {
	reply := &certificattedRpc.String{Value: "pubsubRpc" + args.GetValue()}

	return reply, nil
}

func (p *HelloService) Channel(stream certificattedRpc.HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &certificattedRpc.String{Value: "pubsubRpc " + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

type grpcServer struct {
	auth *token.Authentication
}

func main() {
	certificate, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		log.Panicf("could not load server key pair: %s", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCrt)
	if err != nil {
		log.Panicf("could not read ca certificate: %s", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	certificattedRpc.RegisterHelloServiceServer(grpcServer, new(HelloService))

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("could not listen on %s: %s", port, err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}

}
