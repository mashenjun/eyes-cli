package main

type ErrBody struct {
	Code  int    `json:"ret"`
	Msg   string `json:"msg"`
	Cause string `json:"cause"`
}

//----------------------------------------------------------------------------------------------------------------------
type LoginPostBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code int           `json:"ret"`
	Data LoginRespData `json:"data"`
}

type LoginRespData struct {
	ExpDuration int    `json:"exp_duration"`
	ExpUnit     string `json:"exp_unit"`
	Token       string `json:"token"`
	UserID      int    `json:"user_id"`
}

//----------------------------------------------------------------------------------------------------------------------
type CreateFacelibPostBody struct {
	Name        string `json:"name"`
	Count       int    `json:"count"`
	Description string `form:"description"`
	Source      string `form:"source"`
}

type CreateFacelibResponse struct {
	Code int                   `json:"ret"`
	Data CreateFacelibRespData `json:"data"`
}

type CreateFacelibRespData struct {
	ID             int      `json:"id"`
	GroupID        string   `json:"group_id"`
	HumanInfoNames []string `json:"human_info_names"`
}

//----------------------------------------------------------------------------------------------------------------------
type UploadPostBody struct {
	// TODO
	FacelibID      int    `form:"facelib_id"`
	FacelibGroupID string `json:"facelib_group_id"`
	Name           string `form:"name"`
	Description    string `form:"description"`
	Source         string `form:"source"`
	Gender         uint8  `form:"gender"`
	IDCardNum      string `form:"id_card_num"`
	Count          uint   `form:"count"`
}

type UploadResponse struct {
	Code int            `json:"ret"`
	Data UploadRespData `json:"data"`
}

type UploadRespData struct {
	FacelibID int    `json:"facelib_id"`
	HumanID   int    `json:"human_id"`
	HumanName string `json:"human_name"`
}

//----------------------------------------------------------------------------------------------------------------------
type AddFaceMetaData struct {
	Name      string   `json:"name"`
	IdCardNum string   `json:"id_card_num"`
	Gender    int      `json:"sex"`
	Images    []string `json:"images"`
}

type CreateFaceMetaData struct {
	Name        string `json:"name"`
	Count       int    `json:"size"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

//----------------------------------------------------------------------------------------------------------------------
type EndAddingPostBody struct {
	FacelibID int `json:"facelib_id"`
	Count     int `json:"count"`
}

type EndAddingResponse struct {
	Code int               `json:"ret"`
	Data EndAddingRespData `json:"data"`
}

type EndAddingRespData struct {
	Already int `json:"already"`
}
