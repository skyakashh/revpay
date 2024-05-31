package controller

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbname = "netflix"
const connectionsting = "mongodb+srv://sky:akash@cluster0.g6gxgfm.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const colname = "accounts"
const colname1 = "userid"

var Collection *mongo.Collection
var IdCollection *mongo.Collection

func init() {
	clientoption := options.Client().ApplyURI(connectionsting)

	// connect
	client, err := mongo.Connect(context.TODO(), clientoption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection successfull for mongodb")

	Collection = client.Database(dbname).Collection(colname)
	IdCollection = client.Database(dbname).Collection(colname1)
	// instance ready

	fmt.Println("instance ready")

}
