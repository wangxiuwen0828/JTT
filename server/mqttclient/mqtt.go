//package main
//
//import (
//	"encoding/json"
//	"fmt"
//	mqtt "github.com/eclipse/paho.mqtt.golang"
//	"time"
//)
//
////创建全局的mqtt publish消息处理handler
//var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
//	//打印topic主题名
//	fmt.Printf("publish Client Topic: %s\n", message.Topic())
//	//打印内容
//	fmt.Printf("publish Client Msg: %s\n", message.Payload())
//}
//
////创建全局的mqtt subscribe消息处理handler
//var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
//	//打印topic主题名
//	fmt.Printf("subscribe Client Topic: %s\n", message.Topic())
//	//打印内容
//	fmt.Printf("subscribe Client Msg: %s\n", message.Payload())
//	switch message.Topic() {
//	case "ha":
//		fmt.Println("[subscribe],ha")
//	case "ya":
//		fmt.Println("[subscribe],ya")
//	case "hi":
//		hh := struct {
//			Topic string `json:"topic"`
//			Num   int    `json:"num"`
//		}{}
//
//		if err := json.Unmarshal(message.Payload(), &hh); err != nil {
//			fmt.Println("err:", err)
//		} else {
//			fmt.Println("hi:", hh.Topic, hh.Num)
//		}
//	}
//}
//
//var (
//	disconnectChan    = make(chan bool) //mqtt服务端重启或者断开时用于重新连接
//	connectRetries    = 10              //客户端连接失败重试次数
//	connectRetryDelay = time.Second     //客户端重新连接间隔时间
//)
//
////消息格式，发布消息使用
//type message struct {
//	topic   string      //发布主题
//	payload interface{} //发布内容
//}
//
////客户端管理器
//type MQTTClientManger struct {
//	client   mqtt.Client
//	msgSend  chan message
//	topicSub []string //订阅Subscribe主题topic列表
//}
//
////客户端连接
//func (mg *MQTTClientManger) Connect() bool {
//	var err error
//	for retries := 0; retries < connectRetries; retries++ {
//		if token := mg.client.Connect(); token.Wait() && token.Error() != nil {
//			if retries == connectRetries-1 {
//				err = token.Error()
//			}
//			time.Sleep(connectRetryDelay)
//		} else {
//			return true
//		}
//	}
//
//	fmt.Printf("mqtt connect error: %s\n", err)
//	return false
//}
//
////发布消息
//func (mg *MQTTClientManger) Publish() {
//	for {
//		msg, ok := <-mg.msgSend
//		if ok {
//			//格式化数据，将信息转换为json
//			payload, err := json.Marshal(msg.payload)
//			if err != nil {
//				fmt.Printf("[publish] json.Marshal() error: %s\n", err)
//			}
//			fmt.Printf("%s 打印数据：%s\n", time.Now().Format("2006-01-02 15:04:05"), string(payload))
//			token := mg.client.Publish(msg.topic, 1, false, payload)
//			token.Wait()
//		}
//	}
//}
//
////订阅消息
//func (mg *MQTTClientManger) Subscribe() {
//	for _, topic := range mg.topicSub {
//		token := mg.client.Subscribe(topic, 1, messageSubHandler)
//		token.Wait()
//	}
//}
//
////启动服务
//func (mg *MQTTClientManger) run() {
//	go mg.Subscribe()
//	go mg.Publish()
//}
//
//func NewMQTTClient(server string, topicSub []string) *MQTTClientManger {
//	opts := mqtt.NewClientOptions().AddBroker(server)
//	//设置客户端ID
//	opts.SetClientID("wyd666")
//	opts.SetKeepAlive(30 * time.Second)
//	opts.SetPingTimeout(10 * time.Second)
//	//设置handler
//	opts.SetDefaultPublishHandler(messagePubHandler)
//	opts.SetConnectionLostHandler(func(client mqtt.Client, e error) {
//		fmt.Printf("connect lost,reason:%s\n", e)
//		disconnectChan <- true
//	})
//	//设置是否清空session，这里如果设置为false表示服务器会保留客户端的连接记录，这里设置为true表示每次连接到服务器都以新的身份连接
//	opts.SetCleanSession(false)
//	////设置自动重连
//	//opts.SetAutoReconnect(true)
//
//	return &MQTTClientManger{
//		client:   mqtt.NewClient(opts), //创建客户端连接
//		msgSend:  make(chan message),
//		topicSub: topicSub,
//	}
//}
//
//func main() {
//	topicSub := []string{"ha", "ya", "hi", "$SYS/broker/clients/total"}
//	mqttClient := NewMQTTClient("tcp://127.0.0.1:1883", topicSub)
//	if mqttClient.Connect() {
//		fmt.Println("mqtt connect success")
//		//连接成功
//		mqttClient.run()
//		for num, topic := range topicSub {
//			msg := message{
//				topic: topic,
//				payload: map[string]interface{}{
//					"topic": topic,
//					"num":   num,
//				},
//			}
//			mqttClient.msgSend <- msg
//		}
//	} else {
//		fmt.Println("mqtt connect failed")
//		return
//	}
//	go func() {
//		i := 0
//		ticker := time.NewTicker(3 * time.Second)
//		defer ticker.Stop()
//		for {
//			select {
//			case <-ticker.C:
//				i++
//
//				if i >= 6 {
//					msg := message{
//						topic: "hi",
//						payload: map[string]interface{}{
//							"topic": "hi",
//							"num":   i,
//						},
//					}
//					mqttClient.msgSend <- msg
//				}
//			}
//		}
//	}()
//	for {
//		select {
//		case _, ok := <-disconnectChan:
//			//检测到mqtt服务器重启时，进行重连
//			if ok {
//				if mqttClient.Connect() {
//					mqttClient.Subscribe()
//				}
//			}
//		}
//	}
//}

