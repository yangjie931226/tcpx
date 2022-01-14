package xiface

type IRouter interface {
	//前置处理信息
	PreHandle(request IRequest)
	//处理信息
	Handle(request IRequest)
	//后置处理信息
	PostHandle(request IRequest)
}
