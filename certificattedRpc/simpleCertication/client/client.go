package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpcClient/certificattedRpc"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	clientCrt = "/Users/mac/go/src/grpcDemo/certificattedRpc/simpleCertication/certification/client.crt"
	clientKey = "/Users/mac/go/src/grpcDemo/certificattedRpc/simpleCertication/certification/client.key"
	caCrt     = "/Users/mac/go/src/grpcDemo/certificattedRpc/simpleCertication/certification/ca.crt"
)

func main() {
	certificate, err := tls.LoadX509KeyPair(clientCrt, clientKey)
	if err != nil {
		log.Panicf("could not load client key pair: %s", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCrt)
	if err != nil {
		log.Panicf("could not read ca certificate: %s", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append ca certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ServerName:   "server.io", // NOTE: this is required!
		RootCAs:      certPool,
	})
	conn, err := grpc.Dial(
		"localhost:1107", grpc.WithTransportCredentials(creds),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := certificattedRpc.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &certificattedRpc.String{Value: " grpc"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Println(reply)

	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			num := rand.Intn(100)
			if err := stream.Send(&certificattedRpc.String{Value: "grpc" + strconv.Itoa(num)}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()

	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		fmt.Println(reply.GetValue())

	}
}
