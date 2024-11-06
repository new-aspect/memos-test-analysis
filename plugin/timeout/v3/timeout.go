package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

// 我们先写一个最简单的2个任务，第一个任务执行5s, 第二个任务执行10s, 整体任务超时时间设置为7秒，
// 我们希望得到的结果是整体任务执行7s，超时就退出了
func taskWithTimeout() error {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	// 第一个任务执行 5s
	eg.Go(func() error {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("任务 1 执行中...")
		}
		fmt.Println("任务 1 执行完成")
		return nil
	})

	// 第二个任务执行 10s
	eg.Go(func() error {
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("任务 2 执行中...")
		}
		fmt.Println("任务 2 执行完成")
		return nil
	})

	// 第三个协程来接受报错，通知关闭
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("过了超时时间 %s", err)
	}
	return nil
}

// 任务 2 执行中...
// 任务 1 执行中...
// 任务 2 执行中...
// 任务 1 执行中...
// 任务 1 执行中...
// 任务 2 执行中...
// 任务 2 执行中...
// 任务 1 执行中...
// 任务 1 执行中...
// 任务 1 执行完成
// 任务 2 执行中...
// 任务 2 执行中...
// 任务 2 执行中...
// 任务 2 执行中...
// 任务 2 执行中...
// 任务 2 执行中...
// 任务 2 执行完成
// 过了超时时间 context deadline exceeded
func main() {

	if err := taskWithTimeout(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("所有任务在超时前完成")
	}
}
