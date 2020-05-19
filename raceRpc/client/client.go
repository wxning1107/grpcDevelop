package main

import (
	"context"
	"fmt"
	"github.com/EDDYCJY/go-grpc-example/pkg/gtls"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"google.golang.org/grpc"
	"grpcClient/raceRpc"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	SERVICE_NAME              = "simple_zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://127.0.0.1:9411/api/v1/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
)

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

	conn, err := grpc.Dial("localhost:1107", grpc.WithTransportCredentials(creds), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads())))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := raceRpc.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &raceRpc.String{Value: " grpc"})
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
			if err := stream.Send(&raceRpc.String{Value: "grpc" + strconv.Itoa(num)}); err != nil {
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
