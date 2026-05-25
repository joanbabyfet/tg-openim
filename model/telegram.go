package model

type OpenIMRegisterReq struct {
	Users []OpenIMUser `json:"users"`
}

type OpenIMUser struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
}