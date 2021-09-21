package udp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gitee.com/ictt/JTTM/server"
	"gitee.com/ictt/JTTM/server/ftpclient"
	"gitee.com/ictt/JTTM/server/sertools"
	"gitee.com/ictt/JTTM/tools"
	"gitee.com/ictt/JTTM/tools/logs"
	"gitee.com/ictt/JTTM/tools/sqlDB"
	"github.com/gocsv"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)



//注册回复
func ResRegisterMes(messByte []byte, udpConn *net.UDPConn, udpAddr *net.UDPAddr) {
	messHead10_int := int64(binary.BigEndian.Uint16(messByte[2:4]))
	messHead2 := strconv.FormatInt(messHead10_int, 2)
	vehicleInfo := sqlDB.GetVehicleInfo{}
	vehicleInfo.DeviceIP = udpAddr.IP.String()
	fmt.Println(messHead2)
	if len(messHead2) > 9 {
		for len(messHead2) < 16 {
			messHead2 = "0" + messHead2
		}
		return
	} else {
		vehicleInfo.PhoneNum = hex.EncodeToString(messByte[4:10])
		fmt.Println(vehicleInfo.PhoneNum)
		sequence10 := int64(binary.BigEndian.Uint16(messByte[10:12]))
		fmt.Println(sequence10)
		provinceID := int64(binary.BigEndian.Uint16(messByte[12:14]))
		countyID := int64(binary.BigEndian.Uint16(messByte[14:16]))
		provinceIDStr := strconv.FormatInt(provinceID, 10)
		for len(provinceIDStr) < 2 {
			provinceIDStr = "0" + provinceIDStr
		}
		countyIDStr := strconv.FormatInt(countyID, 10)
		for len(countyIDStr) < 4 {
			countyIDStr = "0" + countyIDStr
		}
		fmt.Println(provinceIDStr + countyIDStr)
		vehicleInfo.ProvinceID = provinceIDStr
		vehicleInfo.CountyID = countyIDStr

		vehicleInfo.ManufacturerID = string(messByte[16:21])

		//vehicleInfo.DeviceModel = string(messByte[21:41])
		//vehicleInfo.DeviceID = string(messByte[41:48])
		////fmt.Println(vehicleInfo.ManufacturerID,vehicleInfo.DeviceModel,vehicleInfo.DeviceID)
		////fmt.Println(messByte[36])
		//color := [8]string{"未上牌", "蓝色", "黄色", "黑色", "白色", "绿色", "黄绿色", "其他"}
		//vehicleInfo.VehiclePlateColor = color[messByte[48]]
		//vehicleInfo.VehicleID = string(messByte[49:])

		vehicleInfo.DeviceModel = string(messByte[21:29])
		vehicleInfo.DeviceID = string(messByte[29:36])
		color := [8]string{"未上牌", "蓝色", "黄色", "黑色", "白色", "绿色", "黄绿色", "其他"}
		vehicleInfo.VehiclePlateColor = color[messByte[36]]
		//vehicleInfo.VehicleID = string(messByte[37:])

		vehicleIDByte,err := tools.GBKToUTF8(messByte[37:])
		if err != nil {
			return
		}
		vehicleInfo.VehicleID = string(vehicleIDByte)
		fmt.Printf("车牌号 %s \n", vehicleInfo.VehicleID)
		//fmt.Println(vehicleInfo.VehicleID)
		vehicleInfo.Status = "ON"
		vehicleInfo.PowerIdentify = tools.GetRandomString(10)
		vehicleInfo.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

		//回复数据编写
		messIDSend, _ := hex.DecodeString("8100")
		//messIDSend := []byte{129
		var sequence uint16 = 1
		sequenceSend := tools.Uint16ToByte(sequence)
		messLength := len(vehicleInfo.PowerIdentify) + 3
		messProperty := tools.Int64ToByte(int64(messLength))
		sendMes := append(messIDSend, messProperty[6:]...)
		sendMes = append(sendMes, messByte[4:10]...)
		sendMes = append(sendMes, sequenceSend...)
		sendMes = append(sendMes, messByte[10:12]...)
		var result byte = 0
		//vehicleList, retBool := sqlDB.QueryFindVehicle()
		//if !retBool {
		//	return
		//}
		//for _, v :=range vehicleList {
		//	if v.DeviceID == vehicleInfo.DeviceID && v.DeviceModel == vehicleInfo.DeviceModel {
		//		result = 3
		//	}
		//	if v.CountyID == vehicleInfo.CountyID && v.VehicleID == vehicleInfo.VehicleID{
		//		result = 1
		//	}
		//}
		sendMes = append(sendMes, result)
		sendMes = append(sendMes, []byte(vehicleInfo.PowerIdentify)...)

		resultSendMes := tools.ParaphraseMess(sendMes)
		if sqlDB.Save(vehicleInfo) {
			tabName := sqlDB.GetTableName(&sqlDB.GetVehicleInfo{})
			fmt.Printf("deviceID=%s saved record to %s's table successfully\n", vehicleInfo.PhoneNum, tabName)
		}


		vehicleList := []*sqlDB.GetVehicleInfo{}
		vehicleList = append(vehicleList, &vehicleInfo) // Add clients
		csvContent, _ := gocsv.MarshalString(&vehicleList) // Get all clients as CSV string
		fmt.Println(csvContent) // Display all clients as CSV string
		logs.BeeLogger.Info("csv info %s", csvContent)
		//data1 := bytes.NewBufferString(csvContent)
		csvName := "device/vehicle_" + vehicleInfo.PhoneNum + ".csv"
		msg := ftpclient.Message{
			FileName: csvName,
			SaveMsg: csvContent,
		}
		ftpclient.FtpMsg <- msg
		//ftpclient.FtpCli.Lock()
		//err := ftpclient.FtpCli.FTPClient.Stor(csvName, data1)
		//ftpclient.FtpCli.Unlock()
		//if err != nil {
		//	fmt.Println("vehicle info to ftp err")
		//}

		fmt.Println(hex.EncodeToString(resultSendMes))
		_, err = udpConn.WriteToUDP(resultSendMes, udpAddr)
		if err != nil {
			fmt.Println("connect error")
			return
		}

		//go func() {
		//	getTemporaryInfo := &sqlDB.GetTemporaryInfo{
		//		PhoneNum:   	vehicleInfo.PhoneNum,
		//		DeviceAddr:   	udpAddr.String(),
		//		Sequence: 		1,
		//		Status:     	"ON",
		//	}
		//
		//	if sqlDB.Save(getTemporaryInfo) {
		//		tabName := sqlDB.GetTableName(&sqlDB.GetTemporaryInfo{})
		//		logs.BeeLogger.Info("PhoneNum=%s saved record to %s's table successfully", getTemporaryInfo.PhoneNum, tabName)
		//	}
		//}()

	}
	fmt.Println("注册登记完成")
	return
}

