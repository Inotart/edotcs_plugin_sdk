package edotcs

// 存储者这个插件sdk的基础类型

type Plugin interface {
	Init_plugin() error                                             // 插件初始化
	Player_Join(player string) error                                // 玩家加入了服务器
	Player_Left(Plugin string) error                                // 玩家离开了服务器
	Player_Message(player string, message string) error             // 玩家消息
	Menu(player string, word []string) error                        // 玩家菜单事件
	Player_Whisper(player string, message string) error             // 玩家私聊
	BlockActorData(Position BlockPos, NBTData map[string]any) error // NBT方块数据更新
	System_Message(
		NeedsTranslation bool,
		SourceName string,
		Message string,
		Parameters []string,
		XUID string,
		PlatformChatID string,
		PlayerRuntimeID string) error // 系统消息
}

type Player struct {
	Name string // 玩家名称
	UUID string // 玩家UUID
}
type BlockPos [3]int32

// X returns the X coordinate of the block position. It is equivalent to BlockPos[0].
func (pos BlockPos) X() int32 {
	return pos[0]
}

// Y returns the Y coordinate of the block position. It is equivalent to BlockPos[1].
func (pos BlockPos) Y() int32 {
	return pos[1]
}

// Z returns the Z coordinate of the block position. It is equivalent to BlockPos[2].
func (pos BlockPos) Z() int32 {
	return pos[2]
}
