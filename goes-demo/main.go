package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"reflect"
)


var client *elastic.Client

func init()  {
	var err error
	client, err = elastic.NewClient(elastic.SetURL("http://es.ops.int-api.xyz"))
	if err != nil {
		checkErr(err)
	}
}


func main() {
	//add()
	//find()
	//update()

	//delete()
	find()

	agg()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


type Item struct {
	Id               int64  `json:"id"`
	Appid            string `json:"appid"`
	AppBAutoId       string `json:"app_b_auto_id"`
}


func add()  {
	//add one
	//client.
	item := Item{Id:int64(21),Appid:fmt.Sprintf("app_%d",21),AppBAutoId:fmt.Sprintf("app_%d",21+200)}
	put, err := client.Index().
		Index("es_test").
		Type("test").
		Id("1").       //这个id也可以指定,不指定的话 es自动生成一个
		BodyJson(item).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println(put)

	//add many
	bulkRequest := client.Bulk()
	for i:=0;i<20;i++ {
		item := Item{Id:int64(i),Appid:fmt.Sprintf("app_%d",i),AppBAutoId:fmt.Sprintf("app_%d",i+200)}
		bulkRequest.Add(elastic.NewBulkIndexRequest().Index("es_test").Type("test").Doc(item))
	}
	bulkRequest.Do(context.TODO())
}

func find()  {
	//find one
	get1, err := client.Get().
		Index("es_test").
		Type("test").
		Id("1").
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}
	var ttyp Item
	json.Unmarshal(*get1.Source,&ttyp)
	fmt.Println("item",ttyp)

	//find many
	searchResult, err := client.Search().
		Index("es_test"). Sort("id", true).
		Type("test").From(0).Size(100).
		Do(context.TODO())
	if err != nil {
		panic(err)
	}
	if searchResult.Hits.TotalHits >0 {
		var ttyp Item
		for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
			t := item.(Item)
			fmt.Println("item ",t)
		}

	}
}


func update() {
	fmt.Println(client.Update().Index("es_test").Type("test").Id("1").
		Doc(map[string]interface{}{"appid": "app_23"}).Do(context.TODO()))
}

func delete()  {
	fmt.Println(client.Delete().Index("es_test").Type("test").Id("1").Do(context.TODO()))
}

func agg() {
	//获取最大的id
	searchResult, err := client.Search().
		Index("es_test").Type("test").
		Aggregation("max_id", elastic.NewMaxAggregation().Field("id")).Size(0).Do(context.TODO())
	if err != nil {
		panic(err)
	}
	var a map[string]float32
	if searchResult != nil {
		if v, found := searchResult.Aggregations["max_id"]; found {
			json.Unmarshal([]byte(*v), &a)
			fmt.Println(a)
		}
	}

	//统计id相同的文档数
	searchResult, err = client.Search().
		Index("es_test").Type("test").
		Aggregation("count", elastic.NewTermsAggregation().Field("id")).Size(0).Do(context.TODO())
	if err != nil {
		panic(err)
	}

	if searchResult != nil {
		if v, found := searchResult.Aggregations["count"]; found {
			var ar elastic.AggregationBucketKeyItems
			err := json.Unmarshal(*v, &ar)
			if err != nil {
				fmt.Printf("Unmarshal failed: %v\n", err)
				return
			}

			for _, item := range ar.Buckets {
				fmt.Printf("id ：%v: count ：%v\n", item.Key, item.DocCount)

			}
		}
	}
}