//鉴权回复
func ResPowerIdentify(messByte []byte, udpConn *net.UDPConn, udpAddr *net.UDPAddr) {
	fmt.Println("收到鉴权请求")
	powerIdentifyMes := server.PowerIdentifyMesReceive{}
	NewMessByte, sign, _ := sertools.HandleMessProperty(&powerIdentifyMes.Header, messByte)
	if !sign {
		fmt.Println("have many package")
		return
	}
	var resultSendMes []byte

	powerIdentifyMes.Data = string(NewMessByte)

	var sequence uint16 = 2
	//获取数据库中的鉴权码
	vehicleInfo := new(sqlDB.GetVehicleInfo)
	if !sqlDB.QueryUserTake(vehicleInfo, map[string]interface{}{"PhoneNum": powerIdentifyMes.Header.PhoneNum}) {
		resultSendMes = sertools.NormalResponse(messByte, sequence, 1)
		sequence++
	} else {
		if powerIdentifyMes.Data != vehicleInfo.PowerIdentify {
			fmt.Println("PowerIdentify error")
			resultSendMes = sertools.NormalResponse(messByte, sequence, 1)
			sequence++
		} else {
			//存储键值对
			resultSendMes = sertools.NormalResponse(messByte, sequence, 0)
			sequence++
			tempClient := &sertools.UdpClient{
				PhoneNum: powerIdentifyMes.Header.PhoneNum,
				UDPConn:  udpConn,
				UDPAddr:  udpAddr,
				Sequence: sequence,
				UpdateTime: time.Now().Unix(),
			}
			sertools.SetUDPClient(vehicleInfo.PhoneNum, tempClient)
			gpsAndPhone := new(sqlDB.GpsAndPhoneNum)
			gpsAndPhone.PhoneNum = powerIdentifyMes.Header.PhoneNum
			gpsAndPhone.GpsTableName = "gps_" + powerIdentifyMes.Header.PhoneNum
			gpsAndPhone.AlarmTableName = "alarm_" + powerIdentifyMes.Header.PhoneNum
			gpsAndPhone.DriverTableName = "drAct_" + powerIdentifyMes.Header.PhoneNum
			sqlDB.Save(gpsAndPhone)
			sqlDB.InitGetLocationInfo(gpsAndPhone.GpsTableName)
			sqlDB.InitGetAlarmInfo(gpsAndPhone.AlarmTableName)
			sqlDB.InitGetDriverInfo(gpsAndPhone.DriverTableName)

		}
	}

	//resultSendMes := normalResponse(messByte)
	fmt.Println(hex.EncodeToString(resultSendMes))
	_, err := udpConn.WriteToUDP(resultSendMes, udpAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}

	fmt.Println("鉴权回复成功")
	if powerIdentifyMes.Data == vehicleInfo.PowerIdentify {
		go func() {
			vehicleInfo := new(sqlDB.GetVehicleInfo)
			fool := sqlDB.QueryUserTake(vehicleInfo, map[string]interface{}{"PhoneNum": powerIdentifyMes.Header.PhoneNum})
			if !fool {
				logs.BeeLogger.Error("查询表格错误")
				return
			}
			vehicleInfo.Status = "ON"
			vehicleInfo.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			sqlDB.Save(vehicleInfo)

			//if server.RemoteClientStatus == "ON" {
			//	go func() {
			//		payload := server.MQTTNotify{}
			//		payload.Protocol = "NotifySub"
			//		payload.Csq = server.Csq
			//		server.Csq ++
			//		payload.Type = "Device"
			//		mqttDevice := server.MQTTDeviceInfo{}
			//		mqttDeviceList := []server.MQTTDeviceInfo{}
			//		mqttDevice.Id = vehicleInfo.PhoneNum
			//		mqttDevice.Name = vehicleInfo.VehicleID
			//		mqttDevice.Status = "ON"
			//
			//		mqttDeviceList = append(mqttDeviceList, mqttDevice)
			//
			//		payload.Data = map[string]interface{}{
			//			"devicelists": mqttDeviceList,
			//		}
			//		msg := server.Message{
			//			Topic:   server.MattPubID,
			//			Payload: payload,
			//		}
			//		server.MqttMess <- msg
			//	}()
			//}

		}()
		go func() {
			QueryStreamProperties(powerIdentifyMes.Header.PhoneNum)
			time.Sleep(time.Second * 10)
			QueryLocationInfo(powerIdentifyMes.Header.PhoneNum)
		}()
		go checkKeepalive(powerIdentifyMes.Header.PhoneNum, udpConn)
		go func() {
			client := sertools.GetTCPClient(powerIdentifyMes.Header.PhoneNum)
			if client != nil {
				//fmt.Println(1111111111)
				sertools.ChangeToUDP.Store(powerIdentifyMes.Header.PhoneNum, "change")
				sertools.DeleteTCPClient(powerIdentifyMes.Header.PhoneNum)
			}
		}()
	}
	return
}

//保活回复
func ResKeepAlive(messByte []byte, udpConn *net.UDPConn, udpAddr *net.UDPAddr) {
	messHead10_int := int64(binary.BigEndian.Uint16(messByte[2:4]))
	messHead2 := strconv.FormatInt(messHead10_int, 2)
	fmt.Println(messHead2)
	if len(messHead2) > 9 {
		for len(messHead2) < 16 {
			messHead2 = "0" + messHead2
		}
		return
	} else {
		PhoneNum := hex.EncodeToString(messByte[4:10])
		client := sertools.GetUDPClient(PhoneNum)
		if client == nil {
			logs.BeeLogger.Error("device disconnected, unable to keepalive")
			fmt.Println("device disconnected, unable to keepalive")
			return
		}

		fmt.Println(PhoneNum)
		sequence10 := int64(binary.BigEndian.Uint16(messByte[10:12]))
		fmt.Println(sequence10)
		//回复数据编写
		resultSendMes := sertools.NormalResponse(messByte, client.Sequence, 0)
		fmt.Println(hex.EncodeToString(resultSendMes))
		_, err := udpConn.WriteToUDP(resultSendMes, udpAddr)
		if err != nil {
			fmt.Println("connect error")
			return
		}
		fmt.Println("保活成功")
		//重置键值对
		client.Sequence++
		client.UpdateTime = time.Now().Unix()
		sertools.SetUDPClient(PhoneNum, client)
		return
	}
}

