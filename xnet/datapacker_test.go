package xnet

import (
	"fmt"
	"testing"
)

func TestDataPacker(t *testing.T) {
	dp := DataPacker{}
	data, err := dp.Pack(2,[]byte("hello test2"))
	if err != nil {
		fmt.Println(err)
		return
	}
	dp2 := DataPacker{}
	mes,err :=dp2.Unpack(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mes.GetRouteId())
	fmt.Println(mes.GetDataLen())
	fmt.Println(string(mes.GetData()))

}

