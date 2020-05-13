package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpcClient/services"
	"log"
)

func main() {
	conn, err := grpc.Dial(":1107", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := services.NewProdServiceClient(conn)
	res, err := client.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 117})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.ProdStock)
}
