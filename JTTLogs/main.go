package main

import (
	"fmt"
	_ "gitee.com/ictt/JTTM/routers"
	"gitee.com/ictt/JTTM/server/ftpclient"
	"gitee.com/ictt/JTTM/server/tcp"
	"gitee.com/ictt/JTTM/server/udp"
	"gitee.com/ictt/JTTM/server/ws"
	"gitee.com/ictt/JTTM/tools/logs"
	"gitee.com/ictt/JTTM/tools/sqlDB"
	"github.com/astaxie/beego"
	"github.com/kqbi/service"
	"os"
	"time"
)

/*
func main() {
	//beego.Run()
	sqlDB.InitDB()
	//udp.ListenUDPServer()
	//tcp.ListenTCPServer()
	//time.Sleep(10*time.Second)
	go udp.ListenUDPServer()
	ws.ListenWebsocketServer()
	//appPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//appPath = strings.Replace(appPath, "\\", "/", -1)
	//beego.BConfig.WebConfig.ViewsPath = appPath + "/views"
	//beego.SetStaticPath("/static", appPath+"/static")
	//beego.Run()

	time.Sleep(time.Second * 60)
	udp.GetAllParameter("017856938601")
	//udp.SendShootNow("017856938601",1, 1,0)
	time.Sleep(time.Second * 10)
	////udp.MultimediaDataUploadSend("017856938601",1,0,"200811000000","200811124320")
	//code, count, list := udp.VideoListGetSend("admin", "017856938601","200818120000","200819184320",0,3,2,1,0)
	//fmt.Println(code, count)
	//for _,v := range list {
	//	fmt.Println(v)
	//}
	//code, url :=udp.StartReplayVideoSend("admin", "017856938601",11,3,1,0,0,0,"200813080000","200813124320")
	//fmt.Println(code, url)
	//time.Sleep(time.Second * 10)
	//code1 := udp.ControlReplayVideoSend("admin","017856938601",11,2,0,"200813090000" )
	//fmt.Println(code1)
	//udp.QueryStreamProperties("017856938601")
	//time.Sleep(time.Second * 5)

	//udp.SetDeviceParameter("017856938601", 1,"000000280400000002")
	//
	//time.Sleep(time.Second * 5)
	//udp.SetDeviceParameter("017856938601", 1,"00000029040000003c")
	//time.Sleep(time.Second * 20)
	//udp.GetAppointParameter("017856938601", 1,"00000076")
	//time.Sleep(time.Second * 10)
	//udp.TrackLocationInfo("017856938601", 15,200)
	//time.Sleep(time.Second * 5)
	//udp.QueryDeviceProperties("017856938601")
	//time.Sleep(time.Second * 10)
	//udp.QueryLocationInfo("017856938601")
	//time.Sleep(time.Second * 60)
	//udp.TrackLocationInfo("017856938601",0,0)
	//udp.SendAcknowledgeAlarm("017856938601",55585,"00000000010000000000000000000000")
	//udp.SendShootNow("017856938601",1, 1,0)
	//for {
	//code, url := udp.RequestRealStream("017856938601", 1, 1, 0)
	//fmt.Println(code, url)
	//time.Sleep(time.Second * 30)
	//code1 := udp.ControlRealStream("017856938601", 1, 0, 0, 1)
	//fmt.Println(code1)
	time.Sleep(time.Second * 30)
	//}
}
*/
type program struct {
	exit chan bool
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	logs.BeeLogger.Info("JTTM Service Start!!!")
	fmt.Printf("%s JTTM Service Start!!!\n", time.Now().Format("2006-01-02 15:04:05"))
	sqlDB.InitDB()
	ftpclient.InitFtpClient()
	go tcp.ListenTCPServer()
	go udp.ListenUDPServer()
	go ws.ListenWebsocketServer()
	go func() {
		for {
			//keepalivedata := bytes.NewBufferString("keepalive")
			//ftpclient.FtpCli.Lock()
			//err := ftpclient.FtpCli.FTPClient.Stor("keepalive.txt", keepalivedata)
			//ftpclient.FtpCli.Unlock()
			//if err != nil {
			//	logs.BeeLogger.Error("ftp err: %s",err)
			//	ftpclient.InitFtpClient()
			//}

			msg := ftpclient.Message{
				FileName: "keepalive.txt",
				SaveMsg: "keepalive",
			}
			ftpclient.FtpMsg <- msg
			time.Sleep(60*time.Second)
		}
	}()
	go ftpclient.FtpFace()
	//go mqttclient.MQTTInit()
	//beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
	//	AllowAllOrigins:  true,
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//	AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	//	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	//	AllowCredentials: true,
	//}))
	beego.Run()
	return
}

func (p *program) Stop(s service.Service) error {
	ftpclient.FtpClient.FtpConn.Quit()
	logs.BeeLogger.Info("JTTM Service Stop!!!")
	fmt.Printf("%s JTTM Service Stop!!!\n", time.Now().Format("2006-01-02 15:04:05"))
	close(p.exit)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "JTTM", //服务显示名称
		DisplayName: "JTTM", //服务名称
		Description: "JTTM", //服务描述
	}

	prg := &program{
		exit: make(chan bool),
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logs.PanicLogger.Fatalln("service.New() error: ", err)
	}

	if len(os.Args) > 1 {
		//install, uninstall, start, stop 的另一种实现方式
		err = service.Control(s, os.Args[1])
		if err != nil {
			logs.PanicLogger.Fatalln(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logs.PanicLogger.Fatalln(err)
	}
}