//心跳保活检测，用于判断蓝方设备是否在线
func checkKeepalive(phoneNum string, udpConn *net.UDPConn) {
	//IPAndPort := tools.StringsJoin(udpAddr.IP.String(), ":", strconv.Itoa(udpAddr.Port))
	for {
		time.Sleep(15 * time.Second)
		_, ok := sertools.ChangeToTCP.Load(phoneNum)
		//fmt.Println(ok)
		if ok {
			logs.BeeLogger.Info("connect change to tcp")
			fmt.Println("connect change to tcp")
			sertools.ChangeToTCP.Delete(phoneNum)
			return
		}

		client := sertools.GetUDPClient(phoneNum)
		if client == nil {
			//心跳超时，设备断开连接
			fmt.Printf("keepalive timeout! deviceID=%s disconnect\n", phoneNum)
			sqlDB.UpdateTableFromPhone(sqlDB.GetVehicleInfo{}, "OFF", time.Now().Format("2006-01-02 15:04:05"), phoneNum)
			sqlDB.UpdateTableFromPhone(sqlDB.GetChannelInfo{}, "OFF", time.Now().Format("2006-01-02 15:04:05"), phoneNum)

			return
		}

		if udpConn != client.UDPConn {
			fmt.Println("重新鉴权了")
			logs.BeeLogger.Error("重新鉴权了")
			return
		}

	}
}

//终端主动上传的定位信息回复
func ResLocateMes(messByte []byte, udpConn *net.UDPConn, udpAddr *net.UDPAddr) {
	fmt.Println("收到位置信息")
	locateDataReceive := server.LocateMesReceive{}

	NewMessByte, sign, _ := sertools.HandleMessProperty(&locateDataReceive.Header, messByte)
	if sign == false {
		return
	}
	sertools.LocationInfo(NewMessByte, &locateDataReceive.Data)
	var addMess = NewMessByte[28:]
	for len(addMess) != 0 {
		addMess = sertools.AdditionalInfo(addMess, &locateDataReceive.Data.AddLocateData)
	}

	locationInfo := sqlDB.GetLocationInfo{
		PhoneNum:    locateDataReceive.Header.PhoneNum,
		InfoType:    locateDataReceive.Data.InfoType,
		AlarmState:  locateDataReceive.Data.AlarmState,
		Latitude:    locateDataReceive.Data.Latitude,
		Longitude:   locateDataReceive.Data.Longitude,
		Altitude:    locateDataReceive.Data.Altitude,
		Speed:       locateDataReceive.Data.Speed,
		Direction:   locateDataReceive.Data.Direction,
		Time:        locateDataReceive.Data.Time,
		Mileage:     locateDataReceive.Data.AddLocateData.Mileage,
		Oil:         locateDataReceive.Data.AddLocateData.Oil,
		SpeedRecode: locateDataReceive.Data.AddLocateData.SpeedRecode,
	}
	//存储最新定位数据
	sertools.NewGPSInfo.Store(locationInfo.PhoneNum, locationInfo)

	gpsftp := []*sqlDB.GetLocationInfo{}
	//缓存10个定位数据再存数据库
	gpsList := sqlDB.GPSList{}
	gps, ok := sertools.GPSInfoListM.Load(locationInfo.PhoneNum)
	if ok {
		gpsList = gps.(sqlDB.GPSList)
		gpsList.Count++
		gpsList.GpsInfo[gpsList.Count] = locationInfo
		if gpsList.Count == 9 {
			sertools.GPSInfoListM.Delete(locationInfo.PhoneNum)
			tableNameInfo := new(sqlDB.GpsAndPhoneNum)
			if !sqlDB.QueryUserTake(tableNameInfo, map[string]interface{}{"PhoneNum": locationInfo.PhoneNum}) {
				fmt.Println("no GPS table")
				logs.BeeLogger.Info("no gps table")
			}
			sqlDB.InsertGPSInfo(tableNameInfo.GpsTableName, gpsList)
			for i := 0; i < 10 ; i++ {
				gpsftp = append(gpsftp, &gpsList.GpsInfo[i]) // Add clients
			}
			csvContent, _ := gocsv.MarshalString(&gpsftp) // Get all clients as CSV string
			fmt.Println(csvContent) // Display all clients as CSV string
			logs.BeeLogger.Info("gps csv info %s", csvContent)
			//gpsdata := bytes.NewBufferString(csvContent)
			csvName := "gps/gps_" + locationInfo.PhoneNum + "_" + tools.GetUUID() + ".csv"
			//ftpclient.FtpCli.Lock()
			//err := ftpclient.FtpCli.FTPClient.Stor(csvName, gpsdata)
			//ftpclient.FtpCli.Lock()
			//if err != nil {
			//	fmt.Println("gps info to ftp err")

			msg := ftpclient.Message{
				FileName: csvName,
				SaveMsg: csvContent,
			}
			ftpclient.FtpMsg <- msg

		} else {
			sertools.GPSInfoListM.Store(locationInfo.PhoneNum, gpsList)
		}
	} else {
		gpsList.Count = 0
		gpsList.GpsInfo[0] = locationInfo
		sertools.GPSInfoListM.Store(locationInfo.PhoneNum, gpsList)
	}

	//缓存10个报警信息再存数据库
	if locateDataReceive.Data.AlarmState != "" {
		alarmInfo := sqlDB.GetAlarmInfo{
			PhoneNum:        locateDataReceive.Header.PhoneNum,
			VehicleAlarm:    locateDataReceive.Data.AlarmState,
			StreamAlarm:     locateDataReceive.Data.AddLocateData.VideoAlarm,
			SignLostChannel: locateDataReceive.Data.AddLocateData.VideoSignalLossAlarm,
			Time:            locateDataReceive.Data.Time,
		}
		alarmList := sqlDB.AlarmList{}
		alarm, ok := sertools.AlarmInfoListM.Load(locationInfo.PhoneNum)
		if ok {
			alarmList = alarm.(sqlDB.AlarmList)
			alarmList.Count++
			alarmList.AlarmInfo[alarmList.Count] = alarmInfo
			if alarmList.Count == 9 {
				sertools.AlarmInfoListM.Delete(locationInfo.PhoneNum)
				tableNameInfo := new(sqlDB.GpsAndPhoneNum)
				if !sqlDB.QueryUserTake(tableNameInfo, map[string]interface{}{"PhoneNum": locationInfo.PhoneNum}) {
					fmt.Println("no ALARM table")
				}
				sqlDB.InsertAlarmInfo(tableNameInfo.AlarmTableName, alarmList)
			} else {
				sertools.AlarmInfoListM.Store(locationInfo.PhoneNum, alarmList)
			}
		} else {
			alarmList.Count = 0
			alarmList.AlarmInfo[0] = alarmInfo
			sertools.AlarmInfoListM.Store(locationInfo.PhoneNum, alarmList)
		}
	}


	//直接上报驾驶行为
	if locateDataReceive.Data.AddLocateData.AbnormalDrivingAlarm.AbnormalDrivingType != "" {
		driverInfo := sqlDB.DrivingAction{
			PhoneNum:            locateDataReceive.Header.PhoneNum,
			AbnormalDrivingType: locateDataReceive.Data.AddLocateData.AbnormalDrivingAlarm.AbnormalDrivingType,
			FatigueDegree:       locateDataReceive.Data.AddLocateData.AbnormalDrivingAlarm.FatigueDegree,
			Time:                locateDataReceive.Data.Time,
		}

		driveftp := []*sqlDB.DrivingAction{}
		driveftp = append(driveftp, &driverInfo) // Add clients
		csvContent, _ := gocsv.MarshalString(&driveftp) // Get all clients as CSV string
		fmt.Println(csvContent) // Display all clients as CSV string
		logs.BeeLogger.Info("driver action csv info %s", csvContent)
		//drivedata := bytes.NewBufferString(csvContent)
		csvName := "action/drAct_" + locationInfo.PhoneNum + "_" + tools.GetUUID() + ".csv"
		//ftpclient.FtpCli.Lock()
		//err := ftpclient.FtpCli.FTPClient.Stor(csvName, drivedata)
		//ftpclient.FtpCli.Unlock()
		//if err != nil {
		//	fmt.Println("driver action info to ftp err")
		//}

		msg := ftpclient.Message{
			FileName: csvName,
			SaveMsg: csvContent,
		}
		ftpclient.FtpMsg <- msg

		//driverActionList := sqlDB.DriverList{}
		//driver, ok := DriverActionListM.Load(locationInfo.PhoneNum)
		//if ok {
		//	driverActionList = driver.(sqlDB.DriverList)
		//	driverActionList.Count++
		//	driverActionList.DriverInfo[driverActionList.Count] = driverInfo
		//	if driverActionList.Count == 9 {
		//		DriverActionListM.Delete(locationInfo.PhoneNum)
		tableNameInfo := new(sqlDB.GpsAndPhoneNum)
		if !sqlDB.QueryUserTake(tableNameInfo, map[string]interface{}{"PhoneNum": locationInfo.PhoneNum}) {
			fmt.Println("no ALARM table")
		}
		sqlDB.InsertDriverActionInfo(tableNameInfo.DriverTableName, driverInfo)
		//	} else {
		//		AlarmInfoListM.Store(locationInfo.PhoneNum, driverActionList)
		//	}
		//} else {
		//	driverActionList.Count = 0
		//	driverActionList.DriverInfo[0] = driverInfo
		//	AlarmInfoListM.Store(locationInfo.PhoneNum, driverActionList)
		//}
	}

	//fmt.Println(locateDataReceive.Data.AddLocateData)
	client := sertools.GetUDPClient(locateDataReceive.Header.PhoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		fmt.Println("device disconnected, unable to keepalive")
		return
	}

	resultSendMes := sertools.NormalResponse(messByte, client.Sequence, 0)
	fmt.Println(hex.EncodeToString(resultSendMes))
	_, err := udpConn.WriteToUDP(resultSendMes, udpAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("位置信息回复成功")
	client.Sequence++
	sertools.SetUDPClient(locateDataReceive.Header.PhoneNum, client)
	return

}

//获取终端所有参数请求
func GetAllParameter(phoneNum string) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}
	sendMes := []byte{0x81, 0x04, 0x00, 0x00}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)
	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送获取终端参数成功")
	return
}

