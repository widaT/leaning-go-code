package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/etcd-io/etcd/clientv3"
)

var client *clientv3.Client

func init() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"172.30.60.8:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
}

//kv操作
func kv() {
	_, err := client.Put(context.Background(), "sample_key", "sample_value")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := clientv3.NewKV(client).Get(context.Background(), "sample_key")
	if n := len(resp.Kvs); n == 0 {
		fmt.Println("not found")
		return
	}
	fmt.Println(string(resp.Kvs[0].Value))

	client.Delete(context.Background(), "sample_key")

	resp, err = clientv3.NewKV(client).Get(context.Background(), "sample_key")
	if n := len(resp.Kvs); n == 0 {
		fmt.Println("not found")
		return
	}
	fmt.Println(string(resp.Kvs[0].Value))
}

//监听 watch
func watch() {
	go func() {
		for {
			ch := client.Watch(context.Background(), "sample_key")
			for wresp := range ch {
				for _, ev := range wresp.Events {
					fmt.Println(ev)
				}
			}
		}
	}()
	_, err := client.Put(context.Background(), "sample_key", "sample_value1")
	if err != nil {
		log.Fatal(err)
	}
	client.Delete(context.Background(), "sample_key")
}

//租约
func lease() {
	lease := clientv3.NewLease(client)
	//设置10秒租约 （过期时间为10秒）
	leaseResp, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		log.Fatal(err)
	}

	leaseID := leaseResp.ID
	fmt.Printf("ttl %d \n", leaseResp.TTL)

	keepRespChan, err := lease.KeepAlive(context.TODO(), leaseID)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
	outer:
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约失效了")
					break outer
				} else {

					fmt.Printf("续租成功 time:%d id:%d,ttl:%d \n", time.Now().Unix(), keepResp.ID, keepResp.TTL)
				}

			}
		}
	}()

	kv := clientv3.NewKV(client)
	putResp, err := kv.Put(context.TODO(), "sample_key_lease", "", clientv3.WithLease(leaseID))
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("写入成功", putResp.Header.Revision)
	}
}

func main() {

	//kv()
	//watch()

	lease()

	time.Sleep(20e9)
	defer client.Close()
}
