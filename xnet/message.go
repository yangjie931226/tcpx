package xnet

import "tcpx/xiface"

type Message struct {
	//处理数据的路由id
	RouteId uint32
	//数据
	Data []byte
	//数据长度
	DataLen uint32
}

//获取数据
func (m *Message) GetData() []byte {
	return m.Data
}

//设置数据
func (m *Message) SetData(data []byte) {
	m.Data = data
}

//获取数据长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

//设置数据长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

func (m *Message) SetRouteId(routeId uint32) {
	m.RouteId = routeId
}

func (m *Message) GetRouteId() uint32 {
	return m.RouteId
}

func NewMessage(data []byte, datalen, routeId uint32) xiface.IMessage {
	message := &Message{
		Data:    data,
		DataLen: datalen,
		RouteId: routeId,
	}
	return message
}
