package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"time"

	pb "git.badhouseplants.net/badhouseplants/postman-service/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedPostmanServer
	senderName     string
	senderPassword string
	receiverEmail  string
	receiverName   string
	smtpHost       string
}

func (s *server) SendEmail(ctx context.Context, in *pb.Email) (*emptypb.Empty, error) {
	// Check for an empty email
	if len(in.SenderEmail) == 0 {
		return nil, fmt.Errorf("email can't be empty")
	}

	// Check for a valid email
	_, err := mail.ParseAddress(in.SenderEmail)
	if err != nil {
		return nil, err
	}

	// Check for an empty name
	if len(in.SenderName) == 0 {
		return nil, fmt.Errorf("name can't be empty")
	}

	var messageData string
	for key, element := range in.GetMessage() {
		messageData = fmt.Sprintf("%s\n%s: %s", messageData, key, element)
	}
	auth := smtp.PlainAuth("", s.senderName, s.senderPassword, s.smtpHost)
	messageTemplate := `To: "%s" <%s>
From: "%s" <%s>
Date: %s
Subject: %s

%s
`
	message := fmt.Sprintf(messageTemplate, s.receiverName, s.receiverEmail, in.SenderName, in.SenderEmail, time.Now().Format("01-02-2006"), in.Subject, messageData)

	if err := smtp.SendMail(s.smtpHost+":25", auth, s.senderName, []string{s.receiverEmail}, []byte(message)); err != nil {
		fmt.Println("Error SendMail: ", err)
		return nil, err
	}
	fmt.Println("Email Sent!")

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	senderName, present := os.LookupEnv("POSTMAN_SENDER_NAME")
	if !present {
		err := fmt.Errorf("POSTMAN_SENDER_NAME variable is not set")
		log.Fatal(err)
	}

	senderPassword, present := os.LookupEnv("POSTMAN_SENDER_PASSWORD")
	if !present {
		err := fmt.Errorf("POSTMAN_SENDER_PASSWORD variable is not set")
		log.Fatal(err)
	}

	receiverEmail, present := os.LookupEnv("POSTMAN_RECEIVER_EMAIL")
	if !present {
		err := fmt.Errorf("POSTMAN_RECEIVER_EMAIL variable is not set")
		log.Fatal(err)
	}

	receiverName, present := os.LookupEnv("POSTMAN_RECEIVER_NAME")
	if !present {
		err := fmt.Errorf("POSTMAN_RECEIVER_NAME variable is not set")
		log.Fatal(err)
	}

	smtpHost, present := os.LookupEnv("POSTMAN_SMTP_HOST")
	if !present {
		err := fmt.Errorf("POSTMAN_SMTP_HOST variable is not set")
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterPostmanServer(s, &server{
		senderName:     senderName,
		senderPassword: senderPassword,
		receiverEmail:  receiverEmail,
		receiverName:   receiverName,
		smtpHost:       smtpHost,
	})

	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
