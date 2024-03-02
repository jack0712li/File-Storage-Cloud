package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
			fmt.Printf("Failed to get data, err: %s", err.Error())
			return
		}
		defer file.Close()
		newFile, err := os.Create("/tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("Failed to create file, err: %s", err.Error())
			return
		}
		defer newFile.Close()
		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err: %s", err.Error())
			return
		}

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

// UploadSucHandler : handles the success of file upload
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Success")
}