//处理终端所有参数
func HandleAllParameter(messByte []byte) (allParameterMes server.AllParameterMes) {
	fmt.Println("收到所有终端参数")

	NewMessByte, sign, _ := sertools.HandleMessProperty(&allParameterMes.Header, messByte)
	if sign == false {
		return
	}

	allParameterMes.Data.SequenceSend = binary.BigEndian.Uint16(NewMessByte[:2])
	fmt.Printf("服务器的消息流水号%v\n", allParameterMes.Data.SequenceSend)
	allParameterMes.Data.ParametersNum = NewMessByte[2]
	fmt.Printf("参数总数%v\n", allParameterMes.Data.ParametersNum)
	var addMess = NewMessByte[3:]
	for len(addMess) != 0 {
		addMess = sertools.DeviceParameters(addMess, allParameterMes.Header.PhoneNum, &allParameterMes.Data)
	}
	return
}

//通用回复
func ResNormalMes(messByte []byte) {
	fmt.Println("收到终端通用回复消息")

	normalMes := server.NormalMesReceive{}
	normalMes.Header.IDReceive = hex.EncodeToString(messByte[:2])
	normalMes.Header.MessProperty.DataLength = int64(binary.BigEndian.Uint16(messByte[2:4]))
	fmt.Printf("消息体长度%v\n", normalMes.Header.MessProperty.DataLength)
	normalMes.Header.PhoneNum = hex.EncodeToString(messByte[4:10])
	fmt.Println(normalMes.Header.PhoneNum)
	normalMes.Header.SequenceReceive = binary.BigEndian.Uint16(messByte[10:12])
	fmt.Printf("终端的消息流水号%v\n", normalMes.Header.SequenceReceive)

	normalMes.Data.SequenceSend = binary.BigEndian.Uint16(messByte[12:14])
	fmt.Printf("服务器的消息流水号%v\n", normalMes.Data.SequenceSend)
	normalMes.Data.IDSend = hex.EncodeToString(messByte[14:16])
	normalMes.Data.Result = messByte[16]
	sertools.MessageClass(messByte[14:17])

	switch normalMes.Data.IDSend {
	case "9101":
		//fmt.Printf("实时音视频传输开启请求")

		streamAndSeq := normalMes.Header.PhoneNum + strconv.FormatInt(int64(normalMes.Data.SequenceSend), 10)
		phoneAndChannel, ok := sertools.StreamChannel.Load(streamAndSeq)
		if ok {
			sertools.StreamChannel.Delete(streamAndSeq)
			if normalMes.Data.Result == 0 {
				sequenceInfo := sertools.StreamInfoData{}
				sequenceInfo.Count = 1
				streamAndChannelStr, _ := phoneAndChannel.(string)
				sequenceInfo.URL = sertools.GetSessionURLFromSTS(streamAndChannelStr)
				sertools.StreamInfo.Store(streamAndChannelStr, sequenceInfo)
			}
		}

	case "9102":
		//fmt.Printf("实时音视频传输控制请求")
		streamAndSeq := normalMes.Header.PhoneNum + strconv.FormatInt(int64(normalMes.Data.SequenceSend), 10)
		streamAndChannel, ok := sertools.StreamChannel.Load(streamAndSeq)
		if ok {
			sertools.StreamChannel.Delete(streamAndSeq)
			if normalMes.Data.Result == 0 {
				streamAndChannelStr, _ := streamAndChannel.(string)
				sertools.StreamInfo.Delete(streamAndChannelStr)
			}
		}
	case "9201":
		//fmt.Println("录像回放请求")
		streamAndSeq := normalMes.Header.PhoneNum + strconv.FormatInt(int64(normalMes.Data.SequenceSend), 10)
		signPhoneAndChannel, ok := sertools.StreamChannel.Load(streamAndSeq)
		if ok {
			sertools.StreamChannel.Delete(streamAndSeq)
			if normalMes.Data.Result == 0 {
				URL := "http://hello.com"
				streamAndChannelStr, _ := signPhoneAndChannel.(string)
				sertools.VideoReply.Store(streamAndChannelStr, URL)
			}
		}

	case "9202":
		//fmt.Println("录像回放控制")
		streamAndSeq := normalMes.Header.PhoneNum + strconv.FormatInt(int64(normalMes.Data.SequenceSend), 10)
		signPhoneAndChannel, ok := sertools.StreamChannel.Load(streamAndSeq)
		if ok {
			sertools.StreamChannel.Delete(streamAndSeq)
			if normalMes.Data.Result == 0 {
				streamAndChannelStr, _ := signPhoneAndChannel.(string)
				sertools.VideoReply.Delete(streamAndChannelStr)
			}
		}

	case "9205":
		fmt.Printf("实时音视频传输状态")
	case "8103":
		fmt.Printf("设置终端参数")
	case "8202":
		fmt.Printf("临时位置跟踪控制")
	case "8203":
		fmt.Printf("人工确认报警控制")
	case "8801":
		fmt.Printf("摄像头立即拍摄命令")
	case "8803":
		fmt.Printf("多媒体数据上传命令")

	}
}

