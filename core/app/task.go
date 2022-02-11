package app

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

// CronTask 定时任务接口，用户的定时任务都必须实现本接口
type CronTask interface {
	GetTaskName() string // 定时任务的名称
	GetCrontabString() string // 返回 crontab 风格的时间格式，如 "* * * * *"
	Entry() // 定时任务的入口函数
}

// 定时任务单元
type taskUint struct {
	EntryID cron.EntryID // 定时任务ID
	Task    CronTask // 定时任务本体
}

// 定时任务队列
var taskQueue []*taskUint

func AddTask(t []CronTask) {
	for _, task := range t {
		ts := &taskUint{Task: task}
		taskQueue = append(taskQueue, ts)
	}
}

// 定时任务总对象
var cr *cron.Cron

func startTask() {
	if len(taskQueue) == 0 {
		return
	}
	cr = cron.New(cron.WithSeconds())

	for _, t := range taskQueue {
		id, err := cr.AddFunc(t.Task.GetCrontabString(), t.Task.Entry)
		if err != nil {
			log.Fatal(fmt.Errorf("定时任务 [%s] 启动出错：%w", t.Task.GetTaskName(), err))
		}
		t.EntryID = id
	}

	cr.Start()
}