package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nsqio/go-nsq"
)

type Nsq_Service struct {
	ip       string
	producer *nsq.Producer
	customer *nsq.Consumer
}

func (this *Nsq_Service) Init(ip string) {
	this.ip = ip
}

// 初始化生产者
func (this *Nsq_Service) InitProducer() {
	var err error
	if this.producer != nil { //重新加载
		this.producer.Stop()
	}
	this.producer, err = nsq.NewProducer(this.ip, nsq.NewConfig())
	if err != nil {
		fmt.Println(err)
	}
}

//发布消息
func (this *Nsq_Service) Publish(topic string, message string) error {
	var err error

	if this.producer != nil {

		if message == "" { //不能发布空串，否则会导致error
			return nil
		}
		err = this.producer.Publish(topic, []byte(message)) // 发布消息
		return err
	}
	return nil
}

type Nsq_Callback func(msg string)
type ConsumerT struct {
	callback Nsq_Callback
}

func (this *ConsumerT) HandleMessage(msg *nsq.Message) error {
	log.Println("receive", msg.NSQDAddress, "message:", string(msg.Body))
	this.callback(string(msg.Body))
	return nil
}

func (this *ConsumerT) SetCallBack(call Nsq_Callback) {
	this.callback = call

}
func (this *Nsq_Service) InitConsumer(topic string, channel string, handeler Nsq_Callback) {
	if this.customer != nil {
		//重新加载
		this.customer.Stop()
	}
	cfg := nsq.NewConfig()
	//cfg.MaxInFlight = 255
	//	cfg.MsgTimeout = time.Second * 10
	cfg.LookupdPollInterval = time.Second          //设置重连时间
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	this.customer = c
	if err != nil {
		log.Println(err)
		return
	}
	this.customer.SetLogger(nil, 0) //屏蔽系统日志
	cust := ConsumerT{}
	cust.SetCallBack(handeler)
	this.customer.AddHandler(&cust) // 添加消费者接口
	ip := strings.Split(this.ip, ",")
	//建立多个nsqd连接
	if err := this.customer.ConnectToNSQDs(ip); err != nil {
		panic(err)
	}

}
