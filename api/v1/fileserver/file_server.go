package fileserver

import (
	"bychat/api/v1/base"
	"bychat/internal/common"
	"bychat/internal/models"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	data := make(map[string]interface{})
	fileType := c.PostForm("fileType")

	var dst = "./api/v1/fileserver/bychat/%s/"
	// 单文件
	// file, _ := c.FormFile("file")

	file, h, err := c.Request.FormFile("file")
	if err != nil {
		base.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	defer file.Close()

	var resURL string
	messageFileType := models.MessageType(fileType)
	switch messageFileType {
	case models.MessageTypeImage:
		config, _, err := image.DecodeConfig(file)
		if err != nil {
			base.Response(c, common.ParameterIllegal, err.Error(), nil)
			return
		}

		width, height := config.Width, config.Height
		data["width"] = width
		data["height"] = height
		resURL = "/fileserver/bychat/" + fileType
		dst = fmt.Sprintf(dst, fileType)
	case models.MessageTypeFile, models.MessageTypeVedio, models.MessageTypeAudio:
		dst = fmt.Sprintf(dst, fileType)
	default:
		base.Response(c, common.ParameterIllegal, "", nil)
		return
	}

	token := time.Now().UnixNano()

	fileName := h.Filename
	size := h.Size

	tokenExtName := fmt.Sprintf("%d%s", token, filepath.Ext(fileName))
	dst = fmt.Sprintf(dst+"%s", tokenExtName)
	logrus.WithField("dst", dst).Info("UploadFile dst")
	// 上传文件至指定的完整文件路径
	c.SaveUploadedFile(h, dst)

	data["token"] = tokenExtName
	data["name"] = fileName
	data["size"] = size
	data["fileType"] = fileType
	data["url"] = fmt.Sprintf("%s/%s", resURL, tokenExtName)
	base.Response(c, common.OK, "", data)
}
