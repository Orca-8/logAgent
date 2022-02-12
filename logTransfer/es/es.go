package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strconv"
	"sync"
	"time"
)

type Msg struct {
	Topic string
	Data  string
}

var (
	client *elastic.Client
	wg     sync.WaitGroup
	esChan chan *Msg
)

// Init 初始化es
func Init(addrs []string) (err error) {
	// 创建client
	client, err = elastic.NewClient(elastic.SetURL(addrs[0]))
	if err != nil {
		return
	}
	fmt.Println("elasticsearch client init succeeded")
	esChan = make(chan *Msg, 100000)
	go sendToEs()
	return
}

// sendToEs 发送消息到es
func sendToEs() {
	for {
		select {
		case msg := <-esChan:
			{
				// 查看目标index--log是否存在，不存在则创建
				exists, err := client.IndexExists("log").Do(context.Background())
				if err != nil {
					// Handle error
					panic(err)
				}
				if !exists {
					mapping := `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"doc":{
			"properties":{
				"topic":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				}
			}
		}
	}
}
`
					_, err := client.CreateIndex("log").Body(mapping).IncludeTypeName(true).Do(context.Background())
					if err != nil {
						// Handle error
						panic(err)
					}
				}
				put1, err := client.Index().
					Index("log").
					Id(strconv.Itoa(int(time.Now().Unix()))).
					BodyJson(Msg{
						Topic: msg.Topic,
						Data:  msg.Data,
					}).
					Do(context.Background())
				if err != nil {
					// Handle error
					panic(err)
				}
				fmt.Printf("Indexed test %v to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

}

// GetEsChan 获取esChan
func GetEsChan() chan *Msg {
	return esChan
}
