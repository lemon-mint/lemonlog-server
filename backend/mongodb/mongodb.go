package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/lemon-mint/lemonlog-server/backend"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongologger struct {
	client     *mongo.Client
	db         string
	collection string
	c          *mongo.Collection
}

//New mongo logger
func New(URI string, DB string, Collection string) backend.LogStore {
	var client, err = mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	col := client.Database(DB).Collection(Collection)
	return mongologger{client: client, db: DB, collection: Collection, c: col}
}

//Put to Mongodb
func (logger mongologger) Put(data []*backend.Log) error {
	datas := make([]interface{}, len(data))
	for i := range data {
		datas[i] = *data[i]
	}
	_, err := logger.c.InsertMany(context.TODO(), datas)
	return err
}

//Del from mongodb
func (logger mongologger) Del(uuid string) error {
	_, err := logger.c.DeleteOne(context.TODO(), bson.M{"uuid": uuid})
	return err
}

func (logger mongologger) Get(uuid string) (*backend.Log, error) {
	var log bson.M
	err := logger.c.FindOne(context.TODO(), bson.M{"uuid": uuid}).Decode(&log)
	if err != nil {
		return nil, err
	}
	output := new(backend.Log)
	output.UUID = log["uuid"].(string)
	output.TimeStamp = log["timestamp"].(int64)
	output.HumanReadableTime = log["humantime"].(string)
	output.LogClass = log["class"].(string)
	output.Category = log["category"].(string)
	output.Body = log["body"].(string)
	return output, nil
}
