package edotcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Inotart/edotcs_plugin_sdk/drpc"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type EDotCS struct {
	Name        string          // 插件名
	Ip          string          // 插件IP,默认 127.0.0.1:8080
	Author      string          // 插件作者
	Version     [3]int64        // 插件版本
	Menu_key    string          // 菜单键,例如 help
	Menu_tip    string          // 菜单提示,例如 帮助菜单
	Description string          // 插件描述
	client      *http.Client    // 客户端连接
	conn        *websocket.Conn // 客户端连接
	Plugins     Plugin          // 插件实例
	// TODO: add fields here
}

func (edotcs *EDotCS) Plugin_init() error {
	//edotcs.plugins = p
	edotcs.client = &http.Client{Timeout: time.Second * 60}

	return nil

}
func (edotcs *EDotCS) Start() {
	err := edotcs.Plugin_init()
	if err != nil {

		log.Println("plugin init failed")
		panic(err)
		return
	}
	edotcs.Listen()
	// TODO: start plugin here
}
func (edotcs *EDotCS) Listen() {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/dotcs/v8/", edotcs.Ip), nil)
	if err != nil {
		log.Println("connect to edotcs failed")
		panic(err)
	}
	edotcs.conn = conn
	defer conn.Close()
	// 发送插件信息
	plugins2 := drpc.Plugin{
		Name:   edotcs.Name,
		Author: edotcs.Author,
		Version: &drpc.Version{
			Major: edotcs.Version[0],
			Minor: edotcs.Version[1],
			Patch: edotcs.Version[2],
		},
		Recommendation:     &edotcs.Description,
		MenuWord:           &edotcs.Menu_key,
		MenuRecommendation: &edotcs.Menu_tip,
	}
	plugins2_data, err := proto.Marshal(&plugins2)
	if err != nil {
		log.Println("marshal plugins2 failed")
		return
	}
	conn.WriteMessage(websocket.TextMessage, bytes.Join([][]byte{
		{0x00},
		plugins2_data,
	}, []byte{}))
	for {
		// 读取消息
		messageType, p, err := conn.ReadMessage()
		fmt.Println(messageType)
		if err != nil {
			log.Fatal(err)
		}
		if messageType != websocket.BinaryMessage {
			continue
		}
		switch p[0] {
		case 1:
			// 接收到插件信息
			message := drpc.Player_Message{}
			err = proto.Unmarshal(p[1:], &message)
			if err != nil {
				log.Println("unmarshal player_message failed")
				continue
			}
			edotcs.Plugins.Player_Message(message.GetPlayer(), message.GetMessage())
		case 2:
			// 接收到玩家进服
			message := drpc.PlayerJoin{}
			err = proto.Unmarshal(p[1:], &message)
			if err != nil {
				log.Println("unmarshal playerjoin failed")
				continue
			}
			edotcs.Plugins.Player_Join(message.GetPlayer())

		case 3:
			// 接收到玩家退出服
			message := drpc.Player_Left{}
			err = proto.Unmarshal(p[1:], &message)
			if err != nil {
				log.Println("unmarshal playerleave failed")
				continue
			}

			edotcs.Plugins.Player_Left(message.GetPlayer())
		case 4:
			// 接受到菜单命令
			message := drpc.Menu{}

			err = proto.Unmarshal(p[1:], &message)
			if err != nil {
				log.Println("unmarshal menu failed")
				continue
			}
			edotcs.Plugins.Menu(message.GetPlayer(), message.GetWord())
		case 5:
			// 接收到私聊消息
			message := drpc.Player_Whisper{}
			err = proto.Unmarshal(p[1:], &message)
			if err != nil {
				log.Println("unmarshal private_message failed")
				continue
			}
			edotcs.Plugins.Player_Whisper(message.GetPlayer(), message.GetMessage())
		case 6:
			// 接收到BlockActorData
			message := drpc.BlockActorData{}
			err = proto.Unmarshal(p[1:], &message)
			if err != nil {
				log.Println("unmarshal blockactordata failed")
				continue
			}
			NBT_Data := map[string]any{}
			err := json.Unmarshal(message.GetNBT_Data(), &NBT_Data)
			if err != nil {
				log.Println("unmarshal nbt_data failed")
				continue
			}
			Position := BlockPos{message.GetPosition().GetX(), message.GetPosition().GetY(), message.GetPosition().GetZ()}
			edotcs.Plugins.BlockActorData(Position, NBT_Data)
		case 7:
			// 接收到系统消息
			message := drpc.System_Message{}
			err = proto.Unmarshal(p[1:], &message)
			if err != nil {
				log.Println("unmarshal system_message failed")
				continue
			}
			edotcs.Plugins.System_Message(message.GetNeedsTranslation(), message.GetSourceName(), message.GetMessage(), message.GetParameters(), message.GetXUID(), message.GetPlatformChatID(), message.GetPlayerRuntimeID())
		default:
			log.Println("unknown message type", p[0])
		}
	}
}
func (edotcs *EDotCS) SendCmd(cmd string) {
	sendcmd := drpc.SendCmd{
		Command: cmd,
	}

	data, err := proto.Marshal(&sendcmd)
	if err != nil {
		log.Println("marshal sendcmd failed")
		return
	}
	edotcs.client.Post(fmt.Sprintf("http://%s/dotcs/v8/cmd", edotcs.Ip), "application/octet-stream", bytes.NewReader(data))
}
func (edotcs *EDotCS) SendWSCmd(cmd string) {
	sendcmd := drpc.SendCmd{
		Command: cmd,
	}

	data, err := proto.Marshal(&sendcmd)
	if err != nil {
		log.Println("marshal sendcmd failed")
		return
	}
	edotcs.client.Post(fmt.Sprintf("http://%s/dotcs/v8/wscmd", edotcs.Ip), "application/octet-stream", bytes.NewReader(data))
}
func (edotcs *EDotCS) Say_To(player string, Message string) {
	sendcmd := drpc.Say_To{
		Player:  player,
		Message: Message,
	}

	data, err := proto.Marshal(&sendcmd)
	if err != nil {
		log.Println("marshal sendcmd failed")
		return
	}
	edotcs.client.Post(fmt.Sprintf("http://%s/dotcs/v8/sayto", edotcs.Ip), "application/octet-stream", bytes.NewReader(data))
}
func (edotcs *EDotCS) Get_Player(player string) (Player, error) {
	Get_Player := drpc.Get_Player_Info{
		Player: player,
	}

	data, err := proto.Marshal(&Get_Player)
	if err != nil {
		log.Println("marshal sendcmd failed")
		return Player{}, err
	}
	resp, err := edotcs.client.Post(fmt.Sprintf("http://%s/dotcs/v8/getplayer", edotcs.Ip), "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return Player{}, err
	}
	resp_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Player{}, err
	}
	player_info := drpc.Return_Player_Info{}
	err = proto.Unmarshal(resp_data, &player_info)
	if err != nil {
		return Player{}, err
	}
	if !player_info.Is_True {
		return Player{}, fmt.Errorf(player_info.GetMessage())
	}
	resp.Body.Close()
	return Player{
		Name: player_info.Player.GetPlayer(),
		UUID: player_info.Player.GetUUID()}, nil
}
func (edotcs *EDotCS) Get_Player_List() ([]Player, error) {
	Get_Player := drpc.Get_Online_Player_Info{}

	data, err := proto.Marshal(&Get_Player)
	if err != nil {
		log.Println("marshal sendcmd failed")
		return []Player{}, err
	}
	resp, err := edotcs.client.Post(fmt.Sprintf("http://%s/dotcs/v8/getplayerlist", edotcs.Ip), "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return []Player{}, err
	}
	resp_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Player{}, err
	}
	player_info := drpc.Return_Online_Player_Info{}
	err = proto.Unmarshal(resp_data, &player_info)
	if err != nil {
		return []Player{}, err
	}
	if !player_info.Is_True {
		return []Player{}, fmt.Errorf(player_info.GetMessage())
	}
	resp.Body.Close()
	players := make([]Player, len(player_info.Players))
	for _, info := range player_info.Players {
		players = append(players, Player{
			Name: info.GetPlayer(),
			UUID: info.GetUUID(),
		})

	}
	return players, nil
}
func NewEDotCS(Name string, Ip string, Author string, Version [3]int64, Menu_key string, Menu_tip string, Description string) *EDotCS {
	edotcs := EDotCS{
		Name:        Name,
		Ip:          Ip,
		Author:      Author,
		Version:     Version,
		Menu_key:    Menu_key,
		Menu_tip:    Menu_tip,
		Description: Description,
	}
	return &edotcs
}

// TODO: add methods here