////向终端发送实时直播的状态
//func SendRealStreamState(phoneNum string, channel byte) {
//	client := GetUDPClient(phoneNum)
//	if client == nil {
//		logs.BeeLogger.Error("device disconnected, unable to keepalive")
//		return
//	}
//	phoneNumByte, err := hex.DecodeString(phoneNum)
//	if err != nil {
//		logs.BeeLogger.Error("phoneNum error")
//		return
//	}
//
//	sendMes := []byte{0x92, 0x05, 0, 2}
//	sendMes = append(sendMes, phoneNumByte...)
//	sequenceByte := tools.Uint16ToByte(client.Sequence)
//	sendMes = append(sendMes, sequenceByte...)
//
//	sendMes = append(sendMes, channel, 0)
//	resultSendMes := tools.ParaphraseMess(sendMes)
//	fmt.Println(hex.EncodeToString(resultSendMes))
//	client.Sequence++
//	setUDPClient(phoneNum, client)
//	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
//	if err != nil {
//		fmt.Println("connect error")
//		return
//	}
//	fmt.Println("发送音视频传输状态通知成功")
//	return
//}

////上行透传数据回复
//func ResDataUplinkPassThrough(messByte []byte, udpConn *net.UDPConn, udpAddr *net.UDPAddr) {
//	fmt.Println("收到透传消息")
//	passThroughMes := server.PassThroughMes{}
//	NewMessByte, sign, _ := HandleMessProperty(&passThroughMes.Header, messByte)
//	if sign == false {
//		return
//	}
//	switch NewMessByte[0] {
//	case 0x00:
//		passThroughMes.Data.CNSSData = hex.EncodeToString(NewMessByte[1:])
//		fmt.Printf("GNSS 模块详细定位数据 %v\n", passThroughMes.Data.CNSSData)
//	case 0x0b:
//		passThroughMes.Data.ICData = hex.EncodeToString(NewMessByte[1:])
//		fmt.Printf("道路运输证 IC 卡信息 %v\n", passThroughMes.Data.ICData)
//	case 0x41:
//		passThroughMes.Data.SerialPort1 = hex.EncodeToString(NewMessByte[1:])
//		fmt.Printf("串口 1 透传 %v\n", passThroughMes.Data.SerialPort1)
//	case 0x42:
//		passThroughMes.Data.SerialPort2 = hex.EncodeToString(NewMessByte[1:])
//		fmt.Printf("串口 2 透传 %v\n", passThroughMes.Data.SerialPort2)
//	default:
//		passThroughMes.Data.Custom.ID = NewMessByte[0]
//		passThroughMes.Data.Custom.Data = hex.EncodeToString(NewMessByte[1:])
//		fmt.Printf("用户自定义透传, id %v, data %v\n", passThroughMes.Data.Custom.ID, passThroughMes.Data.Custom.Data)
//	}
//
//	client := GetUDPClient(passThroughMes.Header.PhoneNum)
//	if client == nil {
//		logs.BeeLogger.Error("device disconnected, unable to keepalive")
//		fmt.Println("device disconnected, unable to keepalive")
//		return
//	}
//	//回复数据编写
//	resultSendMes := NormalResponse(messByte, client.Sequence, 0)
//	fmt.Println(resultSendMes)
//	_, err := udpConn.WriteToUDP(resultSendMes, udpAddr)
//	if err != nil {
//		fmt.Println("connect error")
//		return
//	}
//	fmt.Println("回复数据上行透传成功")
//	//重置键值对
//	client.Sequence++
//	setUDPClient(passThroughMes.Header.PhoneNum, client)
//	return
//}

//设置设备参数
func SetDeviceParameter(phoneNum string, num byte, settings string) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}

	sendMes := []byte{0x81, 0x03}
	settingsLength := len(settings)/2 + 1
	settingsLengthByte := tools.Int64ToByte(int64(settingsLength))
	sendMes = append(sendMes, settingsLengthByte[6:]...)
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	sendMes = append(sendMes, num)
	mes, err := hex.DecodeString(settings)
	if err != nil {
		return
	}
	sendMes = append(sendMes, mes...)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送设置终端参数成功")
	return
}

//获取指定终端参数
func GetAppointParameter(phoneNum string, num byte, parameterID string) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}

	length := uint16(num)*4 + 1
	lengthByte := tools.Uint16ToByte(length)
	sendMes := []byte{0x81, 0x06}
	sendMes = append(sendMes, lengthByte...)
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	allParameter := server.AllParameterDataReceive{}
	allParameter.ParametersNum = num
	sendMes = append(sendMes, allParameter.ParametersNum)
	id, err := hex.DecodeString(parameterID)
	if err != nil {
		fmt.Println("参数id error")
	}
	sendMes = append(sendMes, id...)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送获取指定终端参数成功")
	return
}

//查询终端属性
func QueryDeviceProperties(phoneNum string) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}

	sendMes := []byte{0x81, 0x07, 0, 0}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)
	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送查询终端属性成功")
	return
}

