package errcode

import (
	"fmt"
	"sort"
)

type ErrorMessage struct {
	Code    string
	Message string
	Details []string
}

var codeSet = make(map[string]ErrorMessage, 100)

// 新建错误
func New(errorCode string, errorMessage string) *ErrorMessage {
	if _, ok := codeSet[errorCode]; ok {
		msg := fmt.Sprintf("errorCode exist: %s, use other\n", errorCode)
		panic(msg)
	}

	eMsg := &ErrorMessage{
		Code:    errorCode,
		Message: errorMessage,
	}
	codeSet[errorCode] = *eMsg

	return eMsg
}

// 实现 Error 接口
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("error code: %d, error message :%s\n", e.Code, e.Message)
}

// 添加错误详细描述信息
func (e *ErrorMessage) WithDetails(details ...string) *ErrorMessage {
	eMsg := *e
	eMsg.Details = []string{}
	for _, d := range details {
		eMsg.Details = append(eMsg.Details, d)
	}
	return &eMsg
}

type ErrorMessages []ErrorMessage

// 实现 sort 接口
func (eMsgs ErrorMessages) Len() int {
	return len(eMsgs)
}
func (eMsgs ErrorMessages) Less(i, j int) bool {
	return eMsgs[i].Code < eMsgs[j].Code
}
func (eMsgs ErrorMessages) Swap(i, j int) {
	eMsgs[i], eMsgs[j] = eMsgs[j], eMsgs[i]
}

func GetAllErrorMessage() ErrorMessages {
	eMsgs := make([]ErrorMessage, 0, len(codeSet))

	for _, eMsg := range codeSet {
		eMsgs = append(eMsgs, eMsg)
	}

	sort.Sort(ErrorMessages(eMsgs))

	return eMsgs
}
