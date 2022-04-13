package defines

import (
	"github.com/google/uuid"
	"phoenixbuilder/fastbuilder/uqHolder"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

type LineSource interface {
	Read() string
}

type LineDst interface {
	Write(string)
}

type NoSqlDB interface {
	Get(key string) string
	Delete(key string)
	Commit(key string, v string)
	IterAll(cb func(key string, v string) (stop bool))
	IterWithPrefix(cb func(key string, v string) (stop bool), prefix string)
	IterWithRange(cb func(key string, v string) (stop bool), start, end string)
}

type GameChat struct {
	Name               string
	Msg                []string
	Type               byte
	FrameWorkTriggered bool
	Aux                interface{}
}

type MenuEntry struct {
	Triggers     []string
	ArgumentHint string
	FinalTrigger bool
	Usage        string
}

type GameMenuEntry struct {
	MenuEntry
	OptionalOnTriggerFn func(chat *GameChat) (stop bool)
}

type BackendMenuEntry struct {
	MenuEntry
	OptionalOnTriggerFn func(cmds []string) (stop bool)
}

// CtxProvider 旨在帮助插件发现别的插件主动暴露的接口 GetContext()
// GetUQHolder() 可以获得框架代为维持的信息
type CtxProvider interface {
	GetContext() map[string]interface{}
	GetUQHolder() *uqHolder.UQHolder
}

// ConfigProvider 是帮助一个插件获得和修改别的插件的接口
// 如果仅仅需要自己的配置，这是不必要的
type ConfigProvider interface {
	QueryConfig(name string) interface{}
	GetAllConfigs() []*ComponentConfig
}

// 框架帮忙提供的储存机制，目的在于共享而非沙箱隔离
type StorageAndLogProvider interface {
	GetLogger(topic string) LineDst
	GetNoSqlDB(topic string) NoSqlDB
	GetRelativeFileName(topic string) string
	GetFileData(topic string) ([]byte, error)
	GetJsonData(topic string, data interface{}) error
	WriteFileData(topic string, data []byte) error
	WriteJsonData(topic string, data interface{}) error
}

// 与后端的交互接口
type BackendInteract interface {
	GetBackendDisplay() LineDst
	SetBackendMenuEntry(entry *BackendMenuEntry)
	SetBackendCmdInterceptor(func(cmds []string) (stop bool))
}

// 与游戏的交互接口，通过发出点什么来影响游戏
// 建议扩展该接口以提供更丰富的功能
// 另一种扩展方式是定义新插件并暴露接口
type GameControl interface {
	SayTo(target string, msg string)
	ActionBarTo(target string, msg string)
	TitleTo(target string, msg string)
	SubTitleTo(target string, msg string)
	SendCmd(cmd string)
	SendCmdAndInvokeOnResponse(string, func(output *packet.CommandOutput))
	SendMCPacket(packet.Packet)
	GetPlayerKit(name string) PlayerKit
	GetPlayerKitByUUID(ud uuid.UUID) PlayerKit
	SetOnParamMsg(string, func(chat *GameChat) (catch bool)) error
}

type PlayerKit interface {
	Say(msg string)
	ActionBar(msg string)
	Title(msg string)
	SubTitle(msg string)
	GetRelatedUQ() *uqHolder.Player

	GetViolatedStorage() map[string]interface{}
	GetPersistStorage(k string) string
	CommitPersistStorageChange(k string, v string)

	SetOnParamMsg(func(chat *GameChat) (catch bool)) error
	GetOnParamMsg() func(chat *GameChat) (catch bool)
}

// 与游戏的交互接口，如何捕获和处理游戏的数据包和消息
type GameListener interface {
	SetOnAnyPacketCallBack(func(packet.Packet))
	SetOnTypedPacketCallBack(uint32, func(packet.Packet))
	SetGameMenuEntry(entry *GameMenuEntry)
	SetGameChatInterceptor(func(chat *GameChat) (stop bool))

	AppendOnFirstSeePlayerCallback(cb func(string))
	AppendLoginInfoCallback(cb func(entry protocol.PlayerListEntry))
	AppendLogoutInfoCallback(cb func(entry protocol.PlayerListEntry))
}

// 安全事件发送和处理，比如某插件发现有玩家在恶意修改设置
// 而另一个插件则在 QQ 群里通知这个事件的发生
type SecurityEventIO interface {
	RedAlert(info string)
	RegOnAlertHandler(cb func(info string))
}

type MainFrame interface {
	CtxProvider
	ConfigProvider
	StorageAndLogProvider
	BackendInteract
	SecurityEventIO
	FatalError(err string)
	GetGameControl() GameControl
	GetGameListener() GameListener
}