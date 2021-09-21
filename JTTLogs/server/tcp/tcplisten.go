package tcp

import (
	"encoding/hex"
	"fmt"
	"gitee.com/ictt/JTTM/config"
	"gitee.com/ictt/JTTM/tools"
	"gitee.com/ictt/JTTM/tools/logs"
	"net"
	"time"
)

//func ListenTCPServer() {
//
//	listen, err := net.Listen("tcp", config.TCPAddr)
//	if err != nil {
//		fmt.Println("listen failed, err:", err)
//		return
//	}
//	fmt.Println("tcp server start")
//	defer listen.Close()
//	for {
//
//		conn, err := listen.Accept()
//		if err != nil {
//			fmt.Println("accept failed, err:", err)
//			//continue
//		}else {
//			fmt.Printf("accept success conn =%v, 客户端IP=%v\n", conn, conn.RemoteAddr().String())
//		}
//		//fmt.Println("第二步")
//		go process(conn)
//	}
//}
//
//func process(conn net.Conn) {
//	defer conn.Close()
//		realStream := server.RealStream{}
//		for {
//			headerBuf := make([]byte, 16)
//			n, err := conn.Read(headerBuf)
//			if err != nil {
//				fmt.Println("服务器的读取 err=", err)
//				return
//			}
//			fmt.Println(hex.EncodeToString(headerBuf[:n]))
//			lengthAddr := FormatAnalysis(headerBuf, &realStream.Header)
//			undeterminedBuf2 := make([]byte, lengthAddr+2)
//			n, err = conn.Read(undeterminedBuf2)
//			if err != nil {
//				fmt.Println("服务器的读取 err=", err)
//				return
//			}
//			fmt.Println(hex.EncodeToString(undeterminedBuf2[:n]))
//			UndeterminedField(undeterminedBuf2[:n], &realStream.Header)
//
//			dataBuf := make([]byte, realStream.Header.DataLength)
//			n, err = conn.Read(dataBuf)
//			if err != nil {
//				fmt.Println("服务器的读取 err=", err)
//				return
//			}
//			fmt.Println(n)
//			for n < int(realStream.Header.DataLength) {
//				buf4 := make([]byte, int(realStream.Header.DataLength)-n)
//				n1, err := conn.Read(buf4)
//				if err != nil {
//					fmt.Println("服务器的读取 err=", err)
//					return
//				}
//				dataBuf = append(dataBuf[:n], buf4...)
//				n = n + n1
//				//fmt.Println(buf3)
//			}
//
//			content := hex.EncodeToString(dataBuf)
//			//fmt.Println(content)
//			switch realStream.Header.DataAndPackage.Package {
//				case "0000":
//					fmt.Printf("原子包，不可拆分\n")
//					realStream.Content = content
//					fmt.Println(realStream.Content)
//				case "0001":
//					fmt.Printf("分包处理的第一个包\n")
//					realStream.Content = content
//				case "0010":
//					fmt.Printf("分包处理测最后一个包\n")
//					realStream.Content += content
//					fmt.Println(realStream.Content)
//				case "0011":
//					fmt.Printf("分包处理时的中间包\n")
//					realStream.Content += content
//			}
//
//	}
//}

func ListenTCPServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", config.TCPAddr)
	//fmt.Println(config.TCPAddr)
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("net.ResolveTCPAddr() error: %s", err))
	}
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("net.ListenTCP() error: %s", err))
	}

	logs.BeeLogger.Info("start TCPServer successful!")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "start TCPServer successful!")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer tcpListen.Close()
	for {
		tcpConn, err := tcpListen.AcceptTCP()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			//continue
		}else {
			fmt.Printf("accept success conn =%v, 客户端IP=%v\n", tcpConn, tcpConn.RemoteAddr().String())
		}
		//fmt.Println("第二步")
		go process(tcpConn)
	}
}

