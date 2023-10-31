package Server

import (
	"Chitty-Chat/Logger"
	"Chitty-Chat/gRPC_output"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type messageHandlerServer struct {
	savedMessages []messageUnit
	clients       map[int]gRPC_output.MessageHandler_SendMessageServer
	Timestamp     int64
	mu            sync.Mutex
}
type messageUnit struct {
	ClientName       string
	Message          string
	ClientUniqueCode int
	Timestamp        int64
}

var handler = messageHandlerServer{}

type ChatServer struct {
	gRPC_output.UnimplementedMessageHandlerServer
}

func (s *ChatServer) SendMessage(handlere gRPC_output.MessageHandler_SendMessageServer) error {
	clientUniqueCode := rand.Intn(1e6)
	errorChannel := make(chan error)
	handler.mu.Lock()
	if handler.clients == nil {
		handler.clients = make(map[int]gRPC_output.MessageHandler_SendMessageServer)
	}
	handler.Timestamp++
	handler.clients[clientUniqueCode] = handlere
	handler.mu.Unlock()

	Logger.FileLogger.Println("@Lamport-time", handler.Timestamp, "participant {", clientUniqueCode, "} joined the chat-room")
	for _, element := range handler.clients {
		err := element.Send(&gRPC_output.Message{Timestamp: handler.Timestamp, Message: "User " + strconv.Itoa(clientUniqueCode) + " has joined the chat-room!"})
		if err != nil {
			Logger.ErrorLogger.Println(err)
		}
	}
	go getMessage(handlere, clientUniqueCode, errorChannel)

	go sendMessage(errorChannel)

	return <-errorChannel
}

func getMessage(stream gRPC_output.MessageHandler_SendMessageServer, clientUniqueCode_ int, errorHandler chan error) error {
	for {
		message, err := stream.Recv()
		if err != nil {
			<-errorHandler
		} else {
			handler.mu.Lock()
			handler.savedMessages = append(handler.savedMessages, messageUnit{
				Timestamp:        message.Timestamp,
				Message:          message.Message,
				ClientUniqueCode: clientUniqueCode_,
			})
		}
		fmt.Println("Participant", handler.savedMessages[len(handler.savedMessages)-1].ClientUniqueCode, "wrote:", handler.savedMessages[len(handler.savedMessages)-1].Message)
		handler.mu.Unlock()
	}
}

func sendMessage(errorHandler chan error) error {
	for {
		for {
			time.Sleep(500 * time.Millisecond)
			handler.mu.Lock()

			if len(handler.savedMessages) == 0 {
				handler.mu.Unlock()
				break
			}
			messageTime := handler.savedMessages[0].Timestamp
			message := handler.savedMessages[0].Message
			sender := handler.savedMessages[0].ClientUniqueCode
			handler.mu.Unlock()
			sendMssg := &gRPC_output.Message{}
			handler.Timestamp++
			if message == "#?!((USER//!)LEFT:" {
				Logger.FileLogger.Println("@Lamport-time", messageTime, "participant {", sender, "} left the chat-room")
				sendMssg = &gRPC_output.Message{Timestamp: messageTime, Message: "User " + strconv.Itoa(sender) + " has left the chat-room!"}
				delete(handler.clients, sender)
			} else {
				Logger.FileLogger.Println("@Lamport-time", messageTime, "participant {", sender, "} send the message: {", message, "}")
				sendMssg = &gRPC_output.Message{Timestamp: messageTime, Message: message}
			}
			for uniqueCode, element := range handler.clients {
				if sender != uniqueCode {
					err := element.Send(sendMssg)
					if err != nil {
						<-errorHandler
					}
				}
				handler.mu.Lock()

				if len(handler.savedMessages) > 1 {
					handler.savedMessages = handler.savedMessages[1:]
				} else {
					handler.savedMessages = []messageUnit{}
				}
				handler.mu.Unlock()
			}
		}
	}

}
