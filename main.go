package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"simple-bank/api"
	db "simple-bank/db/sqlc"
	"simple-bank/gapi"
	"simple-bank/pb"
	until "simple-bank/util"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := until.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot open db:", err)
	}
	defer conn.Close()

	// Ping the database to ensure a connection can be established
	err = conn.Ping()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	fmt.Println("Successfully connected to the database")

	if err != nil {
		log.Fatal("cannot connect to db:", err)
		log.Fatalf("cannot connect to db: %v", err)
	}
	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config until.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}

func runGinServer(config until.Config, store db.Store) error {
	server, err := api.NewServer(config, store)
	if err != nil {
		return err
	}
	return server.Start(config.HTTPServerAddress)
}
