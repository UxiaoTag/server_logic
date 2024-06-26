package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"server/csvs"
	"sync"
	"syscall"
	"time"
)

type DBConfig struct {
	DBUser     string `json:"dbuser" `
	DBPassword string `json:"dbpassword" `
}

type ServerConfig struct {
	ServerId      int       `json:"serverid" `
	Host          string    `json:"host" `
	LocalSavePath string    `json:"localsavepath"` //! 本地存储路径
	DBConfig      *DBConfig `json:"database" `
}

type Server struct {
	Wait        sync.WaitGroup
	BanWordBase []string //配置生成
	Lock        *sync.RWMutex

	Config *ServerConfig
}

var server *Server

func GetServer() *Server {
	if server == nil {
		server = new(Server)
		server.Lock = new(sync.RWMutex)
	}
	return server
}

func (self *Server) Start() {
	//读取全局配置
	self.LoadConfig()
	// 加载配置
	rand.Seed(time.Now().Unix())
	csvs.CheckLoadCsv()
	go GetManageBanWord().Run()
	go GetManageHttp().InitData()
	go GetManagePlayer().Run()

	//fmt.Printf("数据测试----start\n")
	playerTest := NewTestPlayer(nil, 10086)
	go playerTest.Run()
	go self.SignalHandle()

	err := http.ListenAndServe(":6666", nil)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}

	//加这条是因为会出现 如果go创建线程慢了导致self.Wait.Wait()判断计时器空直接退出
	time.Sleep(1 * time.Second)

	self.Wait.Wait()
	fmt.Println("服务器关闭成功!")
}

func (self *Server) Close() {
	GetManageBanWord().Close()
}

func (self *Server) AddGo() {
	self.Wait.Add(1)
}

func (self *Server) GoDone() {
	self.Wait.Done()
}

func (self *Server) IsBanWord(txt string) bool {
	self.Lock.RLock()
	defer self.Lock.RUnlock()
	for _, v := range self.BanWordBase {
		match, _ := regexp.MatchString(v, txt)
		if match {
			fmt.Println("发现违禁词:", v)
		}
		if match {
			return match
		}
	}
	return false
}

func (self *Server) UpdateBanWord(banWord []string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	self.BanWordBase = banWord
}

func (self *Server) SignalHandle() {
	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, syscall.SIGINT)

	for {
		select {
		case <-channelSignal:
			fmt.Println("get syscall.SIGINT")
			fmt.Println("get syscall.SIGINT")
			fmt.Println("get syscall.SIGINT")
			fmt.Println("get syscall.SIGINT")
			fmt.Println("get syscall.SIGINT")
			fmt.Println("get syscall.SIGINT")
			fmt.Println("get syscall.SIGINT")
			fmt.Println("get syscall.SIGINT")
			fmt.Println("get syscall.SIGINT")
			self.Close()
		}
	}
}

func (self *Server) LoadConfig() {
	configFile, err := os.ReadFile("../config/config.json")
	if err != nil {
		fmt.Println("error")
		return
	}
	err = json.Unmarshal(configFile, &self.Config)
	if err != nil {
		fmt.Println("error")
		return
	}
	return
}
