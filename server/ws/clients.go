package ws

import "sync"

type Clients struct {
	ClientsMap map[string]*Client
	sync.RWMutex
}

type Client struct {
	sign     string
	username string
	wsConn   *wsConnection
}

//初始化结构体
func newClient() *Clients {
	return &Clients{
		ClientsMap: make(map[string]*Client),
	}
}

//添加
func (c *Clients) set(sign string, client *Client) {
	c.Lock()
	defer c.Unlock()
	c.ClientsMap[sign] = client
}

//获取
func (c *Clients) get(sign string) *Client {
	c.Lock()
	defer c.Unlock()
	return c.ClientsMap[sign]
}

//获取当前所有连接客户端
func (c *Clients) getAll() map[string]*Client {
	c.RLock()
	defer c.RUnlock()

	return c.ClientsMap
}

//删除
func (c *Clients) delete(sign string) {
	c.Lock()
	defer c.Unlock()
	delete(c.ClientsMap, sign)
}
