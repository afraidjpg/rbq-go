package rbq

import (
	"fmt"
	"testing"
)

func TestCQSend(t *testing.T) {
	cqs := newCQSend()
	cqs.AddCQFace(2)
	cqs.AddText("hello")
	cqs.AddText(" world")

	fmt.Println(cqs.cq)
	fmt.Println(cqs.cqm)
}
