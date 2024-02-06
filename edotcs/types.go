package edotcs

// 存储者这个插件sdk的基础类型

type Plugin interface {
	Init_plugin() error                                 // 插件初始化
	Player_Join(player string) error                    // 玩家加入了服务器
	Player_Left(Plugin string) error                    // 玩家离开了服务器
	Player_Message(player string, message string) error // 玩家消息
	Menu(player string, word []string) error            // 玩家菜单事件
}

type Player struct {
	Name string // 玩家名称
	UUID string // 玩家UUID
}
