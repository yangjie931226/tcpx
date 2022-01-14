package xnet

import "tcpx/xiface"

type BaseRoute struct {
}

//前置处理信息
func (r *BaseRoute) PreHandle(request xiface.IRequest) {

}

//处理信息
func (r *BaseRoute) Handle(request xiface.IRequest) {

}

//后置处理信息
func (r *BaseRoute) PostHandle(request xiface.IRequest) {

}
