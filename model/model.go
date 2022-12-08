package model

//定义传输数据格式
type RequestMessage struct {
	FileName     string `json:"filename"`
	UpOrDown     int    `json:"upordown"`
	FileCapacity int64  `json:"filecapacity"`
	//Agent    string `json:"agent"` //预留代理存储
}