package mqttclient

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitee.com/ictt/JTTM/config"
	"gitee.com/ictt/JTTM/server"
	"gitee.com/ictt/JTTM/server/udp"
	"gitee.com/ictt/JTTM/tools/logs"
	"gitee.com/ictt/JTTM/tools/sqlDB"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"time"
)

//创建全局的mqtt publish消息处理handler
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	//打印topic主题名
	fmt.Printf("publish Client Topic: %s\n", message.Topic())
	//打印内容
	fmt.Printf("publish Client Msg: %s\n", message.Payload())
}


var (
	disconnectChan = make(chan bool) //mqtt服务端重启或者断开时用于重新连接
	clientChan = make(chan bool, 2)	 //mqtt另一边客户端断开时重新连接
	MqttClient     *MQTTClient
	mqttSubID = "mqttPub/DASID=" + config.MQTTID
	//MattPubID = "mqttSubReq/clientID=" + config.MQTTID
	mqttHereStatus = "mqttSubClientStatus"	//本地客户端状态
	mqttRemoteStatus = "mqttPubClientStatus"	//远程客户端状态的topic
	//RemoteClientStatus string		//远程客户端状态
	//Csq uint64 = 0
)

////消息格式，发布消息使用
//type Message struct {
//	Topic   string      //发布主题
//	Payload interface{} //发布内容
//}

//MQTT客户端
type MQTTClient struct {
	client   mqtt.Client
	//MsgSend  chan Message
	topicSub []string //订阅Subscribe主题topic列表
}

type mqttRegisterPub struct {
	Id				string	`json:"id"`
	Name			string	`json:"name"`
	Ip				string	`json:"ip"`
	Port			string	`json:"port"`
	TerminalType	int64	`json:"terminaltype"`
	Csq				uint64	`json:"csq"`
	Status			string	`json:"status"`
}


//type MQTTDeviceList struct {
//	Protocol		string	`json:"protocol"`
//	Csq				uint64	`json:"csq"`
//	DeviceLists		[]MQTTDeviceInfo	`json:"devicelists"`
//}
//
//type MQTTDeviceInfo struct {
//	Id				string	`json:"id"`
//	Name			string	`json:"name"`
//	Status			string	`json:"status"`
//}
//
//type MQTTChannelList struct {
//	Protocol		string	`json:"protocol"`
//	Csq				uint64	`json:"csq"`
//	ChannelLists	[]MQTTChannelInfo	`json:"devicelists"`
//}
//
//type MQTTChannelInfo struct {
//	Id				string	`json:"id"`
//	Name			string	`json:"name"`
//	Status			string	`json:"status"`
//	Parentid		string	`json:"parentid"`
//}
//
//type MQTTNotify struct {
//	Protocol		string	`json:"protocol"`
//	Csq				uint64	`json:"csq"`
//	Type			string	`json:"type"`	//不同类型的通知：Device,Channel,Alarm,Record
//	Data 			interface{}	`json:"data"`
//}

