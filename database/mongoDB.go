package database

import (
	"api/config"
	"api/types"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	mdb    *mongo.Database
)

func ConnectMongo() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.Config("MONGO_URL")).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Println(err)
	}
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		fmt.Println(err)
	}
	mdb = client.Database("shrinkr")
	fmt.Println("Connected to MongoDB")
}

func RegisterUser(user *types.User) error {
	collection := mdb.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("User registered")
	return nil
}

func GetUser(email string) (*types.User, error) {
	collection := mdb.Collection("users")
	filter := bson.D{{"username", email}}
	var user types.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// TODO
func AddURL(link *types.LinkDTO, username string) error {
	collection := mdb.Collection("links")
	filter := bson.D{{"key", link.ShortURL}}
	var result types.LinkInfo
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			var linkDoc types.LinkInfo
			linkDoc.Key = link.ShortURL
			linkDoc.LongURL = link.LongURL
			linkDoc.Description = link.Description
			linkDoc.CreatedBy = username
			linkDoc.Created = time.Now().Format("2006-01-02 15:04:05")
			// TODO
			//  also add checks for passcode and clicks and expiration
			_, err := collection.InsertOne(context.Background(), linkDoc)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Link added")
		}
		fmt.Println(err)
	}
	return nil
}

func GetUrlsByUser(username string) ([]types.LinkInfo, error) {
	collection := mdb.Collection("links")
	filter := bson.D{{"createdBy", username}}
	var result []types.LinkInfo
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteLink(shortURL string, username string) error {
	collection := mdb.Collection("links")
	filter := bson.D{{"key", shortURL}, {"createdBy", username}}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetLinkInfo(shortURL string, username string) (*types.LinkInfo, error) {
	collection := mdb.Collection("links")
	filter := bson.D{{"key", shortURL}, {"createdBy", username}}
	var result types.LinkInfo
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
