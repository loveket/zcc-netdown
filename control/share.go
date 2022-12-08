package control

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zcc_netdown/fileupdown"
	"github.com/zcc_netdown/model"
)

func Getsorcelist(c *gin.Context) {
	resList := make([]string, 0)
	for k, _ := range fileupdown.DataList.FileMap {
		resList = append(resList, k)
	}
	//c.JSON(200, gin.H{"sourcelist": resList})
	c.HTML(200, "main.html", gin.H{"sourcelist": resList})
}
func DownCommonSource(c *gin.Context) {
	filename := c.Query("filename")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("conn server err", err)
		return
	}
	rm := &model.RequestMessage{
		FileName: filename,
		UpOrDown: 0,
	}
	result, err := json.Marshal(rm)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	_, err = conn.Write(result)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	resBuf := make([]byte, 2048)
	for {
		_, err := conn.Read(resBuf)
		if err != nil {
			if err == io.EOF {
				c.JSON(200, "down success")
			}
			return
		}
		fileContentDisposition := "attachment;filename=\"" + filename + "\""
		//c.Header("Content-Type", "video/x-mpg") // 这里是压缩文件类型 .mp4
		c.Header("Content-Disposition", fileContentDisposition)
		c.Data(http.StatusOK, "", resBuf)
	}
	// var file *os.File
	// file, err = os.Create(rm.FileName)
	// defer file.Close()
	// if err != nil {
	// 	fmt.Println("create file err", err)
	// 	return
	// }
	// w, err := io.Copy(file, conn)
	// if err != nil {
	// 	fmt.Println("copy file err", err)
	// }
	//fmt.Println("下载完成", w)
	// resBuf := make([]byte, 2048)
	// for {
	// 	m, err := conn.Read(resBuf)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			fmt.Println("--------", resBuf[:m])
	// 			contentType := c.Request.Header.Get("Content-Type")
	// 			fileContentDisposition := "attachment;filename=\"" + filename + "\""
	// 			c.Header("Content-Type", "video/x-mpg") // 这里是压缩文件类型 .mp4
	// 			c.Header("Content-Disposition", fileContentDisposition)
	// 			c.Data(http.StatusOK, contentType, resBuf[:m])
	// 		}
	// 		return
	// 	}
	// }
	// path := fileupdown.DataList.FileMap[filename]
	// file, err := os.Open(path + "/" + filename)
	// defer file.Close()
	// if err != nil {
	// 	c.JSON(201, "file open err")
	// 	return
	// }
	// resBuf := make([]byte, 2048)
	// for {
	// 	_, err := file.Read(resBuf)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			c.JSON(200, "down success")
	// 		}
	// 		return
	// 	}
	// 	fileContentDisposition := "attachment;filename=\"" + filename + "\""
	// 	//c.Header("Content-Type", "video/x-mpg") // 这里是压缩文件类型 .mp4
	// 	c.Header("Content-Disposition", fileContentDisposition)
	// 	c.Data(http.StatusOK, "", resBuf)
	// }
}
func UpCommonSource(c *gin.Context) {
	uploadFile, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "获取文件信息失败!" + err.Error(),
		})
	}
	if uploadFile != nil {
		defer uploadFile.Close()
	}
	r := bufio.NewReader(uploadFile)
	// buf := make([]byte, 2048)
	// for {
	// 	_, err := r.Read(buf)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			//fmt.Println(buf[:n])
	// 			break
	// 		}
	// 		return
	// 	}
	// }
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("conn server err", err)
		return
	}
	rm := &model.RequestMessage{
		FileName: fileHeader.Filename,
		UpOrDown: 1,
	}
	result, err := json.Marshal(rm)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	_, err = conn.Write(result)
	if err != nil {
		fmt.Println("conn write err", err)
		return
	}
	a := make([]byte, 1024)
	n, err := conn.Read(a)
	if err != nil {
		fmt.Println("conn write err", err)
		return
	}
	if "ok" == string(a[:n]) {
		sendFile(conn, fileHeader.Filename, r)
	}
	c.Redirect(http.StatusFound, "/")
	return
}
func sendFile(conn net.Conn, filename string, r *bufio.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("发送文件完成。")
			} else {
				fmt.Println("os.Open err:", err)
			}
			return
		}
		// 写到网络socket中
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return
		}
	}
}
