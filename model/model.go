package model

//定义传输数据格式
type RequestMessage struct {
	FileName string `json:"filename"`
	BackWord string `json:"backword"`
	DataLen  int    `json:"datalen"`
	Data     []byte `json:"data"`
	UpOrDown int    `json:"upordown"`
}
