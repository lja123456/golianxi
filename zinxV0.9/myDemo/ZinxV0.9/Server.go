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
//func (this *PingRouter)PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_,err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
//	if err != nil {
//		fmt.Println("call back before ping error")
//	}
//}
//Test Handle
//在处理conn业务的主方法hook
func (this *PingRouter)Handle(request ziface.IRequest) {
	fmt.Println("Call ping Router Handle...")
	//_,err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping...ping...\n"))
	//if err != nil {
	//	fmt.Println("call back ping...ping...ping error")
	//}
	//先读取客户端的数据，再回写ping..ping..ping
	fmt.Println("recv from client: msgID = ",request.GetMsgID(),
		", data = ",string(request.GetData()))

	err := request.GetConnection().SendMsg(200,[]byte("ping..pign..ping"))
	if err != nil {
		fmt.Println(err)
	}
}
//Test PostHandle
//在处理conn业务之后的钩子方法Hook
//func (this *PingRouter)PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PostHandle...")
//	_,err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
//	if err != nil {
//		fmt.Println("call back after ping error")
//	}
//}

/*
	HelloZinx  test 自定义路由
*/
type HelloZinxRouter struct {
	znet.BaseRouter
}

//Test Handle
//在处理conn业务的主方法hook
func (this *HelloZinxRouter)Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Router Handle...")
	//_,err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping...ping...\n"))
	//if err != nil {
	//	fmt.Println("call back ping...ping...ping error")
	//}
	//先读取客户端的数据，再回写ping..ping..ping
	fmt.Println("recv from client: msgID = ",request.GetMsgID(),
		", data = ",string(request.GetData()))

	err := request.GetConnection().SendMsg(201,[]byte("Hello Wlecome Zinx!!"))
	if err != nil {
		fmt.Println(err)
	}
}
//创建连接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("===> DoConnectionBegin is Called ... ")
	if err := conn.SendMsg(202,[]byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}
}

//连接断开之前的需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("===> DoConnectionLost is Called...")
	fmt.Println("conn ID = ",conn.GetConnID()," is Lost...")
}
func main() {
	//1 创建一个server句柄，使用 Zinx的api
	s := znet.NewServer("[zinx V0.5]")

	//2 注册连接的钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	//3 给当前 Zinx框架添加一个自定义的router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&PingRouter{})
	//4 启动server
	s.Server()

}
