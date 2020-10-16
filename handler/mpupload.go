package handler

import (
	"cloudstorage/cache/redis"
	"cloudstorage/db"
	"cloudstorage/util"
	"fmt"
	red "github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// 分块的初始化信息
type MultipartUploadinfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// 初始化分块上传
func InitialmultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1、解析用户请求信息
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))

	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}
	// 2、redis 连接
	rConn := redis.RedisPool().Get()
	defer rConn.Close()
	// 3、生成分块上传的初始化信息
	upInfo := MultipartUploadinfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().Unix()),
		ChunkSize:  5 * 1024 * 1024, // 5M
		ChunkCount: int(math.Ceil(float64(filesize)) / (5 * 1024 * 1024)),
	}
	// 4、初始化信息保存 redis
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)
	// 5、初始化信息返回客户端
	w.Write(util.NewRespMsg(0, "ok", upInfo).JSONBytes())
}

func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1、解析参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// 2、连接 redis
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	// 3、获取文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)

	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4、更新 redis
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5、返回处理的结果
	w.Write(util.NewRespMsg(0, "ok", nil).JSONBytes())
}

// 通知上传合并
func ComplateUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1、解析参数
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	// 2、redis 连接
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	// 3、通过 uploadid 查询redis 判断是否完成上传
	data, err := red.Values(rConn.Do("HGET", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complate upload failed", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}

	//4、合并分块

	//5、更新唯一文件表及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	db.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	db.OnUserfileUploadFinished(username, filehash, filename, int64(fsize))

	//6、响应处理
	w.Write(util.NewRespMsg(0, "ok", nil).JSONBytes())
}

// 上传文件分块
func Uploadparthandler(w http.ResponseWriter, r *http.Request) {
	// 1、解析参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// 2、获得 redis 连接
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	// 3、获得文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4、更新redis
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5、返回结果
	w.Write(util.NewRespMsg(0, "ok", nil).JSONBytes())
}
