package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"logAgent/conf"
	"logAgent/utils"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	cli *clientv3.Client
)

// 初始化etcd
func Init(endpoints []string) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		return
	}
	//defer cli.Close()
	fmt.Println("etcd init succeeded")
	return
}

// 从etcd读取（get）配置
func Read(key string) (confs []*conf.TailConf, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		return nil, err
	}
	for _, v := range resp.Kvs {
		// fmt.Printf("key:%v   value:%v\n", k, string(v.Value))
		if v.Value == nil {
			// key对应的配置为空
			return nil, err
		}
		err = json.Unmarshal(v.Value, &confs)
		if err != nil {
			return nil, err
		}
	}
	return
}

// 向etcd中put
func Put() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	str := `[{"path":"./log/mysql.log","topic":"mysql_log"},{"path":"./log/redis.log","topic":"redis_log"}]`
	//str := `mysql_log&redis_log`
	_, err := cli.Put(ctx, fmt.Sprintf("collection_%s", utils.GetLocalIP()), str)
	//_, err := cli.Put(ctx, fmt.Sprintf("transfer_%s", utils.GetLocalIP()), str)
	cancel()
	if err != nil {
		fmt.Printf("etcd put faild, err:%v\n", err)
		return
	}
}

// 监控etcd
func WatchConf(key string, newConfChan chan<- []*conf.TailConf) {
	wc := cli.Watch(context.Background(), key)
	for watchChan := range wc {
		for _, ev := range watchChan.Events {
			//fmt.Printf("event: type:%v  key:%v  value:%v\n", ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
			var newConf []*conf.TailConf
			if ev.Kv.Value == nil {
				fmt.Println("配置为空")
			} else {
				err := json.Unmarshal(ev.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("json unmarshal faild, err: %v\n", err)
					return
				}
			}
			newConfChan <- newConf
		}
	}
}
