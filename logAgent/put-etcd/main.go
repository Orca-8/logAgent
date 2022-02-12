package main

import (
	"fmt"
	"logAgent/etcd"
)

func main() {
	// 初始化etcd
	err := etcd.Init([]string{"192.168.106.134:2379"})
	if err != nil {
		fmt.Printf("etcd init faild, err:%v\n", err)
		return
	}
	etcd.Put()
}
