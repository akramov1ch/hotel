package connections

import (
	"api-gateway/api/handler"
	broad "api-gateway/internal/broadcast"
	"api-gateway/internal/controllers/booking"
	hotels "api-gateway/internal/controllers/hotel"
	users "api-gateway/internal/controllers/user"
	redmet "api-gateway/pkg/redis/method"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewBroadcast() *broad.Adjust {
	u := users.UserClinet()
	h:=hotels.Hotel()
	b:=booking.Hotel()
	r := Redis()
	ctx := context.Background()
	return &broad.Adjust{U: u, Ctx: ctx, R: r,H: h,B: b}
}

func NewHandler() *handler.Handler {
	a := NewBroadcast()
	return &handler.Handler{B: a}
}

func Redis() *redmet.Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &redmet.Redis{R: client, Ctx: ctx}
}
