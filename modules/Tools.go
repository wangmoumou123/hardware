package modules

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type ExitCallBack func()

// ExitHandle exit handle 设置信号处理和清理工作
func ExitHandle(callback ExitCallBack) chan struct{} {
	// 创建一个通道，用于接收操作系统信号
	sigChan := make(chan os.Signal, 1)

	// 将指定的信号转发到 sigChan
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个通道用于通知外部代码程序退出
	done := make(chan struct{})

	// 使用 defer 注册清理函数
	go func() {
		defer callback()
		sig := <-sigChan
		fmt.Printf("捕获到信号: %s\n", sig)
		close(done)
	}()
	return done
}
