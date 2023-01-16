package rbq

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCQAt_To(t *testing.T) {
	at := NewCQAt()
	at.To(123456)
	assert.Equal(t, at.String(), "[CQ:at,qq=123456]")

	at.To(1000, 1001, 1002)
	assert.Equal(t, at.String(), "[CQ:at,qq=1000][CQ:at,qq=1001][CQ:at,qq=1002]")

	at.To(0, 2000, 2001)
	assert.Equal(t, at.String(), "[CQ:at,qq=all]")

	at.To(0, -1, 3001, 3002)
	assert.Equal(t, at.String(), "[CQ:at,qq=all]")

	at.To(22, 0)
	assert.Equal(t, at.String(), "[CQ:at,qq=all]")
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
