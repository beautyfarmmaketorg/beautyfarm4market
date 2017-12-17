package entity

type BaseResultEntity struct {
	IsSucess bool `json:"isSucess"`
	Message string `json:"message"`
	Code string `json:"code"`
}

type SendMsgResult struct {
	BaseResultEntity
	Mobile string `json:"mobile"`
}
