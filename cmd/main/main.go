package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/datmaithanh/orderfood/api"
	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/gapi"
	"github.com/datmaithanh/orderfood/pb"
	"github.com/datmaithanh/orderfood/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conn, err := sql.Open(utils.DBDriver, utils.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	runGrpcServer(store)

}



func runGrpcServer(store db.Store) {
	server, err := gapi.NewServer(store)
	if err != nil {
		log.Fatal("Cannot create grpc server: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderFoodServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", utils.GrpcServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener: %w", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start gRPC server: %w", err)
	}
}

func runGinServer(store db.Store) {
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("Cannot run server: %w", err)
	}
	err = server.Start(utils.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

