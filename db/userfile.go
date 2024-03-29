package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

// OnUserFiledUploadFinished: 更新用户文件表
func OnUserFiledUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, e := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`,`file_size`,`upload_at`) values (?,?,?,?,?)")

	if e != nil {
		fmt.Println(e.Error())
		return false
	}
	defer stmt.Close()
	_, e = stmt.Exec(username, filehash, filename, filesize, time.Now())

	if e != nil {
		return false
	}
	return true
}

// QueryUserFileMetas: get user files meta
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, e := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name=? limit ?")
	if e != nil {
		fmt.Println(e.Error())
		return nil, e
	}
	defer stmt.Close()

	rows, e := stmt.Query(username, limit)
	if e != nil {
		return nil, e
	}
	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		e = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if e != nil {
			fmt.Println(e.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}