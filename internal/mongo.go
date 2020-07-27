package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type MongoConfig struct {
	Uri        string
	Database   string
	Collection string
}

type MongoManager struct {
	client *mongo.Client
	config MongoConfig
}

func NewMongoManager(config MongoConfig) MongoManager {
	manager := MongoManager{}
	manager.config = config

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.Uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error during creating MongoManager ", err)
	}
	manager.client = client

	// make sure that client can connect to MongoDB
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Cannot connect to MongoDB ", err)
	}
	return manager
}

func (manager MongoManager) Create(lp LoginPasswordAcls) error {
	collection := manager.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, lp)
	if err != nil {
		log.Print("Cannot insert acl to MongoDB ", err)
	}
	return err
}

func (manager MongoManager) Remove(login Login) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := manager.getCollection().DeleteOne(ctx, bson.M{"login": login.Login})
	if err != nil {
		log.Print("Cannot remove from MongoDB", err)
	}
	return err
}

func (manager MongoManager) GetAll() []LoginPasswordAcls {
	var data []LoginPasswordAcls
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := manager.getCollection().Find(ctx, bson.M{})
	if err != nil {
		log.Print("Error during MongoDB find", err)
	}
	if err = cursor.All(ctx, &data); err != nil {
		log.Print("Error during MongoDB all", err)
	}
	return data
}

func (manager MongoManager) ObserveIfSupported(service ManagerService) {
}
func (manager MongoManager) IsObserveSupported() bool {
	return false
}

func (manager MongoManager) getCollection() *mongo.Collection {
	return manager.client.Database(manager.config.Database).Collection(manager.config.Collection)
}
