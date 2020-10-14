package meta

import (
	"cloudstorage/db"
	"fmt"
)

// FileMeta 文件元信息结构
type FileMeta struct {
	FileShal string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 修改 FileMeta
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileShal] = fmeta
}

// 新增、修改文件元信息到 mysql
func UpdateFileMetaDB(fmeta FileMeta) bool {
	fmt.Println(fmeta)
	return db.OnFileUploadFinished(fmeta.FileShal, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 根据 hash 获取 FileMeta
func GetFileMeta(fileShal string) FileMeta {
	return fileMetas[fileShal]
}

func GetFileMetaDB(filehash string) (FileMeta, error) {
	tfil, err := db.GetFileMeta(filehash)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileShal: tfil.FileHash,
		FileName: tfil.FileName.String,
		FileSize: tfil.FileSize.Int64,
		Location: tfil.FileAddr.String,
	}
	return fmeta, nil
}

// 根据 hash 删除
func RemoveFileMeta(filehash string) {
	delete(fileMetas, filehash)
}
