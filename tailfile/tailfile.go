package tailfile

import (
	"Logagent/commonstruct"
	"Logagent/kafka"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"strings"
)

type tailtask struct {
	path  string
	topic string
	tObj  *tail.Tail
}

func newTailTask(path, topic string) *tailtask {
	tt := &tailtask{
		path:  path,
		topic: topic,
	}
	return tt
}

func (t *tailtask) Init() (err error) {
	// 初始化默认配置
	cfg := tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Follow:    true,
	}
	t.tObj, err = tail.TailFile(t.path, cfg)
	return
}

func (t *tailtask) startCollectLogs() {
	logrus.Info("collect for path:%s is running ...", t.path)
	// tailfile get logs send to client
	for {
		line, ok := <-t.tObj.Lines
		if !ok {
			fmt.Println("tail file close reopen,filename:%s\n", t.path)
			continue
		}

		// 空行就省略
		if len(strings.Trim(line.Text, "\r")) == 0 {
			logrus.Info("出现空行！")
			continue
		}

		// 处理数据:利用通道将同步数据改为异步
		// 将读取的日志封装成kafka日志格式
		msg := &sarama.ProducerMessage{}
		msg.Topic = t.topic
		msg.Value = sarama.StringEncoder(line.Text)
		// 将数据放入管道
		kafka.ToMsgChan(msg)
	}
}

func Init(allconfig []commonstruct.LogsEntry) (err error) {

	for _, conf := range allconfig {
		// 创建一个日志收集任务
		tt := newTailTask(conf.Path, conf.Topic)
		err := tt.Init()
		if err != nil {
			logrus.Errorf("create tailObj failed,path:%s", conf.Path)
			continue
		}

		logrus.Info("create a tail task for path:%s", conf.Path)
		// 启动日志收集
		go tt.startCollectLogs()

	}
	return
}
