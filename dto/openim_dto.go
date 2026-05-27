package dto

type OpenIMRegisterReq struct {
	Users []OpenIMUser `json:"users"`
}

type OpenIMUser struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
}

type OpenIMMessage struct {
	SendID      string `json:"sendID"`
	RecvID      string `json:"recvID"`
	Content     string `json:"content"`
	ContentType int    `json:"contentType"`
}

type TextContent struct {
	Content string `json:"content"`
}