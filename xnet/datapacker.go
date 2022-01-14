package xnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"tcpx/config"
	"tcpx/xiface"
)

// 默认TLV协议
type DataPacker struct{}

//解包
func (dp *DataPacker) Unpack(header []byte) (xiface.IMessage, error) {
	//头部包含 （路由id 数据长度）
	buf := bytes.NewReader(header)
	message := &Message{}
	//路由id  第一部分头部4字节对应路由id字节
	if err := binary.Read(buf, binary.LittleEndian, &message.RouteId); err != nil {
		return nil, err
	}

	//数据长度  第二部分头部4字节对应数据长度
	if err := binary.Read(buf, binary.LittleEndian, &message.DataLen); err != nil {
		return nil, err
	}

	if config.GobalConfig.MaxPackerSize > 0 && message.GetDataLen() > config.GobalConfig.MaxPackerSize {
		return nil, errors.New("too large msg data received")
	}

	return message, nil
}

//封包
func (dp *DataPacker) Pack(routeId uint32, data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	//封装路由id
	if err := binary.Write(&buffer, binary.LittleEndian, routeId); err != nil {
		return nil, err
	}
	//封装数据长度
	if err := binary.Write(&buffer, binary.LittleEndian, uint32(len(data))); err != nil {
		return nil, err
	}

	//封装数据
	if err := binary.Write(&buffer, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (dp *DataPacker) GetHeaderLen() uint32 {
	return 8
}

func NewPataPacker() xiface.IDataPacker {
	dp := &DataPacker{}
	return dp
}
