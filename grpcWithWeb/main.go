package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpcClient/certificattedRpc"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	port      = ":1107"
	serverCrt = "/Users/mac/go/src/grpcDemo/certificattedRpc/simpleCertication/certification/server.crt"
	serverKey = "/Users/mac/go/src/grpcDemo/certificattedRpc/simpleCertication/certification/server.key"
)

type HelloService struct {
}

func (p *HelloService) Hello(ctx context.Context, args *certificattedRpc.String) (*certificattedRpc.String, error) {
	reply := &certificattedRpc.String{Value: "pubsubGrpc" + args.GetValue()}

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

		reply := &certificattedRpc.String{Value: "pubsubGrpc " + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func startServer() {
	creds, err := credentials.NewServerTLSFromFile(serverCrt, serverKey)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	certificattedRpc.RegisterHelloServiceServer(grpcServer, new(HelloService))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "grpc with web")
	})

	http.ListenAndServeTLS(port, serverCrt, serverKey, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.ProtoMajor == 2 && strings.Contains(request.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(writer, request)
		} else {
			print(11111)
			mux.ServeHTTP(writer, request)
		}
	}))

}

func doClientWork() {
	creds, err := credentials.NewClientTLSFromFile(serverCrt, "server.io")
	if err != nil {
		log.Panicf("could not load client key pair: %s", err)
	}
	conn, err := grpc.Dial(
		"localhost"+port, grpc.WithTransportCredentials(creds),
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

func main() {
	go startServer()
	time.Sleep(time.Second)

	doClientWork()
}
