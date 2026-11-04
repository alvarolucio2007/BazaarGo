package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const MongoQueryTimeout time.Duration = 5 * time.Second

// NewConn cretes a new MongoDB connection with [*mongo.Client].
//
// Returns an error if the mongo client cannot connect to the server.
// Returns an error if the [mongo.Client.Ping] health ping fails, and disconnects the client.
func NewConn(addr string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(addr))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}
	return client, nil
}

// ConnToCollection returns a mongoDB collection.
//
// Make sure it has the right databaseName and collectionName parameters.
func ConnToCollection(client *mongo.Client, databaseName, collectionName string) *mongo.Collection {
	return client.Database(databaseName).Collection(collectionName)
}
