package etcd

import (
	"Logagent/commonstruct"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"time"
)

var (
	client *clientv3.Client
)

//初始化
func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Println("Connect etcd failed,err:", err)
		return
	}
	return
}

// 拉取日志收集配置项
func GetConfig(key string) (logsEntryList []commonstruct.LogsEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	rsp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get config from etcd by key:%s failed,err", key, err)
		return
	}

	if len(rsp.Kvs) == 0 {
		logrus.Warnf("get len:0 conf from etcd by key:%s", key)
		return
	}
	ret := rsp.Kvs[0]
	err = json.Unmarshal(ret.Value, &logsEntryList)
	if err != nil {
		logrus.Errorf("json unmarshal failed,err:%v", err)
		return
	}
	return
}