//解析终端回复的属性信息
func HandleDeviceProperties(messByte []byte) {
	fmt.Println("收到终端回复的属性信息")

	deviceProperties := server.DevicePropertiesMes{}
	NewMessByte, sign, _ := sertools.HandleMessProperty(&deviceProperties.Header, messByte)
	if sign == false {
		return
	}

	deviceType := binary.BigEndian.Uint16(NewMessByte[:2])
	deviceProperties.Data.DeviceType = tools.UintTo2byte(uint64(deviceType), 16)
	fmt.Printf("终端类型 %v\n", string(deviceProperties.Data.DeviceType))
	deviceProperties.Data.ManufacturerID = string(NewMessByte[2:7])
	deviceProperties.Data.DeviceModel = string(NewMessByte[7:27])
	deviceProperties.Data.DeviceID = string(NewMessByte[27:34])
	deviceProperties.Data.DeviceSIM = hex.EncodeToString(NewMessByte[34:44])
	fmt.Printf("制造商 ID %v\n终端型号 %v\n终端 ID %v\n终端 SIM 卡 ICCID 号 %v\n ", deviceProperties.Data.ManufacturerID,
		deviceProperties.Data.DeviceModel, deviceProperties.Data.DeviceID, deviceProperties.Data.DeviceSIM)
	n := NewMessByte[44]
	deviceProperties.Data.DeviceHardwareLength = n
	deviceProperties.Data.DeviceHardware = string(NewMessByte[45 : 45+n])
	fmt.Printf("终端硬件版本号长度 %v , 终端硬件版本号 %v\n", deviceProperties.Data.DeviceHardwareLength, deviceProperties.Data.DeviceHardware)
	m := NewMessByte[45+n]
	deviceProperties.Data.DeviceFirmwareLength = m
	deviceProperties.Data.DeviceFirmware = string(NewMessByte[46+n : 46+m+n])
	fmt.Printf("终端固件版本号长度 %v , 终端固件版本号 %v\n", deviceProperties.Data.DeviceFirmwareLength, deviceProperties.Data.DeviceFirmware)
	deviceProperties.Data.GNSSType = tools.UintTo2byte(uint64(NewMessByte[46+m+n]), 8)
	fmt.Printf("GNSS 模块属性 %v\n", string(deviceProperties.Data.GNSSType))
	deviceProperties.Data.CommunicationType = tools.UintTo2byte(uint64(NewMessByte[47+m+n]), 8)
	fmt.Printf("通信模块属性 %v\n", string(deviceProperties.Data.CommunicationType))
}

//处理终端回复的位置信息
func HandleLocationInfo(messByte []byte) {
	fmt.Println("收到位置信息")

	locateDataReceive := server.LocateMesReceive{}
	NewMessByte, sign, _ := sertools.HandleMessProperty(&locateDataReceive.Header, messByte)
	if sign == false {
		return
	}

	locateDataReceive.Data.SequenceSend = binary.BigEndian.Uint16(NewMessByte[:2])
	fmt.Printf("服务器的消息流水号 %v\n", locateDataReceive.Data.SequenceSend)
	sertools.LocationInfo(NewMessByte[2:30], &locateDataReceive.Data)
	var addMess = NewMessByte[30:]
	for len(addMess) != 0 {
		addMess = sertools.AdditionalInfo(addMess, &locateDataReceive.Data.AddLocateData)
	}
	locationInfo := sqlDB.GetLocationInfo{
		PhoneNum:    locateDataReceive.Header.PhoneNum,
		InfoType:    locateDataReceive.Data.InfoType,
		AlarmState:  locateDataReceive.Data.AlarmState,
		Latitude:    locateDataReceive.Data.Latitude,
		Longitude:   locateDataReceive.Data.Longitude,
		Altitude:    locateDataReceive.Data.Altitude,
		Speed:       locateDataReceive.Data.Speed,
		Direction:   locateDataReceive.Data.Direction,
		Time:        locateDataReceive.Data.Time,
		Mileage:     locateDataReceive.Data.AddLocateData.Mileage,
		Oil:         locateDataReceive.Data.AddLocateData.Oil,
		SpeedRecode: locateDataReceive.Data.AddLocateData.SpeedRecode,
	}
	sertools.RealGPSInfo.Store(locateDataReceive.Header.PhoneNum, locationInfo)
	channel := strings.Split(locateDataReceive.Data.AddLocateData.VideoSignalLossAlarm, " ")
	for _, v := range channel[:len(channel)-1] {
		//fmt.Println(v)
		//更新通道信息表格数据
		phoneNumAndChannel := tools.StringsJoin(locateDataReceive.Header.PhoneNum, "_", v)
		sqlDB.UpdateTableFromChannel(sqlDB.GetChannelInfo{}, "Alarm", "Lost", time.Now().Format("2006-01-02 15:04:05"), phoneNumAndChannel)
	}
	return
}

//向终端请求追踪位置信息
func TrackLocationInfo(phoneNum string, timeInterval uint16, totalTime uint32) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}
	sendMes := []byte{0x82, 0x02, 0, 6}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	timeIntervalByte := tools.Uint16ToByte(timeInterval)
	sendMes = append(sendMes, timeIntervalByte...)
	if timeInterval != 0 {
		totalTimeByte := tools.Uint32ToByte(totalTime)
		sendMes = append(sendMes, totalTimeByte...)
	}

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送跟踪位置参数成功")
	return
}

//向终端人工确认报警消息
func SendAcknowledgeAlarm(phoneNum string, sequenceReceive uint16, alarmType string) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}
	sendMes := []byte{0x82, 0x03, 0, 6}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	sequenceReceiveByte := tools.Uint16ToByte(sequenceReceive)
	sendMes = append(sendMes, sequenceReceiveByte...)
	alarmType64 := tools.TwoToInt64(alarmType)
	alarmTypeByte := tools.Int64ToByte(alarmType64)
	sendMes = append(sendMes, alarmTypeByte...)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送人工确认报警成功")
	return
}

//向终端发送控制车辆请求
func SendVehicleControl(phoneNum string, control string) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}
	sendMes := []byte{0x85, 0x00, 0, 1}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	control64 := tools.TwoToInt64(control)
	controlByte := tools.Int64ToByte(control64)
	sendMes = append(sendMes, controlByte[7])

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送车辆控制成功")
	return
}

//解析终端回复的车辆控制应答
func HandleVehicleControlReceive(messByte []byte) {
	fmt.Println("收到车辆控制应答")
	locateDataReceive := server.LocateMesReceive{}
	NewMessByte, sign, _ := sertools.HandleMessProperty(&locateDataReceive.Header, messByte)
	if sign == false {
		return
	}

	locateDataReceive.Data.SequenceSend = binary.BigEndian.Uint16(NewMessByte[:2])
	state := binary.BigEndian.Uint32(NewMessByte[6:10])
	state2 := tools.UintTo2byte(uint64(state), 32)
	switch state2[20] {
	case '0':
		fmt.Println("车门解锁")
	case '1':
		fmt.Println("车门加锁")
	default:
		fmt.Println("参数错误")
	}
	fmt.Printf("服务器的消息流水号 %v\n", locateDataReceive.Data.SequenceSend)

}

