package app

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

type CronTask interface {
	GetTaskName() string
	GetCrontabString() string // 返回 crontab 风格的时间格式，如 "* * * * *"
	Entry()
}

type taskUint struct {
	EntryID cron.EntryID
	Task    CronTask
}

var taskMap []*taskUint

func AddTask(t []CronTask) {
	for _, task := range t {
		ts := &taskUint{Task: task}
		taskMap = append(taskMap, ts)
	}
}

var cr *cron.Cron

func startTask() {
	if len(taskMap) == 0 {
		return
	}
	cr = cron.New(cron.WithSeconds())

	for _, t := range taskMap {
		id, err := cr.AddFunc(t.Task.GetCrontabString(), t.Task.Entry)
		if err != nil {
			log.Fatal(fmt.Errorf("定时任务 [%s] 启动出错：%w", t.Task.GetTaskName(), err))
		}
		t.EntryID = id
	}

	cr.Start()
}