package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	client *mongo.Client
}

func EstablishConn(host, port string) Connection {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://" + host + ":" + port)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	conn := Connection{client}
	return conn
}

func (conn *Connection) Use(dbName, collection string) *mongo.Collection {
	return conn.client.Database(dbName).Collection(collection)

}

func (conn *Connection) DisconnectConn() {
	err := conn.client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
