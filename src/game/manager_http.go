package game

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var manageHttp *ManageHttp

type ManageHttp struct {
}

func GetManageHttp() *ManageHttp {
	if manageHttp == nil {
		manageHttp = new(ManageHttp)
	}
	return manageHttp
}

func (mh *ManageHttp) InitData() {

	http.Handle("/", websocket.Handler(mh.WebsocketHandler))

	http.HandleFunc("/correctname", mh.CorrectName)

}

func (mh *ManageHttp) CorrectName(w http.ResponseWriter, r *http.Request) {
	player.ModPlayer.Name = "修改名称"
}

func (mh *ManageHttp) WebsocketHandler(ws *websocket.Conn) {
	defer ws.Close()

	var player *Player

	for {
		var msg []byte

		//这里是设置websocket3s超时
		ws.SetReadDeadline(time.Now().Add(3 * time.Second))
		err := websocket.Message.Receive(ws, &msg)
		fmt.Println(err)
		if err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				continue
			}
			if player != nil {
				//关闭websocket，存档TODO
				GetManagePlayer().PlayersClose(ws, player.ModPlayer.UserId)
			}
			break
		}

		fmt.Println(string(msg))

		if player == nil {
			var loginMsg MsgLogin
			msgerr := json.Unmarshal(msg, &loginMsg)

			if msgerr != nil {
				//login登录验证TODO
				player = GetManagePlayer().PlayerLoginIn(ws, int(loginMsg.UserId))
				go player.LogicRun()
			}
		} else {
			player.SendLogic(msg)
		}
	}
	return
}
