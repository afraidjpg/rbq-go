package example

import (
	"qq-robot-go/internal/msg"
)

type ExampleTask struct {
}

func (e *ExampleTask) GetTaskName() string {
	return "测试任务"
}
func (e *ExampleTask) GetCrontabString() string {
	// 秒级定时任务
	// 相比linux系统的crontab多了一位，最左位为秒，其他规则与linux的crontab规则一致
	// 当前为，每五秒执行一次
	return "*/5 * * * * *"
}

func (e *ExampleTask) Entry() {
	qq := -1
	msg.SendPrivateMsg(int64(qq), "测试消息")
}

func NewTask() *ExampleTask {
	return new(ExampleTask)
}
