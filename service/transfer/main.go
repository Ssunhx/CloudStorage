package transfer

import (
	"bufio"
	"cloudstorage/config"
	"cloudstorage/db"
	"cloudstorage/mq"
	"cloudstorage/store/oss"
	"encoding/json"
	"os"
)

// 处理文件转移的逻辑
func ProcessTransfer(msg []byte) bool {
	// 解析 msg
	pubData := mq.TransferData{}

	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		return false
	}
	// 根据临时存储文件路径，创建文件句柄
	filed, err := os.Open(pubData.CurLocation)
	if err != nil {
		return false
	}

	// 通过文件句柄将文件读出来并上传到 oss
	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(filed))
	if err != nil {
		return false
	}

	// 更新文件存储路径到文件表
	success := db.UpdateFileLocation(pubData.FileHash, pubData.DestLocation)
	return success
}

func main() {
	mq.StartConsumer(
		config.TransOSSQueueName,
		"transfer-oss",
		ProcessTransfer)
}
