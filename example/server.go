package main

import (
	"fmt"
	"tcpx/xiface"
	"tcpx/xnet"
)

type Router1 struct {
	xnet.BaseRoute
}
func (r *Router1) PreHandle(request xiface.IRequest) {
	request.GetConnection().SetProperty("aaaa","ghghvv")
}
func (r *Router1) Handle(request xiface.IRequest) {
	a,_ := request.GetConnection().GetProperty("aaaa")
	request.GetConnection().SendMessage(0,[]byte(a.(string)))
	request.GetConnection().RemoveProperty("aaaa")
}
func (r *Router1) PostHandle(request xiface.IRequest) {
	_,ok := request.GetConnection().GetProperty("aaaa")
	if !ok {
		fmt.Println("aaaa已删除")
	}
}
type Router2 struct {
	xnet.BaseRoute
}

func (r *Router2) Handle(request xiface.IRequest) {
	fmt.Println("Handle")

	if err := request.GetConnection().SendMessage(1,request.GetMessage().GetData()) ;err != nil {
		fmt.Println(err)
	}

}

func onStartHookFunc (connection xiface.IConnection) {
	fmt.Println("启动钩子")
}
func onCloseHookFunc (connection xiface.IConnection) {
	fmt.Println("关闭钩子")
}
func main()  {
	s := xnet.NewServer()
	s.SetRouter(1,&Router1{})
	s.SetRouter(2,&Router2{})
	s.SetOnConnStartFunc(onStartHookFunc)
	s.SetOnConnCloseFunc(onCloseHookFunc)

	s.Serve()
}
