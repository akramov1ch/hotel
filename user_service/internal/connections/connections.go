package connections

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"user-service/config"
	ntf "user-service/internal/client/notifications"
	intr "user-service/internal/interface"
	"user-service/internal/interface/service"
	adj "user-service/internal/service/adjust"
	grpcmet "user-service/internal/service/methods"
	"user-service/pkg/database/adjust"
	cons "user-service/pkg/kafka/consumer"

	_ "github.com/lib/pq"
)

func NewDatabase() intr.User {
	c := config.Configuration()
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.Database.User, c.Database.Password, c.Database.Host, c.Database.DBname))
	if err != nil {
		log.Println(err)
	}
	if err := db.Ping(); err != nil {
		log.Println(err)
	}
	n := ntf.Hotel()
	return &adjust.Database{Db: db, N: n}
}

func NewService() *service.Service {
	a := NewDatabase()
	return &service.Service{D: a}
}

func NewAdjust() intr.Adjust {
	a := NewService()
	return &adj.Adjust{S: a}
}

func NewAdjustService() *service.Adjust {
	a := NewAdjust()
	return &service.Adjust{A: a}
}

func NewGrpc() *grpcmet.Service {
	a := NewAdjustService()
	return &grpcmet.Service{S: a}
}

func NewConsumer() *cons.Consumer17 {
	a := NewGrpc()
	ctx := context.Background()
	return &cons.Consumer17{C: a, Ctx: ctx}
}
