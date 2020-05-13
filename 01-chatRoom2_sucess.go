package main

import (
	"fmt"
	"strings"
	"time"

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

func MakeMsg(clnt Client,msg string)(buf string){
	buf = "[" + clnt.Addr + "]" + clnt.Name + ": " + msg

	return
}
func HandlerConnect(conn net.Conn) {
	defer conn.Close()
	//创建channel 判断， 用户是否活跃
	hasData := make(chan bool)

	//获取用户 网络地址 IP+port
	netAddr := conn.RemoteAddr().String()
	// 创建新连接用户的结构体信息 默认用户是IP + port
	clnt := Client{make(chan string),netAddr,netAddr}

	//将新连接用户，添加到在线用户map中   key:IP + port  value: client
	onlineMap[netAddr] = clnt

	//创建专门用来给 当前 用户发送消息  的 go程
	go WriteMsgToClient(clnt,conn)
	//发送 用户上线消息到 全局channel中
	//message <- "[" + netAddr + "]" + clnt.Name + "login"
	message <- MakeMsg(clnt,"login")
	//创建一个channel,用来判断用户退出状态
	isQuit := make(chan bool)

	//创建一个匿名go程，专门处理用户发送的消息
	go func () {
		buf := make([]byte,4096)
		for {
			n,err := conn.Read(buf)
			if n == 0 {
				isQuit <- true
				fmt.Printf("检测到客户端: %s 推出\n",clnt.Name)
				return
			}
			if err != nil {
				fmt.Println("conn.Read err: ",err)
				return
			}
			//将读到的用户消息，保存到msg中，string类型
			msg := string(buf[:n-1])
			//提取在线用户列表
			if msg == "who" && len(msg) == 3 {
				conn.Write([]byte("online user list:\n"))
				//遍历当前 map ，获取在线用户
				for _,user := range onlineMap {
					userInfo := user.Addr + ":" + user.Name + "\n"
					conn.Write([]byte(userInfo))
				}
			} else if len(msg) >= 8 && msg[:6] == "rename" {		//rename | a
				newName := strings.Split(msg,"|")[1]			//msg[8:]
				clnt.Name = newName									//修改结构体 name
				onlineMap[netAddr] = clnt							//更新 onlineMap
				conn.Write([]byte("rename sucessful\n"))

			} else {
				//将读到的用户消息，写入到message中
				message <- MakeMsg(clnt,msg)
			}
			//important flush
			hasData <- true
		}
	}()
	//保证不退出
	for {
		//监听 channel 上的数据流动
		select {
		case <- isQuit:
			close(clnt.C) 		//关闭write 协程  十分重要
			delete(onlineMap, clnt.Addr)		//将用户从online移除
			message <- MakeMsg(clnt,"logout")	//写入用户推出消息到全局 channel
			return
		case <-hasData:
			//什么都不做。目的是重置 下面 case 的计时器
		case <- time.After(time.Second * 10 ):
			delete(onlineMap, clnt.Addr)		//将用户从online移除
			message <- MakeMsg(clnt,"logout")	//写入用户推出消息到全局 channel
			return
		}
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
