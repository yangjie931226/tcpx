package xiface

// 服务器接口
type IServer interface {
	// 启动服务tcp服务
	Start()
	// 关闭服务
	Stop()
	//启动其他相关服务
	Serve()
	//设置服务器路由
	SetRouter(uint32, IRouter)
	//设置服务器路由
	GetRouter() map[uint32]IRouter
	//获取连接管理
	GetConnManager() IConnManage
	//设置客户端启动钩子函数
	SetOnConnStartFunc(func(connection IConnection))
	//设置客户端断开钩子函数
	SetOnConnCloseFunc(func(connection IConnection))
	//调用客户端启动钩子函数
	CallOnConnStartFunc(IConnection)
	//调用客户端断开钩子函数
	CallOnConnCloseFunc(IConnection)
}
