package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"pcbook/pb"
	"pcbook/sample"
	"time"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()

	log.Printf("dail server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dail server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	laptop := sample.NewLaptop()
	// laptop.Id = "" // 如果laptop.Id为空，服务器会自动创建
	// laptop.Id = "ee1a0511-5f77-49d2-ae5d-2a348cbb6f8a" // 如果laptop.Id已存在，也没事
	// laptop.Id = "invalid-uuid" // 如果laptop.Id非法，报错
	req := &pb.CreateLaptopRequest{Laptop: laptop}

	// 设置请求超时
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel() // 退出主函数之前调用

	// res, err := laptopClient.CreateLaptop(context.Background(), req)
	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			// not a big deal
			log.Println("laptop already exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return // 在任何错误情况器，我们都必须返回此处
	}
	log.Printf("create laptop with id: %s", res.Id)
}
