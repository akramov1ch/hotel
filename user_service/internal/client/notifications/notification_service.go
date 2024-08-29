package notification

import (
	"log"
	ntf "user-service/pkg/proto/notification"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Hotel() ntf.NotificationClient {
	conn, err := grpc.NewClient("localhost:8084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("notification error",err)
	}
	client := ntf.NewNotificationClient(conn)
	return client
}
