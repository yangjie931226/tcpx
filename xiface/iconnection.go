package xiface

import (
	"context"
	"net"
)

type IConnection interface {
	//启动客户端
	Start()
	//关闭客户端
	Stop()
	//获得当前客户端唯一id
	GetConnID() uint32
	//获得当前客户端链接
	GetConnection() *net.TCPConn
	//获得当前客户端ip地址
	GetRemoteAddr() net.Addr
	//发送信息给客户端
	SendMessage(uint32, []byte) error
	//设置属性
	SetProperty(string, interface{})
	//获取属性
	GetProperty(string) (interface{}, bool)
	//删除属性
	RemoveProperty(string)
	//用户获取连接退出状态
	Context() context.Context
}
