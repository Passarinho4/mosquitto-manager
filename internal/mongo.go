package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (manager MongoManager) Create(lp LoginPasswordAcls) (*string, error) {
	collection := manager.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, lp)
	if err != nil {
		log.Print("Cannot insert acl to MongoDB ", err)
		return nil, err

	} else {
		id := res.InsertedID.(primitive.ObjectID).Hex()
		return &id, nil
	}
}

func (manager MongoManager) Remove(id Id) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objectID, _ := primitive.ObjectIDFromHex(id.Id)
	_, err := manager.getCollection().DeleteOne(ctx, bson.M{"_id": objectID})
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

func (manager MongoManager) Get(id Id) (*LoginPasswordAcls, error) {
	var data LoginPasswordAcls
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objectId, _ := primitive.ObjectIDFromHex(id.Id)
	cursor := manager.getCollection().FindOne(ctx, bson.M{"_id": objectId})
	err := cursor.Err()
	if err != nil {
		return nil, err
	} else {
		err := cursor.Decode(&data)
		if err != nil {
			log.Print("Cannot get from MongoDB", err)
			return nil, err
		}
		return &data, nil
	}
}

func (manager MongoManager) Update(id Id, lp LoginPasswordAcls) error {
	collection := manager.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objectId, _ := primitive.ObjectIDFromHex(id.Id)
	res := collection.FindOneAndReplace(ctx, bson.M{"_id": objectId}, lp)
	if res.Err() != nil {
		log.Print("Cannot update creds to MongoDB ", res.Err())
		return res.Err()

	}
	return nil
}

func (manager MongoManager) ObserveIfSupported(service ManagerService) {
}
func (manager MongoManager) IsObserveSupported() bool {
	return false
}

func (manager MongoManager) getCollection() *mongo.Collection {
	return manager.client.Database(manager.config.Database).Collection(manager.config.Collection)
}
