package xnet

import (
	"tcpx/xiface"
)

type Request struct {
	//客户端连接
	Conn xiface.IConnection
	//数据
	Data xiface.IMessage
}

//获得客户端链接
func (r *Request) GetConnection() xiface.IConnection {
	return r.Conn
}

//获取数据
func (r *Request) GetMessage() xiface.IMessage {
	return r.Data
}


func NewRequest(conn xiface.IConnection, data xiface.IMessage) xiface.IRequest {
	request := &Request{
		Conn: conn,
		Data: data,
	}
	return request
}
