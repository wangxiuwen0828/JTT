package ws

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitee.com/ictt/JTTM/config"
	"gitee.com/ictt/JTTM/server"
	"gitee.com/ictt/JTTM/server/sertools"
	"gitee.com/ictt/JTTM/server/tcp"
	"gitee.com/ictt/JTTM/server/udp"
	"gitee.com/ictt/JTTM/tools"
	"gitee.com/ictt/JTTM/tools/logs"
	"gitee.com/ictt/JTTM/tools/sqlDB"
	"time"
)

//解析基本协议
type BaseMsg struct {
	Data     baseData `json:"data"`
	Sequence uint64   `json:"sequence"`
	Protocol string   `json:"protocol"`
}

type baseData struct {
	Sign string `json:"sign"` //注册时生成，签名
}

//注册请求
type RegisterReq struct {
	Data     registerData `json:"data"`
	Protocol string       `json:"protocol"`
	Sequence uint64       `json:"sequence"`
}

type registerData struct {
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
	Tradeno  string `json:"tradeno"`  //流水号
}

//获取设备列表请求
type GetDeviceListReq struct {
	Data     getDeviceListData `json:"data"`
	Protocol string            `json:"protocol"`
	Sequence uint64            `json:"sequence"`
}

type getDeviceListData struct {
	Sign  string `json:"sign"`
	Page  int  	 `json:"page"`  //页码编号，从1开始
	Limit int  	 `json:"limit"` //分页大小
}

//获取设备通道列表
type GetChannelListReq struct {
	Data     getChannelListData `json:"data"`
	Protocol string             `json:"protocol"`
	Sequence uint64             `json:"sequence"`
}

type getChannelListData struct {
	Sign  string `json:"sign"`
	Page  int  `json:"page"`  //页码编号，从1开始
	Limit int  `json:"limit"` //分页大小
}

//获取GPS信息请求
type GetGPSInfoReq struct {
	Data     getGPSInfoData `json:"data"`
	Protocol string         `json:"protocol"`
	Sequence uint64         `json:"sequence"`
}

type getGPSInfoData struct {
	Sign      string `json:"sign"`
	PhoneNum  string `json:"phoneNum"`
	StartTime string `json:"startTime"` //开始时间
	EndTime   string `json:"endTime"`   //停止时间
}

//实时视频控制请求
type StreamStartReq struct {
	Data     StreamStartData `json:"data"`
	Protocol string          `json:"protocol"`
	Sequence uint64          `json:"sequence"`
}
type StreamStartData struct {
	Sign      string `json:"sign"`
	PhoneNum  string `json:"phoneNum"`  //设备号码
	ChannelID byte   `json:"channelID"` //通道id
	Action    int8   `json:"action"`    //指令行为，1开启，2关闭
	//DataType		byte `json:"dataType"`
	Bitstream byte `json:"bitstream"` //0-主码流， 1-子码流
}

//录像回放请求
type VideoReplayReq struct {
	Data     VideoReplayReqData `json:"data"`
	Protocol string             `json:"protocol"`
	Sequence uint64             `json:"sequence"`
}

type VideoReplayReqData struct {
	Sign      string `json:"sign"`
	PhoneNum  string `json:"phoneNum"`  //设备号码
	ChannelID byte   `json:"channelID"` //通道id
	Action    int8   `json:"action"`    //指令行为，1开启，2关闭
	StartTime string `json:"startTime"` //录像播放开始时间
	EndTime   string `json:"endTime"`   //录像播放停止时间
}

//录像列表请求
type VideoListReq struct {
	Data     VideoListReqData `json:"data"`
	Protocol string           `json:"protocol"`
	Sequence uint64           `json:"sequence"`
}

type VideoListReqData struct {
	Sign      string `json:"sign"`
	PhoneNum  string `json:"phoneNum"`  //设备号码
	ChannelID byte   `json:"channelID"` //通道id
	StartTime string `json:"startTime"` //录像开始时间
	EndTime   string `json:"endTime"`   //录像停止时间
}

