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

func (bp *cake) Player_Message(player string, message string) error {
	fmt.Println("玩家", player, "输入指令", message)
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
