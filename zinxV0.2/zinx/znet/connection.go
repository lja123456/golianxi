package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

/*
	链接模块
 */
type Connection struct {
	//当前连接的 socket  TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的链接状态
	isClosed bool

	//当前连接所绑定的处理业务方法API
	handleAPI  ziface.HandleFunc

	//告知当前连接已经推出的/停止的 channel
	ExitChan chan bool
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32,callback_api ziface.HandleFunc)*Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		handleAPI: callback_api,
		isClosed: false,
		ExitChan: make(chan bool,1),
	}
	return c
}

//连接的读业务方法
func (c *Connection)StartReader() {
	fmt.Println("Reader Groutine is running..")
	defer fmt.Println("connID = ",c.ConnID, " Reader is exit, remote addr is ",c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte,512)
		cnt,err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err",err)
			continue
		}

		//调用当前连接所绑定的HandleAPI
		if c.handleAPI(c.Conn,buf,cnt); err != nil {
			fmt.Println("ConnID ",c.ConnID, " handle is error",err)
			break
		}
	}
}
//启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ",c.ConnID)
	//启动从当前连接的读数据业务

	//TODO 启动从当前连接写数据的业务
}
//停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ",c.ConnID)

	//如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭 socket连接
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)
}
//获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
//获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
//获取远程客户端的TCP 状态  IP port
func (c *Connection) RemoteAddr()  net.Addr {
	return c.Conn.RemoteAddr()
}
//发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {

}