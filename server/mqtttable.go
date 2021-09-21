package server

import (
	//"encoding/json"
	//"fmt"
	"gitee.com/ictt/JTTM/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	//"time"
)

var (
	//disconnectChan = make(chan bool) //mqtt服务端重启或者断开时用于重新连接
	//clientChan = make(chan bool, 2)	 //mqtt另一边客户端断开时重新连接
	//MqttClient     *MQTTClient
	MqttMess = make(chan Message)
	//mqttSubID = "mqttPubRes/clientID=" + config.MQTTID
	MattPubID = "mqttSub/DASID=" + config.MQTTID
	//mqttHereStatus = "mqttSubClientStatus"	//本地客户端状态
	//mqttRemoteStatus = "mqttPubClientStatus"	//远程客户端状态的topic
	RemoteClientStatus string		//远程客户端状态
	Csq uint64 = 0
)

//消息格式，发布消息使用
type Message struct {
	Topic   string      //发布主题
	Payload interface{} //发布内容
}

//MQTT客户端
type MQTTClient struct {
	Client   mqtt.Client
	MsgSend  chan Message
	TopicSub []string //订阅Subscribe主题topic列表
}

type MQTTDeviceList struct {
	Protocol		string              `json:"protocol"`
	Csq				uint64           `json:"csq"`
	DeviceLists		[]MQTTDeviceInfo `json:"devicelists"`
}

type MQTTDeviceInfo struct {
	Id				string	`json:"id"`
	Name			string	`json:"name"`
	Status			string	`json:"status"`
}

type MQTTChannelList struct {
	Protocol		string            `json:"protocol"`
	Csq				uint64         `json:"csq"`
	ChannelLists	[]MQTTChannelInfo `json:"channellists"`
}

type MQTTChannelInfo struct {
	Id				string	`json:"id"`
	Name			string	`json:"name"`
	Status			string	`json:"status"`
	Model			string	`json:"model"`
	Parentid		string	`json:"parentid"`
}

type MQTTNotify struct {
	Protocol		string	`json:"protocol"`
	Csq				uint64	`json:"csq"`
	Type			string	`json:"type"`	//不同类型的通知：Device,Channel,Alarm,Record
	Data 			interface{}	`json:"data"`
}

type MqttRecordListsRes struct {
	Protocol		string               `json:"protocol"`
	Streamid		string               `json:"streamid"`
	Csq				uint64            `json:"csq"`
	Errcode			int64             `json:"errcode"`
	Errmsg			string             `json:"errmsg"`
	Recordlists		[]MqttRecordLists `json:"recordlists"`
}

type MqttRecordLists struct {
	Name			string	`json:"name"`
	Begintime		string	`json:"begintime"`
	Endtime			string	`json:"endtime"`
	Type			string	`json:"type"`
	Filesize		string	`json:"filesize"`
	Filepath 		string 	`json:"filepath"`
}

type MqttRealPlayRes struct {
	Protocol		string	`json:"protocol"`
	Streamid		string	`json:"streamid"`
	Csq				uint64	`json:"csq"`
	Errcode			int64	`json:"errcode"`
	Errmsg			string	`json:"errmsg"`
}