type mqttRecived struct {
	Protocol		string	`json:"protocol"`
}

type mqttRegisterRes struct {
	Protocol		string	`json:"protocol"`
	ErrCode			int16	`json:"errcode"`
	ErrMsg			string	`json:"errmsg"`
}

type mqttRealPlayReceived struct {
	Action			string	`json:"action"` //start开始，stop结束
	Csq				uint64	`json:"csq"`
	Id				string	`json:"id"`
	Ip				string	`json:"ip"`
	Port			string	`json:"port"`
	Protocol		string	`json:"protocol"`
	Type			string	`json:"type"`
}

type mqttRecordPlayReceived struct {
	Action			string	`json:"action"` //start开始，stop结束
	Begintime		string	`json:"begintime"`
	Endtime			string	`json:"endtime"`
	Csq				uint64	`json:"csq"`
	Id				string	`json:"id"`
	Ip				string	`json:"ip"`
	Port			string	`json:"port"`
	Protocol		string	`json:"protocol"`
	Type			string	`json:"type"`
}

type mqttDownLoadReceived struct {
	Action			string	`json:"action"` //start开始，stop结束
	Begintime		string	`json:"begintime"`
	Endtime			string	`json:"endtime"`
	Csq				uint64	`json:"csq"`
	Id				string	`json:"id"`
	Ip				string	`json:"ip"`
	Port			string	`json:"port"`
	Protocol		string	`json:"protocol"`
	Speed 			uint16	`json:"speed"`
	Type			string	`json:"type"`
}

type mqttRecordListsReceived struct {
	Begintime		string	`json:"begintime"`
	Endtime			string	`json:"endtime"`
	Csq				uint64	`json:"csq"`
	Id				string	`json:"id"`
	Protocol		string	`json:"protocol"`
	Type			string	`json:"type"`
}



//客户端连接
func (mg *MQTTClient) Connect() bool {
	i := 0
	for {
		if token := mg.client.Connect(); token.Wait() && token.Error() == nil {
			fmt.Printf("%s mqtt client connection successful\n", time.Now().Format("2006-01-02 15:04:05"))
			return true
		}
		i += 1
		fmt.Println("i=", i)
		time.Sleep(15 * time.Second)
	}

	return false
}

//发布消息
func (mg *MQTTClient) Publish() {
	for {
		msg, ok := <-server.MqttMess
		if ok {
			//格式化数据，将信息转换为json
			payload, err := json.MarshalIndent(msg.Payload,"  ", "  ")
			if err != nil {
				fmt.Printf("[publish] json.Marshal() error: %s\n", err)
			}
			fmt.Printf("%s 打印数据：%s\n", time.Now().Format("2006-01-02 15:04:05"), string(payload))
			token := mg.client.Publish(msg.Topic, 1, false, payload)
			token.Wait()
		}
	}
}

//订阅消息
func (mg *MQTTClient) Subscribe() {
	for _, topic := range mg.topicSub {
		token := mg.client.Subscribe(topic, 1, messageSubHandler)
		token.Wait()
	}
}

//取消订阅
func (mg *MQTTClient) Unsubscribe() {
	for _, topic := range mg.topicSub {
		token := mg.client.Unsubscribe(topic)
		token.Wait()
	}
}

//启动服务
func (mg *MQTTClient) run() {
	go mg.Subscribe()
	go mg.Publish()
}

