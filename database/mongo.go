package database

import (
	"context"
	"github.com/yamamushi/kmud-2020/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type DatabaseHandler struct {
	conf   *config.Config
	host   string
	user   string
	pass   string
	client *mongo.Client
}

func NewDatabaseHandler(config *config.Config) *DatabaseHandler {
	db := &DatabaseHandler{conf: config}
	db.host = db.conf.DB.MongoHost
	db.user = db.conf.DB.MongoUser
	db.pass = db.conf.DB.MongoPass
	return db
}

func (db *DatabaseHandler) Connect() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db.client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+db.host))
	if err != nil {
		return err
	}
	return nil
}

func (db *DatabaseHandler) CheckConnection() (err error) {
	if db.client == nil {
		return db.Connect()
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = db.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return db.Connect()
	}
	return nil
}

func (db *DatabaseHandler) GetCollection(database string, collection string) (mongocollection *mongo.Collection) {
	mongocollection = db.client.Database(database).Collection(collection)
	return mongocollection
}

func (db *DatabaseHandler) Insert(object interface{}, database string, collection string) (err error) {
	err = db.CheckConnection()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	mcollection := db.GetCollection(database, collection)
	_, err = mcollection.InsertOne(ctx, object)
	return err
}

func (db *DatabaseHandler) UpdateOne(filter interface{}, object interface{}, database string, collection string) (err error) {
	err = db.CheckConnection()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	mcollection := db.GetCollection(database, collection)
	_, err = mcollection.UpdateOne(ctx, filter, bson.D{{"$set", object}})
	return err
}

func (db *DatabaseHandler) UpdateMany(filter interface{}, object interface{}, database string, collection string) (err error) {
	err = db.CheckConnection()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	mcollection := db.GetCollection(database, collection)
	_, err = mcollection.UpdateMany(ctx, filter, object)
	return err
}

func (db *DatabaseHandler) DeleteOne(filter interface{}, database string, collection string) (err error) {
	err = db.CheckConnection()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	mcollection := db.GetCollection(database, collection)
	_, err = mcollection.DeleteOne(ctx, filter)
	return err
}

func (db *DatabaseHandler) DeleteMany(filter interface{}, database string, collection string) (err error) {
	err = db.CheckConnection()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	mcollection := db.GetCollection(database, collection)
	_, err = mcollection.DeleteMany(ctx, filter)
	return err
}

func (db *DatabaseHandler) FindOne(filter interface{}, database string, collection string) (output bson.D, err error) {
	err = db.CheckConnection()
	if err != nil {
		return output, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	mcollection := db.GetCollection(database, collection)
	err = mcollection.FindOne(ctx, filter).Decode(&output)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (db *DatabaseHandler) FindAll(filter interface{}, database string, collection string) (results []bson.D, err error) {
	err = db.CheckConnection()
	if err != nil {
		return results, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	mcollection := db.GetCollection(database, collection)

	cur, err := mcollection.Find(ctx, filter)
	if err != nil {
		return results, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			return results, err
		}
		results = append(results, result)
		// do something with result....
	}
	err = cur.Err()
	if err != nil {
		return results, err
	}
	return results, nil
}
