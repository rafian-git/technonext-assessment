package main

import (
	"github.com/redis/go-redis/v9"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ordersv1 "gitlab.com/sample_projects/technonext-assessment/gen/orders/v1"
	"gitlab.com/sample_projects/technonext-assessment/internal/auth"
	"gitlab.com/sample_projects/technonext-assessment/internal/pg"
	"gitlab.com/sample_projects/technonext-assessment/internal/repo"
	"gitlab.com/sample_projects/technonext-assessment/internal/server"
	"gitlab.com/sample_projects/technonext-assessment/internal/service"
)

func main() {
	db, err := pg.Connect()
	if err != nil {
		log.Fatal(err)
	}

	var rdb *redis.Client
	if addr := getenv("REDIS_ADDR", ""); addr != "" {
		rdb = redis.NewClient(&redis.Options{Addr: addr})
	}

	jwtp := auth.NewProvider()
	accessTTL := durEnv("ACCESS_TOKEN_TTL", "120h")
	refreshTTL := durEnv("REFRESH_TOKEN_TTL", "240h")

	authRepo := repo.NewAuthRepo(db)
	orderRepo := repo.NewOrderRepo(db)

	authSvc := service.NewAuthService(authRepo, jwtp, rdb, accessTTL, refreshTTL)
	orderSvc := service.NewOrderService(orderRepo)

	authz := server.NewAuthz(jwtp, authSvc)

	authSrv := server.NewAuthServer(authSvc, jwtp)
	orderSrv := server.NewOrderServer(orderSvc, authz)

	lis, _ := net.Listen("tcp", getenv("GRPC_ADDR", ":50051"))
	s := grpc.NewServer()
	ordersv1.RegisterAuthServiceServer(s, authSrv)
	ordersv1.RegisterOrderServiceServer(s, orderSrv)
	reflection.Register(s)
	log.Println("gRPC listening on", lis.Addr().String())
	log.Fatal(s.Serve(lis))
}

func durEnv(key, def string) time.Duration {
	v := getenv(key, def)
	d, err := time.ParseDuration(v)
	if err != nil {
		d, _ = time.ParseDuration(def)
	}
	return d
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