func NewMQTTClient(servers string, topicSub []string) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker(servers)
	//设置客户端ID
	opts.SetClientID(config.MQTTID)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	//设置handler
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetConnectionLostHandler(func(client mqtt.Client, e error) {
		fmt.Printf("connect lost,reason:%s\n", e)
		disconnectChan <- true
	})
	//设置是否清空session，这里如果设置为false表示服务器会保留客户端的连接记录，这里设置为true表示每次连接到服务器都以新的身份连接
	opts.SetCleanSession(true)
	////设置自动重连
	//opts.SetAutoReconnect(true)
	//设置遗嘱消息

	lastWillData := mqttRegisterPub{
		Id: config.MQTTID,
		Name: "JTSGZ",
		Ip: config.IP,
		TerminalType: 1,
		Csq: server.Csq,
		Status: "OFF",
	}

	//lastWillData := map[string]interface{}{
	//	"serial": config.MQTTID,
	//	"Status":   "OFF",
	//}
	lastWillDataByte, _ := json.MarshalIndent(lastWillData, "", "	")
	opts.SetWill(mqttHereStatus, string(lastWillDataByte), 1, false)

	return &MQTTClient{
		client:   mqtt.NewClient(opts), //创建客户端连接
		//MsgSend:  make(chan Message),
		topicSub: topicSub,
	}
}

func MQTTInit() {
	//该客户端需要订阅的主题，当收到其他客户端的注册请求时进行订阅这些主题
	//订阅主题格式："xxx/clientID=yyy"，每个clientID对应一个mqtt客户端
	//发布主题格1："xxx"，面向所有的mqtt客户端
	//发布主题格式2："xxx/clientID=yyy"，面向指定mqtt客户端发布消息
	topicSub := []string{mqttSubID, mqttRemoteStatus}
	mqttSeverAddr := "tcp://" + config.MQTTAddr
	MqttClient = NewMQTTClient(mqttSeverAddr, topicSub)
	if MqttClient.Connect() {
		MqttClient.Subscribe()
		//fmt.Println("*********1111**********")
		//连接成功，发布消息，允许其他客户端进行注册请求
		go MqttClient.Publish()
		go register()
	}
	//fmt.Println("hahaha")
	for {
		select {
		case _, ok := <-disconnectChan:
			//检测到mqtt服务器重启时，进行重连
			if ok {
				if MqttClient.Connect() {
					payload := mqttRegisterPub{
						Id: config.MQTTID,
						Name: "JTSGZ",
						Ip: config.IP,
						TerminalType: 1,
						Csq: server.Csq,
						Status: "ON",
					}
					server.Csq++
					msg := server.Message{
						Topic:   "mqttSubClientStatus",
						Payload: payload,
					}
					server.MqttMess<- msg
				}
			}
		}
	}
}


