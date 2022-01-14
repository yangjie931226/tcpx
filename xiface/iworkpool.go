package xiface

//所有实现该接口的均实现工作池
type ITask interface {
	DoTask() error
}


type IWorkPool interface {
	Run()
	Submit(ITask)
	
}