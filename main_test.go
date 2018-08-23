package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Nsq() {
	nsq := Nsq_Service{}
	nsq.Init("127.0.0.1:4150")
	nsq.InitProducer()
	msg := MQ_Event{}
	msg.Type = "cmd_batch"
	msg.Data = `ps -ef
	curl http://www.sohu.com`
	c, _ := json.Marshal(&msg)
	fmt.Println(string(c))
	for i := 0; i < 100; i++ {
		//	nsq.Publish("agent_1", fmt.Sprintf("test%d", i))
		nsq.Publish("agent_1", string(c))
	}
}
func Exec() {

}

func TestGet(test *testing.T) {
	Nsq()
	//fmt.Println("-----")
}