//实时gps请求
type RealGPSReq struct {
	Data     RealGPSReqData `json:"data"`
	Protocol string         `json:"protocol"`
	Sequence uint64         `json:"sequence"`
}

type RealGPSReqData struct {
	Sign     string `json:"sign"`
	PhoneNum string `json:"phoneNum"` //设备号码
}

//设置位置汇报时间间隔
type reportParameter struct {
	Protocol string              `json:"protocol"`
	Sequence uint64              `json:"sequence"`
	Data     reportParameterData `json:"data"`
}

type reportParameterData struct {
	Sign     string `json:"sign"`
	PhoneNum string `json:"phoneNum"` //设备号码
	Time 	 uint32	`json:"time"`
}

//处理接收到的websocket客户端数据
func wsReadHandle(msg *wsMessage, wsConn *wsConnection) {
	var baseMsg BaseMsg
	if err := json.Unmarshal(msg.data, &baseMsg); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), err)
		return
	}
	//忽略心跳保活请求打印
	if baseMsg.Protocol != config.JTSKeepaliveReq {
		logs.BeeLogger.Info("remoteAddr=%v, websocketType：%d, websocketData：%s", wsConn.wsSocket.RemoteAddr(), msg.messageType, string(msg.data))
		fmt.Println(string(msg.data))
	}

	if baseMsg.Protocol != config.JTSRegisterReq {
		//如果不是注册请求，需判断所传的sign是否合法，不合法则断开连接
		if client := wsConnAll.get(baseMsg.Data.Sign); client == nil {
			//无效的sign，断开连接
			wsConn.close()
			return
		}
	}

	switch baseMsg.Protocol {
	case config.JTSKeepaliveReq:
		//心跳保活请求
		replyMap := getReplyMap(config.JTSKeepaliveRes, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
	case config.JTSRegisterReq:
		//注册请求
		getRegisterHandle(msg, wsConn, baseMsg)

	case config.JTSGetDeviceListReq:
		//获取设备列表请求
		getDeviceList(msg, wsConn, baseMsg)

	case config.JTSGetChannelListReq:
		//获取设备通道列表请求
		var getChannelListReq GetChannelListReq
		if err := json.Unmarshal(msg.data, &getChannelListReq); err != nil {
			logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
			replyMap := getReplyMap(config.JTSGetChannelListRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
			wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
			return
		}

		if getChannelListReq.Data.Page < 0 || getChannelListReq.Data.Limit < 0 {
			logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, get channelList parameter error", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol)
			replyMap := getReplyMap(config.JTSGetChannelListRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
			wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
			return
		}

		replyMap := getReplyMap(config.JTSGetChannelListRes, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK,baseMsg.Sequence)

		totalCount, _ := sqlDB.TotalCount(&sqlDB.GetChannelInfo{})
		dataMap := replyMap["data"].(map[string]interface{})
		dataMap["totalCount"] = totalCount
		if totalCount == 0 {
			dataMap["channelList"] = make([]interface{}, 0)
		} else {
			var channelList []sqlDB.GetChannelInfo
			if getChannelListReq.Data.Page == 0 && getChannelListReq.Data.Limit == 0 {
				//Page和Limit参数同时为0则查询数据库指定表格的所有数据
				sqlDB.Find(&channelList, &sqlDB.GetChannelInfo{})
			} else {
				start := (getChannelListReq.Data.Page - 1) * getChannelListReq.Data.Limit
				sqlDB.Limit(&channelList, start, getChannelListReq.Data.Limit)
			}

			dataMap["channelList"] = channelList
		}

		replyMap["data"] = dataMap
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)

	case config.JTSGetGPSInfoReq:
		//获取GPS轨迹信息请求
		getGPSList(msg, wsConn, baseMsg)

	case config.JTSStreamActionReq:
		//实时视频控制请求
		go func() {
			getStreamStartOrStop(msg, wsConn, baseMsg)
		}()

	case config.JTSReplayActionReq:
		//录像回放控制请求
		go func() {
			getReplayStartOrStop(msg, wsConn, baseMsg)
		}()

	case config.JTSVideoListReq:
		go func() {
			getVideoList(msg, wsConn, baseMsg)
		}()

	case config.JTSRealGPSReq:
		//获取实时GPS
		getRealGPS(msg, wsConn, baseMsg)

	case config.JTSAlarmListReq:
		//获取报警信息请求
		getAlarmList(msg, wsConn, baseMsg)

	case config.JTSDriverActionReq:
		//获取驾驶行为列表
		getDriverActionList(msg, wsConn, baseMsg)

	case config.JTSSetParameterReq:
		setReportParameter(msg, wsConn, baseMsg)
	}
}

//处理注册请求
func getRegisterHandle(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg) {
	var registerReq RegisterReq
	if err := json.Unmarshal(msg.data, &registerReq); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSRegisterRes, "", config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}

	//生成签名
	sign := tools.SignMD5(tools.StringsJoin(registerReq.Data.Username, registerReq.Data.Password, registerReq.Data.Tradeno))
	wsConn.sign = sign
	client := &Client{
		sign:   sign,
		wsConn: wsConn,
	}

	//存入内存
	wsConnAll.set(sign, client)

	replyMap := getReplyMap(config.JTSRegisterRes, sign, config.JTT_ERROR_SUCCESS_OK,baseMsg.Sequence)
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
}

