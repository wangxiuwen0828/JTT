package config

import (
	"gitee.com/ictt/JTTM/tools/logs"
	"github.com/astaxie/beego"
	"strconv"
)

var (
	IP      		string
	UDPPort 		int64
	UDPAddr 		string //udp服务地址
	TCPPort 		int64
	TCPAddr			string //TCP服务地址
	WebsocketPort 	int64	//websocket服务端口
	WebsocketAddr   string //websocket服务地址
	WSKeepaliveTime int64  //websocket超时时间
	StreamIP        string	//实时流输送到的IP
	StreamPort      int64	//实时流输送到的端口
	//StreamAddr      string
	UrlAddr			string	//请求实时视频url的地址
	MQTTAddr		string	//mqtt的请求地址
	MQTTID			string	//mqtt的id
	//WSKeepaliveOutTime	int64
	FacePath 		string	//人脸图片的保存位置
	FaceWriteTime	int64	//人脸图片写入时间
)

const (
	JTSRegisterReq       = "JTS_Register_req"       //注册请求
	JTSRegisterRes       = "JTS_Register_Res"       //注册回复
	JTSKeepaliveReq      = "JTS_Keepalive_Req"      //心跳保活请求
	JTSKeepaliveRes      = "JTS_Keepalive_Res"      //心跳保活回复
	JTSGetDeviceListReq  = "JTS_GetDeviceList_Req"  //获取设备列表请求
	JTSGetDeviceListRes  = "JTS_GetDeviceList_Res"  //获取设备列表回复
	JTSGetChannelListReq = "JTS_GetChannelList_Req" //获取设备通道列表请求
	JTSGetChannelListRes = "JTS_GetChannelList_Res" //获取设备通道列表回复
	JTSGetGPSInfoReq     = "JTS_GetGPSList_Req"     //获取GPS信息请求
	JTSGetGPSInfoRes     = "JTS_GetGPSList_Res"     //获取GPS信息回复
	JTSStreamActionReq   = "JTS_StreamAction_Req"   //实时视频操作请求
	JTSStreamActionRes   = "JTS_streamAction_Res"   //实时视频操作回复
	JTSReplayActionReq   = "JTS_ReplayAction_Req"   //录像回放操作请求
	JTSReplayActionRes   = "JTS_ReplayAction_Res"   //录像回放操作回复
	JTSVideoListReq      = "JTS_VideoList_Req"      //获取录像列表请求
	JTSVideoListRes      = "JTS_VideoList_Res"      //获取录像列表回复
	JTSAlarmListReq      = "JTS_AlarmList_Req"      //获取报警信息列表请求
	JTSAlarmListRes      = "JTS_AlarmList_Res"      //获取报警信息列表回复
	JTSRealGPSReq        = "JTS_RealGPS_Req"        //实时GPS请求
	JTSRealGPSRes        = "JTS_RealGPS_Res"        //实时GPS回复
	JTSRealAlarmReq      = "JTS_RealAlarm_Req"      //实时报警请求
	JTSRealAlarmRes      = "JTS_RealAlarm_Res"      //实时报警回复
	JTSDriverActionReq   = "JTS_DriverAction_Req"   //驾驶行为列表请求
	JTSDriverActionRes   = "JTS_DriverAction_Res"   //驾驶行为列表回复
	JTSSetParameterReq	 = "JTS_SetParameter_Req"	//设置位置汇报时间间隔请求
	JTSSetParameterRes	 = "JTS_SetParameter_Res"	//设置位置汇报时间间隔回复
	)

func init() {
	//JTSCodeMap = map[int64]string{
	//	JTTM_ERROR_SUCCESS_OK:       "Success OK",
	//	JTTM_ERROR_FAILED:  "Parameter Error",
	//	JTTM_ERROR_SERVER_ERROR: "Device Not Online",
	//}
	//WebsocketAddr = beego.AppConfig.String("websocketAddr")
	//if WebsocketAddr == "" {
	//	logs.PanicLogger.Panicln("init websocketAddr error, websocketAddr cannot be empty")
	//}

	var err error
	WSKeepaliveTime, err = beego.AppConfig.Int64("wsKeepaliveTime")
	if err != nil || WSKeepaliveTime < 0 {
		WSKeepaliveTime = 60
		logs.BeeLogger.Error("init wsKeepaliveTime error, set default value:60, error:%s", err)
	}
	IP = beego.AppConfig.String("ip")
	UDPPort, _ = beego.AppConfig.Int64("udpPort")
	UDPAddr = IP + ":" + strconv.FormatInt(UDPPort, 10)
	if IP == "" || UDPPort == 0 {
		logs.PanicLogger.Panicln("init UDPAddr error, UDP IP or Port cannot be empty")
	}

	TCPPort, _ = beego.AppConfig.Int64("tcpPort")
	TCPAddr = IP + ":" + strconv.FormatInt(TCPPort, 10)
	if TCPPort == 0 {
		logs.PanicLogger.Panicln("init TCPAddr error, TCPPort Port cannot be empty")
	}

	WebsocketPort, _ = beego.AppConfig.Int64("websocketPort")
	WebsocketAddr = IP + ":" + strconv.FormatInt(WebsocketPort, 10)
	if WebsocketPort == 0 {
		logs.PanicLogger.Panicln("init WebsocketPort error, WebsocketPort cannot be empty")
	}

	StreamIP = beego.AppConfig.String("streamIP")
	StreamPort, _ = beego.AppConfig.Int64("streamPort")
	//StreamAddr = IP + ":" + strconv.FormatInt(StreamPort, 10)
	if StreamIP == "" || StreamPort == 0 {
		logs.PanicLogger.Panicln("init StreamAddr error, StreamIP or StreamPort cannot be empty")
	}

	UrlAddr = beego.AppConfig.String("urlAddr")
	if UrlAddr == "" {
		logs.PanicLogger.Panicln("init urlAddr error, urlAddr cannot be empty")
	}

	MQTTAddr = beego.AppConfig.String("mqttAddr")
	MQTTID = beego.AppConfig.String("mqttID")
	if MQTTAddr == "" {
		logs.PanicLogger.Panicln("init mqttAddr error, mqttAddr cannot be empty")
	}

	FacePath = beego.AppConfig.String("facePath")
	if FacePath == "" {
		logs.PanicLogger.Panicln("init facePath error, facePath cannot be empty")
	}
	FaceWriteTime,_ = beego.AppConfig.Int64("faceWriteTime")
	if FaceWriteTime == 0 {
		logs.PanicLogger.Panicln("init FaceWriteTime error, FaceWriteTime cannot be empty")
	}
	//fmt.Println(FacePath)
}

//func init() {
//
//	WSKeepaliveTime,_ = beego.AppConfig.Int64("WSKeepaliveTime")
//	if WSKeepaliveTime == 0 {
//		logs.PanicLogger.Panicln("init WSKeepaliveTime error, WSKeepaliveTime cannot be empty")
//	}
//	WSKeepaliveCount,_ := beego.AppConfig.Int64("wsKeepaliveCount")
//	if WSKeepaliveTime == 0 {
//		logs.PanicLogger.Panicln("init WSKeepaliveCount error, WSKeepaliveCount cannot be empty")
//	}
//	WSKeepaliveOutTime = WSKeepaliveTime * WSKeepaliveCount
//}
