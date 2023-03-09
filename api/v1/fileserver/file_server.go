package fileserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadFile 上传文件
func UploadFile(ctx *gin.Context) {
	// 单文件
	file, _ := ctx.FormFile("file")
	log.Println(file.Filename)

	dst := "./api/v1/fileserver/file/" + file.Filename
	// 上传文件至指定的完整文件路径
	ctx.SaveUploadedFile(file, dst)

	ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
