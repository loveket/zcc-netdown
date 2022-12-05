package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/zcc_netdown/fileupdown"
	"github.com/zcc_netdown/model"
	"github.com/zcc_netdown/pool"
)

var ClientConn sync.Map

//初始化服务器配置
func init() {

}
func ClientConntest() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("cur active conn", ClientConn)
			fmt.Println("cur common source", fileupdown.DataList)
			//fmt.Println("****", fileupdown.DataList.LoadFile("aaa"))
		}
	}
}
func main() {
	go ClientConntest()
	log.Println("net-down server start...")
	service := ":8888" //设置服务器监听的端口
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println("net-down server listener err", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			ClientConn.Delete(conn.RemoteAddr().String())
			log.Println("net-down server from client conn err", err)
			continue
		}
		ClientConn.Store(conn.RemoteAddr().String(), conn)
		go HandleDownType(conn)
	}
}
func HandleDownType(conn net.Conn) {
	flag := pool.Pop()
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		pool.Push(flag)
	// 		log.Println("handle down err", err)
	// 		return
	// 	}
	// }()
	// timeout := time.NewTimer(30 * time.Second)
	// defer timeout.Stop()
	defer ConnClose(conn)
	msgBuf := make([]byte, 1024)
	r := bufio.NewReader(conn)
	var reqMsg *model.RequestMessage
	for {
		_, err := r.Read(msgBuf)
		if err != nil {
			pool.Push(flag)
			log.Println("from client conn Read err", err)
			return
		}
		if io.EOF != nil {
			return
		}
		err = json.Unmarshal(msgBuf, reqMsg)
		if err != nil {
			log.Println("conn json.Unmarshal err", err)
			continue
		}
		if _, ok := fileupdown.DataList.LoadFile(reqMsg.FileName); !ok {
			log.Println("request down file not exist")
			continue
		}
		if reqMsg.UpOrDown == 0 {
			go fileupdown.DownCommonFile(reqMsg, conn, flag)
		} else {

		}
	}

}
func ConnClose(conn net.Conn) {
	conn.Close()
}
