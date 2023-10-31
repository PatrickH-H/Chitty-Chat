package main

//the following code has taken inspiration from https://www.youtube.com/watch?v=pRSKJIt3PYU&t=118s
import (
	"Chitty-Chat/Logger"
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
		Logger.ErrorLogger.Println(err)
	}
	serverID = strings.Trim(serverID, "\r\n")

	log.Println("Connection: " + serverID)
	serverID = "localhost:" + serverID
	conn, err := grpc.Dial(serverID, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		Logger.ErrorLogger.Println(err)
	}
	defer conn.Close()
	client := pb.NewMessageHandlerClient(conn)

	stream, err := client.SendMessage(context.Background())

	if err != nil {
		Logger.ErrorLogger.Println(err)
	}
	ch := clientHandle{stream: stream, lamportTimestamp: 0}

	go ch.sendMessage()
	go ch.receiveMessage()

	bl := make(chan bool)
	<-bl
}
func (ch *clientHandle) sendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			Logger.ErrorLogger.Println(err)
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
				Logger.ErrorLogger.Println(err)
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
			Logger.ErrorLogger.Println(err)
		}
		fmt.Printf("@Lamport-time %d - %s\n", mssg.Timestamp, mssg.Message)

	}
}
