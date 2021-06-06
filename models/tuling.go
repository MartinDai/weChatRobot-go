package models

//图灵机器人响应码常量

const ParamErrCode = 4000   //请求参数格式错误
const NoResultCode = 5000   //无解析结果
const NoApiTimesCode = 4003 //该apikey没有可用请求次数
const TextCode = 10004

const TextResultType = "text"

type ReqParam struct {
	ReqType    int        `json:"reqType"`
	Perception Perception `json:"perception"`
	UserInfo   UserInfo   `json:"userInfo"`
}

type Perception struct {
	InputText InputText `json:"inputText"`
}

type InputText struct {
	Text string `json:"text"`
}

type UserInfo struct {
	ApiKey string `json:"apiKey"`
	UserId string `json:"userId"`
}
