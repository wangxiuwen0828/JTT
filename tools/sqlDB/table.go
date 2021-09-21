package sqlDB

import (
	"gitee.com/ictt/JTTM/config"
)

//设备配置信息表-GetServerInfo
type GetVehicleInfo struct {
	PhoneNum          string            `gorm:"column:PhoneNum;primary_key" json:"phoneNum" csv:"phoneNum"` //终端手机号
	DeviceIP          string            `gorm:"column:DeviceIP" json:"-"`                    //来自设备注册时的UDPAddr
	ProvinceID        string            `gorm:"column:ProvinceID" json:"provinceID" csv:"provinceID"`
	CountyID          string            `gorm:"column:CountyID" json:"countyId" csv:"countyId"`                   //省市县域ID， GB/T 2260
	ManufacturerID    string            `gorm:"column:ManufacturerID" json:"manufacturerID" csv:"manufacturerID"`       //制造商ID
	DeviceModel       string            `gorm:"column:DeviceModel" json:"deviceModel" csv:"deviceModel"`             //终端型号
	DeviceID          string            `gorm:"column:DeviceID" json:"deviceID" csv:"deviceID"`                   //终端ID
	VehiclePlateColor string            `gorm:"column:VehiclePlateColor" json:"vehiclePlateColor"csv:"vehiclePlateColor"` //车牌颜色，0未上牌，1蓝，2黄，3黑，4白，5绿，6黄绿，7其他
	VehicleID         string            `gorm:"column:VehicleID" json:"vehicleID" csv:"vehicleID"`                 //车牌号
	PowerIdentify     string            `gorm:"column:PowerIdentify" json:"-"`
	Status            string            `gorm:"column:Status" json:"status"` //设备在线状态，yes-在线，no-离线
	ChannelCount      byte              `gorm:"column:ChannelCount" json:"channelCount" csv:"channelCount"`
	CreatedAt         config.TimeNormal `gorm:"column:CreatedAt" json:"-"`
	UpdatedAt         string            `gorm:"column:UpdatedAt" json:"-"`
}

//设置表名
func (GetVehicleInfo) TableName() string {
	return "getVehicleInfo"
}

//设备通道信息表-GetChannelInfo
type GetChannelInfo struct {
	PhoneNumAndChannel string            `gorm:"column:PhoneNumAndChannel;primary_key" json:"phoneNumAndChannel"` //终端手机号加通道号
	PhoneNum           string            `gorm:"column:PhoneNum" json:"phoneNum"`                                 //终端手机号
	LogicalChannelID   int64             `gorm:"column:LogicalChannelID" json:"logicalChannelID"`                 //逻辑通道号
	Status             string            `gorm:"column:Status" json:"status"`                                     //通道在线状态，ON-在线，OFF-离线
	Alarm              string            `gorm:"column:Alarm" json:"alarm"`
	CreatedAt          config.TimeNormal `gorm:"column:CreatedAt" json:"-"`
	UpdatedAt          string            `gorm:"column:UpdatedAt" json:"-"`
}

//设置表名
func (GetChannelInfo) TableName() string {
	return "getChannelInfo"
}

//用户信息表
type GetUserInfo struct {
	UserName        string            `gorm:"column:UserName;primary_key" json:"userName"` //用户名
	Password        string            `gorm:"column:Password" json:"password"`
	PermissionLevel string            `gorm:"column:PermissionLevel" json:"permissionLevel"` //密码
	Remarks         string            `gorm:"column:Remarks" json:"remarks"`                 //备注
	CreatedAt       config.TimeNormal `gorm:"column:CreatedAt" json:"-"`
	UpdatedAt       config.TimeNormal `gorm:"column:UpdatedAt" json:"-"`
}

//设置表名
func (GetUserInfo) TableName() string {
	return "getUserInfo"
}

//Gps与设备对应表
type GpsAndPhoneNum struct {
	PhoneNum        string `gorm:"column:PhoneNum;primary_key" json:"phoneNum"` //终端手机号
	GpsTableName    string `gorm:"column:GpsTableName" json:"gpsTableName"`
	AlarmTableName  string `gorm:"column:AlarmTableName" json:"alarmTableName"`
	DriverTableName string `gorm:"column:DriverTableName" json:"driverTableName"`
}

func (GpsAndPhoneNum) TableName() string {
	return "gpsAndPhoneNum"
}

