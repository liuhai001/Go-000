package main

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func NewTracker() *Tracker {
	return &Tracker{
		ch:   make(chan string, 10),
		stop: make(chan struct{}),
	}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}
	fmt.Println("stop before")
	t.stop <- struct{}{}
	fmt.Println("stop after")

}

func (t *Tracker) ShutDown(ctx context.Context) {
	fmt.Println("close before")
	close(t.ch)
	fmt.Println("close after")

	fileName, line, functionName := "?", 0, "?"
	pc, fileName, line, ok := runtime.Caller(1)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		fmt.Println(functionName)
		functionName = filepath.Ext(functionName)
		fmt.Println(functionName)
		functionName = strings.TrimPrefix(functionName, ".")
		fmt.Println(functionName)
	}
	fmt.Printf("fileName:%v , line:%v, functionName:%v\n", fileName, line, functionName)

	select {
	case <-t.stop:
		fmt.Println("stop")
	case <-ctx.Done():
		fmt.Println("time out")
	}
}

func main() {
	tr := NewTracker()
	go tr.Run()
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	tr.ShutDown(ctx)
}
