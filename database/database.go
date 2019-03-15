package database

import (
	"context"
	"github.com/ankurgopher/mongobrain/util"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var client *mongo.Client

func init() {
	client, _ = mongo.NewClient("mongodb+srv://ankurtxns:******@cluster0-ydjii.mongodb.net/test?retryWrites=true")

	// Connect to Atlas instance
	client.Connect(context.TODO())

	// Check if ping to database is working
	err := client.Ping(context.TODO(), nil)

	if err!=nil {
		util.BigError(err.Error())
	}
}

// GetQuizCollection connects to the database & gets the collection
func GetQuizCollection() (*mongo.Collection) {
	coll := client.Database("mongobrain").Collection("quiz")
	return coll
}

func GetUserCollection() (*mongo.Collection) {
	coll := client.Database("mongobrain").Collection("users")
	return coll
}

func GetSessionsCollection() (*mongo.Collection)  {
	coll := client.Database("mongobrain").Collection("sessions")
	return coll
}