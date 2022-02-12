package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

var (
	// 声明一个全局kafka生产者
	producter sarama.SyncProducer
	// 声明一个缓存日志的通道
	dataChan chan *dataOfChan
)

type dataOfChan struct {
	topic, data string
}

// 初始化kafka
func Init(addrs []string, chanMaxSize int) (err error) {
	fmt.Println("开始连接kafka")
	// 配置配置producter
	config := sarama.NewConfig()
	// 返回成功交付的消息
	config.Producer.Return.Successes = true
	// 所有leader和follower都返回ack后才继续生产
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 选择partition的方式为轮询
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	producter, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return err
	}
	fmt.Println("kafka connection succeeded")
	dataChan = make(chan *dataOfChan, chanMaxSize)
	go sendToKafka()
	return
}

// 发送信息
func Send(topic, data string) (err error) {
	// 封装消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	// 发送到kafka
	_, _, err = producter.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Println("send message succeeded")
	return
}

// 发送消息到kafka
func sendToKafka() {
	for {
		select {
		case msg := <-dataChan:
			{
				err := Send(msg.topic, msg.data)
				if err != nil {
					fmt.Printf("send to kafka faild, err: %v\n", err)
					return
				}
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

// 发送消息到缓存通道中
func SendToChan(topic, data string) {
	dataChan <- &dataOfChan{
		topic: topic,
		data:  data,
	}
}