//创建全局的mqtt subscribe消息处理handler
var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	//打印topic主题名
	fmt.Printf("subscribe Client Topic: %s\n", message.Topic())
	//打印内容
	fmt.Printf("subscribe Client Msg: %s\n", message.Payload())
	switch message.Topic() {
	case mqttSubID:
		//fmt.Println("[subscribe],ha")
		reciveData := mqttRecived{}
		err := json.Unmarshal(message.Payload(),&reciveData)
		fmt.Println("解析成功")
		if err !=  nil {
			fmt.Println("收到的数据与本地定义的不匹配")
			logs.BeeLogger.Error("unmarshal MQTT Subscribe Mes err : %s", err)
			return
		}

		switch reciveData.Protocol {
		//
		case "RegisterPub":
			data := mqttRegisterRes{}
			err := json.Unmarshal(message.Payload(), &data)
			if err != nil {
				fmt.Println("收到的数据与本地定义的不匹配")
				logs.BeeLogger.Error("unmarshal MQTT Subscribe Mes err : %s", err)
				return
			}
			if data.ErrCode == 200 {
				clientChan <- true
				//_,ok := mqttCache.Get(config.MQTTID)
				//if !ok {
				//	mqttCache.Set(config.MQTTID, config.MQTTID, cache.DefaultExpiration)
				//}
				//go keepalive()
				server.RemoteClientStatus = "ON"

				//udp.SetDeviceParameter("123456789012",1,"00000014")
				go giveDeviceList()
				go giveChannelList()
			}

		case "RealPlayPub":
			data := mqttRealPlayReceived{}
			err := json.Unmarshal(message.Payload(), &data)
			if err != nil {
				fmt.Println("收到的数据与本地定义的不匹配")
				logs.BeeLogger.Error("unmarshal MQTT Subscribe Mes err : %s", err)
				return
			}
			if data.Action == "start" {
				phoneNum := data.Id[:12]
				channel, _ := hex.DecodeString(data.Id[12:])
				errCode, _ := udp.RequestRealStream(phoneNum, channel[0],1,0)
				if errCode == config.JTT_ERROR_SUCCESS_OK{
					res := server.MqttRealPlayRes{
						Protocol: "RealPlaySub",
						Csq:      data.Csq,
						Streamid: data.Id + "@" + "test",
						Errcode:  config.JTT_ERROR_SUCCESS_OK,
						Errmsg:   config.ErrCodeMap[config.JTT_ERROR_SUCCESS_OK],
					}
					msg := server.Message{
						Topic:   server.MattPubID,
						Payload: res,
					}
					server.MqttMess <- msg
				}
			}


		case "RecordPlayPub":
			data := mqttRecordPlayReceived{}
			err := json.Unmarshal(message.Payload(), &data)
			if err != nil {
				fmt.Println("收到的数据与本地定义的不匹配")
				logs.BeeLogger.Error("unmarshal MQTT Subscribe Mes err : %s", err)
				return
			}

			if data.Action == "start" {
				phoneNum := data.Id[:12]
				channel, _ := hex.DecodeString(data.Id[12:])
				begin := data.Begintime[:4] + data.Begintime[5:7] + data.Begintime[8:10] + data.Begintime[11:13] + data.Begintime[14:16] + data.Begintime[17:]
				end := data.Endtime[:4] + data.Endtime[5:7] + data.Endtime[8:10] + data.Endtime[11:13] + data.Endtime[14:16] + data.Endtime[17:]
				errCode, _ := udp.StartReplayVideoSend(strconv.FormatUint(data.Csq,10), phoneNum, channel[0],1,0,0,0,0, begin, end)
				if errCode == config.JTT_ERROR_SUCCESS_OK{
					res := server.MqttRealPlayRes{
						Protocol: "RecordPlaySub",
						Csq:      data.Csq,
						Streamid: data.Id + "@" + "test",
						Errcode:  config.JTT_ERROR_SUCCESS_OK,
						Errmsg:   config.ErrCodeMap[config.JTT_ERROR_SUCCESS_OK],
					}
					msg := server.Message{
						Topic:   server.MattPubID,
						Payload: res,
					}
					server.MqttMess <- msg
				}
			}


		case "DownLoadPub":
			data := mqttDownLoadReceived{}
			err := json.Unmarshal(message.Payload(), &data)
			if err != nil {
				fmt.Println("收到的数据与本地定义的不匹配")
				logs.BeeLogger.Error("unmarshal MQTT Subscribe Mes err : %s", err)
				return
			}
			if data.Action == "start" {
				res := server.MqttRealPlayRes{
					Protocol: "DownLoadSub",
					Csq:      data.Csq,
					Streamid: data.Id + "@" + "test",
					Errcode:  config.JTT_ERROR_SUCCESS_OK,
					Errmsg:   config.ErrCodeMap[config.JTT_ERROR_SUCCESS_OK],
				}
				msg := server.Message{
					Topic:   server.MattPubID,
					Payload: res,
				}
				server.MqttMess <- msg
			}


		case "RecordListsPub":
			data := mqttRecordListsReceived{}
			err := json.Unmarshal(message.Payload(), &data)
			if err != nil {
				fmt.Println("收到的数据与本地定义的不匹配")
				logs.BeeLogger.Error("unmarshal MQTT Subscribe Mes err : %s", err)
				return
			}
			res := server.MqttRecordListsRes{
				Protocol: "RecordListsSub",
				Csq:      data.Csq,
				Streamid: data.Id + "@" + "test",
				Errcode:  config.JTT_ERROR_SUCCESS_OK,
				Errmsg:   config.ErrCodeMap[config.JTT_ERROR_SUCCESS_OK],
			}
			res.Recordlists = []server.MqttRecordLists{}
			test := server.MqttRecordLists{
				Name: "test",
				Begintime: time.Now().Format("2006-01-02T15:04:05"),
				Endtime: time.Now().Format("2006-01-02T15:04:05"),
				Type: "manual",
				Filesize: "1024",
				Filepath: "/home",
			}
			res.Recordlists = append(res.Recordlists, test)

			msg := server.Message{
				Topic:   server.MattPubID,
				Payload: res,
			}
			server.MqttMess <- msg

		}




		case mqttRemoteStatus:
			data := mqttRegisterPub{}
			err := json.Unmarshal(message.Payload(),&data)
			fmt.Println("解析成功")
			if err !=  nil {
				fmt.Println("收到的数据与本地定义的不匹配")
				logs.BeeLogger.Error("unmarshal MQTT Subscribe Mes err : %s", err)
				return
			}

			if data.Status == "OFF" {
				server.RemoteClientStatus = data.Status
				go register()
			}

	}
}

