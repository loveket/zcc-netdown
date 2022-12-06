package fileupdown

import (
	"log"
	"net"
	"os"

	"github.com/zcc_netdown/model"
)

func DownCommonFile(reqmsg *model.RequestMessage, conn net.Conn, path string) {
	downFile, err := os.ReadFile(path + "\\" + reqmsg.FileName)
	if err != nil {
		log.Println("load local source err", err)
		return
	}
	_, err = conn.Write(downFile)
	if err != nil {
		log.Println("down err", err)
		return
	}
	return
}

func UpCommonFile(reqmsg *model.RequestMessage, conn net.Conn, path string) {
	conn.Write([]byte("ok"))
	allpath := path + "/" + reqmsg.FileName
	file, err := os.Create(allpath)
	defer file.Close()
	if err != nil {
		log.Println("create file err", err)
		return
	}
	DataList.lock.Lock()
	DataList.FileMap[reqmsg.FileName] = path
	DataList.lock.Unlock()
	msgBuf := make([]byte, 2048)
	for {
		n, err := conn.Read(msgBuf)
		if err != nil {
			return
		}
		if n == 0 {
			log.Println("recv file success")
			return
		}
		_, err = file.Write(msgBuf[:n])
		if err != nil {
			log.Println("write file", err)
			return
		}
	}

}
