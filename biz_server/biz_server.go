package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"hero_story.go_server/biz_server/network/broadcaster"
	mywebsocket "hero_story.go_server/biz_server/network/websocket"
	"hero_story.go_server/comm/log"
	"net/http"
	"os"
	"path"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var sessionId int32 = 0

func main() {
	fmt.Println("start bizServer")

	ex, err := os.Executable()

	if nil != err {
		panic(err)
	}

	log.Config(path.Dir(ex) + "/log/biz_server.log")
	log.Info("Hello")

	//
	// 可以用 `ab -n 10000 -c 8 http://127.0.0.1/websocket` 来测试一下性能
	//
	//http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
	//	_, _ = w.Write([]byte("Hello, the World!"))
	//})
	//
	http.HandleFunc("/websocket", webSocketHandshake)
	_ = http.ListenAndServe("127.0.0.1:12345", nil)
}

func webSocketHandshake(w http.ResponseWriter, r *http.Request) {
	if nil == w ||
		nil == r {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if nil != err {
		log.Error("WebSocket upgrade error, %v+", err)
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	log.Info("有新客户端连入")

	sessionId += 1

	ctx := &mywebsocket.CmdContextImpl{
		Conn:      conn,
		SessionId: sessionId,
	}

	// 将指令上下文添加到分组,
	// 当断开连接时移除指令上下文...
	broadcaster.AddCmdCtx(sessionId, ctx)
	defer broadcaster.RemoveCmdCtxBySessionId(sessionId)

	// 循环发送消息
	ctx.LoopSendMsg()
	// 开始循环读取消息
	ctx.LoopReadMsg()
}
