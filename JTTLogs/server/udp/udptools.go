package udp
//
//import (
//	"bytes"
//	"encoding/binary"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"gitee.com/ictt/JTTM/config"
//	"gitee.com/ictt/JTTM/server"
//	"gitee.com/ictt/JTTM/tools"
//	"gitee.com/ictt/JTTM/tools/logs"
//	"gitee.com/ictt/JTTM/tools/sqlDB"
//	"github.com/patrickmn/go-cache"
//	"io/ioutil"
//	"net"
//	"net/http"
//	"strconv"
//	"time"
//)
//
//var (
//	//记录UDP连接信息，当设备注册成功时存储记录，键为Device，值为map键值对
//	//map键值对的键为UDP连接的IP+Port，值为UDPClient
//	udpConnCache *cache.Cache
//	picCache     *cache.Cache
//)
//
//type udpClient struct {
//	PhoneNum string //设备ID
//	//UDPAddrKey string //UDPAddr中的IP+Port组成的字符串
//	UDPConn    *net.UDPConn
//	UDPAddr    *net.UDPAddr
//	UpdateTime int64
//	Sequence   uint16 //记录请求设备配置时生成的随机数，Sequence+2作为我红方向蓝方发送请求信令时使用
//}
//
//type pictureMes struct {
//	PhoneNumAndSeq string
//	//packageNum  uint16
//	newMessByte string
//}
//
//func init() {
//	udpConnCache = cache.New(90*time.Second, 30*time.Second)
//	picCache = cache.New(30*time.Second, 15*time.Second)
//}
//
////存储一个udpClient，初次存储是设备发送注册请求
//func setUDPClient(PhoneNum string, client *udpClient) {
//	//fmt.Println("存储信息")
//	udpConnCache.Set(PhoneNum, client, cache.DefaultExpiration)
//}
//
////存储图片数据
//func setPicture(PhoneNum string, pictureMes *pictureMes) {
//	picCache.Set(PhoneNum, pictureMes, cache.DefaultExpiration)
//}
//
////
//
//////获取一个UDPClient
////func getUDPClient(udpConnKey string) *udpClient {
////	for _, v := range udpConnCache.Items() {
////		client := v.Object.(*udpClient)
////		if client.UDPAddrKey == udpConnKey {
////			return client
////		}
////	}
////
////	return nil
////}
//
////使用PhoneNum筛选出指定的UDPClient
//func GetUDPClient(PhoneNum string) *udpClient {
//	temp, ok := udpConnCache.Get(PhoneNum)
//	if ok {
//		return temp.(*udpClient)
//	}
//
//	logs.BeeLogger.Error("deviceID=%s is disconnect or for other reasons", PhoneNum)
//	return nil
//}
//
////使用PhoneNum获取图片数据
//func GetPicture(PhoneNum string) *pictureMes {
//	temp, ok := picCache.Get(PhoneNum)
//	if ok {
//		return temp.(*pictureMes)
//	}
//
//	logs.BeeLogger.Error("deviceID=%s is disconnect or for other reasons", PhoneNum)
//	return nil
//}
//
////位置信息解析
//func LocationInfo(messByte []byte, locateDataReceive *server.LocateDataReceive) {
//	AlarmSign10 := int64(binary.BigEndian.Uint32(messByte[:4]))
//	AlarmSign2 := strconv.FormatInt(AlarmSign10, 2)
//	AlarmSignDefinition := [32]string{"非法开门", "侧翻预警", "碰撞预警", "车辆非法位移", "车辆非法点火", "车辆被盗",
//		"车辆油量异常", "车辆VSS故障", "路线偏离报警", "路段行驶时间不足/过长", "进出路线", "进出区域", "超时停车",
//		"当天累计驾驶超时", "右转盲区异常报警", "胎压预警", "违规行驶报警", "疲劳驾驶预警", "超速预警", "道路运输证IC卡模块故障",
//		"摄像头故障", "TTS模块故障", "终端LCD或显示屏故障", "终端主电源掉电", "终端主电源欠压", "GNSS天线短路",
//		"GNSS天线未接或被剪短", "GNSS模块发生故障", "危险驾驶", "疲劳驾驶", "超速报警", "紧急报警"}
//	for len(AlarmSign2) < 32 {
//		AlarmSign2 = "0" + AlarmSign2
//	}
//	locateDataReceive.AlarmSign = []byte(AlarmSign2)
//	for i, v := range locateDataReceive.AlarmSign {
//		if v == '1' {
//			fmt.Println(AlarmSignDefinition[i])
//			locateDataReceive.InfoType = "Alarm"
//			locateDataReceive.AlarmState += " " + AlarmSignDefinition[i]
//		}
//	}
//
//	state10 := int64(binary.BigEndian.Uint32(messByte[4:8]))
//	state2 := strconv.FormatInt(state10, 2)
//	stateDefinition1 := [32]string{"", "", "", "", "", "",
//		"", "", "", "", "使用 Galileo 卫星进行定位", "使用 GLONASS 卫星进行定位", "使用北斗卫星进行定位",
//		"使用 GPS 卫星进行定位", "门5开（自定义）", "门4开（驾驶席门）", "门3开（后门）", "门2开（中门）", "门1开（前门）", "车门加锁",
//		"车辆电路断开", "车辆油路断开", "", "", "车道偏移预警", "紧急刹车前撞预警",
//		"经纬度已经保密插件加密", "停运状态", "西经", "南纬", "定位中", "ACC 开"}
//	stateDefinition0 := [32]string{"", "", "", "", "", "",
//		"", "", "", "", "未使用 Galileo 卫星进行定位", "未使用 GLONASS 卫星进行定位", "未使用北斗卫星进行定位",
//		"未使用 GPS 卫星进行定位", "门5关（自定义）", "门4关（驾驶席门）", "门3关（后门）", "门2关（中门）", "门1关（前门）", "车门解锁",
//		"车辆电路正常", "车辆油路正常", "", "", "", "",
//		"经纬度未经保密插件加密", "营运状态", "东经", "北纬", "未定位", "ACC 关"}
//	for len(state2) < 32 {
//		state2 = "0" + state2
//	}
//	locateDataReceive.State = []byte(state2)
//	fmt.Println(string(locateDataReceive.State))
//	for i, v := range locateDataReceive.State {
//		if v == '1' {
//			fmt.Printf("%s	", stateDefinition1[i])
//		} else {
//			fmt.Printf("%s	", stateDefinition0[i])
//		}
//	}
//	switch string(locateDataReceive.State[22:24]) {
//	case "00":
//		fmt.Printf("空车\n")
//	case "10":
//		fmt.Printf("半载\n")
//	case "01":
//		fmt.Printf("\n")
//	case "11":
//		fmt.Printf("满载\n")
//
//	}
//
//	Latitude := float64(binary.BigEndian.Uint32(messByte[8:12])) / 1000000
//	locateDataReceive.Latitude = tools.Float64toDegree(Latitude)
//	switch locateDataReceive.State[29] {
//	case '0':
//		locateDataReceive.Latitude += " N"
//	case '1':
//		locateDataReceive.Latitude += " S"
//	}
//	Longitude := float64(binary.BigEndian.Uint32(messByte[12:16])) / 1000000
//	locateDataReceive.Longitude = tools.Float64toDegree(Longitude)
//	switch locateDataReceive.State[28] {
//	case '0':
//		locateDataReceive.Longitude += " E"
//	case '1':
//		locateDataReceive.Longitude += " W"
//	}
//	locateDataReceive.Altitude = float64(binary.BigEndian.Uint16(messByte[16:18]))
//	locateDataReceive.Speed = float64(binary.BigEndian.Uint16(messByte[18:20])) / 10
//	locateDataReceive.Direction = float64(binary.BigEndian.Uint16(messByte[20:22]))
//	time := hex.EncodeToString(messByte[22:28])
//	locateDataReceive.Time = "20" + time[:2] + "-" + time[2:4] + "-" + time[4:6] + " " + time[6:8] + ":" + time[8:10] + ":" + time[10:12]
//	fmt.Printf("纬度：%v	经度：%v	海拔：%v米\n速度：%vkm/h	方向：%v度	时间：%v\n",
//		locateDataReceive.Latitude, locateDataReceive.Longitude, locateDataReceive.Altitude,
//		locateDataReceive.Speed, locateDataReceive.Direction, locateDataReceive.Time)
//
//}
//
////附加位置信息解析
//func AdditionalInfo(messByte []byte, addLocateData *server.AddLocateMess) []byte {
//	var lenth byte
//	switch messByte[0] {
//	case 0x01:
//		lenth = messByte[1]
//		addLocateData.Mileage = float64(binary.BigEndian.Uint32(messByte[2:2+lenth])) / 10
//		fmt.Printf("行驶里程为 %v km\n", addLocateData.Mileage)
//
//	case 0x02:
//		lenth = messByte[1]
//		addLocateData.Oil = float64(binary.BigEndian.Uint16(messByte[2:2+lenth])) / 10
//		fmt.Printf("油量为 %v L\n", addLocateData.Oil)
//
//	case 0x03:
//		lenth = messByte[1]
//		addLocateData.SpeedRecode = float64(binary.BigEndian.Uint16(messByte[2:2+lenth])) / 10
//		fmt.Printf("速度 %v km/h\n", addLocateData.SpeedRecode)
//
//	case 0x04:
//		lenth = messByte[1]
//		addLocateData.AlarmMesID = int64(binary.BigEndian.Uint16(messByte[2 : 2+lenth]))
//		fmt.Printf("事件id %d \n", addLocateData.AlarmMesID)
//
//	case 0x05:
//		lenth = messByte[1]
//		addLocateData.TirePressure = hex.EncodeToString(messByte[2 : 2+lenth])
//		fmt.Printf("胎压 %v \n", addLocateData.TirePressure)
//
//	case 0x06:
//		lenth = messByte[1]
//		addLocateData.Temperature = int64(binary.BigEndian.Uint16(messByte[2 : 2+lenth]))
//		//messContent2 := strconv.FormatInt(messContent,2)
//		fmt.Printf("车厢温度 %v °C\n", addLocateData.Temperature)
//
//	case 0x11:
//		lenth = messByte[1]
//		addLocateData.SpeedAlarmAdd.LocateType, addLocateData.SpeedAlarmAdd.LocateID = region(messByte)
//	case 0x12:
//		lenth = messByte[1]
//		addLocateData.InAndOutAlarmAdd.LocateType, addLocateData.InAndOutAlarmAdd.LocateID = region(messByte)
//		addLocateData.InAndOutAlarmAdd.Directory = messByte[7]
//	case 0x13:
//		lenth = messByte[1]
//		addLocateData.TimeAlarmAdd.LocateID = int64(binary.BigEndian.Uint32(messByte[2:6]))
//		addLocateData.TimeAlarmAdd.LocateTime = int64(binary.BigEndian.Uint16(messByte[6:8]))
//		addLocateData.TimeAlarmAdd.Result = messByte[8]
//		if addLocateData.TimeAlarmAdd.Result == 0 {
//			fmt.Printf("路段行驶时间不足，id = %v, 时间 = %v\n", addLocateData.TimeAlarmAdd.LocateID, addLocateData.TimeAlarmAdd.LocateTime)
//		} else {
//			fmt.Printf("路段行驶时间过长，id = %v, 时间 = %v\n", addLocateData.TimeAlarmAdd.LocateID, addLocateData.TimeAlarmAdd.LocateTime)
//		}
//
//	case 0x14:
//		lenth = messByte[1]
//		VideoAlarm := binary.BigEndian.Uint32(messByte[2 : 2+lenth])
//		videoAlarm2 := tools.UintTo2byte(uint64(VideoAlarm), 32)
//		videoAlarmData := [32]string{"视频信号丢失", "主存储器故障", "灾备存储单元故障", "其他视频设备故障", "客车超载报警", "异常驾驶行为报警", "特殊报警录像达到存储阈值报警"}
//		for i, v := range videoAlarm2 {
//			if v == '1' {
//				addLocateData.VideoAlarm += videoAlarmData[31-i] + "  "
//			}
//		}
//		fmt.Printf("视频报警标志位%v\n", addLocateData.VideoAlarm)
//
//	case 0x15:
//		lenth = messByte[1]
//		VideoSignalLossAlarm := binary.BigEndian.Uint32(messByte[2 : 2+lenth])
//		VideoSignalLossAlarm2 := tools.UintTo2byte(uint64(VideoSignalLossAlarm), 32)
//		for i, v := range VideoSignalLossAlarm2 {
//			if v == '1' {
//				addLocateData.VideoSignalLossAlarm += strconv.Itoa(32-i) + " "
//				fmt.Printf("第 %v 个逻辑通道视频信号丢失\n", 32-i)
//			}
//		}
//
//	case 0x16:
//		lenth = messByte[1]
//		//VideoSignalShelterAlarm := binary.BigEndian.Uint32(messByte[2:2+lenth])
//		addLocateData.LineNumber = hex.EncodeToString(messByte[2 : 2+lenth])
//		fmt.Printf("线路编码 %v\n", addLocateData.LineNumber)
//		//for i,v := range addLocateData.VideoSignalShelterAlarm {
//		//	if v == '1'{
//		//		fmt.Printf("第 %v 个逻辑通道视频信号遮挡\n", 32 - i)
//		//	}
//		//}
//
//	case 0x17:
//		lenth = messByte[1]
//		//StorageFailureAlarm := messByte[2]
//		addLocateData.BusinessType = messByte[1]
//		fmt.Printf("业务类型 %v\n", addLocateData.BusinessType)
//
//	case 0x18:
//		lenth = messByte[1]
//		AbnormalDrivingType := binary.BigEndian.Uint16(messByte[2:4])
//		AbnormalDrivingType2 := tools.UintTo2byte(uint64(AbnormalDrivingType), 16)
//		AbnormalDrivingTypeData := [16]string{"疲劳", "打电话", "抽烟"}
//		for i, v := range AbnormalDrivingType2 {
//			if v == '1' {
//				addLocateData.AbnormalDrivingAlarm.AbnormalDrivingType += AbnormalDrivingTypeData[15-i] + "  "
//			}
//		}
//		addLocateData.AbnormalDrivingAlarm.FatigueDegree = messByte[4]
//		fmt.Printf("异常驾驶行为类型 %v ,疲劳程度 %v \n", addLocateData.AbnormalDrivingAlarm.AbnormalDrivingType,
//			addLocateData.AbnormalDrivingAlarm.FatigueDegree)
//
//	case 0x25:
//		lenth = messByte[1]
//		messContent := int64(binary.BigEndian.Uint32(messByte[2 : 2+lenth]))
//		vehicleState := strconv.FormatInt(messContent, 2)
//		for len(vehicleState) < 32 {
//			vehicleState = "0" + vehicleState
//		}
//		addLocateData.VehicleState = []byte(vehicleState)
//		fmt.Printf("信号灯状态 %v\n", string(addLocateData.VehicleState))
//	case 0x2A:
//		lenth = messByte[1]
//		messContent := int64(binary.BigEndian.Uint16(messByte[2 : 2+lenth]))
//		IOState := strconv.FormatInt(messContent, 2)
//		for len(IOState) < 16 {
//			IOState = "0" + IOState
//		}
//		addLocateData.IOState = []byte(IOState)
//		fmt.Printf("IO状态 %v\n", string(addLocateData.IOState))
//	case 0x2B:
//		lenth = messByte[1]
//		messContent := int64(binary.BigEndian.Uint32(messByte[2 : 2+lenth]))
//		AD := strconv.FormatInt(messContent, 2)
//		for len(AD) < 32 {
//			AD = "0" + AD
//		}
//		addLocateData.AD = []string{AD[:16], AD[16:]}
//		fmt.Printf("模拟量 AD0 %v , AD1 %v\n", addLocateData.AD[1], addLocateData.AD[0])
//	case 0x30:
//		lenth = messByte[1]
//		addLocateData.SignalIntensity = messByte[2]
//		fmt.Printf("无线通信网络信号强度 %v\n", addLocateData.SignalIntensity)
//	case 0x31:
//		lenth = messByte[1]
//		addLocateData.SatelliteNum = messByte[2]
//		fmt.Printf("GNSS定位卫星数 %v\n", addLocateData.SatelliteNum)
//	default:
//		lenth = messByte[1]
//		custom := server.CustomData{}
//		custom.ID = messByte[0]
//		custom.Data = hex.EncodeToString(messByte[2 : 2+lenth])
//		addLocateData.Custom = append(addLocateData.Custom, custom)
//		fmt.Printf("自定义 id%v data %v\n", custom.ID, custom.Data)
//	}
//
//	return messByte[2+lenth:]
//}
//
////区域信息
//func region(messByte []byte) (locate byte, id int64) {
//	var res string
//	locate = messByte[2]
//	switch locate {
//	case 0:
//		res = "无特定位置"
//	case 1:
//		id = int64(binary.BigEndian.Uint32(messByte[3:7]))
//		res = "圆形区域ID " + strconv.FormatInt(id, 10)
//	case 2:
//		id = int64(binary.BigEndian.Uint32(messByte[3:7]))
//		res = "矩形区域ID " + strconv.FormatInt(id, 10)
//	case 3:
//		id = int64(binary.BigEndian.Uint32(messByte[3:7]))
//		res = "多边形区域ID " + strconv.FormatInt(id, 10)
//	case 4:
//		id = int64(binary.BigEndian.Uint32(messByte[3:7]))
//		res = "路段ID " + strconv.FormatInt(id, 10)
//	}
//	fmt.Println(res)
//	return
//}
//
////通用回答
//func NormalResponse(messByte []byte, sequence uint16, i int) (resultSendMes []byte) {
//	//回复数据编写
//	messIDSend, _ := hex.DecodeString("8001")
//	//messIDSend := []byte{129
//	messHead10_intSend := []byte{0, 5}
//	sendMes := append(messIDSend, messHead10_intSend...)
//	sendMes = append(sendMes, messByte[4:10]...)
//	sequenceByte := tools.Uint16ToByte(sequence)
//	sendMes = append(sendMes, sequenceByte...)
//	sendMes = append(sendMes, messByte[10:12]...)
//	sendMes = append(sendMes, messByte[:2]...)
//	result := []byte{0, 1, 2, 3, 4}
//	sendMes = append(sendMes, result[i])
//	//vehicleInfo.PowerIdentify = "abcdefg"
//	//sendMes = append(sendMes, []byte(vehicleInfo.PowerIdentify)...)
//	resultSendMes = tools.ParaphraseMess(sendMes)
//	return
//	//fmt.Println(resultSendMes)
//}
//
////终端的通用回答解析
//func MessageClass(messByte []byte) {
//	switch hex.EncodeToString(messByte[:2]) {
//	case "9101":
//		fmt.Printf("实时音视频传输开启请求")
//	case "9102":
//		fmt.Printf("实时音视频传输控制请求")
//	case "9201":
//		fmt.Printf("录像回放请求")
//	case "9202":
//		fmt.Printf("录像控制请求")
//	case "9205":
//		fmt.Printf("实时音视频传输状态")
//	case "8103":
//		fmt.Printf("设置终端参数")
//	case "8202":
//		fmt.Printf("临时位置跟踪控制")
//	case "8203":
//		fmt.Printf("人工确认报警控制")
//	case "8801":
//		fmt.Printf("摄像头立即拍摄命令")
//	case "8803":
//		fmt.Printf("多媒体数据上传命令")
//
//	}
//
//	switch messByte[2] {
//	case 0:
//		fmt.Printf("成功/确认\n")
//	case 1:
//		fmt.Printf("失败\n")
//	case 2:
//		fmt.Printf("消息有误\n")
//	case 3:
//		fmt.Printf("不支持\n")
//	}
//}
//
////等待设备返回启动实时直播成功指令，然后从内存取STS返回的sessionURL
//func WaitGetSessionURL(phoneAndChannel string) (int64, string) {
//	idleDelay := time.NewTimer(4 * time.Second)
//	defer idleDelay.Stop()
//
//	for {
//		select {
//		case <-idleDelay.C:
//			//超时
//			return config.JTT_ERROR_SERVER_ERROR, ""
//		default:
//			//每隔100毫秒请求一次数据
//			time.Sleep(100 * time.Millisecond)
//
//			streamInfo, OK := StreamInfo.Load(phoneAndChannel)
//			if OK {
//				streamInfoData, _ := streamInfo.(StreamInfoData)
//				if streamInfoData.URL != "" {
//					return config.JTT_ERROR_SUCCESS_OK, streamInfoData.URL
//				}
//			}
//
//		}
//	}
//}
//
////等待设备返回启动录像回放成功指令，然后从内存取STS返回的sessionURL
//func WaitGetReplayURL(signPhoneAndChannel string) (int64, string) {
//	idleDelay := time.NewTimer(4 * time.Second)
//	defer idleDelay.Stop()
//
//	for {
//		select {
//		case <-idleDelay.C:
//			//超时
//			return config.JTT_ERROR_SERVER_ERROR, ""
//		default:
//			//每隔100毫秒请求一次数据
//			time.Sleep(100 * time.Millisecond)
//
//			videoReply, OK := VideoReply.Load(signPhoneAndChannel)
//			if OK {
//				url, _ := videoReply.(string)
//				if url != "" {
//					return config.JTT_ERROR_SUCCESS_OK, url
//				}
//			}
//
//		}
//	}
//}
//
////等待回复的录像列表
//func WaitReplayList(signPhoneAndChannel string) (int64, uint32, []server.VideoListInfo) {
//	idleDelay := time.NewTimer(4 * time.Second)
//	defer idleDelay.Stop()
//
//	for {
//		select {
//		case <-idleDelay.C:
//			//超时
//			return config.JTT_ERROR_SERVER_ERROR, 0, nil
//		default:
//			//每隔100毫秒请求一次数据
//			time.Sleep(100 * time.Millisecond)
//
//			videoList, OK := VideoList.Load(signPhoneAndChannel)
//			if OK {
//				VideoList.Delete(signPhoneAndChannel)
//				videoListData, _ := videoList.(server.VideoListData)
//				if videoListData.SequenceSend != 0 {
//					return config.JTT_ERROR_SUCCESS_OK, videoListData.VideoCount, videoListData.VideoList
//				}
//			}
//
//		}
//	}
//}
//
////等待设备返实时gps数据
//func WaitGetGPS(phoneNum string) (errCode int64, gps sqlDB.GetLocationInfo) {
//	idleDelay := time.NewTimer(4 * time.Second)
//	defer idleDelay.Stop()
//
//	for {
//		select {
//		case <-idleDelay.C:
//			//超时
//			errCode = config.JTT_ERROR_SERVER_ERROR
//			return
//
//		default:
//			//每隔100毫秒请求一次数据
//			time.Sleep(100 * time.Millisecond)
//			realGps, OK := RealGPSInfo.Load(phoneNum)
//			if OK {
//				RealGPSInfo.Delete(phoneNum)
//				gps, _ = realGps.(sqlDB.GetLocationInfo)
//				errCode = config.JTT_ERROR_SUCCESS_OK
//				return
//			}
//		}
//	}
//}
//
////等待设备返回关闭实时直播成功指令，然后返回结果
//func waitGetStreamEndCode(phoneAndChannel string) int64 {
//	idleDelay := time.NewTimer(4 * time.Second)
//	defer idleDelay.Stop()
//
//	for {
//		select {
//		case <-idleDelay.C:
//			//超时
//			return config.JTT_ERROR_SERVER_ERROR
//		default:
//			//每隔100毫秒请求一次数据
//			time.Sleep(100 * time.Millisecond)
//
//			_, OK := VideoReply.Load(phoneAndChannel)
//			if !OK {
//				return config.JTT_ERROR_SUCCESS_OK
//			}
//		}
//	}
//}
//
////等待设备返回关闭录像回放成功指令，然后返回结果
//func WaitReplayEndCode(phoneAndChannel string) int64 {
//	idleDelay := time.NewTimer(4 * time.Second)
//	defer idleDelay.Stop()
//
//	for {
//		select {
//		case <-idleDelay.C:
//			//超时
//			return config.JTT_ERROR_SERVER_ERROR
//		default:
//			//每隔100毫秒请求一次数据
//			time.Sleep(100 * time.Millisecond)
//
//			_, OK := StreamInfo.Load(phoneAndChannel)
//			if !OK {
//				return config.JTT_ERROR_SUCCESS_OK
//			}
//		}
//	}
//}
//
////处理消息头参数
//func HandleMessProperty(header *server.HeaderReceive, messByte []byte) ([]byte, bool, bool) {
//	sign := true
//	needRepeat := false
//	header.IDReceive = hex.EncodeToString(messByte[:2])
//	dataLength := binary.BigEndian.Uint16(messByte[2:4])
//	messHead2 := strconv.FormatInt(int64(dataLength), 2)
//	for len(messHead2) < 16 {
//		messHead2 = "0" + messHead2
//	}
//	header.MessProperty.Subpackage = messHead2[2]
//	header.MessProperty.Encryption = []byte{messHead2[3], messHead2[4], messHead2[5]}
//	header.MessProperty.DataLength = tools.TwoToInt64(messHead2[6:])
//	fmt.Printf("消息体长度%v\n", header.MessProperty.DataLength)
//	header.PhoneNum = hex.EncodeToString(messByte[4:10])
//	header.SequenceReceive = binary.BigEndian.Uint16(messByte[10:12])
//	var newMessByte []byte
//	switch header.MessProperty.Subpackage {
//	case '0':
//		fmt.Println("package no change")
//		newMessByte = messByte[12:]
//
//	case '1':
//		//fmt.Println("need change")
//		header.PackageCount = binary.BigEndian.Uint16(messByte[12:14])
//		header.PackageNum = binary.BigEndian.Uint16(messByte[14:16])
//		switch header.PackageNum {
//		case 1:
//			sign = false
//			picture := pictureMes{}
//			picture.PhoneNumAndSeq = header.PhoneNum + strconv.FormatInt(int64(header.SequenceReceive), 10)
//			picture.newMessByte = string(messByte[16:])
//			setPicture(picture.PhoneNumAndSeq, &picture)
//
//		case header.PackageCount:
//			sign = true
//			time.Sleep(2 * time.Second)
//			var newMessStr string
//			for i := int64(header.SequenceReceive - header.PackageCount + 1); i < int64(header.SequenceReceive); i++ {
//				phoneNumAndSeq := header.PhoneNum + strconv.FormatInt(i, 10)
//				picture := GetPicture(phoneNumAndSeq)
//				if picture == nil {
//					needRepeat = true
//				}
//				newMessStr += picture.newMessByte
//			}
//			newMessStr += string(messByte[16:])
//			newMessByte = []byte(newMessStr)
//		default:
//			sign = false
//			picture := pictureMes{}
//			picture.PhoneNumAndSeq = header.PhoneNum + strconv.FormatInt(int64(header.SequenceReceive), 10)
//			picture.newMessByte = string(messByte[16:])
//			setPicture(picture.PhoneNumAndSeq, &picture)
//
//		}
//
//	}
//	switch string(header.MessProperty.Encryption) {
//	case "000":
//		fmt.Println("code no change")
//	default:
//		fmt.Println("code need change")
//	}
//	return newMessByte, sign, needRepeat
//}
//
////处理所有终端参数
//func DeviceParameters(messByte []byte, phoneNum string, data *server.AllParameterDataReceive) []byte {
//	lenth := messByte[4]
//	switch hex.EncodeToString(messByte[:4]) {
//	case "00000001":
//
//		data.KeepaliveInterval = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("心跳保活间隔 %d 秒\n", data.KeepaliveInterval)
//
//	case "00000002":
//
//		data.TCPOverTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("TCP 消息应答超时时间 %d s\n", data.TCPOverTime)
//
//	case "00000003":
//
//		data.TCPRepeatNum = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("TCP 消息重传次数 %d \n", data.TCPRepeatNum)
//
//	case "00000004":
//
//		data.UDPOverTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("UDP 消息应答超时时间 %d s\n", data.UDPOverTime)
//
//	case "00000005":
//
//		data.UDPRepeatNum = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("UDP 消息重传次数 %v \n", data.UDPRepeatNum)
//
//	case "00000006":
//
//		data.SMSOverTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("SMS 消息应答超时时间 %d s\n", data.SMSOverTime)
//
//	case "00000007":
//
//		data.SMSRepeatNum = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("SMS 消息重传次数 %v \n", data.SMSRepeatNum)
//	case "00000010":
//
//		data.MainServer.MainServerAPNOrPPP = string(messByte[5 : 5+lenth])
//		fmt.Printf("主服务器 APN %v \n", data.MainServer.MainServerAPNOrPPP)
//	case "00000011":
//
//		data.MainServer.MainServerUsername = string(messByte[5 : 5+lenth])
//		fmt.Printf("主服务器无线通信拨号用户名 %v \n", data.MainServer.MainServerUsername)
//	case "00000012":
//
//		data.MainServer.MainServerPassword = string(messByte[5 : 5+lenth])
//		fmt.Printf("主服务器无线通信拨号密码 %v \n", data.MainServer.MainServerPassword)
//	case "00000013":
//
//		data.MainServer.MainServerIP = string(messByte[5 : 5+lenth])
//		fmt.Printf("主服务器地址 %v \n", data.MainServer.MainServerIP)
//	case "00000014":
//
//		data.BackupServer.BackupServerAPNOrPPP = string(messByte[5 : 5+lenth])
//		fmt.Printf("备份服务器 APN %v \n", data.BackupServer.BackupServerAPNOrPPP)
//	case "00000015":
//
//		data.BackupServer.BackupServerUsername = string(messByte[5 : 5+lenth])
//		fmt.Printf("备份服务器无线通信拨号用户名 %v \n", data.BackupServer.BackupServerUsername)
//	case "00000016":
//
//		data.BackupServer.BackupServerPassword = string(messByte[5 : 5+lenth])
//		fmt.Printf("备份服务器无线通信拨号密码 %v \n", data.BackupServer.BackupServerPassword)
//	case "00000017":
//
//		data.BackupServer.BackupServerIP = string(messByte[5 : 5+lenth])
//		fmt.Printf("备份服务器地址 %v \n", data.BackupServer.BackupServerIP)
//	case "00000018":
//
//		data.TCPPort = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("服务器 TCP 端口 %d \n", data.TCPPort)
//
//	case "00000019":
//
//		data.UDPPort = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("服务器 UDP 端口 %v \n", data.UDPPort)
//
//	case "0000001a":
//
//		data.ICMainServerIP = string(messByte[5 : 5+lenth])
//		fmt.Printf("道路运输证 IC 卡认证主服务器 IP 地址或域名 %v \n", data.ICMainServerIP)
//
//	case "0000001b":
//
//		data.ICTCPPort = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("道路运输证 IC 卡认证主服务器 TCP 端口 %d \n", data.ICTCPPort)
//
//	case "0000001c":
//
//		data.ICUDPPort = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("道路运输证 IC 卡认证主服务器 UDP 端口 %v \n", data.ICUDPPort)
//
//	case "0000001d":
//
//		data.ICBackupServerIP = string(messByte[5 : 5+lenth])
//		fmt.Printf("道路运输证 IC 认证备份服务器 IP 地址或域名，端口同主服\n务器 %v \n", data.ICBackupServerIP)
//
//	case "00000020":
//
//		data.LocationReportStrategy = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		strategy := [3]string{"定时汇报", "定距汇报", "定时定距汇报"}
//		fmt.Printf("位置汇报策略 %v \n", strategy[data.LocationReportStrategy])
//
//	case "00000021":
//
//		data.LocationReportScheme = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		strategy := [2]string{"根据 ACC 状态", "根据登录状态和 ACC状态"}
//		fmt.Printf("位置汇报方案 %v \n", strategy[data.LocationReportScheme])
//
//	case "00000023":
//
//		data.AccompanyServer.AccompanyServerAPNOrPPP = string(messByte[5 : 5+lenth])
//		fmt.Printf("从服务器 APN %v \n", data.AccompanyServer.AccompanyServerAPNOrPPP)
//	case "00000024":
//
//		data.AccompanyServer.AccompanyServerUsername = string(messByte[5 : 5+lenth])
//		fmt.Printf("从服务器无线通信拨号用户名 %v \n", data.AccompanyServer.AccompanyServerUsername)
//	case "00000025":
//
//		data.AccompanyServer.AccompanyServerPassword = string(messByte[5 : 5+lenth])
//		fmt.Printf("从服务器无线通信拨号密码 %v \n", data.AccompanyServer.AccompanyServerPassword)
//	case "00000026":
//
//		data.AccompanyServer.AccompanyServerAddr = string(messByte[5 : 5+lenth])
//		fmt.Printf("从服务器地址 %v \n", data.AccompanyServer.AccompanyServerAddr)
//
//	case "00000022":
//
//		data.ReportTimeInterval.DriverReport = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("驾驶员未登录汇报时间间隔 %v s\n", data.ReportTimeInterval.DriverReport)
//	case "00000027":
//
//		data.ReportTimeInterval.SleepReport = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("休眠时汇报时间间隔 %v s\n", data.ReportTimeInterval.SleepReport)
//	case "00000028":
//
//		data.ReportTimeInterval.AlarmReport = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("紧急报警汇报时间间隔 %v s\n", data.ReportTimeInterval.AlarmReport)
//	case "00000029":
//
//		data.ReportTimeInterval.DefaultReport = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("缺省事件汇报间隔 %v s\n", data.ReportTimeInterval.DefaultReport)
//
//	case "0000002c":
//
//		data.ReportDistanceInterval.DefaultReport = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("缺省事件汇报距离间隔 %v m\n", data.ReportDistanceInterval.DefaultReport)
//	case "0000002d":
//
//		data.ReportDistanceInterval.DriverReport = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("驾驶员未登录汇报距离间隔 %v m\n", data.ReportDistanceInterval.DriverReport)
//	case "0000002e":
//
//		data.ReportDistanceInterval.SleepReport = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("休眠时汇报距离间隔 %v m\n", data.ReportDistanceInterval.SleepReport)
//	case "0000002f":
//
//		data.ReportDistanceInterval.AlarmReport = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("紧急报警汇报距离间隔 %v m\n", data.ReportDistanceInterval.AlarmReport)
//
//	case "00000030":
//
//		data.InflectionPointAngle = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("拐点补传角度 %v 度\n", data.InflectionPointAngle)
//	case "00000031":
//
//		data.IllegalDisplacement = binary.BigEndian.Uint16(messByte[5 : 5+lenth])
//		fmt.Printf("电子围栏半径（非法位移阀值） %v 米\n", data.IllegalDisplacement)
//	case "00000032":
//
//		data.IllegalDrivingPeriod.Start = strconv.Itoa(int(messByte[5])) + ":" + strconv.Itoa(int(messByte[6]))
//		data.IllegalDrivingPeriod.End = strconv.Itoa(int(messByte[7])) + ":" + strconv.Itoa(int(messByte[8]))
//		fmt.Printf("开始时间 %v 结束时间 %v\n", data.IllegalDrivingPeriod.Start, data.IllegalDrivingPeriod.End)
//
//	case "00000040":
//
//		data.PhoneNum.Monitor = string(messByte[5 : 5+lenth])
//		fmt.Printf("监控平台电话号码 %v \n", data.PhoneNum.Monitor)
//	case "00000041":
//
//		data.PhoneNum.Reset = string(messByte[5 : 5+lenth])
//		fmt.Printf("复位电话号码 %v \n", data.PhoneNum.Reset)
//	case "00000042":
//
//		data.PhoneNum.DefaultSetting = string(messByte[5 : 5+lenth])
//		fmt.Printf("恢复出厂设置电话号码 %v \n", data.PhoneNum.DefaultSetting)
//	case "00000043":
//
//		data.PhoneNum.MonitorSMS = string(messByte[5 : 5+lenth])
//		fmt.Printf("监控平台 SMS 电话号码 %v \n", data.PhoneNum.MonitorSMS)
//	case "00000044":
//
//		data.PhoneNum.DeviceSMS = string(messByte[5 : 5+lenth])
//		fmt.Printf("接收终端 SMS 文本报警号码 %v \n", data.PhoneNum.DeviceSMS)
//
//	case "00000045":
//
//		data.DevicePhoneListenStrategy = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		strategy := [2]string{"自动接听", "ACC ON 时自动接听，OFF 时手动接听"}
//		fmt.Printf("终端电话接听策略 %v \n", strategy[data.DevicePhoneListenStrategy])
//	case "00000046":
//
//		data.OnesTalkTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("每次最长通话时间 %v s\n", data.OnesTalkTime)
//	case "00000047":
//
//		data.MonthTalkTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("当月最长通话时间 %v s\n", data.MonthTalkTime)
//
//	case "00000048":
//
//		data.PhoneNum.Listen = string(messByte[5 : 5+lenth])
//		fmt.Printf("监听电话号码 %v \n", data.PhoneNum.Listen)
//	case "00000049":
//
//		data.PhoneNum.MonitorPrivilege = string(messByte[5 : 5+lenth])
//		fmt.Printf("监管平台特权短信号码 %v \n", data.PhoneNum.MonitorPrivilege)
//
//	case "00000050":
//
//		alarmShielding := binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		data.AlarmShielding = tools.UintTo2byte(uint64(alarmShielding), 32)
//		fmt.Printf("报警屏蔽字 %v \n", string(data.AlarmShielding))
//	case "00000051":
//
//		AlarmSMSTxt := binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		data.AlarmSMSTxt = tools.UintTo2byte(uint64(AlarmSMSTxt), 32)
//		fmt.Printf("报警发送文本 SMS 开关 %v \n", string(data.AlarmSMSTxt))
//	case "00000052":
//
//		AlarmShooting := binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		data.AlarmShooting = tools.UintTo2byte(uint64(AlarmShooting), 32)
//		fmt.Printf("报警拍摄开关 %v \n", string(data.AlarmShooting))
//	case "00000053":
//
//		AlarmPhotoSave := binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		data.AlarmPhotoSave = tools.UintTo2byte(uint64(AlarmPhotoSave), 32)
//		fmt.Printf("报警拍摄存储标志 %v \n", string(data.AlarmPhotoSave))
//	case "00000054":
//
//		KeyAlarm := binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		data.KeyAlarm = tools.UintTo2byte(uint64(KeyAlarm), 32)
//		fmt.Printf("关键标志 %v \n", string(data.KeyAlarm))
//
//	case "00000055":
//
//		data.HighestSpeed = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("最高速度 %v km/h\n", data.HighestSpeed)
//	case "00000056":
//
//		data.OverSpeedTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("超速持续时间 %v s\n", data.OverSpeedTime)
//	case "00000057":
//
//		data.KeepDriverTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("连续驾驶时间门限 %v s\n", data.KeepDriverTime)
//	case "00000058":
//
//		data.OneDayDriverTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("当天累计驾驶事件门限 %v s\n", data.OneDayDriverTime)
//	case "00000059":
//
//		data.LeastRestTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("最小休息时间 %v s\n", data.LeastRestTime)
//	case "0000005a":
//
//		data.LongestStopTime = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("最长停车时间 %v s\n", data.LongestStopTime)
//	case "0000005b":
//
//		data.OverSpeedWarningDifference = binary.BigEndian.Uint16(messByte[5:5+lenth]) / 10
//		fmt.Printf("超速报警预警差值 %v km/h\n", data.OverSpeedWarningDifference)
//	case "0000005c":
//
//		data.FatigueDriverWarningDifference = binary.BigEndian.Uint16(messByte[5 : 5+lenth])
//		fmt.Printf("疲劳驾驶预警差值 %v s\n", data.FatigueDriverWarningDifference)
//	case "0000005d":
//
//		data.CollisionAlarmParameters.Collision = messByte[6]
//		data.CollisionAlarmParameters.CollisionAcceleration = messByte[5]
//		fmt.Printf("疲劳驾驶预警差值:碰撞事件 %v 单位 4ms,碰撞加速度 %v 单位0.1g\n",
//			data.CollisionAlarmParameters.Collision, data.CollisionAlarmParameters.CollisionAcceleration)
//	case "0000005e":
//
//		data.RolloverAlarmAngle = binary.BigEndian.Uint16(messByte[5 : 5+lenth])
//		fmt.Printf("侧翻报警参数设置 %v 度\n", data.RolloverAlarmAngle)
//
//	case "00000064":
//
//		timingCameraControl := binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		TimingCameraControl2 := tools.UintTo2byte(uint64(timingCameraControl), 32)
//		data.TimingCameraControl.Interval = tools.TwoToInt64(string(TimingCameraControl2[:15]))
//		data.TimingCameraControl.OtherSign = TimingCameraControl2[15:]
//		fmt.Printf("定时拍照控制,定时事件间隔 %v , 其他标志 %v\n", string(data.TimingCameraControl.Interval), string(data.TimingCameraControl.OtherSign))
//	case "00000065":
//
//		DistanceCameraControl := binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		DistanceCameraControl2 := tools.UintTo2byte(uint64(DistanceCameraControl), 32)
//		data.DistanceCameraControl.Interval = tools.TwoToInt64(string(DistanceCameraControl2[:15]))
//		data.DistanceCameraControl.OtherSign = DistanceCameraControl2[15:]
//		fmt.Printf("定距拍照控制,定时事件间隔 %v , 其他标志 %v\n", string(data.DistanceCameraControl.Interval), string(data.DistanceCameraControl.OtherSign))
//
//	case "00000070":
//
//		data.ImageOrVideoInstruction = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("图像/视频指令 %v \n", data.ImageOrVideoInstruction)
//	case "00000071":
//
//		data.Brightness = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("亮度 %v \n", data.Brightness)
//	case "00000072":
//
//		data.ContrastRatio = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("对比度 %v \n", data.ContrastRatio)
//	case "00000073":
//
//		data.Saturation = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("饱和度 %v \n", data.Saturation)
//	case "00000074":
//
//		data.Chroma = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("色度 %v \n", data.Chroma)
//
//	case "00000075":
//		data.AudioAndVideoParameters.RealTimeStreamCoding = messByte[5]
//		data.AudioAndVideoParameters.RealTimeStreamResolutionRatio = messByte[6]
//		data.AudioAndVideoParameters.RealTimeStreamKeyFrameInterval = binary.BigEndian.Uint16(messByte[7:9])
//		data.AudioAndVideoParameters.RealTimeStreamTargetFrame = messByte[9]
//		data.AudioAndVideoParameters.RealTimeStreamTargetKbps = binary.BigEndian.Uint32(messByte[10:14])
//		data.AudioAndVideoParameters.StorageStreamCoding = messByte[14]
//		data.AudioAndVideoParameters.StorageStreamResolutionRatio = messByte[15]
//		data.AudioAndVideoParameters.StorageStreamKeyFrameInterval = binary.BigEndian.Uint16(messByte[16:18])
//		data.AudioAndVideoParameters.StorageStreamTargetFrame = messByte[18]
//		data.AudioAndVideoParameters.StorageStreamTargetKbps = binary.BigEndian.Uint32(messByte[19:23])
//		OSDSubtitleOverlaySettings := binary.BigEndian.Uint16(messByte[23:25])
//		data.AudioAndVideoParameters.OSDSubtitleOverlaySettings = tools.UintTo2byte(uint64(OSDSubtitleOverlaySettings), 16)
//		data.AudioAndVideoParameters.UseAudioOrNot = messByte[25]
//		fmt.Printf("音视频参数%v\n", data.AudioAndVideoParameters)
//	case "00000076":
//		data.AudioAndVideoChannelList.AudioAndVideoChannelCount = messByte[5]
//		data.AudioAndVideoChannelList.AudioChannelCount = messByte[6]
//		data.AudioAndVideoChannelList.VideoChannelCount = messByte[7]
//		num := (lenth - 3) / 4
//		channelComparisonTable := server.ChannelComparisonTableData{}
//		for i := 0; i < int(num); i++ {
//			channelComparisonTable.PhysicalChannelID = messByte[i*4+8]
//			channelComparisonTable.LogicalChannelID = messByte[i*4+9]
//			channelComparisonTable.ChannelType = messByte[i*4+10]
//			channelComparisonTable.ConnectedToPTZ = messByte[i*4+11]
//			data.AudioAndVideoChannelList.ChannelComparisonTable = append(data.AudioAndVideoChannelList.ChannelComparisonTable, channelComparisonTable)
//		}
//		fmt.Printf("通道参数%v\n", data.AudioAndVideoChannelList)
//
//	case "00000077":
//		data.SingleChannelVideoParameter.ChannelCount = messByte[5]
//		num := (lenth - 1) / 21
//		singleVideoParametersData := server.SingleVideoParametersData{}
//		for i := 0; i < int(num); i++ {
//			singleVideoParametersData.ChannelID = messByte[21*i+6]
//			singleVideoParametersData.RealTimeStreamCoding = messByte[21*i+7]
//			singleVideoParametersData.RealTimeStreamResolutionRatio = messByte[21*i+8]
//			singleVideoParametersData.RealTimeStreamKeyFrameInterval = binary.BigEndian.Uint16(messByte[21*i+9 : 21*i+11])
//			singleVideoParametersData.RealTimeStreamTargetFrame = messByte[21*i+11]
//			singleVideoParametersData.RealTimeStreamTargetKbps = binary.BigEndian.Uint32(messByte[21*i+12 : 21*i+16])
//			singleVideoParametersData.StorageStreamCoding = messByte[21*i+16]
//			singleVideoParametersData.StorageStreamResolutionRatio = messByte[21*i+17]
//			singleVideoParametersData.StorageStreamKeyFrameInterval = binary.BigEndian.Uint16(messByte[21*i+18 : 21*i+20])
//			singleVideoParametersData.StorageStreamTargetFrame = messByte[21*i+20]
//			singleVideoParametersData.StorageStreamTargetKbps = binary.BigEndian.Uint32(messByte[21*i+21 : 21*i+25])
//			OSDSubtitleOverlaySettings := binary.BigEndian.Uint16(messByte[21*i+25 : 21*i+27])
//			singleVideoParametersData.OSDSubtitleOverlaySettings = tools.UintTo2byte(uint64(OSDSubtitleOverlaySettings), 16)
//			data.SingleChannelVideoParameter.SingleVideoParameters = append(data.SingleChannelVideoParameter.SingleVideoParameters, singleVideoParametersData)
//		}
//		fmt.Printf("单独通道视频参数%v\n", data.SingleChannelVideoParameter)
//	case "00000079":
//		data.SpecialAlarmRecording.Threshold = messByte[5]
//		data.SpecialAlarmRecording.Duration = messByte[6]
//		data.SpecialAlarmRecording.StartTime = messByte[7]
//	case "0000007a":
//		VideoAlarmScreenWord := binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		data.VideoAlarmScreenWord = tools.UintTo2byte(uint64(VideoAlarmScreenWord), 32)
//		fmt.Printf("视频相关报警屏蔽字%v\n", string(data.VideoAlarmScreenWord))
//	case "0000007b":
//		data.VideoAnalysisAlarm.LoadNum = messByte[5]
//		data.VideoAnalysisAlarm.FatigueThreshold = messByte[6]
//		fmt.Printf("视频分析报警参数定义及说明%v\n", data.VideoAnalysisAlarm)
//	case "0000007c":
//		data.DeviceWakeupMode.WakeupMode = messByte[5]
//		data.DeviceWakeupMode.WakeupCondition = messByte[6]
//		data.DeviceWakeupMode.TimingWakeupDay = messByte[7]
//		data.DeviceWakeupMode.UseTimingWakeup = messByte[8]
//		data.DeviceWakeupMode.TimeStart1 = hex.EncodeToString(messByte[9:10]) + ":" + hex.EncodeToString(messByte[10:11])
//		data.DeviceWakeupMode.TimeEnd1 = hex.EncodeToString(messByte[11:12]) + ":" + hex.EncodeToString(messByte[12:13])
//		data.DeviceWakeupMode.TimeStart2 = hex.EncodeToString(messByte[13:14]) + ":" + hex.EncodeToString(messByte[14:15])
//		data.DeviceWakeupMode.TimeEnd2 = hex.EncodeToString(messByte[15:16]) + ":" + hex.EncodeToString(messByte[16:17])
//		data.DeviceWakeupMode.TimeStart3 = hex.EncodeToString(messByte[17:18]) + ":" + hex.EncodeToString(messByte[18:19])
//		data.DeviceWakeupMode.TimeEnd3 = hex.EncodeToString(messByte[19:20]) + ":" + hex.EncodeToString(messByte[20:21])
//		data.DeviceWakeupMode.TimeStart4 = hex.EncodeToString(messByte[21:22]) + ":" + hex.EncodeToString(messByte[22:23])
//		data.DeviceWakeupMode.TimeEnd4 = hex.EncodeToString(messByte[23:24]) + ":" + hex.EncodeToString(messByte[24:25])
//		fmt.Printf("终端休眠唤醒模式设置数据%v\n", data.DeviceWakeupMode)
//
//	case "00000080":
//
//		data.Odometer = binary.BigEndian.Uint32(messByte[5:5+lenth]) / 10
//		fmt.Printf("车辆里程表读数 %v km\n", data.Odometer)
//
//	case "00000081":
//
//		data.ProvinceID = binary.BigEndian.Uint16(messByte[5 : 5+lenth])
//		fmt.Printf("车辆所在的省域 ID %v \n", data.ProvinceID)
//
//	case "00000082":
//
//		data.CityID = binary.BigEndian.Uint16(messByte[5 : 5+lenth])
//		fmt.Printf("车辆所在的市域 ID %v \n", data.CityID)
//	case "00000083":
//
//		data.VehicleID = string(messByte[5 : 5+lenth])
//		fmt.Printf("公安交通管理部门颁发的机动车号牌 %v \n", data.VehicleID)
//	case "00000084":
//
//		data.VehiclePlateColor = messByte[5]
//		fmt.Printf("车牌颜色 %v \n", data.VehiclePlateColor)
//
//	case "00000090":
//
//		GNSSPositioningMode := uint64(messByte[5])
//		data.GNSSPositioningMode = tools.UintTo2byte(GNSSPositioningMode, 8)
//		GNSSPositioningModeMes := [8]string{"", "", "", "", "galileo定位", "GLONASS 定位", "北斗定位", "GPS定位"}
//		for i, v := range data.GNSSPositioningMode {
//			if v == '1' {
//				fmt.Printf("GNSS 定位模式 %v \n", GNSSPositioningModeMes[i])
//			}
//		}
//	case "00000091":
//
//		data.GNSSBps = messByte[5]
//		GNSSBpsMes := [6]int64{4800, 9600, 19200, 38400, 57600, 115200}
//		fmt.Printf("GNSS 波特率 %v \n", GNSSBpsMes[data.GNSSBps])
//	case "00000092":
//
//		data.GNSSOutputFrequency = messByte[5]
//		GNSSBpsMes := [5]int64{500, 1000, 2000, 3000, 4000}
//		fmt.Printf("GNSS 模块详细定位数据输出频率 %v ms\n", GNSSBpsMes[data.GNSSOutputFrequency])
//	case "00000093":
//
//		data.GNSSAcquisitionFrequency = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("GNSS 模块详细定位数据采集频率 %v s\n", data.GNSSAcquisitionFrequency)
//	case "00000094":
//
//		data.GNSSUploadMethod = messByte[5]
//		GNSSBpsMes := [16]string{"本地存储", "按时间间隔上传", "按距离间隔上传", "", "", "", "", "",
//			"", "", "", "按累计时间上传", "按累计距离上传", "按累计条数上传", "", ""}
//		fmt.Printf("GNSS 模块详细定位数据上传方式 %v \n", GNSSBpsMes[data.GNSSUploadMethod])
//	case "00000095":
//
//		data.GNSSUploadSettings = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		GNSSBpsMes := [16]string{"", "秒", "米", "", "", "", "", "",
//			"", "", "", "秒", "米", "条", "", ""}
//		fmt.Printf("GNSS 模块详细定位数据上传设置 %v ,单位 %v \n", data.GNSSUploadSettings, GNSSBpsMes[data.GNSSUploadMethod])
//
//	case "00000100":
//
//		data.CANAcquisitionInterval1 = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("CAN 总线通道 1 采集时间间隔 %v ms \n", data.CANAcquisitionInterval1)
//	case "00000101":
//
//		data.CANUploadInterval1 = binary.BigEndian.Uint16(messByte[5 : 5+lenth])
//		fmt.Printf("CAN 总线通道 1 上传时间间隔 %v s \n", data.CANUploadInterval1)
//	case "00000102":
//
//		data.CANAcquisitionInterval2 = binary.BigEndian.Uint32(messByte[5 : 5+lenth])
//		fmt.Printf("CAN 总线通道 2 采集时间间隔 %v ms \n", data.CANAcquisitionInterval2)
//	case "00000103":
//
//		data.CANUploadInterval2 = binary.BigEndian.Uint16(messByte[5 : 5+lenth])
//		fmt.Printf("CAN 总线通道 2 上传时间间隔 %v s \n", data.CANUploadInterval2)
//	case "00000110":
//		CANIDCollectionSettings := binary.BigEndian.Uint64(messByte[5 : 5+lenth])
//		CANIDCollectionSettings2 := tools.UintTo2byte(CANIDCollectionSettings, 64)
//		data.CANIDCollectionSettings.AcquisitionInterval = binary.BigEndian.Uint32(messByte[5:9])
//		data.CANIDCollectionSettings.Channel = CANIDCollectionSettings2[32]
//		data.CANIDCollectionSettings.FrameType = CANIDCollectionSettings2[33]
//		data.CANIDCollectionSettings.CollectionMode = CANIDCollectionSettings2[34]
//		data.CANIDCollectionSettings.ID = tools.TwoToInt64(string(CANIDCollectionSettings2[35:]))
//		fmt.Printf("CAN 总线 ID 单独采集设置 %v  \n", data.CANIDCollectionSettings)
//
//	default:
//		id := binary.BigEndian.Uint32(messByte[:4])
//		if 0x0111 <= id && id <= 0x01ff {
//			settings := server.CANIDCollectionSettingsData{}
//			CANIDCollectionSettings := binary.BigEndian.Uint64(messByte[5 : 5+lenth])
//			CANIDCollectionSettings2 := tools.UintTo2byte(CANIDCollectionSettings, 64)
//			settings.AcquisitionInterval = binary.BigEndian.Uint32(messByte[5:9])
//			settings.Channel = CANIDCollectionSettings2[32]
//			settings.FrameType = CANIDCollectionSettings2[33]
//			settings.CollectionMode = CANIDCollectionSettings2[34]
//			settings.ID = tools.TwoToInt64(string(CANIDCollectionSettings2[35:]))
//			data.OtherCANIDCollectionSettings = append(data.OtherCANIDCollectionSettings, settings)
//			fmt.Printf("其他CAN 总线 ID 单独采集设置 %v  \n", data.CANIDCollectionSettings)
//		} else {
//			custom := server.CustomData1{}
//			custom.ID = hex.EncodeToString(messByte[:4])
//			if lenth > 0 {
//				custom.Data = hex.EncodeToString(messByte[5 : 5+lenth])
//			} else {
//				custom.Data = "no mes"
//			}
//
//			data.Custom = append(data.Custom, custom)
//			fmt.Printf("自定义 id%v data %v\n", custom.ID, custom.Data)
//		}
//
//	}
//
//	return messByte[5+lenth:]
//}
//
////处理音视频资源列表
//func VideoListInfoHandle(messByte []byte, videoListData *server.VideoListData) []byte {
//	videoInfo := server.VideoListInfo{}
//	videoInfo.ChannelID = messByte[0]
//	startTime := hex.EncodeToString(messByte[1:7])
//	videoInfo.StartTime = startTime[:2] + "-" + startTime[2:4] + "-" + startTime[4:6] + " " + startTime[6:8] + ":" + startTime[8:10] + ":" + startTime[10:12]
//	endTime := hex.EncodeToString(messByte[7:13])
//	videoInfo.EndTime = endTime[:2] + "-" + endTime[2:4] + "-" + endTime[4:6] + " " + endTime[6:8] + ":" + endTime[8:10] + ":" + endTime[10:12]
//	alarmType := binary.BigEndian.Uint64(messByte[13:21])
//	videoInfo.AlarmType = string(tools.UintTo2byte(alarmType, 64))
//	videoInfo.MediaType = messByte[21]
//	videoInfo.Bitstream = messByte[22]
//	videoInfo.StorageType = messByte[23]
//	size := float64(binary.BigEndian.Uint32(messByte[24:28]))
//	sizeM := size / 1024 / 1024
//	sizeStr := strconv.FormatFloat(sizeM, 'f', 3, 64)
//	//float保留两位小数
//	videoInfo.Size, _ = strconv.ParseFloat(sizeStr, 64)
//	videoListData.VideoList = append(videoListData.VideoList, videoInfo)
//	return messByte[28:]
//}
//
////获取实时视频播放URL，第一次内存中无此URL，向设备发送start指令，等设备返回status=200时向STS服务发送POST请求URL
////当内存中有URL时，直接向STS服务发送POST请求URL，不直接返回是防止URL地址因为STS服务器断开等原因而失效
//func GetSessionURLFromSTS(phoneAndChannel string) string {
//	client := &http.Client{
//		Transport: &http.Transport{
//			Dial: func(netw, addr string) (net.Conn, error) {
//				conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
//				if err != nil {
//					return nil, err
//				}
//				conn.SetDeadline(time.Now().Add(time.Second * 3)) //设置发送接受数据超时
//				return conn, nil
//			},
//			ResponseHeaderTimeout: time.Second * 4,
//		},
//	}
//	data := struct {
//		SessionID string `json:"sessionId"`
//	}{
//		SessionID: phoneAndChannel,
//	}
//
//	sessionByte, err := json.MarshalIndent(data, "  ", "    ")
//	if err != nil {
//		logs.BeeLogger.Error("getSessionURLFromSTS() => xml.MarshalIndent() error: %s", err)
//		return ""
//	}
//
//	request, err := http.NewRequest("POST", tools.StringsJoin("http://", config.UrlAddr, "/jt/StartStream"), bytes.NewReader(sessionByte))
//	if err != nil {
//		logs.BeeLogger.Error("getSessionURLFromSTS() => http.NewRequest() error: %s", err)
//		return ""
//	}
//
//	//增加Header选项
//	request.Header.Set("Content-Type", "application/json")
//
//	response, err := client.Do(request)
//	if err != nil {
//		logs.BeeLogger.Error("getSessionURLFromSTS() => http.Do() error: %s", err)
//		return ""
//	}
//	defer response.Body.Close()
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		logs.BeeLogger.Error("getSessionURLFromSTS() => ioutil.ReadAll() error: %s", err)
//		return ""
//	}
//	logs.BeeLogger.Debug("getSessionURLFromSTS() body is: %s", string(body))
//	bodyData := struct {
//		ErrorCode    int64  `json:"errorCode"`
//		ErrorCodeStr string `json:"errorCodeStr"`
//		PushUrl      string `json:"pushUrl"`
//	}{}
//	err = json.Unmarshal(body, &bodyData)
//	if err != nil {
//		logs.BeeLogger.Error("getSessionURLFromSTS() => json.Unmarshal() error: %s", err)
//		return ""
//	}
//	if bodyData.ErrorCode == 0 {
//		return bodyData.PushUrl
//	}
//
//	return ""
//}
