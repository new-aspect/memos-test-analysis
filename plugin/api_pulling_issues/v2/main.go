package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"strings"
	"time"
)

var (
	summaryTaskCounts = 100             // 总任务两是100
	concurrencyCount  = 10              // 并发数是10
	timeout           = 3 * time.Minute // 超时时间是3分钟
)

func pullingLogic(ctx context.Context, taskId int) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			result := getTaskResult(taskId)
			if strings.Contains(result, "任务未完成") {
				time.Sleep(5 * time.Second)
				continue
			}
			fmt.Printf("任务ID %d 完成\n", taskId)
			return nil
		}
	}
}

func getTaskResult(taskId int) string {
	rand.Seed(time.Now().UnixNano() + int64(taskId))
	if rand.Intn(2) == 0 {
		return "任务未完成"
	}

	// 这里是模拟第三方接口偶尔会出现的一个接口无论如何都执行不完，即使超时24小时
	if taskId == 22 || taskId == 56 || taskId == 77 {
		return "任务未完成"
	}
	return "任务完成"
}

func main() {
	// 总任务循环
	for i := 0; i < summaryTaskCounts; i += concurrencyCount {
		// 并发执行拉去
		// 每次重试都创建一个带有超时的 context
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// 使用 errorgroup 来管理这次重试的协程
		eg, ctx := errgroup.WithContext(ctx)
		for j := 0; j < concurrencyCount; j++ {
			taskId := i + j // 捕获循环变量
			eg.Go(func() error {
				return pullingLogic(ctx, taskId)
			})
		}
		if err := eg.Wait(); err != nil {
			fmt.Printf("并发执行出错: %v \n", err)
		}
	}
	fmt.Printf("所有任务执行完成")
}
