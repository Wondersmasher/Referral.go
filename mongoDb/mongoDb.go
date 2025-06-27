package mongodb

import (
	"context"
	"fmt"
	"time"

	// "github.com/Wondersmasher/Referral/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var UserCollection *mongo.Collection

func InitDb() {
	// fmt.Println("Connecting to MongoDB...")
	// fmt.Println(env.MONGO_DB_COLLECTION, env.MONGO_DB_DATABASE, env.MONGO_DB_URL)
	// client, err := mongo.Connect(options.Client().ApplyURI(env.MONGO_DB_URL))

	// if err != nil {
	// 	panic(err)
	// }

	// UserCollection = client.Database(env.MONGO_DB_DATABASE).Collection(env.MONGO_DB_COLLECTION)

	// fmt.Println(client)
	// err = EnsureUserEmailUniqueIndex(UserCollection)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Connected to MongoDB...")
	fmt.Println("Connecting to MongoDB...")
	client, _ := mongo.Connect(options.Client().ApplyURI("mongodb+srv://wondersmasher:.E.79kFRqzt57pW@golang.ydf6sqc.mongodb.net/?retryWrites=true&w=majority&appName=Golang"))
	UserCollection = client.Database("Golang").Collection("user")
	err := UserCollection.Drop(context.TODO())
	if err != nil {
		fmt.Println("Error dropping db", err)
	}
	err = EnsureUserEmailUniqueIndex(UserCollection)

	if err != nil {
		fmt.Println("Error connecting to MongoDB for unique email", err)
	}
	fmt.Println("Connected to MongoDB...")
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
