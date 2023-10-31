package main

import (
	"Chitty-Chat/Logger"
	ServerStruct "Chitty-Chat/Server"
	"Chitty-Chat/gRPC_output"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {

	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000"
	}
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		Logger.ErrorLogger.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	Logger.FileLogger.Println("Listening @ : " + Port)
	grpcServer := grpc.NewServer()

	server := ServerStruct.ChatServer{}

	gRPC_output.RegisterMessageHandlerServer(grpcServer, &server)
	err = grpcServer.Serve(listen)
	if err != nil {
		Logger.ErrorLogger.Fatalf("Failed to start gRPC server :: %v", err)
	}
}
