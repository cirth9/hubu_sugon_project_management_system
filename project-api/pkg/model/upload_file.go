package model

type UploadFileReq struct {
	TaskCode         string `form:"taskCode"`
	ProjectCode      string `form:"projectCode"`
	ProjectName      string `form:"projectName"`
	TotalChunks      int    `form:"totalChunks"`
	RelativePath     string `form:"relativePath"`
	Filename         string `form:"filename"`
	ChunkNumber      int    `form:"chunkNumber"`
	ChunkSize        int    `form:"chunkSize"`
	CurrentChunkSize int    `form:"currentChunkSize"`
	TotalSize        int    `form:"totalSize"`
	Identifier       string `form:"identifier"`
}
