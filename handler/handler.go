package handler

import (
	"cloudstorage/db"
	"cloudstorage/meta"
	"cloudstorage/store/oss"
	"cloudstorage/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
		//meta.UpdateFileMeta(fileMeta)

		// 写入 ceph
		//newFile.Seek(0, 0)
		//data, _ := ioutil.ReadAll(newFile)
		//bucket := ceph.GetCephBucket("userfile")
		//cephPath := "/ceph/" + fileMeta.FileShal
		//_ = bucket.Put(cephPath, data, "object-stream", s3.PublicRead)
		//fileMeta.Location = cephPath

		// oss 存储
		// oss 路径中不能以 / 开头
		osspath := "oss/" + fileMeta.FileShal
		err = oss.Bucket().PutObject(osspath, newFile)
		if err != nil {
			fmt.Println(err.Error())
			w.Write([]byte("upload oss failed"))
		}
		fileMeta.Location = osspath

		_ = meta.UpdateFileMetaDB(fileMeta)
		r.ParseForm()
		username := r.Form.Get("username")
		suc := db.OnUserfileUploadFinished(username, fileMeta.FileShal, fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		} else {
			w.Write([]byte("Upload failed"))
		}
		//http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传已完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finish")
}

// 获取文件信息
func GetFileMetahandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)

	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 文件批量查询
func FileQueryhandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")

	userFiles, err := db.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")

	filemeta := meta.GetFileMeta(filehash)
	fmt.Println(filemeta)
	f, err := os.Open(filemeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment 表示文件会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+filemeta.FileName+"\"")
	w.Write(data)
}

// 文件修改
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileShal := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	// 只有 0 才可以修改
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 只能是 POST 请求
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 修改名字
	curFileMeta := meta.GetFileMeta(fileShal)
	curFileMeta.FileName = newFileName
	//meta.UpdateFileMeta(curFileMeta)
	_ = meta.UpdateFileMetaDB(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 返回数据
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 文件删除
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileShal := r.Form.Get("filehash")
	// 删除磁盘文件
	fMeta := meta.GetFileMeta(fileShal)
	os.Remove(fMeta.Location)
	// 在 dict 中删除
	meta.RemoveFileMeta(fileShal)
	w.WriteHeader(http.StatusOK)
}

// 秒传接口
func TryFastUploadhandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 1、解析参数
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 2、文件表中查询相同 hash 的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3、查不到记录
	_fileMeta := meta.FileMeta{}
	if fileMeta == _fileMeta {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}

	suc := db.OnUserfileUploadFinished(username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败",
		Data: nil,
	}
	w.Write(resp.JSONBytes())
	return
}

// oss 文件下载
func OSSDownloadURLHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")

	row, _ := db.GetFileMeta(filehash)

	signedURL := oss.DownloadURL(row.FileAddr.String)

	w.Write([]byte(signedURL))
}
