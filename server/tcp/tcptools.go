package tcp
//
//import (
//	"encoding/binary"
//	"fmt"
//	"gitee.com/ictt/JTTM/server"
//	"gitee.com/ictt/JTTM/tools"
//	"gitee.com/ictt/JTTM/tools/logs"
//	"github.com/patrickmn/go-cache"
//	"net"
//	"strconv"
//	"time"
//)
//
//var (
//	//记录tcp连接信息，当设备注册成功时存储记录，键为Device，值为map键值对
//	//map键值对的键为UDP连接的IP+Port，值为tcpClient
//	tcpConnCache *cache.Cache
//)
//
//type tcpClient struct {
//	PhoneNum string //设备ID
//	TCPAddr	   string
//	TCPConn    *net.TCPConn
//	UpdateTime int64
//	Sequence   uint16 //记录请求设备配置时生成的随机数，Sequence+2作为我红方向蓝方发送请求信令时使用
//}
//
//func init() {
//	tcpConnCache = cache.New(90*time.Second, 30*time.Second)
//}
//
////存储一个udpClient，初次存储是设备发送注册请求
//func setTCPClient(PhoneNum string, client *tcpClient) {
//	//fmt.Println("存储信息")
//	tcpConnCache.Set(PhoneNum, client, cache.DefaultExpiration)
//}
//
////使用PhoneNum筛选出指定的TCPClient
//func getTCPClient(PhoneNum string) *tcpClient {
//	temp, ok := tcpConnCache.Get(PhoneNum)
//	if ok {
//		return temp.(*tcpClient)
//	}
//
//	logs.BeeLogger.Error("deviceID=%s is disconnect or for other reasons", PhoneNum)
//	return nil
//}
//
//
////通过conn获取phoneNum
//func getPhoneFromConn(tcpAddr  string) string {
//	for _, v := range tcpConnCache.Items() {
//		client := v.Object.(*tcpClient)
//		if client.TCPAddr == tcpAddr {
//			return client.PhoneNum
//		}
//	}
//	return ""
//}
//
////处理消息体属性
//func handleMessProperty(messByte []byte) (int64, bool) {
//	header := server.HeaderReceive{}
//	sign := true
//	//needRepeat := false
//	//header.IDReceive = hex.EncodeToString(messByte[:2])
//	dataLength := binary.BigEndian.Uint16(messByte)
//	messHead2 := strconv.FormatInt(int64(dataLength), 2)
//	for len(messHead2) < 16 {
//		messHead2 = "0" + messHead2
//	}
//	header.MessProperty.Subpackage = messHead2[2]
//	header.MessProperty.Encryption = []byte{messHead2[3], messHead2[4], messHead2[5]}
//	header.MessProperty.DataLength = tools.TwoToInt64(messHead2[6:])
//	fmt.Printf("消息体长度%v\n", header.MessProperty.DataLength)
//	switch header.MessProperty.Subpackage {
//	case '0':
//		fmt.Println("package no change")
//
//	case '1':
//		sign = false
//		fmt.Println("need change")
//		header.MessProperty.DataLength += 4
//	}
//	return header.MessProperty.DataLength, sign
//}
//
//func checkTCPConnect (phoneNum string, tcpConn *net.TCPConn) {
//
//}
//
