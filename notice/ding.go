package notice

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"openapiai/utils/request"
)

type DingMsgText struct {
	Content string `json:"content"`
}

type DingMsgAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingMsg struct {
	MsgType string      `json:"msgtype"`
	Text    DingMsgText `json:"text"`
	At      DingMsgAt   `json:"at"`
}

// 通用发送消息
func SendDingMsg(api string, content string) {
	ip := request.GetIntranetIp()
	content = content + " [cosumer_msg_ip] : " + ip + "\n"
	msgByte := generateMsg(content)
	req := request.NewRequest()
	// api notice
	noticeApi, ok := DingTokensMap[api]
	if !ok {
		// service notice
		noticeApi = ServiceDingTokensMap[api]
	}

	runmode := beego.AppConfig.String("runmode")
	if noticeApi != "" && runmode == "prod" {
		req.DoHttpPost(noticeApi, msgByte)
	}
}

func generateMsg(content string) []byte {
	text := DingMsgText{
		Content: content,
	}

	at := DingMsgAt{
		IsAtAll: true,
	}

	msg := DingMsg{
		MsgType: "text",
		Text:    text,
		At:      at,
	}

	msgByte, _ := json.Marshal(msg)
	return msgByte
}

func ProcessDingNoticeMsgByApi(api string, content string) {
	SendDingMsg(api, content)
}

func ApiNotice(api string, appKey string, traceId string, rpcId string, errMsg string) string {
	ip := request.GetIntranetIp()
	return fmt.Sprintf(" [api] : %s\n [appKey] : %s\n [traceId] : %s\n [rpcId] : %s failed.\n [errMsg] : %s\n [produce_msg_ip] : %s\n", api, appKey, traceId, rpcId, errMsg, ip)
}
