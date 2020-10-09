package main

import (
	"cloudstorage/handler"
	"fmt"
	"net/http"
)

func main() {
	// 文件上传
	http.HandleFunc("/file/upload", handler.UploadHandler)
	// 上传成功
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	// 获取文件信息
	http.HandleFunc("/file/meta", handler.GetFileMetahandler)
	// 文件下载
	http.HandleFunc("/file/download", handler.DownloadHandler)
	err := http.ListenAndServe(":9234", nil)
	if err != nil {
		fmt.Println("failed start server")
	}
}
