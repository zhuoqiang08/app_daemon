package main

import (
	"encoding/json"
	"log"

	"gopkg.in/ini.v1"
)

type MQ_Event struct {
	Id   string `json:"id"` //uuid
	Type string `json:"type"`
	Data string `json:"data"`
}
type MQ_Ret struct {
	Id      string `json:"id"`
	success bool   `json:"success`
	Msg     string `json:"msg"`   //cmd Ret
	Error   string `json:"error"` //cmd error
}

//命令执行后 返回
func Msg_Handle(content string) {
	msg := MQ_Event{}
	log.Println("MSG_HANDLE", content)
	json.Unmarshal([]byte(content), &msg)

	if msg.Type == "cmd" {
		ExecCmd_Nsq(msg.Data, msg.Id, nsq_client)
	}
	if msg.Type == "cmd_batch" {
		ExecCmd_File_Nsq(msg.Data, msg.Id, nsq_client)
	}

}

var nsq_client *Nsq_Service

const RET_TOPIC = "AGENT_RET"

func Nsq_Init() {
	//nsq配置以及
	//消息标题为本机的IP
	if nsq_client == nil {
		nsq_client = &Nsq_Service{}
	}
	conf, err := ini.Load("app.conf") //加载配置文件
	if err != nil {
		log.Println("load config file fail!")
		return
	}
	ip, _ := conf.Sections()[0].GetKey("nsq")
	agent_id, _ := conf.Sections()[0].GetKey("agent_id")
	log.Println("Nsq", ip.String(), agent_id)
	nsq_client.Init(ip.String())
	nsq_client.InitProducer()
	//go nsq_client.InitConsumer("agent_"+agent_id.String(), "agent", Msg_Handle)
	go nsq_client.InitConsumer(RET_TOPIC, "agent", Msg_Handle)
}

//通过http配置启动参数。
//通过nsq来执行指定命令行.
func main() {
	Nsq_Init()
	endRunning := make(chan bool, 1)
	<-endRunning
	println("-----")
}
