package tcpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/zcc_netdown/common"
	"github.com/zcc_netdown/fileupdown"
	"github.com/zcc_netdown/model"
	"github.com/zcc_netdown/pool"
	"github.com/zcc_netdown/utils"
)

var ClientConn sync.Map

func StartServer() {
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
	defer pool.Push(flag)
	defer fmt.Println("hahaha")
	msgBuf := make([]byte, 2048)
	var reqMsg model.RequestMessage
	for {
		n, err := conn.Read(msgBuf)
		if err != nil {
			log.Println("from client conn Read err", err)
			return
		}
		err = json.Unmarshal(msgBuf[:n], &reqMsg)
		if err != nil {
			log.Println("conn json.Unmarshal err", err)
			continue
		}

		fmt.Println("----", string(msgBuf))
		fmt.Println("*****", reqMsg)
		if reqMsg.UpOrDown == 0 {
			path, ok := fileupdown.DataList.LoadFile(reqMsg.FileName)
			if !ok {
				log.Println("request down file not exist")
				continue
			}
			go fileupdown.DownCommonFile(&reqMsg, conn, path)
		} else {
			lastname := utils.GetFileLastName(reqMsg.FileName)
			upAllPath := common.CommonSourceAddr + "/" + lastname + "-file"
			isExist, err := utils.IsExistPath(upAllPath)
			if err != nil {
				log.Println("check path err", err)
				return
			}
			if !isExist {
				err := os.Mkdir(upAllPath, os.ModePerm)
				if err != nil {
					log.Println("create dir err", err)
					return
				}
			}
			fmt.Println("---------", reqMsg)
			fileupdown.UpCommonFile(&reqMsg, conn, upAllPath)
		}
	}

}
func ConnClose(conn net.Conn) {
	ClientConn.Delete(conn.RemoteAddr().String())
	conn.Close()
}
