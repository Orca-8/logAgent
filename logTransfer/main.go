package main

import (
	"fmt"
	"github.com/hell-6/logTransfer/conf"
	"github.com/hell-6/logTransfer/es"
	"github.com/hell-6/logTransfer/etcd"
	"github.com/hell-6/logTransfer/kafka"
	"github.com/hell-6/logTransfer/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"sync"
)

func main() {
	// 读取配置文件
	file, err := ioutil.ReadFile("./conf/conf.yaml")
	if err != nil {
		fmt.Printf("read file of config faild, err: %v\n", err)
		return
	}
	var confs conf.Conf
	yaml.Unmarshal(file, &confs)
	// 初始化etcd
	err = etcd.Init(confs.Etcd.Endpoint)
	if err != nil {
		fmt.Printf("init etcd faild, err: %v\n", err)
		return
	}
	// 从etcd读transfer_ip
	topics, err := etcd.Read(fmt.Sprintf(confs.Etcd.Key, utils.GetLocalIP()))
	// 初始化es
	err = es.Init(confs.Es.Address)
	if err != nil {
		log.Printf("init es faild, err: %v\n", err)
		return
	}
	// 初始化kafka
	err = kafka.Init(confs.Kafka.Addr, topics)
	if err != nil {
		fmt.Printf("init kafka faild, err； %v\n", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
