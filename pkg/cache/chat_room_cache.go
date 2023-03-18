package cache

import (
	"bychat/internal/models"
)

var roomCache = make(map[uint32][]*models.UserOnline)

// SetChatRoomUser 设置缓存
func SetChatRoomUser(roomID uint32, user *models.UserOnline) {
	roomCache[roomID] = append(roomCache[roomID], user)
}

// DelChatRoomUser 设置缓存
func DelChatRoomUser(roomID, userID uint32) {
	for i, r := range roomCache[roomID] {
		if r.ID == userID {
			roomCache[roomID] = append(roomCache[roomID][:i], roomCache[roomID][i+1:]...)
			break
		}
	}
}

// GetChatRoomUser 设置缓存
func GetChatRoomUser(roomID uint32) []*models.UserOnline {
	return roomCache[roomID]
}

// GetChatRoomID 获取房间ID
func GetChatRoomID(userID uint32) uint32 {
	for roomID, list := range roomCache {
		for _, r := range list {
			if r.ID == userID {
				return roomID
			}
		}
	}
	return 0
}
