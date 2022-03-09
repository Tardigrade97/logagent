package main

import (
	"Logagent/etcd"
	"Logagent/kafka"
	"Logagent/tailfile"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	KafkaConfig `ini:"kafka"`
	Logs        `ini:"logs"`
	EtcdConfig  `int:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `int:"chan_size"`
}

type Logs struct {
	Logpath string `ini:"log_path"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

func starService() {
	select {}
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

	// 初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed,err:", err)
		return
	}
	// 从etcd中拉取要收集的配置项目
	allConf, err := etcd.GetConfig(configObj.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed,err:", err)
		return
	}
	fmt.Println(allConf)

	// 根据配置中的日志路径初始化
	err = tailfile.Init(allConf) //把etcd中获取的配置项传入init
	if err != nil {
		logrus.Errorf("tailfile init failed,err:", err)
		return
	}
	fmt.Println("init tailfile sucess!")
	starService()
}