//向终端立即拍照请求
func SendShootNow(phoneNum string, channel byte, photoNum uint16, timeInterval uint16) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}
	sendMes := []byte{0x88, 0x01, 0, 12}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	sendMes = append(sendMes, channel)
	photoNumByte := tools.Uint16ToByte(photoNum)
	sendMes = append(sendMes, photoNumByte...)
	timeIntervalByte := tools.Uint16ToByte(timeInterval)
	sendMes = append(sendMes, timeIntervalByte...)
	sendMes = append(sendMes, 1, 4, 1, 0, 0, 0, 0)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送摄像头立即拍摄命令成功")
	return
}

//查询终端音视频属性
func QueryStreamProperties(phoneNum string) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}
	sendMes := []byte{0x90, 0x03, 0, 0}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)
	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送查询终端音视频属性成功")
	return
}

//终端上传音视频属性解析
func HandleStreamProperties(messByte []byte) {
	fmt.Println("收到终端音视频属性应答")
	deviceStreamPropertiesMes := server.DeviceStreamPropertiesMes{}

	NewMessByte, sign, _ := sertools.HandleMessProperty(&deviceStreamPropertiesMes.Header, messByte)
	if sign == false {
		return
	}

	deviceStreamPropertiesMes.Data.AudioCoding = NewMessByte[0]
	deviceStreamPropertiesMes.Data.AudioChannelCount = NewMessByte[1]
	deviceStreamPropertiesMes.Data.AudioSamplingFrequency = NewMessByte[2]
	deviceStreamPropertiesMes.Data.AudioSamplingBits = NewMessByte[3]
	deviceStreamPropertiesMes.Data.AudioFrameLength = binary.BigEndian.Uint16(NewMessByte[4:6])
	deviceStreamPropertiesMes.Data.AllowAudioOutput = NewMessByte[6]
	deviceStreamPropertiesMes.Data.VideoCoding = NewMessByte[7]
	deviceStreamPropertiesMes.Data.MaxAudioChannelCount = NewMessByte[8]
	deviceStreamPropertiesMes.Data.MaxVideoChannelCount = NewMessByte[9]
	vehicleInfo := new(sqlDB.GetVehicleInfo)
	if !sqlDB.QueryUserTake(vehicleInfo, map[string]interface{}{"PhoneNum": deviceStreamPropertiesMes.Header.PhoneNum}) {
		logs.BeeLogger.Error("vehicle db err")
		return
	}
	vehicleInfo.ChannelCount = deviceStreamPropertiesMes.Data.MaxVideoChannelCount
	sqlDB.Save(vehicleInfo)
	channelInfo := sqlDB.GetChannelInfo{}
	//mqttChannelInfo := server.MQTTChannelInfo{}
	//mqttChannelList := []server.MQTTChannelInfo{}

	for i := 1; i <= int(deviceStreamPropertiesMes.Data.MaxVideoChannelCount); i++ {
		channelInfo.PhoneNumAndChannel = deviceStreamPropertiesMes.Header.PhoneNum +"_" + strconv.Itoa(i)
		channelInfo.PhoneNum = deviceStreamPropertiesMes.Header.PhoneNum
		channelInfo.Status = "ON"
		channelInfo.LogicalChannelID = int64(i)
		channelInfo.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		if sqlDB.Save(channelInfo) {
			logs.BeeLogger.Info("save channel[%v] info success",i)
		}

		//mqttChannelInfo.Id = channelInfo.PhoneNumAndChannel
		//mqttChannelInfo.Parentid = channelInfo.PhoneNum
		//mqttChannelInfo.Name = strconv.Itoa(i)
		//mqttChannelInfo.Status = "ON"
		//
		//mqttChannelList = append(mqttChannelList, mqttChannelInfo)

	}

	//if server.RemoteClientStatus == "ON" {
	//	payload := server.MQTTNotify{}
	//	payload.Protocol = "NotifySub"
	//	payload.Csq = server.Csq
	//	server.Csq++
	//	payload.Type = "Channel"
	//	payload.Data = map[string]interface{}{
	//		"channellists": mqttChannelList,
	//	}
	//	msg := server.Message{
	//		Topic:   server.MattPubID,
	//		Payload: payload,
	//	}
	//	server.MqttMess <- msg
	//}

	fmt.Printf("终端上传音视频属性%v\n", deviceStreamPropertiesMes.Data)
}

//解析运营登记信息并返回
func ResOperationRegistrationMes(messByte []byte, udpConn *net.UDPConn, udpAddr *net.UDPAddr) {
	fmt.Println("收到运营登记信息")
	operationRegistrationMes := server.OperationRegistrationMes{}
	NewMessByte, sign, _ := sertools.HandleMessProperty(&operationRegistrationMes.Header, messByte)
	if sign == false {
		return
	}
	operationRegistrationMes.Data.LineNumber = binary.BigEndian.Uint32(NewMessByte[:4])
	operationRegistrationMes.Data.EmployeeNumber = hex.EncodeToString(NewMessByte[4:])
	fmt.Printf("线路编号 %v, 员工编号 %v\n", operationRegistrationMes.Data.LineNumber, operationRegistrationMes.Data.EmployeeNumber)

	client := sertools.GetUDPClient(operationRegistrationMes.Header.PhoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		fmt.Println("device disconnected, unable to keepalive")
		return
	}

	resultSendMes := sertools.NormalResponse(messByte, client.Sequence, 0)
	fmt.Println(hex.EncodeToString(resultSendMes))
	_, err := udpConn.WriteToUDP(resultSendMes, udpAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("位置信息回复成功")
	client.Sequence++
	sertools.SetUDPClient(operationRegistrationMes.Header.PhoneNum, client)
	return
}

//多媒体数据检索
func MultimediaDataSearch(phoneNum string, channel, alarmType byte, startTime, endTime string) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}
	sendMes := []byte{0x88, 0x02, 0, 15}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	sendMes = append(sendMes, 0, channel, alarmType)
	startByte, _ := hex.DecodeString(startTime)
	endByte, _ := hex.DecodeString(endTime)
	sendMes = append(sendMes, startByte...)
	sendMes = append(sendMes, endByte...)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送存储多媒体数据检索命令成功")
	return
}

//多媒体数据上传
func MultimediaDataUploadSend(phoneNum string, channel, alarmType byte, startTime, endTime string) {
	client := sertools.GetUDPClient(phoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		return
	}
	phoneNumByte, err := hex.DecodeString(phoneNum)
	if err != nil {
		logs.BeeLogger.Error("phoneNum error")
		return
	}
	sendMes := []byte{0x88, 0x03, 0, 16}
	sendMes = append(sendMes, phoneNumByte...)
	sequenceByte := tools.Uint16ToByte(client.Sequence)
	sendMes = append(sendMes, sequenceByte...)

	sendMes = append(sendMes, 0, channel, alarmType)
	startByte, _ := hex.DecodeString(startTime)
	endByte, _ := hex.DecodeString(endTime)
	sendMes = append(sendMes, startByte...)
	sendMes = append(sendMes, endByte...)
	sendMes = append(sendMes, 0)

	resultSendMes := tools.ParaphraseMess(sendMes)
	fmt.Println(hex.EncodeToString(resultSendMes))
	client.Sequence++
	sertools.SetUDPClient(phoneNum, client)
	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("发送存储多媒体数据上传命令成功")
	return
}

