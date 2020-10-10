package meta

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

// 根据 hash 获取 FileMeta
func GetFileMeta(fileShal string) FileMeta {
	return fileMetas[fileShal]
}

// 根据 hash 删除
func RemoveFileMeta(filehash string) {
	delete(fileMetas, filehash)
}
