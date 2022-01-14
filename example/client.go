package main

import (
	"fmt"
	"net"
	"tcpx/xnet"
	"time"
)


func main()  {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	for  {
		dp := xnet.DataPacker{}
		data,err := dp.Pack(2,[]byte("test tcpx"))
		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err := conn.Write(data); err != nil {
			return
		}

		buf := make([]byte, 512)
		cnt,err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