//多媒体数据上传解析
func MultimediaDataUploadHandle(messByte []byte, udpConn *net.UDPConn, udpAddr *net.UDPAddr) {
	fmt.Println("收到多媒体数据上传")
	multimediaMes := server.MultimediaMes{}
	NewMessByte, sign, needRepeat := sertools.HandleMessProperty(&multimediaMes.Header, messByte)
	if sign == false {
		return
	}
	multimediaMes.Data.ID = binary.BigEndian.Uint32(NewMessByte[:4])
	multimediaMes.Data.MediaType = NewMessByte[4]
	multimediaMes.Data.MediaCoding = NewMessByte[5]
	multimediaMes.Data.IncidentCoding = NewMessByte[6]
	multimediaMes.Data.ChannelID = NewMessByte[7]
	sertools.LocationInfo(NewMessByte[8:36], &multimediaMes.Data.LocateData)

	client := sertools.GetUDPClient(multimediaMes.Header.PhoneNum)
	if client == nil {
		logs.BeeLogger.Error("device disconnected, unable to keepalive")
		fmt.Println("device disconnected, unable to keepalive")
		return
	}
	sendMess := []byte{0x88, 0x00}
	var packageID []byte
	if needRepeat {
		length := 5 + 2*multimediaMes.Header.PackageCount
		lengthByte := tools.Uint16ToByte(length)
		sendMess = append(sendMess, lengthByte...)
		for i := uint16(1); i <= multimediaMes.Header.PackageCount; i++ {
			packageID = append(packageID, tools.Uint16ToByte(i)...)
		}
	} else {
		sendMess = append(sendMess, 0, 4)
		channelID := strconv.FormatInt(int64(multimediaMes.Data.ChannelID), 10)
		pictureType := strconv.FormatInt(int64(multimediaMes.Data.MediaType), 10)
		pictureID := strconv.FormatInt(int64(multimediaMes.Data.ID), 10)
		appPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		appPath = strings.Replace(appPath, "\\", "/", -1)
		filepath := appPath + "/picture/" + multimediaMes.Header.PhoneNum + "_" + channelID + "_" + pictureType + "_" + pictureID + ".jpg"
		multimediaMes.Data.DataPackage = NewMessByte[36:]
		err2 := ioutil.WriteFile(filepath, multimediaMes.Data.DataPackage, 0666)
		if err2 != nil {
			return
		}
	}

	phoneNum, _ := hex.DecodeString(multimediaMes.Header.PhoneNum)
	sendMess = append(sendMess, phoneNum...)
	sequence := tools.Uint16ToByte(client.Sequence)
	sendMess = append(sendMess, sequence...)
	id := tools.Uint32ToByte(multimediaMes.Data.ID)
	sendMess = append(sendMess, id...)

	if needRepeat {
		packageCount := tools.Uint16ToByte(multimediaMes.Header.PackageCount)
		sendMess = append(sendMess, packageCount[1])
		sendMess = append(sendMess, packageID...)
	}

	resultSendMes := tools.ParaphraseMess(sendMess)
	fmt.Println(hex.EncodeToString(resultSendMes))
	_, err := udpConn.WriteToUDP(resultSendMes, udpAddr)
	if err != nil {
		fmt.Println("connect error")
		return
	}
	fmt.Println("多媒体数据上传回复成功")
	client.Sequence++
	sertools.SetUDPClient(multimediaMes.Header.PhoneNum, client)
	return

}

//处理返回的录像资源列表
func VideoListGetHandle(messByte []byte) {
	fmt.Println("收到录像资源列表")
	videoListMes := server.VideoListMes{}
	NewMessByte, sign, _ := sertools.HandleMessProperty(&videoListMes.Header, messByte)
	if sign == false {
		return
	}
	videoListMes.Data.SequenceSend = binary.BigEndian.Uint16(NewMessByte[:2])
	videoListMes.Data.VideoCount = binary.BigEndian.Uint32(NewMessByte[2:6])
	videoInfoByte := NewMessByte[6:]
	for len(videoInfoByte) > 0 {
		videoInfoByte = sertools.VideoListInfoHandle(videoInfoByte, &videoListMes.Data)
	}
	//
	streamAndSeq := videoListMes.Header.PhoneNum + strconv.FormatInt(int64(videoListMes.Data.SequenceSend), 10)
	signPhoneAndChannel, ok := sertools.StreamChannel.Load(streamAndSeq)
	if ok {
		sertools.StreamChannel.Delete(streamAndSeq)
		signPhoneAndChannelStr, _ := signPhoneAndChannel.(string)
		sertools.VideoList.Store(signPhoneAndChannelStr, videoListMes.Data)
	}

	return
}

//func RequestActualPassword(phoneNum string){
//	client := GetUDPClient(phoneNum)
//	if client == nil {
//		logs.BeeLogger.Error("device disconnected, unable to keepalive")
//		return
//	}
//	phoneNumByte,err := hex.DecodeString(phoneNum)
//	if err != nil {
//		logs.BeeLogger.Error("phoneNum error")
//		return
//	}
//
//	sendMes := []byte{0x17, 0x00, 0x00, 0x1c}
//	sendMes = append(sendMes,phoneNumByte...)
//	sequenceByte := tools.Uint16ToByte(client.Sequence)
//	sendMes = append(sendMes, sequenceByte...)
//	vehicleInfo := new(sqlDB.GetVehicleInfo)
//	if !sqlDB.QueryUserTake(vehicleInfo, map[string]interface{}{"PhoneNum": phoneNum}){
//		fmt.Println("数据库中无数据")
//		return
//	}
//	for i := 0; i < 15; i++ {
//		sendMes = append(sendMes, 0)
//	}
//	sendMes = append(sendMes, []byte(vehicleInfo.VehicleID)...)
//	sendMes = append(sendMes, 1)
//
//	resultSendMes := tools.ParaphraseMess(sendMes)
//	fmt.Println(hex.EncodeToString(resultSendMes))
//	client.Sequence ++
//	setUDPClient(phoneNum, client)
//	_, err = client.UDPConn.WriteToUDP(resultSendMes, client.UDPAddr)
//	if err != nil {
//		fmt.Println("connect error")
//		return
//	}
//	fmt.Println("发送获取终端参数成功")
//	return
//}
