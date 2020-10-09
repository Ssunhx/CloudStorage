package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

type ShalStream struct {
	_shal hash.Hash
}

func (obj *ShalStream) Update(data []byte) {
	if obj._shal == nil {
		obj._shal = sha1.New()
	}
	obj._shal.Write(data)
}

func (obj *ShalStream) Sum() string {
	return hex.EncodeToString(obj._shal.Sum([]byte("")))
}

func Shal(data []byte) string {
	_shal := sha1.New()
	_shal.Write(data)
	return hex.EncodeToString(_shal.Sum([]byte("")))
}

func FileShal(file *os.File) string {
	_shal := sha1.New()
	io.Copy(_shal, file)
	return hex.EncodeToString(_shal.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

// 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 获取文件大小
func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