func register() {
	server.Csq++
	payload := mqttRegisterPub{
		Id: config.MQTTID,
		Name: "JTSGZ",
		Ip: config.IP,
		TerminalType: 1,
		Csq: server.Csq,
		Status: "ON",
	}
	server.Csq++
	msg :=server. Message{
		Topic:   "mqttSubClientStatus",
		Payload: payload,
	}
	server.MqttMess <- msg

	t1 := time.NewTicker(60 * time.Second)
	defer t1.Stop()

	for {
		select {
		case _, ok := <-clientChan:
			//检测到mqtt服务器重启时，进行重连
			if ok {
				//fmt.Println("qq")
				return
			}
		case <-t1.C:
			payload.Csq = server.Csq
			server.Csq ++
			msg := server.Message{
				Topic:   "mqttSubClientStatus",
				Payload: payload,
			}
			server.MqttMess <- msg
		}
	}
}

//func keepalive() {
//	for {
//		payload := PayloadData{}
//		payload.IVMS.Header.Protocol = "iVMS_Msg_Type_KeepAliveReq"
//		payload.IVMS.Body = map[string]interface{}{
//			"keepalive": "keepalive",
//		}
//		msg := Message{
//			Topic:   MattPubID,
//			Payload: payload,
//		}
//		MqttClient.MsgSend <- msg
//		time.Sleep(15*time.Second)
//		_,ok := mqttCache.Get(config.MQTTID)
//		if !ok {
//			go register()
//			return
//		}
//	}
//}

func giveDeviceList() {
	var deviceList []sqlDB.GetVehicleInfo
	sqlDB.Find(&deviceList, &sqlDB.GetVehicleInfo{})
	//fmt.Println(deviceList)
	mqttDevice  := server.MQTTDeviceInfo{}
	//var mqttDeviceList  []MQTTDevice
	mqttDeviceList  := server.MQTTDeviceList{}
	mqttDeviceList.DeviceLists = []server.MQTTDeviceInfo{}
	mqttDeviceList.Protocol = "DeviceListsSub"
	mqttDeviceList.Csq = server.Csq
	server.Csq++

	for _,v := range deviceList {
		mqttDevice.Id = v.PhoneNum
		mqttDevice.Name = v.VehicleID
		mqttDevice.Status = v.Status
		mqttDeviceList.DeviceLists = append(mqttDeviceList.DeviceLists, mqttDevice)
	}

	msg := server.Message{
		Topic:   server.MattPubID,
		Payload: mqttDeviceList,
	}
	server.MqttMess <- msg
}

func giveChannelList() {
	var channelList []sqlDB.GetChannelInfo
	sqlDB.Find(&channelList, &sqlDB.GetChannelInfo{})
	//fmt.Println(channelList)
	mqttChannelInfo := server.MQTTChannelInfo{}
	mqttChannelList := server.MQTTChannelList{}
	mqttChannelList.ChannelLists = []server.MQTTChannelInfo{}
	mqttChannelList.Protocol = "ChannelListsSub"
	mqttChannelList.Csq = server.Csq
	server.Csq++

	for _,v := range channelList {
		mqttChannelInfo.Id = v.PhoneNumAndChannel
		mqttChannelInfo.Parentid = v.PhoneNum
		mqttChannelInfo.Name = strconv.FormatInt(v.LogicalChannelID, 10)
		mqttChannelInfo.Status = v.Status
		mqttChannelInfo.Model = "车载摄像头"

		mqttChannelList.ChannelLists = append(mqttChannelList.ChannelLists, mqttChannelInfo)
	}

	msg := server.Message{
		Topic:   server.MattPubID,
		Payload: mqttChannelList,
	}
	server.MqttMess <- msg
}



//type RegisterBodyData struct {
//	Serial	string `json:"serial"`
//}