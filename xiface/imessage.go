package xiface

type IMessage interface {
	//获取数据
	GetData() []byte
	//设置数据
	SetData([]byte)
	//获取数据长度
	GetDataLen() uint32
	//设置数据长度
	SetDataLen(uint32)
	//设置路由id
	SetRouteId(uint32)
	//获取路由idet
	GetRouteId() uint32
}
