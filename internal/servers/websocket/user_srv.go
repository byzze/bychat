package websocket

func UserList(appID uint32) []*Client {
	c := clientManager.GetUserClients()
	return c
}
