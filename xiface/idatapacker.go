package xiface

type IDataPacker interface {
	//解包
	Unpack(header []byte) (IMessage, error)
	//封包
	Pack(routeId uint32, data []byte) ([]byte, error)
}
