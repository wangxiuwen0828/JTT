package udp

import (
	"encoding/hex"
	"fmt"
	"gitee.com/ictt/JTTM/config"
	"gitee.com/ictt/JTTM/tools"
	"gitee.com/ictt/JTTM/tools/logs"
	"net"
	"time"
)

func ListenUDPServer() {

	udpAddr, err := net.ResolveUDPAddr("udp", config.UDPAddr)
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("net.ResolveUDPAddr() error: %s", err))
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("net.ListenUDP() error: %s", err))
	}

	logs.BeeLogger.Info("start UDPServer successful!")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "start UDPServer successful!")
	defer udpConn.Close()

	for {
		buf := make([]byte, 4096)
		n, udpAddr, err := udpConn.ReadFromUDP(buf)

		if err != nil {
			logs.BeeLogger.Error("udpConn.ReadFromUDP() error: %s", err)
			continue
		}

		logs.BeeLogger.Debug("UDP原始数据(udp raw data)：%s", hex.EncodeToString(buf[:n]))

		go UdpHandle(buf[:n], udpConn, udpAddr)
	}
}

func UdpHandle(bufData []byte, udpConn *net.UDPConn, udpAddr *net.UDPAddr) {
	fmt.Println(hex.EncodeToString(bufData))
	messByte, err := tools.Changebodymess(bufData)
	if err != "" {
		return
	}
	mess16ID := hex.EncodeToString(messByte[:2])
	switch mess16ID {
	case "0001":
		//处理收到终端发送的通用回复
		ResNormalMes(messByte)
	case "0002":
		//收到保活请求并回复
		ResKeepAlive(messByte, udpConn, udpAddr)
	case "0100":
		//收到注册请求并处理
		ResRegisterMes(messByte, udpConn, udpAddr)
	case "0102":
		//收到鉴权请求
		ResPowerIdentify(messByte, udpConn, udpAddr)
	case "0104":
		//收到终端参数
		HandleAllParameter(messByte)
	case "0107":
		HandleDeviceProperties(messByte)
	case "0200":
		//收到位置信息并回复
		ResLocateMes(messByte, udpConn, udpAddr)
	case "0201":
		//收到位置信息并解析
		HandleLocationInfo(messByte)
	case "0500":
		//收到车辆控制应答并解析
		HandleVehicleControlReceive(messByte)
	case "0801":
		//收到图片信息并解析
		MultimediaDataUploadHandle(messByte, udpConn, udpAddr)
	case "0b01":
		//收到运营登记信息并回复
		ResOperationRegistrationMes(messByte, udpConn, udpAddr)
	//case "0900":
	//	ResDataUplinkPassThrough(messByte, udpConn, udpAddr)
	case "1003":
		//收到终端上传音视频属性
		HandleStreamProperties(messByte)
	case "1205":
		//收到录像列表
		VideoListGetHandle(messByte)

	}

}
