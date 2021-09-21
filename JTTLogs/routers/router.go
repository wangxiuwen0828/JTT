package routers

import (
	"gitee.com/ictt/JTTM/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	beego.Router("/index/hook/on_stream_none_reader", &controllers.MainController{}, "post:HooStop")
	beego.Router("/api/3rd/gas", &controllers.MainController{}, "post:Gas")
	beego.Router("/api/3rd/ehome/mpdcdata", &controllers.MainController{}, "post:MPDCData")
	//beego.Router("/api/3rd/ehome/faces", &controllers.MainController{}, "post:Faces")
	//beego.Router("/api/3rd/ehome/alarm/110", &controllers.MainController{}, "post:PassengerFlowTime")
	//beego.Router("/api/3rd/ehome/alarm/111", &controllers.MainController{}, "post:PassengerFlowFrame")
	//ns := beego.NewNamespace("/api/v1",
	//	//页面跳转
	//	beego.NSNamespace("/home",
	//		//添加设备界面
	//		beego.NSRouter("/addDevice",&controllers.MainController{},"get:AddDevice"),
	//		//检测设备信息
	//		//beego.NSRouter("/checkDevice",&controllers.MainController{},"post:CheckDevice"),
	//		//获取设备数据显示到页面中
	//		beego.NSRouter("/getChannelData",&controllers.MainController{},"get:GetChannelData"),
	//		beego.NSRouter("/getDeviceData",&controllers.MainController{},"get:GetDeviceData"),
	//		//查询人数
	//		//beego.NSRouter("/queryPeople",&controllers.MainController{},"post:QueryPeople"),
	//		//删除设备
	//		//beego.NSRouter("/delDevice",&controllers.MainController{},"post:DelDevice"),
	//		//添加设备信息
	//		//beego.NSRouter("/addDeviceInfo",&controllers.MainController{},"post:AddDeviceInfo"),
	//		//修改通道信息
	//		//beego.NSRouter("/updateChannelInfo",&controllers.MainController{},"post:UpdateChannelInfo"),
	//	),
	//)
	//beego.AddNamespace(ns)
	//关闭实时直播或者录像回放
	//nsHook := beego.NewNamespace("/index",
	//	beego.NSNamespace("/hook",
	//		beego.NSRouter("/on_stream_none_reader", &controllers.HookController{}, "post:HookStop"),
	//		//beego.NSRouter("/on_publish", &controllers.HookController{}, "post:HookPublish"),
	//		//beego.NSRouter("/on_record_mp4", &controllers.HookController{}, "post:HookRecordMP4"),
	//		//beego.NSRouter("/on_http_access", &controllers.HookController{}, "post:HookHttpAccess"),
	//	),
	//)
	//beego.AddNamespace(nsHook)
}
