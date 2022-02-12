package kafka

import (
	"fmt"
	"github.com/hell-6/logTransfer/es"
	"time"

	"github.com/Shopify/sarama"
)

var (
	// 声明一个全局kafka消费者
	consumer sarama.Consumer
)

// 初始化kafka
func Init(addrs []string, topics []string) (err error) {
	fmt.Println("开始连接kafka")
	consumer, err = sarama.NewConsumer(addrs, sarama.NewConfig())
	if err != nil {
		return err
	}
	fmt.Println("kafka connection succeeded")
	for _, topic := range topics {
		fmt.Println(topic)
		go getMsg(topic)
	}
	return
}

// 读取消息
func getMsg(topic string) {
	// 获取topic对应的分区列表
	partitions, err := consumer.Partitions(topic)
	var partitionConsumer sarama.PartitionConsumer
	for _, partition := range partitions {
		// 创建每个分区的消费者
		partitionConsumer, err = consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("getMsg from kafka faild, err； %v\n", err)
			return
		}
		go func() {
			// 获取esChan
			esChan := es.GetEsChan()
			// 从分区消费者中取出消息发往es
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					{
						//fmt.Printf("Consumed message offset %d\n", msg.Offset)
						fmt.Printf("%v:Consumed message value: %v\n", topic, string(msg.Value))
						// 把消息包装成es.Msg结构体，然后存到esChan通道
						esChan <- &es.Msg{
							Topic: topic,
							Data:  string(msg.Value),
						}
					}
				default:
					time.Sleep(time.Millisecond * 50)
				}
			}
		}()
	}
}
