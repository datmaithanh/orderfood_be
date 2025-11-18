package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/datmaithanh/orderfood/api"
	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/gapi"
	"github.com/datmaithanh/orderfood/pb"
	"github.com/datmaithanh/orderfood/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	conn, err := sql.Open(utils.DBDriver, utils.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	go runGatewayServer(store)
	runGrpcServer(store)
	// runGinServer(store)

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

func runGatewayServer(store db.Store) {
	server, err := gapi.NewServer(store)
	if err != nil {
		log.Fatal("Cannot create HTTP gateway server: %w", err)
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, called := context.WithCancel(context.Background())
	defer called()
	err = pb.RegisterOrderFoodServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("Cannot register gateway server: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", utils.ServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener: %w", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Cannot start HTTP gateway server: %w", err)
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
