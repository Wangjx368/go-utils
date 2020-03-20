package wslock

import (
	"github.com/gorilla/websocket"
	"sync"
)

// WSer ws长链接的带锁对象
type WSer struct {
	IsClosed bool
	WS       *websocket.Conn
	Lock     *sync.RWMutex
}

// Write WSer对象向长链接写响应
func (wser *WSer) Write(res interface{}) error {
	wser.Lock.Lock()
	defer wser.Lock.Unlock()
	err := wser.WS.WriteJSON(res)
	return err
}

// Close 安全关断ws链接
func (wser *WSer) Close() {
	wser.Lock.Lock()
	defer wser.Lock.Unlock()
	if !wser.IsClosed {
		wser.WS.Close()
		wser.IsClosed = true
	}
}
