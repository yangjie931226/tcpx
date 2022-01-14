package xnet

import (
	"context"
	"fmt"
	"net"
	"sync"
	"tcpx/config"
	"tcpx/xiface"
)

type Server struct {
	//服务器名称
	Name string
	//服务器ip
	IP string
	//服务器端口
	Port int
	//ip类型
	IPVersion string
	//消息管理
	MsgHanlder xiface.IMsgHandler
	//客户端连接管理
	ConnManager xiface.IConnManage
	//客户端连接钩子函数
	onConnStartHookFunc func(c xiface.IConnection)
	//客户端断开连接钩子函数
	onConnCloseHookFunc func(c xiface.IConnection)
	//上下文管理
	ctx context.Context
	//通知连接关闭的方法
	ctxCancel context.CancelFunc
	//路由处理方法
	ApiRouters map[uint32]xiface.IRouter
	//工作池
	WorkerPool xiface.IWorkPool
	mutex sync.RWMutex
}

// 启动服务tcp服务
func (s *Server) Start() {
	fmt.Println("[Server Start]", s.IPVersion, s.Name)
	//获取监听的地址
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("Resolve TCP Addr error: ", err)
		return
	}
	//启动监听
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("Listen TCP error: ", err)
		return
	}
	var connID uint32
	connID = 1
	//启动循环接收连接进来的客户端
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept TCP error: ", err)
			continue
		}
		//超过最大长度限制断开
		if s.ConnManager.GetConnSize() > config.GobalConfig.MaxConnLen {
			conn.Close()
			fmt.Println("客户端超过服务器最大连接数")
			continue
		}
		//启动一个goroutine处理每个连接进来的客户端
		dealConn := NewConnection(connID, conn, s, s.WorkerPool)
		//添加到连接管理
		s.ConnManager.AddConn(connID, dealConn)
		go dealConn.Start()
		connID++
	}

}

// 关闭服务
func (s *Server) Stop() {
	fmt.Println("[Server Stop]", s.IPVersion, s.Name)
	//关闭服务器 回收资源
	s.ConnManager.Clear()
	//关闭协程池
	s.MsgHanlder.Stop()
	//上下文通知服务器关闭
	s.ctxCancel()
}

//启动其他相关服务
func (s *Server) Serve() {
	go s.Start()

	//启动协程池处理业务
	go s.WorkerPool.Run()
	for {
		select {
		case <-s.ctx.Done(): //获取上下文管理管道关闭信号，退出循环
			return
		}
	}
}

//设置服务器路由
func (s *Server) SetRouter(routeId uint32, router xiface.IRouter) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ApiRouters[routeId] = router
	//s.MsgHanlder.SetRouter(routeId, router)
}

//设置服务器路由
func (s *Server) GetRouter() map[uint32]xiface.IRouter {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ApiRouters
}

//设置客户端启动钩子函数
func (s *Server) SetOnConnStartFunc(onConnStartHookFunc func(connection xiface.IConnection)) {
	s.onConnStartHookFunc = onConnStartHookFunc
}

//设置客户端断开钩子函数
func (s *Server) SetOnConnCloseFunc(onConnCloseHookFunc func(connection xiface.IConnection)) {
	s.onConnCloseHookFunc = onConnCloseHookFunc
}

//调用客户端启动钩子函数
func (s *Server) CallOnConnStartFunc(connection xiface.IConnection) {
	if s.onConnStartHookFunc != nil {
		s.onConnStartHookFunc(connection)
	}
}

//调用客户端断开钩子函数
func (s *Server) CallOnConnCloseFunc(connection xiface.IConnection) {
	if s.onConnCloseHookFunc != nil {
		s.onConnCloseHookFunc(connection)
	}
}

//获取连接管理
func (s *Server) GetConnManager() xiface.IConnManage {
	return s.ConnManager
}



func NewServer() xiface.IServer {
	ctx, ctxCancel := context.WithCancel(context.Background())
	s := &Server{
		Name:        config.GobalConfig.Name,
		IPVersion:   "tcp",
		IP:          config.GobalConfig.IP,
		Port:        config.GobalConfig.Port,
		ctx:         ctx,
		ctxCancel:   ctxCancel,
		MsgHanlder:  NewMsgHandler(),
		ConnManager: NewConnManager(),
		ApiRouters: map[uint32]xiface.IRouter{},
		WorkerPool:NewWorkPool(config.GobalConfig.WorkPoolSize),
	}
	return s
}
