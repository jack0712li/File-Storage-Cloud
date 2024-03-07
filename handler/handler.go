package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	dblayer "filestore-server/db"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// UploadHandler : handles the upload of files
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := os.ReadFile("static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internal Server Error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err: %s\n", err.Error())
			io.WriteString(w, "Failed to get data")
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: os.TempDir() + "\\" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),

		}
	
		tempFilePath := fileMeta.Location
		newFile, err := os.Create(tempFilePath)
		if err != nil {
			fmt.Printf("Failed to create file, err: %s\n", err.Error())
			io.WriteString(w, "Failed to create file")
			return
		}


		defer newFile.Close()
	
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err: %s\n", err.Error())
			io.WriteString(w, "Failed to save data into file")
			return
		}
		
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		// meta.UpdateFileMeta(fileMeta)
		_= meta.UpdateFileMetaDB(fileMeta)
		fmt.Print(fileMeta.FileSha1)
	
		// 更新用户文件记录
		r.ParseForm()
		username := r.Form.Get("username")
		isSuc := dblayer.OnUserFiledUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if isSuc {
			// 上传完成，跳转到home页面
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		} else {
			w.Write([]byte("upload failed:  更新用户文件表记录失败"))
		}
	}

}

// UploadSucHandler : handles the success of file upload
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Success")
}

// GetFileMetaHandler : get file metadata
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash:= r.Form["filehash"][0]
	// fMeta := meta.GetFileMeta(filehash)
	fMeta,err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data,err:= json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1:= r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f,err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data , err := io.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// FileMetaUpdateHandler : update file metadata
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType:= r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// FileDeleteHandler : delete file
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fileSha1)
	os.Remove(fm.Location)
	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}
