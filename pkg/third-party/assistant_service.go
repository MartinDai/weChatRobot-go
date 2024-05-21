package third_party

type AssistantService interface {
	ProcessText(fromUserName, toUserName, content string) interface{}
}
