package main

type Uuid string

type PluginRegister struct {
	Name            string   `json:"name"`
	UrlToReset      string   `json:"urlToReset"`
	UrlToRun        string   `json:"urlToRun"`
	UrlStatus       string   `json:"urlStatus"`
	CronjobSchedule string   `json:"cronjobSchedule"`
	FilExtension    []string `json:"filExtension"`
}

type PluginStatus struct {
	IsRunning bool `json:"isRunning"`
}

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

type RequestPluginRun struct {
	PluginName string `json:"pluginName"`
	LogMessage string `json:"logMessage"`
	Data       string `json:"data"`
}

type RequestPluginRunResponse struct {
	Status string `json:"status"`
}
