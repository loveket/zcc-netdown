package fileupdown

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/zcc_netdown/model"
)

func DownCommonFile(reqmsg *model.RequestMessage, conn net.Conn, path string) {
	file, err := os.Open(path + "\\" + reqmsg.FileName)
	if err != nil {
		log.Println("load local source err", err)
		return
	}
	buf := make([]byte, 2048)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("load file data success")
				conn.Close()
			}
			return
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Println("down err", err)
			return
		}
	}
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
	msgBuf := make([]byte, 2048)
	for {
		n, _ := conn.Read(msgBuf)
		if n == 0 {
			log.Println("recv file success")
			DataList.lock.Lock()
			DataList.FileMap[reqmsg.FileName] = path
			DataList.lock.Unlock()
			return
		}
		_, err = file.Write(msgBuf[:n])
		if err != nil {
			log.Println("write file", err)
			return
		}

	}

}
