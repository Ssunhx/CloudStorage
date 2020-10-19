package main

import (
	"cloudstorage/handler"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// 文件上传
	http.HandleFunc("/file/upload", handler.UploadHandler)
	// 上传成功
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	// 获取文件信息
	http.HandleFunc("/file/meta", handler.GetFileMetahandler)
	// 文件下载
	http.HandleFunc("/file/download", handler.DownloadHandler)
	// 文件更新
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	// 文件删除
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	// 文件查询
	http.HandleFunc("/file/query", handler.FileQueryhandler)

	// 分块上传
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.InitialmultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complate", handler.HTTPInterceptor(handler.ComplateUploadHandler))

	// 注册
	http.HandleFunc("/user/signup", handler.SignUPHandler)
	// 用户登录
	http.HandleFunc("/user/signin", handler.SignInHandler)
	// 用户信息
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))
	// 秒传
	http.HandleFunc("/file/fastupload", handler.TryFastUploadhandler)

	// oss download
	http.HandleFunc("/file/downloadurl", handler.HTTPInterceptor(handler.OSSDownloadURLHandler))

	err := http.ListenAndServe(":9234", nil)
	if err != nil {
		fmt.Println("failed start server")
	}
}
