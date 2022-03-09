package main

import (
	"Logagent/kafka"
	"Logagent/tailfile"
	"fmt"
	"github.com/Shopify/sarama"
	"gopkg.in/ini.v1"
)

type Config struct {
	KafkaConfig `ini:"kafka"`
	Logs        `ini:"logs"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `int:"chan_size"`
}

type Logs struct {
	Logpath string `ini:"log_path"`
}

func startService() (err error) {
	// tailfile get logs send to client
	for {
		line, ok := <-tailfile.TailObj.Lines
		if !ok {
			fmt.Println("tail file close reopen,filename:%s\n", tailfile.TailObj.Filename)
			continue
		}
		// 处理数据:利用通道将同步数据改为异步
		// 将读取的日志封装成kafka日志格式
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder(line.Text)
		// 将数据放入管道
		kafka.ToMsgChan(msg)
	}
	return
}

func main() {
	// 读取配置文件
	configObj := new(Config)
	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		fmt.Println("load config fialed,err:", err)
		return
	}

	// 连接kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.ChanSize)
	if err != nil {
		fmt.Println("kafka init filed,err:", err)
		return
	}
	fmt.Println("kafka init success")

	// 根据配置中的日志路径初始化
	err = tailfile.Init(configObj.Logpath)
	if err != nil {
		fmt.Println("tailfile init failed,err:", err)
		return
	}
	fmt.Println("init tailfile sucess!")

	// 获取并处理log
	err = startService()
	if err != nil {
		fmt.Println("startService failed!err:", err)
		return
	}
}
