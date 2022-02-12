package test

import (
	"fmt"
	"github.com/hell-6/logTransfer/es"
	"testing"
)

func TestSendToEs(t *testing.T) {
	addrs := []string{"http://192.168.106.134:9200"}
	err := es.Init(addrs)
	if err != nil {
		fmt.Printf("es init faild, err: %v\n", err)
	}
	//es.SendToEs()
}
