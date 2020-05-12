package main

import (
	"fmt"
	"net"
)
//创建用户结构体类型
type Client struct {
	C chan string
	Name string
	Addr string
}

//创建全局map，存储在线用户
var onlineMap map [string]Client

//创建全局channel传递用户消息
var message = make(chan string)

func WriteMsgToClient(clnt Client,conn net.Conn) {
	//监听 用户自带Channel上是否有消息
	for msg := range clnt.C {
		conn.Write([]byte(msg + "\n"))
	}
}
func HandlerConnect(conn net.Conn) {
	defer conn.Close()

	//获取用户 网络地址 IP+port
	netAddr := conn.RemoteAddr().String()
	// 创建新连接用户的结构体信息 默认用户是IP + port
	clnt := Client{make(chan string),netAddr,netAddr}

	//将新连接用户，添加到在线用户map中   key:IP + port  value: client
	onlineMap[netAddr] = clnt

	//创建专门用来给 当前 用户发送消息  的 go程
	go WriteMsgToClient(clnt,conn)
	//发送 用户上线消息到 全局channel中
	message <- "[" + netAddr + "]" + clnt.Name + "login"

	//保证不退出
	for {
		;
	}

}

func Manager() {
	//初始化 onlineMap
	onlineMap = make(map[string]Client)

	//监听全局channel中是否有数据,有数据存储至msg，无数据阻塞
	for {

		msg := <-message

		//循环发送消息给所有在线用户。
		for _,clnt :=range onlineMap {
			clnt.C <- msg
		}
	}

}
func main() {
	//创建监听套接字
	listener,err := net.Listen("tcp","127.0.0.1:8000")
	if err != nil {
		fmt.Println("Listen err",err)
		return
	}
	defer listener.Close()

	//创建管理者go程，管理map和全局channel
	go Manager()
	//循环监听客户端连接请求
	for {
		conn,err := listener.Accept()
		if err != nil {
			fmt.Println("accept err",err)
			return
		}
		//启动go程处理客户端数据请求
		go HandlerConnect(conn)
	}
}
