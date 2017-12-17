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

func GetBaseSucessRes() BaseResultEntity {
	return BaseResultEntity{IsSucess:true,Message:"响应成功",Code:"200"}
}

func GetBaseFailRes() BaseResultEntity {
	return BaseResultEntity{IsSucess:false,Message:"响应失败",Code:"300"}
}
