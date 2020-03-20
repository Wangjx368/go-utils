package ws

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"net/http"
	"openapiai/utils/response"
	"openapiai/utils/wslock"
	"sync"
)

const (
	readBufSize  = 102400
	writeBufSize = 102400
)

func WsConnect(ctx *context.Context) (wsSafe *wslock.WSer) {
	logs.Notice("端上发起的请求建立长链接消息为： ", ctx.Request)
	ws, err := websocket.Upgrade(ctx.ResponseWriter, ctx.Request, nil, readBufSize, writeBufSize)
	wsSafe = NewWSer(ws)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(ctx.ResponseWriter, "Not a websocket handshake", 400)
		errMsg := "WebSocket建立连接时，请求协议不是WebSocket握手。"
		logs.Error(errMsg)
		// myErr.WsRespError(wsSafe, ctx, 16000, errMsg, err)
		return nil
	} else if err != nil {
		errMsg := "WebSocket建立连接时，无法设置WebSocket连接。"
		response.WsRespError(wsSafe, ctx, 16001, errMsg, err)
		return nil
	}
	//ws.WriteMessage(1, []byte("dial ws success")) // ToDO
	return wsSafe
}

func NewWSer(curWS *websocket.Conn) *wslock.WSer {
	wser := new(wslock.WSer)
	wser.WS = curWS
	wser.Lock = new(sync.RWMutex)
	wser.IsClosed = false
	return wser
}
