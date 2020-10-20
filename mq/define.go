package mq

import "cloudstorage/common"

// 写入 rabbitmq 的消息结构体
type TransferData struct {
	FileHash      string
	CurLocation   string
	DestLocation  string
	DestStoreType common.StoreType
}
