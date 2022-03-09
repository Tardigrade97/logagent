package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func Init(address []string, chansize int64) (err error) {
	// 1.生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 2.连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		fmt.Println("producer closed,err:", err)
		return
	}
	// 初始化msgchan
	msgChan = make(chan *sarama.ProducerMessage, chansize)
	go sendMsg()
	return
}

//从通道msgchan中读取数据并发送到kafka中
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send message failed,err:", err)
				return
			}
			logrus.Infof("send msg to kafka sucess,pid:%v offset:%v", pid, offset)
		}
	}
}

func ToMsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
