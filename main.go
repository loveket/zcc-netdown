package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/zcc_netdown/httpserver"
	"github.com/zcc_netdown/tcpserver"
)

//初始化服务器配置
func init() {

}
func ClientConntest() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			//fmt.Println("cur active conn", ClientConn)
			//fmt.Println("cur common source", fileupdown.DataList)
			//fmt.Println("****", fileupdown.DataList.LoadFile("aaa"))
		}
	}
}
func MyselfCheck() {
	http.ListenAndServe(":8889", nil)
}
func main() {
	go ClientConntest()
	go MyselfCheck()
	go httpserver.StartHttpServer()
	tcpserver.StartServer()
}
