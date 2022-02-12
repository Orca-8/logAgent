package main

import (
	"fmt"
	"io/ioutil"
	"logAgent/conf"
	"logAgent/etcd"
	"logAgent/kafka"
	"logAgent/tailLog"
	"logAgent/utils"
	"sync"

	"gopkg.in/yaml.v2"
)

func main() {
	// 读配置文件
	b, err := ioutil.ReadFile("./conf/conf.yaml")
	if err != nil {
		fmt.Println("read config file faild, err:", err)
		return
	}
	cf := conf.Conf{}
	// ymal解析
	yaml.Unmarshal(b, &cf)
	// 初始化kafka
	err = kafka.Init(cf.Kafka.Addr, cf.Kafka.ChanMaxSize)
	if err != nil {
		fmt.Printf("kafka init faild, err:%v\n", err)
		return
	}

	// 初始化etcd
	err = etcd.Init(cf.Etcd.Endpoint)
	if err != nil {
		fmt.Printf("etcd init faild, err:%v\n", err)
		return
	}

	// 从etcd中根据配置文件取出key对应的配置
	var confs []*conf.TailConf
	confs, err = etcd.Read(fmt.Sprintf(cf.Etcd.Key, utils.GetLocalIP()))
	fmt.Printf("localkey: %v\n", fmt.Sprintf(cf.Etcd.Key, utils.GetLocalIP()))
	if err != nil {
		fmt.Printf("read from etcd faild, err:%v\n", err)
		return
	}

	// 根据confs中的配置创建任务管理者
	//var taskMsr *tailLog.TaskMsr
	_, err = tailLog.CreateTailMsr(confs)
	if err != nil {
		fmt.Printf("create taskMsr faild, err: %v\n", err)
		return
	}

	// 监听etcd
	go etcd.WatchConf(fmt.Sprintf(cf.Etcd.Key, utils.GetLocalIP()), tailLog.GetNewConfChan())

	for _, v := range confs {
		fmt.Printf("[conf]--path: %v  topic: %v\n", v.Path, v.Topic)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}
