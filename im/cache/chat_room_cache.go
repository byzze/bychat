package cache

import (
	"bychat/im/models"

	"github.com/sirupsen/logrus"
)

// TODO 后续修改为redis
var roomCache = make(map[uint32][]uint32)

// AddChatRoomUser 设置缓存
func AddChatRoomUser(roomID, userID uint32) {
	roomCache[roomID] = append(roomCache[roomID], userID)
}

// DelChatRoomUser 设置缓存
func DelChatRoomUser(roomID, userID uint32) {
	for i, id := range roomCache[roomID] {
		if id == userID {
			roomCache[roomID] = append(roomCache[roomID][:i], roomCache[roomID][i+1:]...)
			break
		}
	}
}

// GetChatRoomUser 设置缓存
func GetChatRoomUser(roomID uint32) []*models.UserOnline {
	var list = roomCache[roomID]
	var res = make([]*models.UserOnline, 0)
	for _, v := range list {
		u, err := GetUserOnlineInfo(v)
		if err != nil {
			logrus.WithError(err).Error("GetUserOnlineInfo")
			continue
		}
		if u != nil {
			res = append(res, u)
		}
	}
	return res
}

// GetChatRoomID 获取房间ID
func GetChatRoomID(userID uint32) uint32 {
	for roomID, list := range roomCache {
		for _, id := range list {
			if id == userID {
				return roomID
			}
		}
	}
	return 0
}
