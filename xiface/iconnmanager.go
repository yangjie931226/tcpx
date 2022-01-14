package xiface

type IConnManage interface {
	//获取链接总数
	GetConnSize() int
	//添加链接
	AddConn(uint32,IConnection)
	//删除链接
	RemoveConn(uint32)
	//根据ConnID获取连接
	GetConnById(connId uint32) (IConnection, bool)
	//清除所有连接
	Clear()

}
