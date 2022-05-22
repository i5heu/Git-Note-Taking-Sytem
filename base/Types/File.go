package types



type FileGet struct {
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
}

type FileGetResponse struct {
	FileName     string `json:"fileName"`
	FilePath     string `json:"filePath"`
	FilExtension string `json:"filExtension"`
	FileContent  string `json:"fileContent"` // base64
	LastModified string `json:"lastModified"`
}

type FileInfoGet struct {
	FileName     string `json:"fileName"`
	FilePath     string `json:"filePath"`
	FilExtension string `json:"filExtension"`
	FileContent  string `json:"fileContent"` // base64
	LastModified string `json:"lastModified"`
}

type FileListResponse struct {
	List []FileInfoGet `json:"list"`
}

type FilePut struct {
	FileName     string `json:"fileName"`
	FilePath     string `json:"filePath"`
	FileContent  string `json:"fileContent"` // base64
	LastModified string `json:"lastModified"`
}

type FilePutResponse struct {
	Status string `json:"status"`
}

type FileDelete struct {
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
}

type FileDeleteResponse struct {
	Status string `json:"status"`
}
