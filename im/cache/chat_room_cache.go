package cache

var roomCache = make(map[uint32][]uint32)

// SetChatRoomUser 设置缓存
func SetChatRoomUser(roomID, userID uint32) {
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
func GetChatRoomUser(roomID uint32) []uint32 {
	return roomCache[roomID]
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
