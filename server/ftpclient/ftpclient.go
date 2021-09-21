package ftpclient

import (
	"bytes"
	"fmt"
	"gitee.com/ictt/JTTM/tools/logs"
	"github.com/astaxie/beego"
	"github.com/jlaffaye/ftp"
	"time"
)
var (
	//FTPClient *ftp.ServerConn
	//err error
	FtpMsg  = make(chan Message, 1000)
	FtpClient *FTPClient
)



type FTPClient struct {
	FtpConn	*ftp.ServerConn
}
type Message struct {
	FileName	string
	SaveMsg		string
}

//var FtpCli = struct{
//	sync.Mutex
//	FTPClient *ftp.ServerConn
//}{}

func InitFtpClient() {
	FtpUser := beego.AppConfig.String( "ftp::ftpUser")
	FtpPWD := beego.AppConfig.String( "ftp::ftpPWD")
	FtpAddr := beego.AppConfig.String("ftp::ftpAddr")
	if FtpAddr == "" || FtpPWD == "" || FtpUser == "" {
		logs.PanicLogger.Panicln("init ftp error, ftpUser or ftpPWD or ftpAddr cannot be empty")
	}
	//fmt.Println(FtpPWD,FtpUser,FtpAddr)

	FtpConn, err := ftp.Dial(FtpAddr, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		logs.BeeLogger.Error("init ftp connect error: %s",err)
	}

	err = FtpConn.Login(FtpUser, FtpPWD)
	if err != nil {
		logs.BeeLogger.Error("init ftp login error: %s",err)
	}
	logs.BeeLogger.Info(fmt.Sprintf("successful connection to ftp"))
	fmt.Printf("%s successful connection to ftp\n", time.Now().Format("2006-01-02 15:04:05"))

	FtpClient = &FTPClient{
		FtpConn: FtpConn,
	}
	go FtpClient.Store()
}

func (ft *FTPClient) Store()  {
	for {
		msg, ok := <-FtpMsg
		if ok {
			data := bytes.NewBufferString(msg.SaveMsg)
			//fmt.Println(1)
			err := ft.FtpConn.Stor(msg.FileName, data)
			//fmt.Println(2)
			if err != nil {
				//fmt.Println("cuo wu le")
				logs.BeeLogger.Error("save message %s to ftp error : %s", msg.FileName, err)
				FtpMsg <- msg
				go InitFtpClient()
				break
			}
		}
	}
}