//获取设备列表
func getDeviceList(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg) {
	var getDeviceListReq GetDeviceListReq
	if err := json.Unmarshal(msg.data, &getDeviceListReq); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSGetDeviceListRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	if getDeviceListReq.Data.Page < 0 || getDeviceListReq.Data.Limit < 0 {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, get deviceList parameter error", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol)
		replyMap := getReplyMap(config.JTSGetDeviceListRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}

	replyMap := getReplyMap(config.JTSGetDeviceListRes, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK,baseMsg.Sequence)

	totalCount, _ := sqlDB.TotalCount(&sqlDB.GetVehicleInfo{})

	dataMap := replyMap["data"].(map[string]interface{})
	dataMap["totalCount"] = totalCount
	if totalCount == 0 {
		dataMap["deviceList"] = make([]interface{}, 0)
	} else {
		var deviceList []sqlDB.GetVehicleInfo
		if getDeviceListReq.Data.Page == 0 && getDeviceListReq.Data.Limit == 0 {
			//Page和Limit参数同时为0则查询数据库指定表格的所有数据
			sqlDB.Find(&deviceList, &sqlDB.GetVehicleInfo{})
		} else {
			start := (getDeviceListReq.Data.Page - 1) * getDeviceListReq.Data.Limit
			sqlDB.Limit(&deviceList, start, getDeviceListReq.Data.Limit)
		}

		dataMap["deviceList"] = deviceList
	}

	replyMap["data"] = dataMap
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
}

//处理实时视频控制
func getStreamStartOrStop(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg) {
	var streamStartReq StreamStartReq
	if err := json.Unmarshal(msg.data, &streamStartReq); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSStreamActionRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}

	if streamStartReq.Data.Action != 1 && streamStartReq.Data.Action != 2 {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, get channelList parameter error", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol)
		replyMap := getReplyMap(config.JTSStreamActionRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}

	//replyMap := getReplyMap(config.JTSStreamActionReq, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK)
	var retCode int64
	var URL string //实时直播返回的sessionURL
	//var channelListInfo sqlDB.GetChannelInfo
	//phoneNumAndChannel := streamStartReq.Data.PhoneNum + "_" + strconv.Itoa(int(streamStartReq.Data.ChannelID))

	//switch sqlDB.First(&channelListInfo, map[string]interface{}{"PhoneNumAndChannel": phoneNumAndChannel}) {
	//case 1:
	//记录存在
	//return udp.SendInviteStream(channelListInfo.DeviceID, channelID)
	switch streamStartReq.Data.Action {
		case 1:
		//开启实时直播
			//retCode, URL = udp.RequestRealStream(streamStartReq.Data.PhoneNum, streamStartReq.Data.ChannelID, 1, streamStartReq.Data.Bitstream)
			client := sertools.GetUDPClient(streamStartReq.Data.PhoneNum)
			if client !=  nil {
				retCode, URL = udp.RequestRealStream(streamStartReq.Data.PhoneNum, streamStartReq.Data.ChannelID, 1, streamStartReq.Data.Bitstream)
			} else {
				retCode, URL = tcp.RequestRealStream(streamStartReq.Data.PhoneNum, streamStartReq.Data.ChannelID, 1, streamStartReq.Data.Bitstream)
			}
		case 2:
		//关闭实时直播
		//retCode = udp.ControlRealStream(streamStartReq.Data.PhoneNum, streamStartReq.Data.ChannelID,0,0,0)
			retCode = config.JTT_ERROR_SUCCESS_OK
		default:
			retCode = config.JTT_ERROR_PARAMETER_ERROR
	}
	//case -1:
	//查询语句出错
	//retCode = config.JTT_ERROR_DB_ERR
	//case 0:
	//记录不存在
	//retCode = config.JTT_ERROR_CHANNEL_NOT_FOUNT
	//}

	replyMap := getReplyMap(config.JTSStreamActionRes, baseMsg.Data.Sign, retCode,baseMsg.Sequence)
	dataMap := replyMap["data"].(map[string]interface{})
	dataMap["phoneNum"] = streamStartReq.Data.PhoneNum
	dataMap["channelID"] = streamStartReq.Data.ChannelID
	dataMap["action"] = streamStartReq.Data.Action
	//if URL != "" {
	//实时直播返回正确的sessionURL
	dataMap["url"] = URL
	//}
	replyMap["data"] = dataMap
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
	//实时视频操作请求，type类型：1-开始；2-关闭
}

//处理录像回放控制
func getReplayStartOrStop(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg) {
	var videoReplayReq VideoReplayReq
	if err := json.Unmarshal(msg.data, &videoReplayReq); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSReplayActionRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}

	if len(videoReplayReq.Data.StartTime) != 14 || len(videoReplayReq.Data.EndTime) != 14 {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, get channelList parameter error", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol)
		replyMap := getReplyMap(config.JTSReplayActionRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}

	//replyMap := getReplyMap(config.JTSStreamActionReq, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK)
	var retCode int64
	var URL string //录像回放返回的sessionURL
	//var channelListInfo sqlDB.GetChannelInfo
	//phoneNumAndChannel := videoReplayReq.Data.PhoneNum + "_" + strconv.Itoa(int(videoReplayReq.Data.ChannelID))
	//
	//switch sqlDB.First(&channelListInfo, map[string]interface{}{"PhoneNumAndChannel": phoneNumAndChannel}) {
	//case 1:
	//记录存在
	//return udp.SendInviteStream(channelListInfo.DeviceID, channelID)
	client := wsConnAll.get(videoReplayReq.Data.Sign)
	switch videoReplayReq.Data.Action {
	case 1:
		//开启实时直播
		retCode, URL = udp.StartReplayVideoSend(client.username, videoReplayReq.Data.PhoneNum, videoReplayReq.Data.ChannelID, 2, 1, 0, 0, 0, videoReplayReq.Data.StartTime[2:], videoReplayReq.Data.EndTime[2:])
	case 2:
		//关闭实时直播
		retCode = udp.ControlReplayVideoSend(client.username, videoReplayReq.Data.PhoneNum, videoReplayReq.Data.ChannelID, 2, 0, videoReplayReq.Data.StartTime[2:])
	}
	//case -1:
	//查询语句出错
	//retCode = config.JTT_ERROR_DB_ERR
	//case 0:
	//记录不存在
	//	retCode = config.JTT_ERROR_CHANNEL_NOT_FOUNT
	//}

	replyMap := getReplyMap(config.JTSReplayActionRes, baseMsg.Data.Sign, retCode,baseMsg.Sequence)
	dataMap := replyMap["data"].(map[string]interface{})
	dataMap["phoneNum"] = videoReplayReq.Data.PhoneNum
	dataMap["channelID"] = videoReplayReq.Data.ChannelID
	dataMap["action"] = videoReplayReq.Data.Action
	if URL != "" {
		//实时直播返回正确的sessionURL
		dataMap["url"] = URL
	}
	replyMap["data"] = dataMap
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
	//实时视频操作请求，type类型：1-开始；2-关闭
}

//处理录像列表
func getVideoList(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg) {
	var videoListReq VideoListReq
	if err := json.Unmarshal(msg.data, &videoListReq); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSVideoListRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}

	if len(videoListReq.Data.StartTime) != 14 || len(videoListReq.Data.EndTime) != 14 {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, get channelList parameter error", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol)
		replyMap := getReplyMap(config.JTSVideoListRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}

	//replyMap := getReplyMap(config.JTSStreamActionReq, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK)
	var retCode int64
	var count uint32
	var videoList []server.VideoListInfo

	//var channelListInfo sqlDB.GetChannelInfo
	//phoneNumAndChannel := videoListReq.Data.PhoneNum + "_" + strconv.Itoa(int(videoListReq.Data.ChannelID))
	//
	//switch sqlDB.First(&channelListInfo, map[string]interface{}{"PhoneNumAndChannel": phoneNumAndChannel}) {
	//case 1:
	//记录存在
	client := wsConnAll.get(videoListReq.Data.Sign)
	retCode, count, videoList = udp.VideoListGetSend(client.username, videoListReq.Data.PhoneNum, videoListReq.Data.StartTime[2:], videoListReq.Data.EndTime[2:], 0, videoListReq.Data.ChannelID, 2, 0, 0)

	//case -1:
	//	//查询语句出错
	//	retCode = config.JTT_ERROR_DB_ERR
	//case 0:
	//	//记录不存在
	//	retCode = config.JTT_ERROR_CHANNEL_NOT_FOUNT
	//}

	replyMap := getReplyMap(config.JTSVideoListRes, baseMsg.Data.Sign, retCode,baseMsg.Sequence)
	dataMap := replyMap["data"].(map[string]interface{})
	dataMap["phoneNum"] = videoListReq.Data.PhoneNum
	dataMap["channelID"] = videoListReq.Data.ChannelID
	dataMap["count"] = count
	if retCode == config.JTT_ERROR_SUCCESS_OK {
		dataMap["videoList"] = videoList
	}
	replyMap["data"] = dataMap
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
	//实时视频操作请求，type类型：1-开始；2-关闭
}

//获取GPS轨迹信息请求
func getGPSList(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg) {
	var getGPSInfoReq GetGPSInfoReq
	if err := json.Unmarshal(msg.data, &getGPSInfoReq); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSGetGPSInfoRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	phoneAndTable := new(sqlDB.GpsAndPhoneNum)
	if !sqlDB.QueryUserTake(phoneAndTable, map[string]interface{}{"PhoneNum": getGPSInfoReq.Data.PhoneNum}) {
		logs.BeeLogger.Error("cannot find %s gps table", getGPSInfoReq.Data.PhoneNum)
		replyMap := getReplyMap(config.JTSGetGPSInfoRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	gpsList, retBool := sqlDB.QueryFindGPS(phoneAndTable.GpsTableName)
	if !retBool {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol)
		replyMap := getReplyMap(config.JTSGetGPSInfoRes, baseMsg.Data.Sign, config.JTT_ERROR_DB_ERR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	startFormatTime, _ := time.Parse("20060102150405", getGPSInfoReq.Data.StartTime)
	startTime := startFormatTime.Unix()
	endFormatTime, _ := time.Parse("20060102150405", getGPSInfoReq.Data.EndTime)
	endTime := endFormatTime.Unix()
	var newGpsList []sqlDB.GetLocationInfo
	for _, v := range gpsList {
		FormatTime, _ := time.Parse("2006-01-02 15:04:05", v.Time)
		Time := FormatTime.Unix()
		if startTime <= Time && Time <= endTime {
			newGpsList = append(newGpsList, v)
		}
	}
	replyMap := getReplyMap(config.JTSGetGPSInfoRes, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK,baseMsg.Sequence)
	dataMap := replyMap["data"].(map[string]interface{})

	dataMap["gpsInfo"] = newGpsList
	dataMap["totalCount"] = len(newGpsList)
	replyMap["data"] = dataMap
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
}

//获取实时GPS
func getRealGPS(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg) {
	var realGPSReq RealGPSReq
	if err := json.Unmarshal(msg.data, &realGPSReq); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSRealGPSRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	//fmt.Println(realGPSReq)
	//retCode, realGps := udp.QueryLocationInfo(realGPSReq.Data.PhoneNum)
	//
	//replyMap := getReplyMap(config.JTSGetGPSInfoRes, baseMsg.Data.Sign, retCode)
	realGps, ok := sertools.NewGPSInfo.Load(realGPSReq.Data.PhoneNum)
	var retCode int64
	if ok {
		retCode = config.JTT_ERROR_SUCCESS_OK
	} else {
		retCode = config.JTT_ERROR_DEVICELOST
	}

	replyMap := getReplyMap(config.JTSRealGPSRes, baseMsg.Data.Sign, retCode,baseMsg.Sequence)
	dataMap := replyMap["data"].(map[string]interface{})
	dataMap["gpsInfo"] = realGps
	replyMap["data"] = dataMap
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
}

//获取报警信息列表
func getAlarmList(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg) {
	var getAlarmInfoReq GetGPSInfoReq
	if err := json.Unmarshal(msg.data, &getAlarmInfoReq); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSAlarmListRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	phoneAndTable := new(sqlDB.GpsAndPhoneNum)
	if !sqlDB.QueryUserTake(phoneAndTable, map[string]interface{}{"PhoneNum": getAlarmInfoReq.Data.PhoneNum}) {
		logs.BeeLogger.Error("cannot find %s gps table", getAlarmInfoReq.Data.PhoneNum)
		replyMap := getReplyMap(config.JTSAlarmListRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	alarmList, retBool := sqlDB.QueryFindAlarm(phoneAndTable.AlarmTableName)
	if !retBool {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol)
		replyMap := getReplyMap(config.JTSAlarmListRes, baseMsg.Data.Sign, config.JTT_ERROR_DB_ERR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	startFormatTime, _ := time.Parse("20060102150405", getAlarmInfoReq.Data.StartTime)
	startTime := startFormatTime.Unix()
	endFormatTime, _ := time.Parse("20060102150405", getAlarmInfoReq.Data.EndTime)
	endTime := endFormatTime.Unix()
	var newAlarmList []sqlDB.GetAlarmInfo
	for _, v := range alarmList {
		FormatTime, _ := time.Parse("2006-01-02 15:04:05", v.Time)
		Time := FormatTime.Unix()
		if startTime <= Time && Time <= endTime {
			newAlarmList = append(newAlarmList, v)
		}
	}
	replyMap := getReplyMap(config.JTSAlarmListRes, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK,baseMsg.Sequence)
	dataMap := replyMap["data"].(map[string]interface{})

	dataMap["alarmInfo"] = newAlarmList
	dataMap["totalCount"] = len(newAlarmList)
	replyMap["data"] = dataMap
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
}

//获取驾驶行为列表
func getDriverActionList(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg) {
	var getDriverInfoReq GetGPSInfoReq
	if err := json.Unmarshal(msg.data, &getDriverInfoReq); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSDriverActionRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	phoneAndTable := new(sqlDB.GpsAndPhoneNum)
	if !sqlDB.QueryUserTake(phoneAndTable, map[string]interface{}{"PhoneNum": getDriverInfoReq.Data.PhoneNum}) {
		logs.BeeLogger.Error("cannot find %s gps table", getDriverInfoReq.Data.PhoneNum)
		replyMap := getReplyMap(config.JTSDriverActionRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	driverActionList, retBool := sqlDB.QueryFindDriverAction(phoneAndTable.DriverTableName)
	if !retBool {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol)
		replyMap := getReplyMap(config.JTSDriverActionRes, baseMsg.Data.Sign, config.JTT_ERROR_DB_ERR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	startFormatTime, _ := time.Parse("20060102150405", getDriverInfoReq.Data.StartTime)
	startTime := startFormatTime.Unix()
	endFormatTime, _ := time.Parse("20060102150405", getDriverInfoReq.Data.EndTime)
	endTime := endFormatTime.Unix()
	var newDriverList []sqlDB.DrivingAction
	for _, v := range driverActionList {
		FormatTime, _ := time.Parse("2006-01-02 15:04:05", v.Time)
		Time := FormatTime.Unix()
		if startTime <= Time && Time <= endTime {
			newDriverList = append(newDriverList, v)
		}
	}
	replyMap := getReplyMap(config.JTSDriverActionRes, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK,baseMsg.Sequence)
	dataMap := replyMap["data"].(map[string]interface{})

	dataMap["driverActionList"] = newDriverList
	dataMap["totalCount"] = len(newDriverList)
	replyMap["data"] = dataMap
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
}

//设置位置汇报时间间隔
func setReportParameter(msg *wsMessage, wsConn *wsConnection, baseMsg BaseMsg)  {
	var parameter reportParameter
	if err := json.Unmarshal(msg.data, &parameter); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, wsReadHandle() ==> protocol=%s, json.Unmarshal() error: %s", wsConn.wsSocket.RemoteAddr(), baseMsg.Protocol, err)
		replyMap := getReplyMap(config.JTSSetParameterRes, baseMsg.Data.Sign, config.JTT_ERROR_PARAMETER_ERROR,baseMsg.Sequence)
		wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
		return
	}
	timeStr := hex.EncodeToString(tools.Uint32ToByte(parameter.Data.Time))
	fmt.Println(timeStr)
	udp.SetDeviceParameter(parameter.Data.PhoneNum, 1, timeStr)
	replyMap := getReplyMap(config.JTSSetParameterRes, baseMsg.Data.Sign, config.JTT_ERROR_SUCCESS_OK,baseMsg.Sequence)
	wsReplyWrite(msg, wsConn, baseMsg.Protocol, replyMap)
}
//构造websocket回复Map
func getReplyMap(protocol, sign string, errCode int64, sequence uint64) map[string]interface{} {
	replyMap := map[string]interface{}{
		"errCode": errCode,
		"errMsg":  config.ErrCodeMap[errCode],
		"sequence": sequence,
		"data": map[string]interface{}{
			"sign": sign,
		},
		"protocol": protocol,
	}

	return replyMap
}

//构造websocket回复函数
func wsReplyWrite(msg *wsMessage, wsConn *wsConnection, protocol string, replyMap map[string]interface{}) {
	replyByte, _ := json.MarshalIndent(&replyMap, " ", "  ")
	if err := wsConn.wsWrite(msg.messageType, replyByte); err != nil {
		logs.BeeLogger.Error("remoteAddr=%v, protocol=%s ---> websocket write message to queue error:%s", wsConn.wsSocket.RemoteAddr(), protocol, err)
		wsConn.close()
	}
}
