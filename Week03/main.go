package main
import (
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//不用errgroup 返回的ctx
	g, _ := errgroup.WithContext(ctx)
	//http1
	http1 := NewHttpServer(":8080")
	//http2
	http2 := NewHttpServer(":8081")

	//启动1
	g.Go(func() error {
		if err := http1.start(); err != nil {
			//一个退出，要全部退出，调用cancel()
			cancel()
			return err
		}
		return nil
	})

	//启动2
	g.Go(func() error {
		if err := http2.start(); err != nil {
			//一个退出，要全部退出，调用cancel()
			cancel()
			return err
		}
		return nil
	})

	// 监听sig信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
			select {
			case s := <-quit:
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					cancel()
				default:
				}
			}
	}()

	// context取消后，关闭http server
	go func() {
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			//等待http server shutdown
			if err := http1.shutdown(ctx); err != nil {
				log.Println("http1 shutdown err: ", err)
			}
			if err := http2.shutdown(ctx); err != nil {
				log.Println("http2 shutdown err: ", err)
			}
		}

	}()

	//等待所有http server退出
	if err := g.Wait(); err != nil {
		log.Println("all exit: ", err)
	}
}

type httpServer struct {
	server http.Server
}

func NewHttpServer(addr string) *httpServer {
	return &httpServer{
		server: http.Server{
			Addr: addr,
		},
	}
}

func (h *httpServer) start() error {
	log.Println("http server start!")
	return h.server.ListenAndServe()
}

func (h *httpServer) shutdown(ctx context.Context) error {
	log.Println("http server shutdown!")
	return h.server.Shutdown(ctx)
}