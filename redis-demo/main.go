package main

import (
	"fmt"
	"github.com/go-redis/redis"
)


var client *redis.Client

func init()  {
	client = redis.NewClient(&redis.Options{
		Addr:     "172.30.60.8:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func ping()  {
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}


func main() {

	ping()
	/*go func() {

		for {

			keyv()
		}
	}()

	go func() {
		for {

			keyv()
		}
	}()
*/

	//list()
	//hash()
	//pipeline()
	//set()
/*	sortset()
	time.Sleep(10e10)*/

}

func keyv()  {
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

func list()  {
	client.LPush("wida","ddddd")
	ret,_ := client.LPop("wida").Result()
	fmt.Println(ret)
	_,err := client.LPop("wida").Result()
	if err == redis.Nil {
		fmt.Println("empty list")
	}
}

func set()  {
	client.SAdd("set","wida")
	client.SAdd("set","wida1")
	client.SAdd("set","wida2")
	ret,_:= client.SMembers("set").Result()
	fmt.Println(ret)
}

func sortset()  {
	client.ZAdd("page_rank", redis.Z{10 ,"google.com"})
	client.ZAdd("page_rank", redis.Z{9 ,"baidu.com"},redis.Z{8 ,"bing.com"})
	ret,_:=client.ZRangeWithScores("page_rank",0,-1).Result()
	fmt.Println(ret)
}

func hash()  {
	client.HSet("hset","wida",1)
	ret ,_:=client.HGet("hset","wida").Result()
	fmt.Println(ret)
	ret ,err:=client.HGet("hset","wida3").Result()
	if err == redis.Nil {
		fmt.Println("key not found")
	}
	client.HSet("hset","wida2",2)
	r ,_:=client.HGetAll("hset").Result()
	fmt.Println(r)
}

func pipeline()  {
	pipe := client.Pipeline()
	pipe.HSet("hset2","wida",1)
	pipe.HSet("hset2","wida2",2)

	ret := pipe.HGetAll("hset2")
	fmt.Println(pipe.Exec())
	fmt.Println(ret.Result())
}