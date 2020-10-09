package handler

import (
	"cloudstorage/meta"
	"cloudstorage/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { //
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Println("fail to get data, err: ", err)
		}

		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			return
		}

		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("fail to save file, err:%s ", err.Error())
			return
		}

		// 生成文件 hash
		newFile.Seek(0, 0)
		fileMeta.FileShal = util.FileShal(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传已完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finish")
}
