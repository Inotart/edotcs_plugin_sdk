package main

import (
	"fmt"
	"os"

	"github.com/Inotart/edotcs_plugin_sdk/edotcs"
)

type cake struct {
	edotcs.BasePlugin
	EDotCS *edotcs.EDotCS
}

// 用户输入菜单命令后执行,这里的例子是 .eat art
func (bp *cake) Menu(player string, menu []string) error {
	bp.EDotCS.Say_To(player, "你想吃什么大饼？") // 发送消息给玩家
	return nil
}
func main() {
	fmt.Println("基本的饼插件演示")
	sdk := edotcs.EDotCS{
		Name:        "饼插件",
		Version:     [3]int64{0, 0, 1},
		Ip:          os.Args[1],
		Author:      "Inotart",
		Menu_key:    "eat",
		Menu_tip:    "吃大饼",
		Description: "这是一个饼插件，可以吃大饼。",
	}
	sdk.Plugins = func() edotcs.Plugin {
		return &cake{EDotCS: &sdk}
	}()
	sdk.Start()
}
