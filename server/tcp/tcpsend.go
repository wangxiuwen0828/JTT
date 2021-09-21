package tcp

import (
	"encoding/hex"
	"fmt"
	"gitee.com/ictt/JTTM/config"
	"gitee.com/ictt/JTTM/server"
	"gitee.com/ictt/JTTM/server/sertools"
	"gitee.com/ictt/JTTM/tools"
	"gitee.com/ictt/JTTM/tools/logs"
	"gitee.com/ictt/JTTM/tools/sqlDB"
	"strconv"
)

/*
向终端请求开始录像播放
username——websocket传入的客户唯一标识符，phoneNum——设备号码，channel——设备逻辑通道号
dataType——音视频类型：0：音视频，1：音频，2：视频，3：视频或音视频
Bitstream——码流类型：0：主码流或子码流，1：主码流，2：子码流；如果此通道只传输音频，此字段置 0
storageType——存储器类型：0：主存储器或灾备存储器，1：主存储器，2：灾备存储器
PlaybackMode——回放方式：0：正常回放；1：快进回放；2：关键帧快退回放；3：关键帧播放；4、单帧上传；
speed——快进或快退倍数：回放方式为 1 和 2 时，此字段内容有效，否则置 0，  0：无效；  1：1 倍；  2：2 倍；  3：4 倍；  4：8 倍；  5：16 倍
startTime， endTime——开始和结束时间，省略年的前两位，如2018年9月6号13点14分12秒 在此应输入 180906131412 ，2018省略了20

返回字段 int64是errCode ,string是播放的url
*/
func StartReplayVideoSend(username, phoneNum string, channel, dataType, Bitstream, storageType, PlaybackMode, speed byte, startTime, endTime string) (int64, string) {

	client := sertools.GetTCPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return config.JTT_ERROR_DEVICELOST, ""
	}

	phoneAndSeq := phoneNum + strconv.FormatInt(int64(client.Sequence), 10)
	signPhoneAndChannel := username + phoneNum + hex.EncodeToString([]byte{channel})
	sertools.StreamChannel.Store(phoneAndSeq, signPhoneAndChannel)

	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return config.JTT_ERROR_DEVICELOST, ""
	}
	n := int64(len(config.IP))
	mesLength := tools.Int64ToByte(n + 23)
	ipLength := tools.Int64ToByte(n)
	sendMes := append([]byte{0x92, 0x01}, mesLength[6:]...)
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	sendMes = append(sendMes, ipLength[7])
	sendMes = append(sendMes, []byte(config.StreamIP)...)
	tcpPortByte := tools.Int64ToByte(config.StreamPort)
	sendMes = append(sendMes, tcpPortByte[6:]...)
	udpPortByte := tools.Int64ToByte(0)
	sendMes = append(sendMes, udpPortByte[6:]...)
	sendMes = append(sendMes, channel, dataType, Bitstream, storageType, PlaybackMode, speed)
	startByte, _ := hex.DecodeString(startTime)
	endByte, _ := hex.DecodeString(endTime)
	sendMes = append(sendMes, startByte...)
	sendMes = append(sendMes, endByte...)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetTCPClient(phoneNum, client)
	_, err = client.TCPConn.Write(resultSendMes)
	if err != nil {
		//StreamChannel.Delete(phoneAndSeq)
		fmt.Println("connect error")
		return config.JTT_ERROR_DEVICELOST, ""
	}

	fmt.Println("发送录像回放指令成功")
	return sertools.WaitGetReplayURL(signPhoneAndChannel)
	//return waitGetSessionURL(phoneAndChannel)
}

/*
向终端发送控制录像指令
username——websocket传入的客户唯一标识符，phoneNum——设备号码，channel——设备逻辑通道号
control——回放控制：0：开始传输；1：暂停传输；2：结束传输；3：快进传输；4：关键帧快退传输；5：拖动传输；6：关键帧传输
speed——快进或快退倍数：传输控制为 3 和 4 时，此字段内容有效，否则置 0，  0：无效；  1：1 倍；  2：2 倍；  3：4 倍；  4：8 倍；  5：16 倍
time——拖动回放位置，回放控制为 5 时，此字段有效。省略年的前两位，如2018年9月6号13点14分12秒 在此应输入 180906131412 ，2018省略了20

返回字段 int64是errCode
*/
func ControlReplayVideoSend(username, phoneNum string, channel, control, speed byte, time string) int64 {
	client := sertools.GetTCPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return config.JTT_ERROR_DEVICELOST
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return config.JTT_ERROR_DEVICELOST
	}
	phoneAndSeq := phoneNum + strconv.FormatInt(int64(client.Sequence), 10)
	signPhoneAndChannel := username + phoneNum + hex.EncodeToString([]byte{channel})
	sertools.StreamChannel.Store(phoneAndSeq, signPhoneAndChannel)

	sendMes := []byte{0x92, 0x02, 0, 9}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)
	sendMes = append(sendMes, channel, control, speed)
	timeByte, _ := hex.DecodeString(time)
	sendMes = append(sendMes, timeByte...)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetTCPClient(phoneNum, client)
	_, err = client.TCPConn.Write(resultSendMes)
	if err != nil {
		fmt.Println("connect error")
		return config.JTT_ERROR_DEVICELOST
	}
	fmt.Println("发送控制录像回放指令成功")

	return sertools.WaitReplayEndCode(signPhoneAndChannel)
}

