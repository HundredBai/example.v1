package scheduler

import (
	"net/http"
)

type Scheduler interface {
	// Init 用于初始化调度器。
	// 参数 requestArgs 代表请求相关的参数
	// 参数 dataArgs 代表数据相关的参数
	// 参数 componentArgs 代表组件相关的参数
	Init(requestArgs RequestArgs,
		dataArgs DataArgs,
		componentArgs ComponentArgs) error
	// Start 用于启动调度器并执行爬取流程
	// 参数 firstHTTPReq 即代表首次请求，调度器会以此为起始点开始执行爬取流程
	Start(firstHTTPReq *http.Request) error
	// Stop 用于停止调度器的运行
	// 所有处理模块执行的流程都会被中止
	Stop() error
	// Status 用于获取调度器的状态
	Status() Status
	// ErrorChan 用于获得错误通道
	// 调度器以及各个处理模块运行过程中出现的所有错误都会被发送到该通道
	// 若结果值为 nil，则说明错误通道不可用或调度器已被停止
	ErrorChan() <-chan error
	// Idle 用于判断所有处理模块是否都处于空闲状态
	Idle() bool
	// Summary 用于获取摘要实例
	Summary() SchedSummary
}

type SchedulerImpl struct {

}
