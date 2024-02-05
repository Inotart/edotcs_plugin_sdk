package edotcs

type BasePlugin struct {
	EDotCS *EDotCS
}

// Init_plugin 初始化插件,此方法必须重写
func (bp *BasePlugin) Init_plugin() error {
	// 此方法必须重新写
	return nil
}

// Player_join 玩家加入服务器时调用,此方法可重写
func (bp *BasePlugin) Player_Join(player string) error {
	return nil
}

// Player_left 玩家离开服务器时调用,此方法可重写
func (bp *BasePlugin) Player_Left(player string) error {
	return nil
}

// Player_Message  玩家发送消息时调用,此方法可重写
func (bp *BasePlugin) Player_Message(player string, message string) error {
	return nil
}

// Menu 菜单点击时调用,此方法可重写
func (bp *BasePlugin) Menu(player string, menu []string) error {
	return nil
}
