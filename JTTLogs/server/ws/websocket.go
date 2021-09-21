package ws

import (
	"errors"
	"fmt"
	"gitee.com/ictt/JTTM/config"
	"gitee.com/ictt/JTTM/tools/logs"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	wsConnAll *Clients //ws的所有连接
	maxConnId uint64
	upGrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//允许所有的CORS跨域请求，正式环境可以关闭
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

//客户端读写消息
type wsMessage struct {
	//websocket.TextMessage 消息类型
	messageType int
	data        []byte
}

//客户端连接
type wsConnection struct {
	wsSocket *websocket.Conn //底层websocket
	inChan   chan *wsMessage //读队列
	outChan  chan *wsMessage //写队列

	mutex     sync.Mutex //避免重复关闭管道,加锁处理
	isClosed  bool
	closeChan chan byte //关闭通知
	sign      string
}

func wsHandler(resp http.ResponseWriter, req *http.Request) {
	//应答客户端告知升级连接为websocket
	wsSocket, err := upGrader.Upgrade(resp, req, nil)
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("upgrade websocket error: %s", err))
	}

	wsConn := &wsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *wsMessage, 1000),
		outChan:   make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
		sign:      "",
	}

	logs.BeeLogger.Info("%s connected websocket successfully", wsConn.wsSocket.RemoteAddr())
	//处理器，可在此处进行心跳保活操作
	go wsConn.processLoop()
	//读协程
	go wsConn.wsReadLoop()
	//写协程
	go wsConn.wsWriteLoop()
}

//处理队列中的消息
func (wsConn *wsConnection) processLoop() {
	//处理消息队列中的消息
	//获取到消息队列中的消息，处理完成后，发送消息给客户端
	for {
		msg, err := wsConn.wsRead()
		if err != nil {
			logs.BeeLogger.Error("remoteAddr=%v, read message error: %s", wsConn.wsSocket.RemoteAddr(), err)
			wsConn.close()
			break
		}
		//处理接收到的数据，将要发送到客户端的数据写入消息队列中
		wsReadHandle(msg, wsConn)
	}
}

//处理消息队列中的消息
func (wsConn *wsConnection) wsReadLoop() {
	for {
		wsConn.wsSocket.SetReadDeadline(time.Now().Add(time.Duration(config.WSKeepaliveTime) * time.Second))
		//读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			//判断是不是超时
			if netErr, ok := err.(net.Error); ok {
				if netErr.Timeout() {
					logs.BeeLogger.Error("ReadMessage timeout remote: %v", wsConn.wsSocket.RemoteAddr())
					wsConn.close()
					return
				}
			}
			//其他错误，如果是 1001 和 1000 就不打印日志
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logs.BeeLogger.Error("ReadMessage other remote:%v error: %v", wsConn.wsSocket.RemoteAddr(), err)
			}

			//切断服务
			wsConn.close()
			return
		}
		req := &wsMessage{
			messageType: msgType,
			data:        data,
		}

		//放入请求队列,消息入栈
		select {
		case wsConn.inChan <- req:
		case <-wsConn.closeChan:
			wsConn.close()
			return
		}
	}
}

//发送消息给客户端
func (wsConn *wsConnection) wsWriteLoop() {
	for {
		select {
		//取一个应答
		case msg := <-wsConn.outChan:
			//写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				logs.BeeLogger.Error("remoteAddr=%v, websocket.WriteMessage() error: %s", wsConn.wsSocket.RemoteAddr(), err)
				//切断服务
				wsConn.close()
				return
			}
		case <-wsConn.closeChan:
			//获取到关闭通知
			return
		}
	}
}

//写入消息到队列中
func (wsConn *wsConnection) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("websocket connection closed")
	}
	return nil
}

//读取消息队列中的消息
func (wsConn *wsConnection) wsRead() (*wsMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		// 获取到消息队列中的消息
		return msg, nil
	case <-wsConn.closeChan:

	}
	return nil, errors.New("websocket connection closed")
}

//关闭连接
func (wsConn *wsConnection) close() {
	logs.BeeLogger.Info("close websocket connection!!!")
	//线程安全，可多次调用
	wsConn.wsSocket.Close()
	//利用标记，让closeChan只关闭一次
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		// 删除这个连接的变量
		wsConnAll.delete(wsConn.sign)
		close(wsConn.closeChan)
	}
}

//启动程序
func ListenWebsocketServer() {
	//初始化
	wsConnAll = newClient()

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "start websocketServer successful!")
	logs.BeeLogger.Info("start websocketServer successful!")
	http.HandleFunc("/ws", wsHandler)
	err := http.ListenAndServe(config.WebsocketAddr, nil)
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "start websocketServer failed, trigger panic!")
		logs.PanicLogger.Panicln(fmt.Sprintf("start websocketServer error: %s", err))
	}
}
