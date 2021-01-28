package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"runtime"
	"runtime/debug"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}

	for {
		//监听端口
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		// 开始goroutine监听连接
		Go(func() {
			handleConn(conn)
		})
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	message := make(chan string, 1)
	Go(func() {
		handleWrite(conn, ctx, message)
	})
	rd := bufio.NewReader(conn)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || strings.Contains(line, "EOF") { //停止
			log.Printf("read error: %v\n", err)
			break
		}
		message <- line
	}
	fmt.Println("reader done")
}

func handleWrite(conn net.Conn, ctx context.Context, msg chan string) {
	defer conn.Close()

	wr := bufio.NewWriter(conn)
	for {
		select {
		case <-ctx.Done():
			log.Printf("writer ctx err %+v", ctx.Err())
			log.Printf("writer done")
			log.Printf("Number of active goroutines %d", runtime.NumGoroutine())
			return
		case line := <-msg:
			wr.Write([]byte(line))
			wr.Flush()
		}
	}
}

//封装一个go关键字
func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("panic: err:%v, stack:%v\n", err, string(debug.Stack()))
			}
		}()
		f()
	}()
}
