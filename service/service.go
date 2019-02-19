package service

import (
	"fmt"
	"log"
	"mg-coms/coms"
	"mg-coms/config"
	"mg-coms/notice"
	"mg-coms/websocketio"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

//Service 人机交互服务
type Service struct {
	coms   *coms.Coms
	config config.Config
	state  bool
}

//ServiceBean 返回结构体
type ServiceBean struct {
	Error   bool        `json:"error"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

//Start　启动服务
func (s *Service) Start() {
	// s.ServiceLoad()
	r := mux.NewRouter()

	r.HandleFunc("/message", s.SendAllMessage).Methods("GET")
	r.HandleFunc("/message/{filename}", s.SendMessage).Methods("GET")

	r.HandleFunc("/restart", s.ServicRestart).Methods("GET")

	r.HandleFunc("/process", s.SearchProcess).Methods("GET")

	r.HandleFunc("/users", s.ModifiUser).Methods("PUT")
	r.HandleFunc("/users", s.AppendUser).Methods("POST")
	r.HandleFunc("/users", s.SearchAllUsers).Methods("GET")
	r.HandleFunc("/users/{name}", s.SearchUser).Methods("GET")
	r.HandleFunc("/users/{name}/{password}", s.LoginUser).Methods("GET")
	r.HandleFunc("/users/{name}", s.RemoveUser).Methods("DELETE")
	r.HandleFunc("/coms/commands", s.Command).Methods("PUT")
	r.HandleFunc("/commands", s.SearchAllPotocol).Methods("GET")
	r.HandleFunc("/coms/channels/{name}", s.Channel).Methods("GET")
	r.HandleFunc("/coms/devices", s.AllDevices).Methods("GET")
	r.HandleFunc("/coms/devices/{channelname}/{deviceid}", s.OneDevice).Methods("GET")

	r.HandleFunc("/version", s.CreateVersion).Methods("POST")
	r.HandleFunc("/version", s.SearchVersion).Methods("GET")
	r.HandleFunc("/version", s.ModifiVersion).Methods("PUT")
	r.HandleFunc("/version/{name}", s.RemoveVersion).Methods("DELETE")

	r.HandleFunc("/templates", s.ModifiTemplate).Methods("PUT")
	r.HandleFunc("/templates", s.AppendTemplate).Methods("POST")
	r.HandleFunc("/templates", s.SearchAllTemplate).Methods("GET")
	r.HandleFunc("/templates/{templatename}", s.RemoveTemplate).Methods("DELETE")
	r.HandleFunc("/templates/{templatename}", s.SearchOneTemplate).Methods("GET")

	r.HandleFunc("/pointaddr", s.SearchPointAddr).Methods("GET")
	r.HandleFunc("/pointaddr/{cname}", s.SearchNewPointAddr).Methods("GET")
	r.HandleFunc("/pointaddr", s.ModifiPointAddr).Methods("POST")
	r.HandleFunc("/pointaddr", s.DeletePointAddr).Methods("PUT")

	r.HandleFunc("/channels", s.ModifiChannel).Methods("PUT")
	r.HandleFunc("/channels", s.AppendChannel).Methods("POST")
	r.HandleFunc("/channels", s.SearchAllChannel).Methods("GET")
	r.HandleFunc("/channels/{channelname}", s.SearchOneChannel).Methods("GET")
	r.HandleFunc("/channels/{channelname}", s.RemoveChannel).Methods("DELETE")

	r.HandleFunc("/devices", s.ModifiDevice).Methods("PUT")
	r.HandleFunc("/devices", s.AppendDevice).Methods("POST")
	r.HandleFunc("/devices/{channelname}", s.SearchAllDevice).Methods("GET")
	r.HandleFunc("/devices/{channelname}/{deviceid}", s.SearchOneDevice).Methods("GET")
	r.HandleFunc("/devices/{channelname}/{deviceid}", s.RemoveDevice).Methods("DELETE")

	r.HandleFunc("/points", s.ConfigPoint).Methods("PUT", "POST")
	r.HandleFunc("/points/{channelname}/{deviceid}", s.SearchPoint).Methods("GET")
	r.HandleFunc("/points/{channelname}/{deviceid}/{pointid}/{kind}", s.RemovePoint).Methods("DELETE")
	// static := http.FileServer(http.Dir("/home/jiaweisi/Documents/coms/dist/"))
	static := http.FileServer(http.Dir("../dist/"))
	// http.Handle("/", http.StripPrefix("/", static))
	r.PathPrefix("/").Handler(http.StripPrefix("/", static))
	http.Handle("/", r)
	http.Handle("/ws", websocket.Handler(websocketio.WsHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAddServer: ", err)
	}
}

//New 创建服务
func New() *Service {
	return new(Service).Init()
}

//InitComs 初始化COMS
func (s *Service) InitComs() {
	coms, err := coms.New()

	if err != nil {
		fmt.Println("启动失败", err)
		s.state = false
	}
	s.coms = coms
	s.state = true
}

//ServiceLoad 加载服务配置
func (s *Service) ServiceLoad() {
	notice.SetFileCallBack(websocketio.NoticeData)
	s.InitComs()
	// fmt.Println("DeviceMap=======", s.coms.DeviceMap)
	s.coms.Start()
	s.config.Load()
	websocketio.ChannelMap = s.coms.ChannelMap
	websocketio.Server()
}

//ServicRestart 重启通讯管理机
func (s *Service) ServicRestart(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s.Restart()
}

func (s *Service) Restart() {
	s.coms.Stop()
	s.config.Empty()
	time.Sleep(time.Second * 2)
	s.ServiceLoad()
}

//Init 初始化Service
func (s *Service) Init() *Service {
	s.ServiceLoad()

	// s.Start()
	return s
}
