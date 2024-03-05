package main

import (
	"context"
	"fmt"
	"github.com/kleczynski/log-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient
	// create context in order to disconnect

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	app := Config{
		Models: data.New(client),
	}
	// Register RPC server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Println("error registering RPC ", err)
	}
	go func() {
		err := app.rpcListen()
		if err != nil {
			log.Println("error starting RPC", err)
		}
	}()

	go app.gRPCListen()

	log.Println("starting service on port ", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting: ", err)
	}
	log.Println("Connected to mongo!")
	return c, nil
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}

	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {

		}
	}(listen)

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)

	}

}