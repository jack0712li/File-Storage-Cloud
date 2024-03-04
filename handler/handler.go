package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"filestore-server/meta"
	"time"
	"filestore-server/util"
	"encoding/json"
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
		meta.UpdateFileMeta(fileMeta)
		fmt.Print(fileMeta.FileSha1)
	
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
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
	fMeta := meta.GetFileMeta(filehash)
	data,err:= json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}