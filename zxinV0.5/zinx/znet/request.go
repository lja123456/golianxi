package znet

import "zinx/ziface"

type Request struct {
	//已经和客户端建立好的链接
	conn ziface.IConnection
	//客户端请求的数据
	//data []byte
	msg ziface.IMessage
}

//得到当前链接
func (r *Request)GetConnection()  ziface.IConnection {
	return r.conn
}

//得到请求的消息数据
func (r *Request)GetData() []byte {
	//return  r.data
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}