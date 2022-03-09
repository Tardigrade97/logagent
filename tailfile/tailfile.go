package tailfile

import (
	"fmt"
	"github.com/hpcloud/tail"
)

var (
	TailObj *tail.Tail
)

func Init(filepath string) (err error) {
	// 初始化默认配置
	config := tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Follow:    true,
	}
	// 打开文件开始读取数据
	TailObj, err = tail.TailFile(filepath, config)
	if err != nil {
		fmt.Println("tail init failed,err:", err)
		return
	}
	return
}
