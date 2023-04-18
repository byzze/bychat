package fileserver

import (
	"bychat/api/base"
	"bychat/internal/domain/message"
	"bychat/pkg/common"
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

	var dstDir = "./api/v1/fileserver/bychat/%s/"
	// 单文件
	file, h, err := c.Request.FormFile("file")
	if err != nil {
		base.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	defer file.Close()

	dstDir = fmt.Sprintf(dstDir, fileType)

	var reqURL string
	messageType := message.MsgType(fileType)
	switch messageType {
	case message.MsgTypeImage:
		config, _, err := image.DecodeConfig(file)
		if err != nil {
			logrus.WithError(err).Error("DecodeConfig Img")
			base.Response(c, common.ParameterIllegal, err.Error(), nil)
			return
		}

		width, height := config.Width, config.Height
		data["width"] = width
		data["height"] = height
		reqURL = "/fileserver/bychat/" + fileType
	// case models.MessageTypeFile, models.MessageTypeVedio, models.MessageTypeAudio:
	default:
		base.Response(c, common.ParameterIllegal, "", nil)
		return
	}

	token := time.Now().UnixNano()

	fileName := h.Filename
	size := h.Size

	tokenExtName := fmt.Sprintf("%d%s", token, filepath.Ext(fileName))
	dstDir = fmt.Sprintf(dstDir+"%s", tokenExtName)

	logrus.WithField("dst", dstDir).Info("UploadFile dst")
	// 上传文件至指定的完整文件路径
	err = c.SaveUploadedFile(h, dstDir)
	if err != nil {
		logrus.WithError(err).Error("SaveUploadedFile")
		base.Response(c, common.OperationFailure, err.Error(), nil)
		return
	}

	data["token"] = tokenExtName
	data["name"] = fileName
	data["size"] = size
	data["fileType"] = fileType
	data["url"] = fmt.Sprintf("%s/%s", reqURL, tokenExtName)
	base.Response(c, common.OK, "", data)
}
