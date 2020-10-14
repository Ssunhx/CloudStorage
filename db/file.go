package db

import (
	mydb "cloudstorage/db/mysql"
	"database/sql"
	"fmt"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// 文件上传成功，保存 meta 到数据库
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DB.Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("failed to prepare statment, err: " + err.Error())
		return false
	}

	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("failed with hash %s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

// 获取文件元信息 mtea
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DB.Prepare("select file_sha1, file_name, file_size, file_addr from tbl_file " +
		"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}
