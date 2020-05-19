package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//iSerer 的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前的Server 添加一个router, server注册的连接对应的处理业务
	Router ziface.IRouter
}


//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s,listenner at IP: %s,Port: %d is starting...\n",utils.GlobalObject.Name,
		utils.GlobalObject.Host,utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPackeetSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)
	fmt.Printf("[Start] Server Listennre at IP : %s,Port %d, is starting\n",s.IP,s.Port)

	go func() {
		//1 获取一个TCP的Addr
		addr,err := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ",err)
			return
		}

		//2 监听服务器的地址
		listenner,err := net.ListenTCP(s.IPVersion,addr)
		if err != nil {
			fmt.Println("listen ",s.IPVersion, "err ",err)
			return
		}

		fmt.Println("start Zinx server succ, ",s.Name, "succ,Listenning...")

		var cid uint32
		cid = 0
		//3 阻塞的等待客户端连接，处理客户端连接业务
		for {
			//如果有客户端连接过来，阻塞会返回
			conn,err :=listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err",err)
				continue
			}
			//将处理新连接的业务方法  和 conn 进行绑定 得到我们的连接模块
			dealConn := NewConnection(conn,cid, s.Router)
			cid++

			//启动 当前的连接业务处理
			go dealConn.Start()

		}
	}()

}
//停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源，状态或者一些已经开辟的链接信息 进行停止或回收
}
//运行服务器
func (s *Server) Server() {
	//启动server的功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {

	}
}

//路由功能： 给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Succ!!")
}
/*
	初始化Server 模块的方法
 */
func NewServer(name string) ziface.IServer{
	s := &Server {
		//Name: name,
		Name: utils.GlobalObject.Name,
		IPVersion: "tcp4",
		//IP: "0.0.0.0",
		IP: utils.GlobalObject.Host,
		//Port: 8999,
		Port: utils.GlobalObject.TcpPort,
		Router: nil,
	}
	return s
}