//车辆位置信息表-GetLocationInfo
type GetLocationInfo struct {
	PhoneNum string `gorm:"column:PhoneNum" json:"phoneNum" csv:"phoneNum"` //终端手机号
	InfoType string `gorm:"column:InfoType" json:"infoType" csv:"infoType"` //消息类型
	//AlarmSign				string				`gorm:"column:AlarmSign" json:"alarmSign"`	 //报警标志
	//VehicleState			string				`gorm:"column:VehicleState" json:"vehicleState"`	 //车辆状态
	AlarmState  string  `gorm:"column:AlarmState" json:"alarmState" csv:"alarmState"`   //报警状态
	Latitude    string  `gorm:"column:Latitude" json:"latitude" csv:"latitude"`       //纬度
	Longitude   string  `gorm:"column:Longitude" json:"longitude" csv:"longitude"`     //经度
	Altitude    float64 `gorm:"column:Altitude" json:"altitude" csv:"altitude"`       //海拔
	Speed       float64 `gorm:"column:Speed" json:"speed" csv:"speed"`             //速度
	Direction   float64 `gorm:"column:Direction" json:"direction" csv:"direction"`     //方向 0-359，正北为 0，顺时针
	Time        string  `gorm:"column:Time" json:"time" csv:"time"`               //时间
	Mileage     float64 `gorm:"column:Mileage" json:"mileage" csv:"mileage"`         //行驶里程
	Oil         float64 `gorm:"column:Oil" json:"oil" csv:"oil"`                 //油量
	SpeedRecode float64 `gorm:"column:SpeedRecode" json:"speedRecode" csv:"speedRecode"` //里程表上的车速
	//AlarmMesID				int64				`gorm:"column:AlarmMesID" json:"alarmMesID"`	 //人工确认报警ID
	//TirePressure			string
	//Temperature				int64
	//SpeedAlarmAdd			SpeedAlarmAddMes
	//InAndOutAlarmAdd		InAndOutAlarmAddMes
	//TimeAlarmAdd			TimeAlarmAddMes
	//VehicleAddState			string				`gorm:"column:VehicleAddState" json:"vehicleAddState"`	 //车辆扩展状态位
	//IOState					string				`gorm:"column:IOState" json:"iOState"`	 //IO状态位
	//AD0 					string				`gorm:"column:AD0" json:"aD0"`	 //模拟量
	//AD1						string				`gorm:"column:AD1" json:"aD1"`	 //模拟量
	//SignalIntensity			byte				`gorm:"column:SignalIntensity" json:"signalIntensity"`	 //无线通信网络信号强度
	//SatelliteNum			byte				`gorm:"column:SatelliteNum" json:"satelliteNum"`	 //GNSS 定位卫星数
	//VideoAlarm				string				`gorm:"column:VideoAlarm" json:"videoAlarm"`	 //视频相关报警
	//VideoSignalLossAlarm	string				`gorm:"column:VideoSignalLossAlarm" json:"videoSignalLossAlarm"`	 //视频丢失通道
	//LineNumber				string				`gorm:"column:LineNumber" json:"lineNumber"`	 //线路编码
	//BusinessType			byte				`gorm:"column:BusinessType" json:"businessType"`	 //业务类型
	//AbnormalDrivingType		string				`gorm:"column:AbnormalDrivingAlarm" json:"abnormalDrivingAlarm"`	 //驾驶员异常行为报警
	//FatigueDegree			byte				`gorm:"column:FatigueDegree" json:"fatigueDegree"`	 //疲劳程度
	//Status     			string            	`gorm:"column:Status" json:"status"` 					 //设备在线状态，yes-在线，no-离线
}

////设置表名
//func (GetLocationInfo) TableName(tableName string) string {
//	return tableName
//}

//设备配置信息表-GetAlarmInfo
type GetAlarmInfo struct {
	PhoneNum        string `gorm:"column:PhoneNum" json:"phoneNum" csv:"phoneNum"`               //终端手机号
	VehicleAlarm    string `gorm:"column:VehicleAlarm" json:"vehicleAlarm" csv:"vehicleAlarm"`       //车辆报警状态信息
	StreamAlarm     string `gorm:"column:StreamAlarm" json:"streamAlarm" csv:"streamAlarm"`         //视频报警标志信息
	SignLostChannel string `gorm:"column:SignLostChannel" json:"signLostChannel" csv:"signLostChannel"` //视频信号丢失的通道信息
	Time            string `gorm:"column:Time" json:"time" csv:"time"`                       //时间
}

//设置表名
//func (GetAlarmInfo) TableName() string {
//	return "getAlarmInfo"
//}

type DrivingAction struct {
	PhoneNum            string `gorm:"column:PhoneNum" json:"phoneNum" csv:"phoneNum"`                       //终端手机号
	AbnormalDrivingType string `gorm:"column:AbnormalDrivingType" json:"abnormalDrivingType" csv:"abnormalDrivingType"` //驾驶员异常行为报警
	FatigueDegree       byte   `gorm:"column:FatigueDegree" json:"fatigueDegree" csv:"fatigueDegree"`             //疲劳程度
	Time                string `gorm:"column:Time" json:"time" csv:"time" `                               //时间
}

////录像列表
//type VideoListInfo struct {
//	DeviceID		string				`gorm:"column:DeviceID;unique_index" json:"deviceid"`
//	StartTime		string				`gorm:"column:StartTime;unique_index" json:"startTime"`
//	EndTime			string				`gorm:"column:EndTime;unique_index" json:"endTime"`
//	RecordedFile	[]string			`gorm:"column:RecordedFile" json:"recordedFile"`
//}
//
//type StreamSession struct {
//	DeviceChannelID     string            `gorm:"column:DeviceChannelID unique_index:kv" json:"lowerdomain"`
//	Session     		string            `gorm:"column:Session" json:"session"`
//}
//func (StreamSession) TableName() string {
//	return "streamSession"
//}

type GPSList struct {
	Count   int8
	GpsInfo [10]GetLocationInfo
}

type AlarmList struct {
	Count     int8
	AlarmInfo [10]GetAlarmInfo
}

type DriverList struct {
	Count      int8
	DriverInfo [10]DrivingAction
}
