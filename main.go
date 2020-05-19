package main

import "fmt"

//func main() {
//	conn, err := grpc.Dial(":1107", grpc.WithInsecure())
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
//
//	client := services.NewProdServiceClient(conn)
//	res, err := client.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 117})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(res.ProdStock)
//}

func bubble(s []int) {
	n := len(s)
	for sorted := false; !sorted; n-- {
		for j := 0; j < n-1; j++ {
			if s[j+1] < s[j] {
				s[j+1], s[j] = s[j], s[j+1]
				sorted = true
			}
		}
		sorted = !sorted
	}
}

func fibonacci(n int) int {
	f, g := 0, 1
	for ; n > 0; n-- {
		g = g + f
		f = g - f
	}
	return g
}

func main() {
	s := []int{4, 6, 8, 2, 1, 7, 3, 3, 5}
	bubble(s)
	fmt.Println(s)
	fmt.Println(fibonacci(5))
}
