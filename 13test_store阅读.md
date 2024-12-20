# test/store阅读笔记

```
├── store
│   ├── memo_test.go
│   ├── store.go
│   ├── system_setting_test.go
│   └── user_test.go
└── test.go
```


我们先看momo_test.go里面有一个TestMemoStore，首先我们整理他的意图，他做一下事情
1. 打开测试数据库
2. 在测试数据库创建一名用户
3. 调用数据库的创建Memo函数
4. 调用数据库的修改Memo函数
5. 调用数据库的查询Memo函数
6. 调用数据库的删除Memo函数

测试数据库是获取测试的配置，根据测试的配置打开数据库，我注意到store和db是两个结构，所以将一个细节，就是在store目录有func New(db *sql.DB) *Store这样的函数被调用了3次，1次是测试时创建一个测试数据库给他，1次是服务启动时调用，一次是命令行选择setup 时可以初始化数据库，他通过这种方法解耦。
同样我注意到他的Server里面也有Store，也就是他调用api的时候可以向下面这样描述他的操作数据库的意图

```go
func (s *Server) registerSystemRoutes(g *echo.Group) {
    g.GET("/status", func(c echo.Context) error {
        hostUser, err := s.Store.FindUser(ctx, &hostUserFind)
    }
}
```

然后我发现他创建或删除的时候，传入的是一个结构体，而且这样的写法感觉更容易阅读
```go
memoCreate := &api.MemoCreate{
    CreatorID:  user.ID,
    Content:    "test_content",
    Visibility: api.Public,
}
memo, err := store.CreateMemo(ctx, memoCreate)
```