/*
发送查询录像资源列表
username——websocket传入的客户唯一标识符，phoneNum——设备号码，channel——设备逻辑通道号
startTime， endTime——开始和结束时间，省略年的前两位，如2018年9月6号13点14分12秒 在此应输入 180906131412 ，2018省略了20
alarmType ——报警标志，bit0-31 见 JT/T 808-2011 表 18 报警标志位定义，bit32-63 见 表 11；全 0 表示无报警类型条件
mediaType——音视频资源类型：0：音视频，1：音频，2：视频，3：视频或音视频
Bitstream——码流类型：0：所有码流，1：主码流，2：子码流；如果此通道只传输音频，此字段置 0
storageType——存储器类型：0：所有存储器，1：主存储器，2：灾备存储器

返回数据 int64——是errCode，uint32是指录像文件的数量，[]server.VideoListInfo是指具体的录像文件信息
*/
func VideoListGetSend(username, phoneNum string, startTime, endTime string, alarmType uint64, channel, mediaType, Bitstream, storageType byte) (int64, uint32, []server.VideoListInfo) {
	client := sertools.GetTCPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return config.JTT_ERROR_DEVICELOST, 0, nil
	}

	phoneAndSeq := phoneNum + strconv.FormatInt(int64(client.Sequence), 10)
	signPhoneAndChannel := username + phoneNum + hex.EncodeToString([]byte{channel})
	sertools.StreamChannel.Store(phoneAndSeq, signPhoneAndChannel)

	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return config.JTT_ERROR_PARAMETER_ERROR, 0, nil
	}

	sendMes := []byte{0x92, 0x05, 0, 24}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	sendMes = append(sendMes, channel)
	startByte, _ := hex.DecodeString(startTime)
	endByte, _ := hex.DecodeString(endTime)
	sendMes = append(sendMes, startByte...)
	sendMes = append(sendMes, endByte...)
	alarmTypeByte := tools.Uint64ToByte(alarmType)
	sendMes = append(sendMes, alarmTypeByte...)
	sendMes = append(sendMes, mediaType, Bitstream, storageType)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetTCPClient(phoneNum, client)
	_, err = client.TCPConn.Write(resultSendMes)
	if err != nil {
		fmt.Println("connect error")
		return config.JTT_ERROR_DEVICELOST, 0, nil
	}
	fmt.Println("发送查询录像资源列表命令成功")
	return sertools.WaitReplayList(signPhoneAndChannel)
}

/*
向终端发送控制实时直播指令
phoneNum——设备号码，channel——设备逻辑通道号
control——控制指令：0：关闭音视频传输指令  1：切换码流（增加暂停和继续）  2：暂停该通道所有流的发送  3：恢复暂停前流的发送，与暂停前的流类型一致
audioAndVideo——关闭音视频类型：0：关闭该通道有关的音视频数据  1：只关闭该通道有关的音频，保留该通道有关的视频  2：只关闭该通道有关的视频，保留该通道有关的音频
Bitstream——码流类型：0：主码流，1：子码流

返回字段 int64是errCode
*/
func ControlRealStream(phoneNum string, channel, control, audioAndVideo, Bitstream byte) int64 {
	client := sertools.GetTCPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return config.JTT_ERROR_DEVICELOST
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return config.JTT_ERROR_PARAMETER_ERROR
	}
	phoneAndSeq := phoneNum + strconv.FormatInt(int64(client.Sequence), 10)
	phoneAndChannel := phoneNum + hex.EncodeToString([]byte{channel})
	//_,ok := StreamInfo.Load(phoneAndChannel)
	//if !ok {
	//	return config.JTT_ERROR_SUCCESS_OK
	//} //else {
	//	streamInfoData,_ := streamInfo.(StreamInfoData)
	//if streamInfoData.Count > 1 {
	//	streamInfoData.Count--
	//  StreamInfo.Store(phoneAndChannel, streamInfoData)
	//  return config.JTT_ERROR_SUCCESS_OK
	//}
	//}
	sertools.StreamChannel.Store(phoneAndSeq, phoneAndChannel)

	sendMes := []byte{0x91, 0x02, 0, 4}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	sendMes = append(sendMes, channel, control, audioAndVideo, Bitstream)
	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetTCPClient(phoneNum, client)
	_, err = client.TCPConn.Write(resultSendMes)
	if err != nil {
		fmt.Println("connect error")
		return config.JTT_ERROR_DEVICELOST
	}
	fmt.Println("发送控制实时音视频成功")
	logs.BeeLogger.Info("发送 %s", hex.EncodeToString(resultSendMes))
	return config.JTT_ERROR_SUCCESS_OK
	//return waitGetStreamEndCode(phoneAndChannel)
}

