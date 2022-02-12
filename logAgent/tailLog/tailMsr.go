package tailLog

import (
	"fmt"
	"logAgent/conf"
	"time"
)

var taskMsr *TaskMsr

type TaskMsr struct {
	//confs    []*conf.TailConf
	TailObjs    map[string]*TailTask
	NewConfChan chan []*conf.TailConf
}

// 创建tail管理者
func CreateTailMsr(confs []*conf.TailConf) (*TaskMsr, error) {
	taskMsr = &TaskMsr{
		TailObjs:    make(map[string]*TailTask, 16),
		NewConfChan: make(chan []*conf.TailConf),
	}
	for _, v := range confs {
		tailTask, err := NewTailObj(v.Path, v.Topic)
		if err != nil {
			return nil, err
		}
		taskMsr.TailObjs[fmt.Sprintf("%s_%s", v.Topic, v.Path)] = tailTask
	}
	go taskMsr.watchNewConfChan()
	return taskMsr, nil
}

// 监控NewConfChan
func (taskMsr *TaskMsr) watchNewConfChan() {
	for {
		select {
		// 有了新配置
		case newConf := <-taskMsr.NewConfChan:
			{
				var newKeys []string
				// 如果新配置为空，则将所有tailObj关闭
				if newConf == nil {
					for _, o := range taskMsr.TailObjs {
						taskMsr.TailObjs[fmt.Sprintf("%s_%s", o.Topic, o.Path)].Cancel()
						delete(taskMsr.TailObjs, fmt.Sprintf("%s_%s", o.Topic, o.Path))
					}
					continue
				}
				for _, v := range newConf {
					newKeys = append(newKeys, fmt.Sprintf("%s_%s", v.Topic, v.Path))
					_, ok := taskMsr.TailObjs[fmt.Sprintf("%s_%s", v.Topic, v.Path)]
					if !ok {
						// 没有这个tail，就创建
						tailObj, err := NewTailObj(v.Path, v.Topic)
						if err != nil {
							fmt.Printf("add tail faild, err；%v\n", err)
							return
						}
						taskMsr.TailObjs[fmt.Sprintf("%s_%s", v.Topic, v.Path)] = tailObj
					}
				}

				// 如果当前有的tailobj不在新配置中，则删除
				for curkey := range taskMsr.TailObjs {
					var isDelete = true
					for _, newkey := range newKeys {
						if curkey == newkey {
							isDelete = false
						}
					}
					if isDelete {
						taskMsr.TailObjs[curkey].Cancel()
						delete(taskMsr.TailObjs, curkey)
					}
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func GetNewConfChan() chan<- []*conf.TailConf {
	return taskMsr.NewConfChan
}
