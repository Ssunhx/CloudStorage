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
	err := http.ListenAndServe(":9234", nil)
	if err != nil {
		fmt.Println("failed start server")
	}
}
