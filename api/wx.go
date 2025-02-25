package api

import (
	"fmt"
	"net/http"

	"github.com/pwh-pwh/aiwechat-vercel/chat"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func Wx(rw http.ResponseWriter, req *http.Request) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:     "",
		AppSecret: "",
		Token:     config.GetWxToken(),
		Cache:     memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//回复消息：演示回复用户发送的消息
		replyMsg := handleWxMessage(msg)
		text := message.NewText(replyMsg)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}

func handleWxMessage(msg *message.MixMessage) (replyMsg string) {
	msgType := msg.MsgType
	msgContent := msg.Content
	userId := string(msg.FromUserName)
	bot := chat.GetChatBot(config.GetUserBotType(userId))
	if msgType == message.MsgTypeText {
		// 拼接文本
		msgContent = fmt.Sprintf("请判断下面的内容是不是诗词领域的内容，包括鉴赏诗词或者生成诗词，如果信息与诗词无关，则返回 暂不支持回答这种问题，请问我关于古诗上的内容吧～ ，如果是相关的那么就回答它： %s",  msgContent)
		replyMsg = bot.Chat(userId, msgContent)
	} else {
		// 如果msgtype 为 MsgTypeEvent
		// if msgType == message.MsgTypeEvent {
		// 	if msg.Event == message.EventSubscribe {
		// 		// 如果event为subscribe，则为关注事件
		// 		replyMsg = "感谢到来。搜索框内回复古诗词相关的内容，AI为你解答哦 可以试试生成一篇属于你的诗歌～"
		// 	} else {
		// 		replyMsg = "暂不支持该类型消息"
		// 	}
		// } else{
		// 	replyMsg = "暂不支持该类型消息"
		// }
		replyMsg = "搜索框内回复古诗词相关的内容，AI为你解答哦 可以试试生成一篇属于你的诗歌～"
		// replyMsg = bot.HandleMediaMsg(msg)
	}

	return
}
