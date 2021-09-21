package controllers

import (
	"gitee.com/ictt/JTTM/models"
)

//type MainController struct {
//	beego.Controller
//}

//初始页面
func (self *MainController) DeviceInfo() {
	self.TplName = "deviceInfo.html"
}

//添加设备页面
func (self *MainController) AddDevice() {
	self.TplName = "addDevice.html"
}

////添加设备信息
//func (self *MainController) AddDeviceInfo() {
//
//	self.Data["json"] = models.AddDeviceInfo(self.Ctx.Input.RequestBody)
//	self.ServeJSON()
//
//}

//同步通道信息信息到主页
func (self *MainController) GetChannelData() {
	start, _ := self.GetInt64("start")
	limit, _ := self.GetInt64("limit")

	self.Data["json"] = models.GetChannelData(start, limit)
	self.ServeJSON()
}

//同步设备信息信息到主页
func (self *MainController) GetDeviceData() {
	self.Data["json"] = models.GetDeviceData()
	self.ServeJSON()
}

////删除设备信息
//func (self *MainController) DelDevice()  {
//	self.Data["json"] = models.DelDevice(self.Ctx.Input.RequestBody)
//	self.ServeJSON()
//}

////查询总人数
//func (self *MainController) QueryPeople()  {
//	self.Data["json"] = models.QueryPeople(self.Ctx.Input.RequestBody)
//	self.ServeJSON()
//}

////检查设备信息是否正确
//func (self *MainController) CheckDevice()  {
//	self.Data["json"] = models.CheckDevice(self.Ctx.Input.RequestBody)
//	self.ServeJSON()
//}
//func (c *MainController) Get() {
//	c.Data["Website"] = "beego.me"
//	c.Data["Email"] = "astaxie@gmail.com"
//	c.TplName = "index.tpl"
//}
