package ziface

/*
	将请求的消息封装到一个Message中，定义抽象接口
 */

type IMessage interface {
	GetMsgID() uint32	//獲取消息的ID
	//獲取消息的長度
	GetMsgLen() uint32
	//获取消息的内容
	GetData() []byte
	//设置消息的ID
	SetMsgID(uint32)
	//设置消息的内容
	SetData([]byte)
	//设置消息的长度
	SetDataLen(uint32)
}