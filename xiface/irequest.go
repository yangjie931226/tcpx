package xiface

type IRequest interface {
	//获得客户端链接
	GetConnection() IConnection
	//获取数据
	GetMessage() IMessage
}
