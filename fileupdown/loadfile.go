package fileupdown

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zcc_netdown/common"
)

var DataList *AllDataMessage

type AllDataMessage struct {
	FileMap map[string]string
	lock    sync.RWMutex
}

//存储
func (adm *AllDataMessage) SaveFile(filename, path string) {
	adm.lock.Lock()
	defer adm.lock.Unlock()
	adm.FileMap[filename] = path
}

//读取
func (adm *AllDataMessage) LoadFile(filename string) (path string, isexist bool) {
	adm.lock.RLock()
	defer adm.lock.RUnlock()
	path = adm.FileMap[filename]
	if path == " " || len(path) == 0 {
		return path, false
	}
	return path, true
}

//读取公共资源所有数据
func init() {
	DataList = &AllDataMessage{
		FileMap: make(map[string]string),
	}
	var files []string
	err := filepath.Walk(common.CommonSourceAddr, func(path string, info os.FileInfo, err error) error {
		//files = append(files, path)
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Println("load filepath err", err)
		return
	}
	for _, file := range files {
		//fmt.Println("*******", file)
		// strList:=strings.Split(file,"\\")
		// strListLen:=len(strList)
		// strList[strListLen-1]=
		lastPos := strings.LastIndex(file, "\\")
		len := len(file)
		filename := file[lastPos+1 : len]
		path := file[0:lastPos]
		DataList.FileMap[filename] = path
	}
}
