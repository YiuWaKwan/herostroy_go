package websocket

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/reflect/protoreflect"
	"hero_story.go_server/biz_server/handler"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
	"hero_story.go_server/comm/main_thread"
	"time"
)

const oneSecond = 1000
const readMsgCountPerSecond = 16

// CmdContextImpl 就是 MyCmdContext 的 WebSocket 实现
type CmdContextImpl struct {
	userId       int64
	clientIpAddr string
	Conn         *websocket.Conn
	sendMsgQ     chan protoreflect.ProtoMessage // BlockingQueue
	SessionId    int32
}

func (ctx *CmdContextImpl) BindUserId(val int64) {
	ctx.userId = val
}

func (ctx *CmdContextImpl) GetUserId() int64 {
	return ctx.userId
}

func (ctx *CmdContextImpl) GetClientIpAddr() string {
	return ctx.clientIpAddr
}

func (ctx *CmdContextImpl) Write(msgObj protoreflect.ProtoMessage) {
	if nil == msgObj ||
		nil == ctx.Conn ||
		nil == ctx.sendMsgQ {
		return
	}

	ctx.sendMsgQ <- msgObj // queue.push
}

func (ctx *CmdContextImpl) SendError(errorCode int, errorInfo string) {
}

func (ctx *CmdContextImpl) Disconnect() {
	if nil != ctx.Conn {
		_ = ctx.Conn.Close()
	}
}

// LoopSendMsg 循环发送消息,
// 内部通过协程来实现
func (ctx *CmdContextImpl) LoopSendMsg() {
	// 首先构建发送队列
	ctx.sendMsgQ = make(chan protoreflect.ProtoMessage, 64)

	go func() { // new Thread().start(90 -> { ... })
		for {
			msgObj := <-ctx.sendMsgQ // queue.pop

			if nil == msgObj {
				continue
			}

			byteArray, err := msg.Encode(msgObj)

			if nil != err {
				log.Error("%+v", err)
				return
			}

			if err := ctx.Conn.WriteMessage(websocket.BinaryMessage, byteArray); nil != err {
				log.Error("%+v", err)
			}
		}
	}() // 相当于启动一个线程, 专门负责发送消息
}

// LoopReadMsg 循环读取消息
func (ctx *CmdContextImpl) LoopReadMsg() {
	if nil == ctx.Conn {
		return
	}

	// 设置读取字节数限制
	ctx.Conn.SetReadLimit(64 * 1024)

	t0 := int64(0)
	counter := 0

	for {
		_, msgData, err := ctx.Conn.ReadMessage()

		if nil != err {
			log.Error("%+v", err)
			break
		}

		t1 := time.Now().UnixMilli()

		if (t1 - t0) > oneSecond {
			t0 = t1
			counter = 0
		}

		if counter >= readMsgCountPerSecond {
			log.Error("消息过于频繁")
			continue
		}

		counter++

		msgCode := binary.BigEndian.Uint16(msgData[2:4])
		newMsgX, err := msg.Decode(msgData[4:], int16(msgCode))

		if nil != err {
			log.Error(
				"消息解码错误, msgCode = %d, error = %+v",
				msgCode, err,
			)
			continue
		}

		log.Info(
			"收到客户端消息, msgCode = %d, msgName = %s",
			msgCode,
			newMsgX.Descriptor().Name(),
		)

		// 创建指令处理器
		cmdHandler := handler.CreateCmdHandler(msgCode)

		if nil == cmdHandler {
			log.Error(
				"未找到指令处理器, msgCode = %d",
				msgCode,
			)
			continue
		}

		main_thread.Process(func() {
			cmdHandler(ctx, newMsgX)
		})
	}

	// 处理用户离线逻辑
	handler.OnUserQuitHandler(ctx)
}
