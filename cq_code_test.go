package rbq

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCQAt_To(t *testing.T) {
	at := NewCQAt()
	at.To(123456)
	assert.Equal(t, at.String(), "[CQ:at,qq=123456,name=123456]")

	at.To(1000, 1001, 1002)
	assert.Equal(t, at.String(), "[CQ:at,qq=1000,name=1000][CQ:at,qq=1001,name=1001][CQ:at,qq=1002,name=1002]")

	at.To(0, 2000, 2001)
	assert.Equal(t, at.String(), "[CQ:at,qq=all,name=全体成员]")

	at.To(0, -1, 3001, 3002)
	assert.Equal(t, at.String(), "[CQ:at,qq=all,name=全体成员]")
	//
	at.To(22, 0)
	assert.Equal(t, at.String(), "[CQ:at,qq=all,name=全体成员]")

	at.ToWithNotExistName([]string{"不存在", "不存在2号"}, []int64{4001, 4002})
	assert.Equal(t, at.String(), "[CQ:at,qq=4001,name=不存在][CQ:at,qq=4002,name=不存在2号]")

	at.ToWithNotExistName([]string{"不存在", "不存在2号", ""}, []int64{4001, 4002, 0})
	assert.Equal(t, at.String(), "[CQ:at,qq=all,name=全体成员]")

	at.ToWithNotExistName([]string{"不存在", "不存在2号"}, []int64{4001, 4002, 5001})
	assert.Equal(t, at.Errors()[0].Error(), "at: name和userId长度不一致")
}

func TestCQFace_Id(t *testing.T) {
	face := NewCQFace()
	face.Id(1)
	assert.Equal(t, face.String(), "[CQ:face,id=1]")

	face.Id(-1)
	assert.Equal(t, face.String(), "")
	assert.Equal(t, face.Errors()[0].Error(), "face: id 必须在 0-221 之间")

	face.Id(1, 5, 114)
	assert.Equal(t, face.String(), "[CQ:face,id=1][CQ:face,id=5][CQ:face,id=114]")

	face.Id(222)
	assert.Equal(t, face.String(), "")
	assert.Equal(t, face.Errors()[0].Error(), "face: id 必须在 0-221 之间")

	face.Id(13, 2, -1)
	assert.Equal(t, face.String(), "")
	assert.Equal(t, face.Errors()[0].Error(), "face: id 必须在 0-221 之间")
}

func TestCQRecord_File(t *testing.T) {
	record := NewCQRecord()
	record.File("http://downsc.chinaz.net/Files/DownLoad/sound1/201906/11582.mp3")
	assert.Equal(t, record.String(), "[CQ:record,file=http://downsc.chinaz.net/Files/DownLoad/sound1/201906/11582.mp3]")

	record.AllOption("rename.mp3", 0, "http://downsc.chinaz.net/Files/DownLoad/sound1/201906/11582.mp3", -1, -1, 30)
	assert.Equal(t, record.String(), "[CQ:record,file=rename.mp3,magic=0,url=http://downsc.chinaz.net/Files/DownLoad/sound1/201906/11582.mp3,timeout=30]")
}
