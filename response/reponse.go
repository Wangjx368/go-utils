package response

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"io"
	myErr "openapiai/utils/errors"
	"openapiai/utils/meta"
	"openapiai/utils/request"
	"openapiai/utils/wslock"
	"os"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResponseMergency struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Sid  string      `json:"sid"`
	Data interface{} `json:"data"`
}

type InternalResponse struct {
	Stat    int    `json:"stat" description:"返回状态码 1-正常,其他异常"`
	Message string `json:"message" description:"返回信息描述"`
}

type SendMsg struct {
	Flag      string
	Err       error
	ErrorCode int
	Msg       interface{}
}

func NewResponse(ctx *context.Context, code int, msg string, data interface{}) *Response {
	res := &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	res.SetResponseContext(ctx)
	return res
}

func NewResponseMergency(ctx *context.Context, code int, msg string, data interface{}) *ResponseMergency {
	sid := ""
	if ctx.Input.GetData("Sid") != nil {
		sid = ctx.Input.GetData("Sid").(string)
	}
	res := &ResponseMergency{
		Code: code,
		Msg:  msg,
		Sid:  sid,
		Data: data,
	}
	res.SetResponseContextMergency(ctx)
	return res
}

func NewInternalResponse(ctx *context.Context, code int, msg string, data interface{}) *InternalResponse {
	res := &InternalResponse{
		Stat:    code,
		Message: msg,
	}
	return res
}

func (res *Response) SetResponseContext(ctx *context.Context) {
	IsSuccessful := (res.Code == 0)
	rc := meta.ResponseContext{
		Code:         res.Code,
		IsSuccessful: IsSuccessful,
	}
	request.SetInputDataRespContext(ctx, rc)
}

func (res *ResponseMergency) SetResponseContextMergency(ctx *context.Context) {
	IsSuccessful := (res.Code == 0)
	rc := meta.ResponseContext{
		Code:         res.Code,
		IsSuccessful: IsSuccessful,
	}
	request.SetInputDataRespContext(ctx, rc)
}

func (res *Response) ServeJSON(ctx *context.Context) {
	resJson, err := json.Marshal(res)
	if err != nil {
		res.Code = -15000
	}
	ctx.WriteString(string(resJson))
}

func (res *Response) WsServeJSON(ws *websocket.Conn, ctx *context.Context) {
	resJson, err := json.Marshal(res)
	if err != nil {
		res.Code = -15000
	}
	ws.WriteJSON(resJson)
}

func SendChan(ch chan<- SendMsg, flag string, theErr error, errorCode int, msg interface{}) {
	sm := SendMsg{
		flag,
		theErr,
		errorCode,
		msg,
	}
	ch <- sm
}

func WritePump(wser *wslock.WSer, ch <-chan SendMsg, ctx *context.Context) {
	defer func() {
		logs.Notice("write close")
		wser.Close()
	}()
	for {
		select {
		case sm := <-ch:
			errMsg := myErr.ErrorsMap[sm.ErrorCode]
			switch sm.Flag {
			case "success":
				res := NewResponse(ctx, sm.ErrorCode, errMsg, sm.Msg)
				wserr := wser.Write(res)
				if wserr != nil {
					wser.Close()
				}
			case "error":
				errLog := errMsg + fmt.Sprintf(" 错误详情为：%s ", sm.Err)
				logs.Error(errLog)
				res := NewResponse(ctx, sm.ErrorCode, errLog, nil)
				wserr := wser.Write(res)
				if wserr != nil {
					wser.Close()
				}
			case "fatal":
				errLog := errMsg + fmt.Sprintf(" 错误详情为：%s ", sm.Err)
				logs.Error(errLog)
				res := NewResponse(ctx, sm.ErrorCode, errLog, nil)
				wser.Write(res)
				wser.Close()
			}
		}
	}
}

