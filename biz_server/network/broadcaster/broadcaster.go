package broadcaster

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"hero_story.go_server/biz_server/base"
)

var innerMap = make(map[int32]base.MyCmdContext)

// AddCmdCtx 添加指令上下文分组
func AddCmdCtx(sessionId int32, ctx base.MyCmdContext) {
	if nil == ctx {
		return
	}

	innerMap[sessionId] = ctx
}

// RemoveCmdCtxBySessionId 移除指令上下文分组
func RemoveCmdCtxBySessionId(sessionId int32) {
	if sessionId <= 0 {
		return
	}

	delete(innerMap, sessionId)
}

// Broadcast 广播消息
func Broadcast(msgObj protoreflect.ProtoMessage) {
	if nil == msgObj {
		return
	}

	for _, ctx := range innerMap {
		if nil != ctx {
			ctx.Write(msgObj)
		}
	}
}
