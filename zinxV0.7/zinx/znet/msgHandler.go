package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

/*
	消息处理模块的实现
 */

type MsgHandle struct {
	//存放每个MsgHandle 所对应的处理方法
	Apis map[uint32] ziface.IRouter
}

//初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		make(map[uint32] ziface.IRouter),
	}
}
//调度/执行对应的router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 从REquest中找到msgID
	handler,ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ",request.GetMsgID(), " is NOT FOUNT! Need register!")
	}
	//2 根据MsgID 调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32,router ziface.IRouter) {
	//1 判断 当前msg绑定的API处理方法是否已经存在
	if _,ok := mh.Apis[msgID]; ok {
		//id已经注册了
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与API绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ",msgID, " succ!")

}