// WritePumpIgnoreError 老鼠英语专用
func WritePumpIgnoreError(wser *wslock.WSer, ch <-chan SendMsg, ctx *context.Context) {
	defer func() {
		logs.Notice("write close")
		wser.Close()
	}()
	for {
		select {
		case sm := <-ch:
			errMsg := myErr.ErrorsMap[sm.ErrorCode]
			switch sm.Flag {
			case "nlp": // 此时，ErrorCode为待识别的tts文本个数
				res := NewResponse(ctx, 0, myErr.ErrorsMap[0], sm.Msg)
				wserr := wser.Write(res)
				// 更新ctx
				SetCtxRespNLPNum(ctx, sm)
				if wserr != nil {
					wser.Close()
				}

			case "tts":
				res := NewResponse(ctx, 0, myErr.ErrorsMap[0], sm.Msg)
				wserr := wser.Write(res)
				// 更新ctx
				SetCtxRespTTSNum(ctx, sm)
				// 向数据库插入trace数据
				InsertCtxTrace(ctx, sm)
				if wserr != nil {
					wser.Close()
				}

			case "error":
				saveRespFiles(ctx, sm.ErrorCode, errMsg, sm.Err) // 记录错误响应信息
				if sm.ErrorCode != 16017 && sm.ErrorCode != 16106 && sm.ErrorCode != 16109 && sm.ErrorCode != 16110 {
					errLog := errMsg + fmt.Sprintf(" 错误详情为：%s ", sm.Err)
					logs.Error(errLog)
					res := NewResponseMergency(ctx, sm.ErrorCode, errMsg, nil)
					// 向数据库插入trace数据
					InsertCtxTrace(ctx, sm)
					wserr := wser.Write(res)
					if wserr != nil {
						wser.Close()
					}
				}

			case "fatal":
				saveRespFiles(ctx, sm.ErrorCode, errMsg, sm.Err) // 记录错误响应信息
				errLog := errMsg + fmt.Sprintf(" 错误详情为：%s ", sm.Err)
				logs.Error(errLog)
				res := NewResponseMergency(ctx, sm.ErrorCode, errMsg, nil)
				// 向数据库插入trace数据
				InsertCtxTrace(ctx, sm)
				wser.Write(res)
				wser.Close()
			}
		}
	}
}

func saveRespFiles(ctx *context.Context, errorCode int, errMsg string, err error) {
	if errorCode == 16102 || errorCode == 16017 {
		return
	}
	// init
	usrID := ""
	sid := ""
	dialogStatus := ""
	if ctx.Input.GetData("UsrId") != nil {
		usrID = ctx.Input.GetData("UsrId").(string)
		sid = ctx.Input.GetData("Sid").(string)
		dialogStatusTmp := ctx.Input.GetData("NlpDialogStatus").(int)
		dialogStatus = fmt.Sprintf("%d", dialogStatusTmp)
	}

	// 需要向数据库中插入的数据
	respInfoFile := "./savefiles/" + "resp" + ".info"
	fileObj, err := os.OpenFile(respInfoFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		logs.Debug("Failed to open the file", err.Error())
		os.Exit(2)
	}
	defer fileObj.Close()
	respInfo := fmt.Sprintf("%d", errorCode) + "\t" + errMsg + fmt.Sprintf("%s", err) + "\t" + dialogStatus + "\t" + usrID + "\t" + sid + "\n"
	io.WriteString(fileObj, respInfo)
}

// WsRespError WebSocket返回给调用方报错信息
func WsRespError(wser *wslock.WSer, ctx *context.Context, errorCode int, errMsg string, err error) {
	errLog := errMsg + fmt.Sprintf(" 错误详情为：%s ", err)
	logs.Error(errLog)
	res := NewResponse(ctx, errorCode, errMsg, nil)
	wserr := wser.Write(res)
	if wserr != nil {
		wser.Close()
	}
}

// WsRespFatalError WebSocket返回给调用方报错信息,并且关闭连接
func WsRespFatalError(wser *wslock.WSer, ctx *context.Context, errorCode int, errMsg string, err error) {
	errLog := errMsg + fmt.Sprintf(" 错误详情为：%s ", err)
	logs.Error(errLog)
	res := NewResponse(ctx, errorCode, errMsg, nil)
	wser.Write(res)
	wser.Close()
}
