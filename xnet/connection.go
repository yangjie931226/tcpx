package xnet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"tcpx/xiface"
)

type Connection struct {
	//客户端唯一id
	ConnID uint32
	//客户端隶属服务
	TcpServer xiface.IServer
	//客户端当前连接
	Conn *net.TCPConn
	//消息管理
	MsgHanlder xiface.IMsgHandler
	//上下文管理
	ctx context.Context
	//通知连接关闭的方法
	ctxCancel context.CancelFunc
	//发送信息通道
	sendMessageChan chan []byte
	//连接属性
	property map[string]interface{}
	//管理连接的锁
	connMutex sync.RWMutex
	//是否关闭
	isClose    bool
	WorkerPool xiface.IWorkPool
}

//启动客户端
func (c *Connection) Start() {
	fmt.Println("Start ConnID", c.ConnID)
	//调用连接客户端钩子函数
	c.TcpServer.CallOnConnStartFunc(c)
	//启动读客户端goroutine
	go c.startReader()
	//启动写客户端goroutine
	go c.startWriter()
}

//关闭客户端
func (c *Connection) Stop() {
	fmt.Println("Stop ConnID", c.ConnID)

	//调用关闭客户端钩子函数
	c.TcpServer.CallOnConnCloseFunc(c)
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	if c.isClose {
		return
	}

	//从链接管理删除
	c.TcpServer.GetConnManager().RemoveConn(c.GetConnID())
	//上下文发送退出信号
	c.ctxCancel()
	//关闭客户端连接
	c.Conn.Close()
	//关闭发送信息管道
	close(c.sendMessageChan)
	//停止标志
	c.isClose = true
}

//获得当前客户端唯一id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获得当前客户端链接
func (c *Connection) GetConnection() *net.TCPConn {
	return c.Conn
}

//获得当前客户端ip地址
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送信息给客户端
func (c *Connection) SendMessage(routeId uint32, message []byte) error {
	c.connMutex.RLock()
	defer c.connMutex.RUnlock()
	if c.isClose {
		return errors.New("Connection is closed when Send Message")
	}
	dp := DataPacker{}
	data, err := dp.Pack(routeId, message)
	if err != nil {
		fmt.Println("Pack message error:", err)
		return errors.New("Pack message error")
	}
	//往sendMessageChan管道发送信息
	c.sendMessageChan <- data
	return nil
}

//设置属性
func (c *Connection) SetProperty(key string, val interface{}) {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}
	c.property[key] = val
}

//获取属性
func (c *Connection) GetProperty(key string) (interface{}, bool) {
	c.connMutex.RLock()
	defer c.connMutex.RUnlock()
	val, ok := c.property[key]
	return val, ok
}

//删除属性
func (c *Connection) RemoveProperty(key string) {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	delete(c.property, key)
}

//用户获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) startReader() {
	fmt.Println("startReader ConnID:", c.ConnID)
	defer c.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 读取头部信息
			dp := DataPacker{}
			buf := make([]byte, dp.GetHeaderLen())
			_, err := io.ReadFull(c.Conn, buf)
			if err != nil {
				fmt.Println("Conn Read error: ", err)
				return
			}
			fmt.Println(buf)
			message, err := dp.Unpack(buf)
			if err != nil {
				fmt.Println("Unpack error:", err)
				return
			}
			// 读取数据部分信息
			var databuf []byte
			if message.GetDataLen() > 0 {
				databuf = make([]byte, message.GetDataLen())
				_, err = io.ReadFull(c.Conn, databuf)
				if err != nil {
					fmt.Println("Conn Read error: ", err)
					return
				}
			}

			//设置message数据部分信息
			message.SetData(databuf)

			request := NewRequest(c, message)
			task := &RequestTask{
				Request:    request,
				ApiRouters: c.TcpServer.GetRouter(),
			}
			c.WorkerPool.Submit(task)
			//把信息发入协程池处理客户端发送过来的信息
			//c.MsgHanlder.SendRequestTask(request)
		}

	}
}

func (c *Connection) startWriter() {
	fmt.Println("startWriter ConnID:", c.ConnID)
	//启动循环监听发送信息管道
	for {
		select {
		case message := <-c.sendMessageChan:
			if _, err := c.Conn.Write(message); err != nil {
				fmt.Println("Conn Write error: ", err)
			}
		case <-c.ctx.Done(): //获取上下文管理管道关闭信号，退出循环
			return
		}
	}
}

func NewConnection(connID uint32, conn *net.TCPConn, tcpServer xiface.IServer, workPool xiface.IWorkPool) *Connection {
	ctx, ctxCancel := context.WithCancel(context.Background())
	c := &Connection{
		ConnID:          connID,
		Conn:            conn,
		TcpServer:       tcpServer,
		ctx:             ctx,
		ctxCancel:       ctxCancel,
		sendMessageChan: make(chan []byte, 1024),
		//MsgHanlder:      msgHandler,
		property:   nil,
		WorkerPool: workPool,
	}
	return c
}