/*
向终端请求实时直播
phoneNum——设备号码，channel——设备逻辑通道号
dataType——数据类型：0：音视频，1：视频，2：双向对讲，3：监听，4：中心广播，5：透传
Bitstream——码流类型：0：主码流，  1：子码流

返回字段 int64是errCode ,string是播放的url
*/
func RequestRealStream(phoneNum string, channel, dataType, Bitstream byte) (int64, string) {
	client := sertools.GetTCPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return config.JTT_ERROR_DEVICELOST, ""
	}

	phoneAndSeq := phoneNum + strconv.FormatInt(int64(client.Sequence), 10)
	phoneAndChannel := phoneNum + hex.EncodeToString([]byte{channel})
	streamInfo, ok := sertools.StreamInfo.Load(phoneAndChannel)
	if ok {

		streamInfoData, _ := streamInfo.(sertools.StreamInfoData)
		streamInfoData.Count++
		sertools.StreamInfo.Store(phoneAndChannel, streamInfoData)
		if streamInfoData.URL != "" {
			url := sertools.GetSessionURLFromSTS(phoneAndChannel)
			if url == "" {
				sertools.StreamInfo.Delete(phoneAndChannel)
				return config.JTT_ERROR_OVERTIME, ""
			} else {
				return config.JTT_ERROR_SUCCESS_OK, url
			}
		}
	}
	sertools.StreamChannel.Store(phoneAndSeq, phoneAndChannel)

	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return config.JTT_ERROR_DEVICELOST, ""
	}
	n := int64(len(config.StreamIP))
	mesLength := tools.Int64ToByte(n + 8)
	ipLength := tools.Int64ToByte(n)
	sendMes := append([]byte{0x91, 0x01}, mesLength[6:]...)
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	sendMes = append(sendMes, ipLength[7])
	sendMes = append(sendMes, []byte(config.StreamIP)...)
	tcpPortByte := tools.Int64ToByte(config.StreamPort)
	sendMes = append(sendMes, tcpPortByte[6:]...)
	udpPortByte := tools.Int64ToByte(config.StreamPort)
	sendMes = append(sendMes, udpPortByte[6:]...)
	sendMes = append(sendMes, channel, dataType, Bitstream)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetTCPClient(phoneNum, client)
	logs.BeeLogger.Info("send to tcp: %s", hex.EncodeToString(resultSendMes))
	_, err = client.TCPConn.Write(resultSendMes)
	if err != nil {
		sertools.StreamChannel.Delete(phoneAndSeq)
		fmt.Println("connect error")
		return config.JTT_ERROR_DEVICELOST, ""
	}

	fmt.Println("发送获取实时音视频成功")
	return sertools.WaitGetSessionURL(phoneAndChannel)
}

/*
主动查询实时位置信息
phoneNum——设备号码，

返回字段 int64是errCode ,gpsInfo 是实时gps数据
*/
func QueryLocationInfo(phoneNum string) (errCode int64, gpsInfo sqlDB.GetLocationInfo) {
	client := sertools.GetTCPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		errCode = config.JTT_ERROR_DEVICELOST
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		errCode = config.JTT_ERROR_PARAMETER_ERROR
		return
	}
	sendMes := []byte{0x82, 0x01, 0, 0}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)
	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetTCPClient(phoneNum, client)
	_, err = client.TCPConn.Write(resultSendMes)
	if err != nil {
		fmt.Println("connect error")
		errCode = config.JTT_ERROR_DEVICELOST
		return
	}
	fmt.Println("发送查询位置参数成功")
	errCode, gpsInfo = sertools.WaitGetGPS(phoneNum)
	return
}

///*
//设置终端的圆形区域
//
// */
//
//func SetRoundRegionSend (phoneNum string, ) (errCode int64) {
//
//}