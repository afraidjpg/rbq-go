package app

import "sync"

func StartApp() {
	startListening()
	startTask()
	waiting()
}

func waiting() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

