package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
	基于Zinx框架来开发的 服务器端应用程序
 */

/*
	ping  test 自定义路由
 */
type PingRouter struct {
	znet.BaseRouter
}

//Test PreRouter
func (this *PingRouter)PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}
//Test Handle
//在处理conn业务的主方法hook
func (this *PingRouter)Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping...ping...\n"))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}
//Test PostHandle
//在处理conn业务之后的钩子方法Hook
func (this *PingRouter)PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}


func main() {
	//1 创建一个server句柄，使用 Zinx的api
	s := znet.NewServer("[zinx V0.2]")

	//2 给当前 Zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	//3 启动server
	s.Server()

}
