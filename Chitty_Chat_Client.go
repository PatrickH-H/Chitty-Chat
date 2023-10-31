package main

import (
	pb "Chitty-Chat/gRPC_output"
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strings"
	"time"
)

type clientHandle struct {
	stream           pb.MessageHandler_SendMessageClient
	clientName       string
	lamportTimestamp int64
}

func main() {
	fmt.Println("Enter port number")

	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read from console %v", err)
	}
	serverID = strings.Trim(serverID, "\r\n")

	log.Println("Connection: " + serverID)
	serverID = "localhost:" + serverID
	conn, err := grpc.Dial(serverID, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server %v", err)
	}
	defer conn.Close()

	client := pb.NewMessageHandlerClient(conn)

	stream, err := client.SendMessage(context.Background())

	if err != nil {
		log.Fatalf("Failed to call SendMessage %v", err)
	}
	ch := clientHandle{stream: stream, lamportTimestamp: 0}
	//ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	bl := make(chan bool)
	<-bl

}

func (ch *clientHandle) clientConfig() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read from console %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")
}

func (ch *clientHandle) sendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)

		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from console %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")
		if clientMessage == "quit" {
			clientMessageBox := &pb.Message{
				Message:   "#?!((USER//!)LEFT:", //Message indicating user left. The message has no significance other than a message a user is NOT likely to write themselves
				Timestamp: ch.lamportTimestamp,
			}
			err = ch.stream.Send(clientMessageBox)
			time.Sleep(500 * time.Millisecond)
			os.Exit(1)
		} else {
			clientMessageBox := &pb.Message{
				Message:   clientMessage,
				Timestamp: ch.lamportTimestamp,
			}
			err = ch.stream.Send(clientMessageBox)

			if err != nil {
				log.Printf("Error while sending message %v", err)
			}
		}
		ch.lamportTimestamp++

	}
}
func (ch *clientHandle) receiveMessage() {

	for {
		mssg, err := ch.stream.Recv()
		if mssg.Timestamp > ch.lamportTimestamp {
			ch.lamportTimestamp = mssg.Timestamp
		}
		ch.lamportTimestamp++
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}
		fmt.Printf("@Lamport-time %d - %s\n", mssg.Timestamp, mssg.Message)

	}
}
