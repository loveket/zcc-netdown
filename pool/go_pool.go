package pool

import "github.com/zcc_netdown/common"

var goPool = make(chan int, common.MaxDownActiveNum)

func init() {
	for i := 1; i <= common.MaxDownActiveNum; i++ {
		goPool <- i
	}
}
func Push(i int) {
	goPool <- i
}

func Pop() int {
	return <-goPool
}
