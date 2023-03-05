package cache

import (
	"bychat/internal/models"
)

var roomCache = make(map[uint32][]*models.UserOnline)

func SetRoomUser(roomID uint32, user *models.UserOnline) {

	roomCache[roomID] = append(roomCache[roomID], user)
}

func DelRoomUser(roomID, userID uint32) {
	for i, r := range roomCache[roomID] {
		if r.ID == userID {
			roomCache[roomID] = append(roomCache[roomID][:i], roomCache[roomID][i+1:]...)
			break
		}
	}
}

func GetRoomUser(roomID uint32) []*models.UserOnline {

	return roomCache[roomID]
}
