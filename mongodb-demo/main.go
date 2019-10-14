package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"reflect"
)

var client *mongo.Client
var collection *mongo.Collection
type Student struct {
	Name string
	Age  int
	No string
}

func init()  {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://172.30.60.8:27017"))
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("testing").Collection("students")
}

func ping()  {
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}


//curd
func add()  {
	wida := Student{"wida", 32, "8001"}
	amy := Student{"amy", 25, "8002"}
	sunny := Student{"sunny", 35, "8003"}

	ret, err := collection.InsertOne(nil, wida)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("写入一个文档", ret.InsertedID)
	student := []interface{}{amy, sunny}

	ret2, err := collection.InsertMany(nil, student)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("写入多个文档 ", ret2.InsertedIDs)
}

func find()  {
	var result Student
	err := collection.FindOne(nil,  bson.D{{"name", "amy"}, {"age", 25}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", result)
	findOptions := options.Find()
	findOptions.SetLimit(3)

	var results []*Student

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(nil)

	for _,r := range results {
		fmt.Printf("%+v\n", r)
	}

}
func update()  {
	filter := bson.D{{"no", "8001"}}
	update := bson.D{
		{"$set", bson.D{
			{"age", 33},
			{"name", "wida2"},
		}},
	}
	ret, err := collection.UpdateOne(nil, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", ret.MatchedCount, ret.ModifiedCount)
}

func del()  {
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents\n", deleteResult.DeletedCount)
}


func index()  {
	indexView := collection.Indexes()
	ret,err :=indexView.CreateOne(context.Background(), mongo.IndexModel{
		Keys:  bsonx.Doc{{"name", bsonx.Int32(-1)}},
		Options: options.Index().SetName("testname").SetUnique(true), //这边设置了唯一限定，不设定默认不是唯一的
	})

	fmt.Println(ret,err)
}

func main()  {




	index()
	add()
	update()
	find()
	del()
}
