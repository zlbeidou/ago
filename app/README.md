# 概述
app 模块是快速搭建程序的简易框架, 将程序划分成若干 service 简化程序设计.
service 注册之后只需要关注各个 service 即可, 框架会监听系统信号管理 service 的生命周期.

# 术语
## service
包括 onceService 和 backendService
### onceService
只有 Init 方法用于初始化, Init 执行完成服务即结束.
### backendService
含有 Init 和 Run 方法, Init 用于初始化, Run 是服务是实际运行逻辑.
每个 backendService 会启一个 goroutine 来运行 Run.
Run 可以监听 ctx.Done 判断程序是否收到结束信号, 结束前可以处理资源回收的事情.
## rear
service 结束后, 程序退出前执行的殿后服务, 主要用于资源回收
存在 rear 的原因是应对有些工具模块本身比较轻量, 逻辑上不希望当做 service 来看, 但是仍然需要在程序结束前进行资源回收的情况. 例如一些单例模式的工具, 在程序结束前希望进行资源回收.

# 使用方法
具体的案例可以参考 example.
## main.go
```
app.Init()
app.Run()
```
执行 app.Init() 和 app.Run() 即可.

## service 注册
app.RegisterService 可以进行服务注册.
通常在 init() 函数进行服务注册, 在 main.go 中以匿名导入的方式完成注册.

## rear 注册
### 注册 rear 函数
直接注册 rear 函数
### 将注册耦合在逻辑中
难以将资源回收提炼成函数的情况, 例如回收一个 goroutine.
```
app.RearStarted()
go func() {
    defer app.RearStopped()

    for {
        select {
        case <-app.Done():
            return
        case sth:
            // do sth
        }
    }
}()
```