//func process(conn *net.TCPConn) {
//	defer conn.Close()
//	for {
//		headerBuf := make([]byte, 15)
//		_, err := conn.Read(headerBuf)
//		if err != nil {
//			fmt.Println("服务器的读取 err=", err)
//			return
//		}
//		headerBuf = bytes.ReplaceAll(headerBuf,[]byte{0x7d,0x02},[]byte{0x7e})
//		headerBuf = bytes.ReplaceAll(headerBuf,[]byte{0x7d,0x01},[]byte{0x7d})
//
//		//fmt.Println(hex.EncodeToString(headerBuf[:n]))
//		length,_ := handleMessProperty(headerBuf[3:5])
//		trueLength := (length + 12 + 3) - int64(len(headerBuf))
//		for trueLength != 0 {
//			DataBuf := make([]byte, trueLength)
//			_, err := conn.Read(DataBuf)
//			if err != nil {
//				fmt.Println("服务器的读取 err=", err)
//				return
//			}
//			DataBuf = append(headerBuf[len(headerBuf)-1:], DataBuf...)
//			DataBuf = bytes.ReplaceAll(DataBuf,[]byte{0x7d,0x02},[]byte{0x7e})
//			DataBuf = bytes.ReplaceAll(DataBuf,[]byte{0x7d,0x01},[]byte{0x7d})
//			trueLength = trueLength + 1 - int64(len(DataBuf))
//			headerBuf = append(headerBuf[:len(headerBuf)-1], DataBuf...)
//		}
//		fmt.Println(hex.EncodeToString(headerBuf))
//		//go tcpHandle(headerBuf,conn)
//	}
//}

func process(conn *net.TCPConn) {
	defer conn.Close()
	var moreBuf []byte
	for {
		headerBuf := make([]byte, 1024)
		n, err := conn.Read(headerBuf)
		if err != nil {
			//phoneNum := getPhoneFromConn(conn.RemoteAddr().String())
			//fmt.Println(phoneNum)
			//if phoneNum != "" {
			//	connectChange := server.ConnectChangeData {
			//		PhoneNum: phoneNum,
			//		Change: "tcp",
			//	}
			//	server.ConnectChangeChan <- connectChange
			//}
			fmt.Println("服务器的读取 err=", err)
			logs.BeeLogger.Error("服务器的读取 err= %s", err)
			return
		}
		headerBuf = append(moreBuf, headerBuf[:n]...)
		var start int
		for i, v := range headerBuf {
			if v == 0x7e {
				if i != 0 && headerBuf[i-1] != 0x7e {
					go tcpHandle(headerBuf[start:i +1], conn)
					start = i + 1
				}
			}
		}
		moreBuf = headerBuf[start:]
		//fmt.Println(hex.EncodeToString(moreBuf))
	}
}

func tcpHandle(bufData []byte, tcpConn *net.TCPConn) {
	fmt.Printf("tcp recived : %s \n", hex.EncodeToString(bufData))
	logs.BeeLogger.Info("tcp recived : %s", hex.EncodeToString(bufData))
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
		ResKeepAlive(messByte, tcpConn)
	case "0100":
		//收到注册请求并处理
		ResRegisterMes(messByte, tcpConn)
	case "0102":
		//收到鉴权请求
		ResPowerIdentify(messByte, tcpConn)
	case "0104":
		//收到终端参数
		HandleAllParameter(messByte)
	case "0107":
		HandleDeviceProperties(messByte)
	case "0200":
		//收到位置信息并回复
		ResLocateMes(messByte, tcpConn)
	case "0201":
		//收到位置信息并解析
		HandleLocationInfo(messByte)
	//case "0500":
	//	//收到车辆控制应答并解析
	//	HandleVehicleControlReceive(messByte)
	//case "0801":
	//	//收到图片信息并解析
	//	MultimediaDataUploadHandle(messByte, tcpConn)
	//case "0b01":
	//	//收到运营登记信息并回复
	//	ResOperationRegistrationMes(messByte, tcpConn)
	case "0900":
		ResDataUplinkPassThrough(messByte, tcpConn)
	case "1003":
		//收到终端上传音视频属性
		HandleStreamProperties(messByte)
	case "1205":
		//收到录像列表
		VideoListGetHandle(messByte)

	}

}
