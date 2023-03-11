package fileserver

import (
	"bychat/api/v1/base"
	"bychat/internal/common"
	"bychat/internal/models"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	data := make(map[string]interface{})
	fileType := c.PostForm("fileType")

	var dst = "./api/v1/fileserver/file/%s/"

	switch fileType {
	case string(models.MessageTypeFile), string(models.MessageTypeImg),
		string(models.MessageTypeVedio), string(models.MessageTypeSound):
		dst = fmt.Sprintf(dst, fileType)
	default:
		base.Response(c, common.ParameterIllegal, "", nil)
		return
	}
	// 单文件
	file, _ := c.FormFile("file")

	token := time.Now().UnixNano()

	fileName := file.Filename
	size := file.Size

	tokenExtName := fmt.Sprintf("%d%s", token, filepath.Ext(fileName))
	dst = fmt.Sprintf(dst+"%s", tokenExtName)
	logrus.WithField("dst", dst).Info("UploadFile dst")
	// 上传文件至指定的完整文件路径
	c.SaveUploadedFile(file, dst)

	data["token"] = tokenExtName
	data["name"] = fileName
	data["size"] = size
	data["fileType"] = fileType
	base.Response(c, common.OK, "", data)
}
