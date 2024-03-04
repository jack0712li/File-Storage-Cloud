package meta

// FileMeta : file metadata structure
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}


var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta : add/update file metadata
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// GetFileMeta : get file metadata
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

