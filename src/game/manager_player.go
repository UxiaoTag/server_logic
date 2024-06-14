package game

import (
	"fmt"
	"math"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var managePlayer *ManagePlayer

type ManagePlayer struct {
	Players map[int]*Player
	lock    *sync.RWMutex
}

func GetManagePlayer() *ManagePlayer {
	if managePlayer == nil {
		managePlayer = new(ManagePlayer)
		managePlayer.Players = make(map[int]*Player)
		managePlayer.lock = new(sync.RWMutex)
	}
	return managePlayer
}

func (mp *ManagePlayer) PlayerLoginIn(ws *websocket.Conn, userId int) *Player {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	playerinfo, ok := mp.Players[userId]

	if ok {
		if playerinfo.ws != ws {
			oldws := playerinfo.ws
			playerinfo.ws = ws
			if oldws != nil {
				//顶号
				playerinfo.ws.Write([]byte("账号别处登录"))
				playerinfo.ws.Close()
			}

		}
	} else {
		playerinfo = NewTestPlayer(ws, userId)
		mp.Players[userId] = playerinfo
	}
	playerinfo.exitTime = math.MaxInt64
	return playerinfo
}

func (mp *ManagePlayer) PlayersClose(ws *websocket.Conn, userId int) {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	playerinfo, ok := mp.Players[userId]

	if ok {
		//确定同一个连接
		if playerinfo.ws == ws {
			playerinfo.ws = nil
			playerinfo.exitTime = time.Now().Unix() + 10 //这里应该是十分钟好一点
			fmt.Println("websocket断开")
			//断开连接
		}
	}
}

func (mp *ManagePlayer) Run() {

	ticker := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-ticker.C:
			mp.CheckPlayerOff()
		}
	}
}
func (mp *ManagePlayer) CheckPlayerOff() {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	for k, p := range mp.Players {
		if p.exitTime <= time.Now().Unix() {
			delete(mp.Players, k)
		}
	}
}
