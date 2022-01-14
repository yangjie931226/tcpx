package xnet

import (
	"fmt"
	"sync"
	"tcpx/xiface"
)

type ConnManager struct {
	ConnMap map[uint32]xiface.IConnection
	cmgrMutex sync.RWMutex
}

//获取链接总数
func (cmgr *ConnManager) GetConnSize() int {
	cmgr.cmgrMutex.RLock()
	defer cmgr.cmgrMutex.RUnlock()
	return len(cmgr.ConnMap)
}

//添加链接
func (cmgr *ConnManager) AddConn(connId uint32, conn xiface.IConnection) {
	cmgr.cmgrMutex.Lock()
	defer cmgr.cmgrMutex.Unlock()
	cmgr.ConnMap[connId] = conn
	fmt.Println("添加链接", cmgr.ConnMap)
}

//删除链接
func (cmgr *ConnManager) RemoveConn(connId uint32) {
	cmgr.cmgrMutex.Lock()
	defer cmgr.cmgrMutex.Unlock()
	delete(cmgr.ConnMap, connId)
	fmt.Println("删除链接", cmgr.ConnMap)
}

//根据ConnID获取连接
func (cmgr *ConnManager) GetConnById(connId uint32) (xiface.IConnection, bool) {
	cmgr.cmgrMutex.RLock()
	defer cmgr.cmgrMutex.RLock()
	conn, ok := cmgr.ConnMap[connId]
	return conn, ok
}


//清除所有连接
func (cmgr *ConnManager) Clear() {
	cmgr.cmgrMutex.Lock()
	defer cmgr.cmgrMutex.Unlock()
	for k, v := range cmgr.ConnMap {
		v.Stop()
		delete(cmgr.ConnMap, k)
	}
}

func NewConnManager() xiface.IConnManage {
	cmgr := &ConnManager{
		ConnMap: make(map[uint32]xiface.IConnection),
	}
	return cmgr
}
