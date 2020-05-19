package znet

import (
	"errors"
	"fmt"
	"io"
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
//	handleAPI  ziface.HandleFunc

	//告知当前连接已经推出的/停止的 channel
	ExitChan chan bool

	//该连接处理的方法Router
	//Router ziface.IRouter
	//消息的管理MsgID 和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32,msgHandler ziface.IMsgHandle)*Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		//zinx v0.2  handleAPI: callback_api,
		MsgHandle:msgHandler,
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
		//buf := make([]byte,512)
		//buf := make([]byte,utils.GlobalObject.MaxPackageSize)
		//_,err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err",err)
		//	continue
		//}
		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 8个字节
		headData := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(c.GetTCPConnection(),headData); err != nil {
			fmt.Println("read msg head error",err)
			break
		}

		//拆包，得到msgID 和  msgDatalen 放在 msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error",err)
			break
		}

		//根据dataLen 再次读取Data, 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte,msg.GetMsgLen())
			if _,err := io.ReadFull(c.GetTCPConnection(),data); err != nil {
				fmt.Println("read msg data error ",err)
				break
			}
		}
		msg.SetData(data)
		//得到当前conn 数据的Request请求数据
		req := Request{
			c,
			msg,
		}
		//go func(request ziface.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)
		//从路由中，找到注册绑定的Conn对应的router调用
		//c.Router.PreHandle(&req)
		//根据绑定号的MsgId 找到对应处理api业务
		go c.MsgHandle.DoMsgHandler(&req)
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
//func (c *Connection) Send(data []byte) error {
//	return nil
//}
//提供一个SendMsg 方法， 将我们要发送给客户端的数据，先进性封包，在发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包，  MsgDataLen | MsgID | Data
	dp := NewDataPack()
	binaryMsg,err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ",msgId)
		return errors.New("Pack error msg")
	}
	//将数据发送给客户端
	if _,err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id: ",msgId," error: ",err)
		return errors.New("conn Write error")
	}
	return nil
}
