package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var UserCollection *mongo.Collection

func InitDb() {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}

	UserCollection = client.Database("test").Collection("test")

	fmt.Println(client)
	err = EnsureUserEmailUniqueIndex(UserCollection)
	if err != nil {
		panic(err)
	}
}

func EnsureUserEmailUniqueIndex(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
