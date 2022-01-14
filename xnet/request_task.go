package xnet

import (
	"fmt"
	"tcpx/xiface"
)

type RequestTask struct {
	Request xiface.IRequest
	ApiRouters map[uint32]xiface.IRouter
}



func (rt *RequestTask)DoTask() error {
	router,ok := rt.ApiRouters[rt.Request.GetMessage().GetRouteId()]
	if !ok {
		return fmt.Errorf("无此路由id %v",rt.Request.GetMessage().GetRouteId())
	}
	router.PreHandle(rt.Request)
	router.Handle(rt.Request)
	router.PostHandle(rt.Request)
	return nil
}
