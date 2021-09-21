package controllers

import (
	//"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitee.com/ictt/JTTM/server/ftpclient"
	"gitee.com/ictt/JTTM/server/sertools"
	"gitee.com/ictt/JTTM/server/tcp"
	"gitee.com/ictt/JTTM/server/udp"
	"gitee.com/ictt/JTTM/tools"
	"gitee.com/ictt/JTTM/tools/logs"
	"github.com/astaxie/beego"
	"github.com/gocsv"
	"strconv"
	"time"
)

type MainController struct {
	beego.Controller
}

//var passengerFlowM sync.Map
//实时直播-直播流停止
func (this *MainController) HooStop() {
	fmt.Println("Hook接收关闭流时间: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("hook_on_stream_none_reader: ")
	logs.BeeLogger.Info("hook_on_stream_none_reader: %s",string(this.Ctx.Input.RequestBody))
	fmt.Println(string(this.Ctx.Input.RequestBody))
	data := struct {
		App    string `json:"app"`    //流应用名
		Schema string `json:"schema"` //rtsp或rtmp
		Stream string `json:"stream"` //流ID
		Vhost  string `json:"vhost"`  //流虚拟主机
		//Params string `json:"params"` //token
	}{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)

	if err != nil {
		logs.BeeLogger.Error("analysis stream stop error: %s", err)
		this.Data["json"] = map[string]interface{}{
			"close": false,
			"code":  0,
		}
		this.ServeJSON()
		return
	}

	switch data.App {
	case "live":
		//关闭实时直播
		for len(data.Stream) < 14 {
			data.Stream = "0" + data.Stream
		}
		phoneNum := data.Stream[:12]
		channel, _ := hex.DecodeString(data.Stream[12:])
		fmt.Println("执行关闭实时直播！！！")
		//fmt.Println(channel)
		//fmt.Println(phoneNum)
		client := sertools.GetUDPClient(phoneNum)
		if client !=  nil {
			udp.ControlRealStream(phoneNum, channel[0], 0, 0, 0)
		} else {
			tcp.ControlRealStream(phoneNum, channel[0], 0, 0, 0)
		}

		this.Data["json"] = map[string]interface{}{
			"close": "true",
			"code":  1,
		}
		this.ServeJSON()
		//case "record":
		//	//关闭录像回放
		//	fmt.Println("执行关闭录像回放！！！")
		//	tcpServers.ListenPlaybackStop(data.Stream[:20], data.Stream, data.Params)
		//}
	}
	return
}

type dangerousThings struct {
	//Code int64 `json:"code"`
	AlertAdd1 		string `json:"alert_Add1"`
	AlertAdd2 		string `json:"alert_Add2"`
	AlertAdd3 		string `json:"alert_Add3"`
	AlertAdd4 		string `json:"alert_Add4"`
	AlertAdd5 		string `json:"alert_Add5"`
	AlertAdd6 		string `json:"alert_Add6"`
	AlertAdd7 		string `json:"alert_Add7"`
	AlertAdd8		string `json:"alert_Add8"`
	AlertAdd9 		string `json:"alert_Add9"`
	AlertAdd10 		string `json:"alert_Add10"`
	Alert_ID		int64 `json:"alert_ID"`
	Alert_DWBH		string `json:"alert_DWBH"`
	Alert_XLBH		string `json:"alert_XLBH"`
	Alert_CH		string `json:"alert_CH"`
	Alert_RQ		string `json:"alert_RQ"`
	Alert_SJ		string `json:"alert_SJ"`
	Alert_JSY		string `json:"alert_JSY"`
	Alert_Exp		string `json:"alert_Exp"`
	Alert_PACNO		int64 `json:"alert_PACNO"`
	Alert_Value1	int64 `json:"alert_Value1"`
	Alert_Value2	int64 `json:"alert_Value2"`
	Alert_Value3	int64 `json:"alert_Value3"`
	Alert_Value4	int64 `json:"alert_Value4"`
	Alert_Value5	int64 `json:"alert_Value5"`
	Alert_Value6	int64 `json:"alert_Value6"`
	Alert_Value7	int64 `json:"alert_Value7"`
	Alert_Value8	int64 `json:"alert_Value8"`
	Alert_Value9	int64 `json:"alert_Value9"`
	Alert_Value10	int64 `json:"alert_Value10"`
	Alert_Warn1		int64 `json:"alert_Warn1"`
	Alert_Warn2		int64 `json:"alert_Warn2"`
	Alert_Warn3		int64 `json:"alert_Warn3"`
	Alert_Warn4		int64 `json:"alert_Warn4"`
	Alert_Warn5		int64 `json:"alert_Warn5"`
	Alert_Warn6		int64 `json:"alert_Warn6"`
	Alert_Warn7		int64 `json:"alert_Warn7"`
	Alert_Warn8		int64 `json:"alert_Warn8"`
	Alert_Warn9		int64 `json:"alert_Warn9"`
	Alert_Warn10	int64 `json:"alert_Warn10"`
	Alert_CJRQ		string `json:"alert_CJRQ"`
	Alert_CJSJ		string `json:"alert_CJSJ"`
	Alert_JD		int64 `json:"alert_JD"`
	Alert_WD		int64 `json:"alert_WD"`
}

type dangerousData struct {
	Code int64               `json:"code"`
	IsSuccess bool           `json:"isSuccess"`
	Massage string           `json:"massege"`
	Result []dangerousThings `json:"result"`
}

type dangerousTable struct {
	AlertAdd 		string `csv:"alert_Add"`
	Alert_ID		int64 `csv:"alert_ID"`
	Alert_DWBH		string `csv:"alert_DWBH"`
	Alert_XLBH		string `csv:"alert_XLBH"`
	Alert_CH		string `csv:"alert_CH"`
	Alert_RQ		string `csv:"alert_RQ"`
	Alert_SJ		string `csv:"alert_SJ"`
	Alert_JSY		string `csv:"alert_JSY"`
	Alert_Exp		string `csv:"alert_Exp"`
	Alert_PACNO		int64 `csv:"alert_PACNO"`
	Alert_Value		int64 `csv:"alert_Value"`
	Alert_Warn		int64 `csv:"alert_Warn"`
	Alert_CJRQ		string `csv:"alert_CJRQ"`
	Alert_CJSJ		string `csv:"alert_CJSJ"`
	Alert_JD		int64 `csv:"alert_JD"`
	Alert_WD		int64 `csv:"alert_WD"`
}

func (this *MainController) Gas() {
	//fmt.Println("get gas danger info: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("get gas danger info:")
	logs.BeeLogger.Info("get gas danger info: %s",string(this.Ctx.Input.RequestBody))
	fmt.Println(string(this.Ctx.Input.RequestBody))
	data := dangerousData{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)

	if err != nil {
		logs.BeeLogger.Error("analysis stream stop error: %s", err)
		fmt.Println(err)
		this.Data["json"] = map[string]interface{}{
			"close": false,
			"code":  0,
		}
		this.ServeJSON()
		return
	}

	switch data.Code {
	case 200:
		//关闭实时直播
		for i := 0; i < len(data.Result); i++ {
			addr := [10]string{data.Result[i].AlertAdd1, data.Result[i].AlertAdd2, data.Result[i].AlertAdd3,
				data.Result[i].AlertAdd4, data.Result[i].AlertAdd5, data.Result[i].AlertAdd6, data.Result[i].AlertAdd7,
				data.Result[i].AlertAdd8, data.Result[i].AlertAdd9, data.Result[i].AlertAdd10}
			value := [10]int64{data.Result[i].Alert_Value1,data.Result[i].Alert_Value2, data.Result[i].Alert_Value3,
				data.Result[i].Alert_Value4, data.Result[i].Alert_Value5, data.Result[i].Alert_Value6, data.Result[i].Alert_Value7,
				data.Result[i].Alert_Value8, data.Result[i].Alert_Value9, data.Result[i].Alert_Value10}
			warn := [10]int64{data.Result[i].Alert_Warn1, data.Result[i].Alert_Warn2, data.Result[i].Alert_Warn3,
				data.Result[i].Alert_Warn4, data.Result[i].Alert_Warn5, data.Result[i].Alert_Warn6, data.Result[i].Alert_Warn7,
				data.Result[i].Alert_Warn8, data.Result[i].Alert_Warn9, data.Result[i].Alert_Warn10}

			driveftp := []*dangerousTable{}
			for  j := 0; j < 10; j++ {
				if addr[j] != "" {
					danger := dangerousTable{
						AlertAdd:		addr[j],
						Alert_ID:		data.Result[i].Alert_ID,
						Alert_DWBH:		data.Result[i].Alert_DWBH,
						Alert_XLBH:		data.Result[i].Alert_XLBH,
						Alert_CH:		data.Result[i].Alert_CH,
						Alert_RQ:		data.Result[i].Alert_RQ,
						Alert_SJ:		data.Result[i].Alert_SJ,
						Alert_JSY:		data.Result[i].Alert_JSY,
						Alert_Exp:		data.Result[i].Alert_Exp,
						Alert_PACNO:	data.Result[i].Alert_PACNO,
						Alert_Value:	value[j],
						Alert_Warn:		warn[j],
						Alert_CJRQ:		data.Result[i].Alert_CJRQ,
						Alert_CJSJ:		data.Result[i].Alert_CJSJ,
						Alert_JD:		data.Result[i].Alert_JD,
						Alert_WD:		data.Result[i].Alert_WD,
					}
					driveftp = append(driveftp, &danger) // Add clients
					}
				}
			csvContent, _ := gocsv.MarshalString(&driveftp) // Get all clients as CSV string
			fmt.Println(csvContent) // Display all clients as CSV string
			logs.BeeLogger.Info("danger data csv info %s", csvContent)
			//dangerdata := bytes.NewBufferString(csvContent)
			csvName := "danger/danger_" + data.Result[i].Alert_XLBH + "_" + data.Result[i].Alert_CH + ".csv"
			//ftpclient.FtpCli.Lock()
			//err := ftpclient.FtpCli.FTPClient.Stor(csvName, dangerdata)
			//ftpclient.FtpCli.Unlock()
			//if err != nil {
			//	fmt.Println("danger info to ftp err")
			//	logs.BeeLogger.Error("danger info to ftp err")
			//}
			msg := ftpclient.Message{
				FileName: csvName,
				SaveMsg: csvContent,
			}
			ftpclient.FtpMsg <- msg

		}
		this.Data["json"] = map[string]interface{}{
			"close": "true",
			"code":  1,
		}
		this.ServeJSON()
	}
	return
}

type passengerFlowTimeData struct {
	DeviceID	string 	`json:"DeviceID" csv:"deviceID"`
	AlarmType 	int64 	`json:"AlarmType" csv:"alarmType"`
	AlarmDescribe	string `json:"AlarmDescribe" csv:"alarmDescribe"`
	StartTime 	string 	`json:"StartTime" csv:"startTime"`
	StopTime 	string 	`json:"StopTime" csv:"stopTime"`
	EnterNum	int64 	`json:"EnterNum" csv:"enterNum"`
	LeaveNum	int64	`json:"LeaveNum" csv:"leaveNum"`
	VideoChannel int8	`json:"VideoChannel" csv:"videoChannel"`
}

//按时间接受客流信息
func (this *MainController) PassengerFlowTime() {
	//fmt.Println("get gas danger info: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("get PassengerFlow info:")
	logs.BeeLogger.Info("get PassengerFlow info: %s", string(this.Ctx.Input.RequestBody))
	fmt.Println(string(this.Ctx.Input.RequestBody))
	data := passengerFlowTimeData{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.BeeLogger.Error("PassengerFlowTime stop error: %s", err)
		fmt.Println(err)
		this.Data["json"] = map[string]interface{}{
			"close": false,
			"code":  0,
		}
		this.ServeJSON()
		return
	}

	defer func() {
		this.Data["json"] = map[string]interface{}{
			"close": "true",
			"code":  1,
		}
		this.ServeJSON()
	}()

	psssenger := []*passengerFlowTimeData{}
	psssenger = append(psssenger, &data)
	csvContent, _ := gocsv.MarshalString(&psssenger) // Get all clients as CSV string
	fmt.Println(csvContent) // Display all clients as CSV string
	logs.BeeLogger.Info("PassengerFlowTime  data csv info %s", csvContent)
	//passengerdata := bytes.NewBufferString(csvContent)
	csvName := "passengerFlow/byTime_" +  data.DeviceID + "_"+ strconv.Itoa(int(data.VideoChannel)) + "_" + tools.GetUUID() + ".csv"
	//err = ftpclient.FTPClient.Stor(csvName, passengerdata)
	//ftpclient.FtpCli.Lock()
	//err = ftpclient.FtpCli.FTPClient.Stor(csvName, passengerdata)
	//ftpclient.FtpCli.Unlock()
	//if err != nil {
	//	fmt.Println("PassengerFlowTime info to ftp err:", err)
	//}
	msg := ftpclient.Message{
		FileName: csvName,
		SaveMsg: csvContent,
	}
	ftpclient.FtpMsg <- msg
	return
}

type passengerFlowFrameData struct {
	DeviceID	string 	`json:"DeviceID" csv:"deviceID"`
	AlarmType 	int64 	`json:"AlarmType" csv:"alarmType"`
	AlarmDescribe	string `json:"AlarmDescribe" csv:"alarmDescribe"`
	StartTime 	string 	`json:"StartTime" csv:"startTime"`
	EnterNum	int64 	`json:"EnterNum" csv:"enterNum"`
	LeaveNum	int64	`json:"LeaveNum" csv:"leaveNum"`
	VideoChannel int8	`json:"VideoChannel" csv:"videoChannel"`
}



//按每张图接受客流信息
func (this *MainController) PassengerFlowFrame() {
	//fmt.Println("get gas danger info: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("get PassengerFlowFrame info:")
	logs.BeeLogger.Info("get PassengerFlowFrame info: %s", string(this.Ctx.Input.RequestBody))
	fmt.Println(string(this.Ctx.Input.RequestBody))
	data := passengerFlowFrameData{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.BeeLogger.Error("analysis PassengerFlowFrame error: %s", err)
		fmt.Println(err)
		this.Data["json"] = map[string]interface{}{
			"close": false,
			"code":  0,
		}
		this.ServeJSON()
		return
	}

	defer func() {
		this.Data["json"] = map[string]interface{}{
			"close": "true",
			"code":  1,
		}
		this.ServeJSON()
	}()
	psssenger := []*passengerFlowFrameData{}
	//passengerFlowi, ok := passengerFlowM.Load(data.DeviceID)

	//if ok {
	//	passengerFlows := passengerFlowi.(passengerFlowFrameData)
	//	passengerFlowM.Delete(data.DeviceID)
	//	psssenger = append(psssenger, &passengerFlows)
		psssenger = append(psssenger, &data)

		csvContent, _ := gocsv.MarshalString(&psssenger) // Get all clients as CSV string
		fmt.Println(csvContent) // Display all clients as CSV string
		logs.BeeLogger.Info("PassengerFlowFrame  data csv info %s", csvContent)
		//passengerdata := bytes.NewBufferString(csvContent)
		csvName := "passengerFlow/byFrame_" + data.DeviceID  + "_" + tools.GetUUID() + ".csv"
		//err = ftpclient.FTPClient.Stor(csvName, passengerdata)
		//ftpclient.FtpCli.Lock()
		//err = ftpclient.FtpCli.FTPClient.Stor(csvName, passengerdata)
		//ftpclient.FtpCli.Unlock()
		//if err != nil {
		//	fmt.Println("PassengerFlowFrame info to ftp err:",err)
		//}
		msg := ftpclient.Message{
			FileName: csvName,
			SaveMsg: csvContent,
		}
		ftpclient.FtpMsg <- msg

	//} else {
	//	passengerFlowM.Store(data.DeviceID,data)
	//}

	return
}

type facePic struct {
	DeviceID	string	`json:"DeviceID" csv:"deviceID"`
	PictureTime	string	`json:"PictureTime" csv:"pictureTime"`
	FileName	string	`json:"FileName" csv:"fileName"`
	PictureData	string	`json:"PictureData" csv:"pictureData"`
}

//接受人脸图片
func (this *MainController) Faces() {
	//fmt.Println("get gas danger info: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("get faces info:")
	logs.BeeLogger.Info("get Faces info: %s", string(this.Ctx.Input.RequestBody))
	fmt.Println(string(this.Ctx.Input.RequestBody))
	data := facePic{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.BeeLogger.Error("analysis Faces error: %s", err)
		fmt.Println(err)
		this.Data["json"] = map[string]interface{}{
			"close": false,
			"code":  0,
		}
		this.ServeJSON()
		return
	}

	defer func() {
		this.Data["json"] = map[string]interface{}{
			"close": "true",
			"code":  1,
		}
		this.ServeJSON()
	}()
	picDataStr,err := base64.StdEncoding.DecodeString(data.PictureData)
	if err != nil {
		logs.BeeLogger.Info("decode picture err:", err)
	}
	//picData := bytes.NewBuffer(picDataStr)
	csvName := "faces/" + data.DeviceID + "_"+ data.PictureTime + "_" +data.FileName + "_" + tools.GetUUID()  + ".jpg"
	//err = ftpclient.FTPClient.Stor(csvName, picData)
	//ftpclient.FtpCli.Lock()
	//err = ftpclient.FtpCli.FTPClient.Stor(csvName, picData)
	//ftpclient.FtpCli.Unlock()
	//if err != nil {
	//	fmt.Println("Faces info to ftp err:",err)
	//}
	msg := ftpclient.Message{
		FileName: csvName,
		SaveMsg: string(picDataStr),
	}
	ftpclient.FtpMsg <- msg
	return
}

type mpdcDataMes struct {
	DeviceID	string 	`json:"DeviceID" csv:"deviceID"`
	Index 		string	`json:"Index" csv:"index"`
	StartTime 	string 	`json:"StartTime" csv:"startTime"`
	StopTime 	string 	`json:"StopTime" csv:"stopTime"`
	EnterNum	string 	`json:"EnterNum" csv:"enterNum"`
	LeaveNum	string	`json:"LeaveNum" csv:"leaveNum"`
	VideoChannel string	`json:"VideoChannel" csv:"videoChannel"`
	Count		string 	`json:"Count" csv:"count"`
	Level		string	`json:"Level" csv:"level"`
	RetransFlag int64	`json:"RetransFlag" csv:"retransFlag"`
}

//按时间接受客流信息
func (this *MainController) MPDCData() {
	//fmt.Println("get gas danger info: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("get PassengerFlow info:")
	logs.BeeLogger.Info("get PassengerFlow info: %s", string(this.Ctx.Input.RequestBody))
	fmt.Println(string(this.Ctx.Input.RequestBody))
	data := mpdcDataMes{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.BeeLogger.Error("PassengerFlowTime stop error: %s", err)
		fmt.Println(err)
		this.Data["json"] = map[string]interface{}{
			"close": false,
			"code":  0,
		}
		this.ServeJSON()
		return
	}

	defer func() {
		this.Data["json"] = map[string]interface{}{
			"close": "true",
			"code":  1,
		}
		this.ServeJSON()
	}()

	psssenger := []*mpdcDataMes{}
	psssenger = append(psssenger, &data)
	csvContent, _ := gocsv.MarshalString(&psssenger) // Get all clients as CSV string
	fmt.Println(csvContent) // Display all clients as CSV string
	logs.BeeLogger.Info("PassengerFlowTime  data csv info %s", csvContent)
	//passengerdata := bytes.NewBufferString(csvContent)
	csvName := "passengerFlow/" +  data.DeviceID + "_"+ data.VideoChannel + "_" + tools.GetUUID() + ".csv"
	msg := ftpclient.Message{
		FileName: csvName,
		SaveMsg: csvContent,
	}
	ftpclient.FtpMsg <- msg
	return
}

//func (self *HookController) HookPublish() {
//	fmt.Println("hook_on_publish: ", string(self.Ctx.Input.RequestBody))
//	replyData := map[string]interface{}{
//		"code":       0,
//		"enableHls":  true,
//		"enableMP4":  false,
//		"enableRtxp": true,
//		"msg":        "success",
//	}
//	//解析接收的内容
//	data := struct {
//		App string `json:"app"` //流应用名
//	}{}
//	err := json.Unmarshal(self.Ctx.Input.RequestBody, &data)
//	if err != nil {
//		BeeLogger.Error("analysis hook on_publish error: %s", err)
//	} else {
//		if data.App == "download" {
//			replyData["enableMP4"] = true
//		}
//	}
//
//	self.Data["json"] = replyData
//	self.ServeJSON()
//}
//
//func (self *HookController) HookRecordMP4() {
//	fmt.Println("hook_on_record_mp4: ", string(self.Ctx.Input.RequestBody))
//	data := struct {
//		App    string `json:"app"`    //流应用名
//		Stream string `json:"stream"` //流ID
//		URL    string `json:"url"`
//		Vhost  string `json:"vhost"`
//	}{}
//	err := json.Unmarshal(self.Ctx.Input.RequestBody, &data)
//	if err != nil {
//		BeeLogger.Error("analysis hook on_record_mp4 error: %s", err)
//		return
//	}
//
//	if data.App == "download" {
//		//用户token(用户token是32的UUID，国标设备编码是20位)
//		if len(data.Stream) < 52 {
//			return
//		}
//		userToken := data.Stream[len(data.Stream)-32:]
//		ip, port := tcpServers.GetStreamIPAndPort("hls")
//		videoURL := utils.StringsJoin("http://", ip, ":", port, "/", data.URL, "?token=", userToken)
//
//		SetPlaybackDownloadVideo(userToken, data.Stream, videoURL)
//	}
//}
//
////保存录像下载的内容
//func SetPlaybackDownloadVideo(userToken, streamID, videoURL string) {
//	userAuth, found := authCache.AuthCache.Get(userToken)
//	if found {
//		userAuthTemp := userAuth.(*authCache.UserAuth)
//		userAuthTemp.Download[streamID] = videoURL
//		authCache.AuthCache.Set(userToken, userAuthTemp, cache.DefaultExpiration)
//	}
//}
//
////访问http文件服务器上hls之外的文件时触发
//func (self *HookController) HookHttpAccess() {
//	fmt.Println("hook_on_http_access: ", string(self.Ctx.Input.RequestBody))
//	data := struct {
//		Params string `json:"params"`
//	}{}
//	err := json.Unmarshal(self.Ctx.Input.RequestBody, &data)
//	if err != nil {
//		BeeLogger.Error("analysis hook on_http_access error: %s", err)
//		return
//	}
//
//	_, found := authCache.AuthCache.Get(data.Params[len("token="):])
//	if found {
//		//有效的用户token
//		self.Data["json"] = map[string]interface{}{
//			"code":   0,
//			"err":    "",
//			"path":   "",
//			"second": 3600,
//		}
//		self.ServeJSON()
//	}
//}
