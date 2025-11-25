package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/datmaithanh/orderfood/api"
	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/gapi"
	"github.com/datmaithanh/orderfood/pb"
	"github.com/datmaithanh/orderfood/utils"
	"github.com/datmaithanh/orderfood/worker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	utils.LoadConfig()
	conn, err := sql.Open(utils.DBDriver, utils.DBSource)
	if err != nil {
		log.Fatal().Msgf("cannot connect to db: %s", err)
	}

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr:    utils.Redis_Addr,
		Password: utils.Redis_Password,
		TLSConfig: &tls.Config{
		},
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store)
	go runGatewayServer(store, taskDistributor)
	runGrpcServer(store, taskDistributor)
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	processor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := processor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

}

func runGrpcServer(store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("Cannot create grpc server: %s", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLoger)

	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterOrderFoodServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", utils.GrpcServerAddress)
	if err != nil {
		log.Fatal().Msgf("Cannot create listener: %s", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msgf("Cannot start gRPC server: %s", err)
	}
}

func runGatewayServer(store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("Cannot create HTTP gateway server: %s", err)
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
		log.Fatal().Msgf("Cannot register gateway server: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", utils.ServerAddress)
	if err != nil {
		log.Fatal().Msgf("Cannot create listener: %s", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	handler := gapi.HTTPLoger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msgf("Cannot start HTTP gateway server: %s", err)
	}
}

func runGinServer(store db.Store) {
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal().Msgf("Cannot run server: %s", err)
	}
	err = server.Start(utils.ServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot start server: %s", err)
	}
}
