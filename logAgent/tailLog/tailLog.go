package tailLog

import (
	"context"
	"fmt"
	"logAgent/kafka"
	"time"

	"github.com/hpcloud/tail"
)

// tail任务
type TailTask struct {
	Path    string
	Topic   string
	TailObj *tail.Tail
	Ctx     context.Context
	Cancel  context.CancelFunc
}

// 创建tail任务
func NewTailObj(path, topic string) (tailTask *TailTask, err error) {
	// 配置tail
	config := &tail.Config{
		// 重新打开
		ReOpen: true,
		// 是否跟随
		Follow: true,
		// 从哪儿开始读
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		// 当文件不存在时是否报错
		MustExist: false,
		Poll:      true,
	}
	tailObj, err := tail.TailFile(path, *config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	tailTask = &TailTask{
		Path:    path,
		Topic:   topic,
		TailObj: tailObj,
		Ctx:     ctx,
		Cancel:  cancel,
	}
	fmt.Println("tail create succeeded")
	go tailTask.SendToChan()
	return
}

// 发送消息到缓存通道中
func (tailTask *TailTask) SendToChan() error {
	logChan := tailTask.TailObj.Lines
	for {
		select {
		// 如果有日志，则发送到缓存区,否则停一秒再继续
		case line := <-logChan:
			/* err := Send(tailTask.Topic, line.Text)
			if err != nil {
				return err
			} */
			kafka.SendToChan(tailTask.Topic, line.Text)
		case <-tailTask.Ctx.Done():
			return nil
		default:
			time.Sleep(time.Second)
		}
	}
}
