package server

//实时直播消息头
type RealStreamHeader struct {
	FrameHeader           string
	Vpxcc                 VpxccMean
	Mpt                   MptMean
	PackageNum            uint16
	SIM                   string
	Channel               byte
	DataAndPackage        DataAndPackageMean
	TimeStamp             uint64
	Last_I_Frame_Interval uint16
	Last_Frame_Interval   uint16
	DataLength            uint16
}

//vpxcc
type VpxccMean struct {
	V  int64
	P  int64
	X  int64
	CC int64
}

//mpt
type MptMean struct {
	M  int64
	PT int64
}

//分包与数据类型
type DataAndPackageMean struct {
	Data    string
	Package string
}

//实时直播tcp流
type RealStream struct {
	Header  RealStreamHeader
	Content string
}

//收到的通用消息体
type NormalDataReceive struct {
	SequenceSend uint16
	IDSend       string
	Result       byte
}

//收到的通用消息头
type HeaderReceive struct {
	IDReceive       string
	PhoneNum        string
	SequenceReceive uint16
	PackageCount    uint16
	PackageNum      uint16
	MessProperty    MessPropertyReceive
}

//收到的通用回复
type NormalMesReceive struct {
	Header HeaderReceive
	Data   NormalDataReceive
}

//收到的鉴权消息
type PowerIdentifyMesReceive struct {
	Header HeaderReceive
	Data   string
}

//type RequestRealStreamSend struct {
//	Header 			HeaderSend
//	Data			string
//}

//发送的通用消息头
type HeaderSend struct {
	IDReceive    string
	DataLength   uint16
	PhoneNum     string
	SequenceSend uint16
}

//收到的消息体属性
type MessPropertyReceive struct {
	Subpackage byte
	Encryption []byte
	DataLength int64
}

//发送的通用消息
type NormalMesSend struct {
	Header HeaderSend
	Data   NormalDataSend
}

//发送的通用消息体
type NormalDataSend struct {
	SequenceReceive uint16
	IDReceive       string
	Result          byte
}

//收到的位置消息
type LocateMesReceive struct {
	Header HeaderReceive
	Data   LocateDataReceive
}

//收到的位置消息体
type LocateDataReceive struct {
	SequenceSend  uint16
	InfoType      string
	AlarmState    string
	AlarmSign     []byte
	State         []byte
	Latitude      string
	Longitude     string
	Altitude      float64
	Speed         float64
	Direction     float64
	Time          string
	AddLocateData AddLocateMess
}

//位置消息体中的附加消息
type AddLocateMess struct {
	Mileage              float64
	Oil                  float64
	SpeedRecode          float64
	AlarmMesID           int64
	TirePressure         string
	Temperature          int64
	SpeedAlarmAdd        SpeedAlarmAddMes
	InAndOutAlarmAdd     InAndOutAlarmAddMes
	TimeAlarmAdd         TimeAlarmAddMes
	VehicleState         []byte
	IOState              []byte
	AD                   []string
	SignalIntensity      byte
	SatelliteNum         byte
	VideoAlarm           string
	VideoSignalLossAlarm string
	LineNumber           string
	BusinessType         byte
	AbnormalDrivingAlarm AbnormalDrivingAlarmData
	Custom               []CustomData
}

//异常驾驶行为报警详细描述
type AbnormalDrivingAlarmData struct {
	AbnormalDrivingType string
	FatigueDegree       byte
}

//用户自定义消息
type CustomData struct {
	ID   byte
	Data string
}

//超速警报信息
type SpeedAlarmAddMes struct {
	LocateType byte
	LocateID   int64
}

//出入区域警报信息
type InAndOutAlarmAddMes struct {
	LocateType byte
	LocateID   int64
	Directory  byte
}

//时间警报信息
type TimeAlarmAddMes struct {
	LocateTime int64
	LocateID   int64
	Result     byte
}

//收到的所有终端参数消息
type AllParameterMes struct {
	Header HeaderReceive
	Data   AllParameterDataReceive
}

//收到的所有终端参数消息体
type AllParameterDataReceive struct {
	SequenceSend                   uint16
	ParametersNum                  byte
	KeepaliveInterval              uint32
	TCPOverTime                    uint32
	TCPRepeatNum                   uint32
	UDPOverTime                    uint32
	UDPRepeatNum                   uint32
	SMSOverTime                    uint32
	SMSRepeatNum                   uint32
	MainServer                     MainServerData
	BackupServer                   BackupServerData
	TCPPort                        uint32
	UDPPort                        uint32
	ICMainServerIP                 string
	ICTCPPort                      uint32
	ICUDPPort                      uint32
	ICBackupServerIP               string
	LocationReportStrategy         uint32
	LocationReportScheme           uint32
	ReportTimeInterval             ReportTimeIntervalData
	AccompanyServer                AccompanyServerData
	ReportDistanceInterval         ReportDistanceIntervalData
	InflectionPointAngle           uint32
	IllegalDisplacement            uint16
	IllegalDrivingPeriod           IllegalDrivingPeriodData
	DevicePhoneListenStrategy      uint32
	OnesTalkTime                   uint32
	MonthTalkTime                  uint32
	PhoneNum                       AllPhoneNumData
	AlarmShielding                 []byte
	AlarmSMSTxt                    []byte
	AlarmShooting                  []byte
	AlarmPhotoSave                 []byte
	KeyAlarm                       []byte
	HighestSpeed                   uint32
	OverSpeedTime                  uint32
	KeepDriverTime                 uint32
	OneDayDriverTime               uint32
	LeastRestTime                  uint32
	LongestStopTime                uint32
	OverSpeedWarningDifference     uint16
	FatigueDriverWarningDifference uint16
	CollisionAlarmParameters       CollisionAlarmParametersData
	RolloverAlarmAngle             uint16
	TimingCameraControl            TimingCameraControlData
	DistanceCameraControl          DistanceCameraControlData
	ImageOrVideoInstruction        uint32
	Brightness                     uint32
	ContrastRatio                  uint32
	Saturation                     uint32
	Chroma                         uint32
	Odometer                       uint32
	ProvinceID                     uint16
	CityID                         uint16
	VehicleID                      string
	VehiclePlateColor              byte
	GNSSPositioningMode            []byte
	GNSSBps                        byte
	GNSSOutputFrequency            byte
	GNSSAcquisitionFrequency       uint32
	GNSSUploadMethod               byte
	GNSSUploadSettings             uint32
	CANAcquisitionInterval1        uint32
	CANUploadInterval1             uint16
	CANAcquisitionInterval2        uint32
	CANUploadInterval2             uint16
	CANIDCollectionSettings        CANIDCollectionSettingsData
	OtherCANIDCollectionSettings   []CANIDCollectionSettingsData
	AudioAndVideoParameters        AudioAndVideoParametersData
	AudioAndVideoChannelList       AudioAndVideoChannelListData
	SingleChannelVideoParameter    SingleChannelVideoParameterData
	SpecialAlarmRecording          SpecialAlarmRecordingData
	VideoAlarmScreenWord           []byte
	VideoAnalysisAlarm             VideoAnalysisAlarmData
	DeviceWakeupMode               DeviceWakeupModeData
	Custom                         []CustomData1
}

//用户自定义参数
type CustomData1 struct {
	ID   string
	Data string
}

//主服务器参数
type MainServerData struct {
	MainServerAPNOrPPP string
	MainServerUsername string
	MainServerPassword string
	MainServerIP       string
}

//备用服务器参数
type BackupServerData struct {
	BackupServerAPNOrPPP string
	BackupServerUsername string
	BackupServerPassword string
	BackupServerIP       string
}

//从服务器参数
type AccompanyServerData struct {
	AccompanyServerAPNOrPPP string
	AccompanyServerUsername string
	AccompanyServerPassword string
	AccompanyServerAddr     string
}

//时间间隔汇报参数
type ReportTimeIntervalData struct {
	DriverReport  uint32
	SleepReport   uint32
	AlarmReport   uint32
	DefaultReport uint32
}

//距离间隔汇报参数
type ReportDistanceIntervalData struct {
	DriverReport  uint32
	SleepReport   uint32
	AlarmReport   uint32
	DefaultReport uint32
}

//所有电话参数
type AllPhoneNumData struct {
	Monitor          string
	Reset            string
	DefaultSetting   string
	MonitorSMS       string
	DeviceSMS        string
	Listen           string
	MonitorPrivilege string
}

//总线参数
type CANIDCollectionSettingsData struct {
	AcquisitionInterval uint32
	Channel             byte
	FrameType           byte
	CollectionMode      byte
	ID                  int64
}

//非法行驶路段参数
type IllegalDrivingPeriodData struct {
	Start string
	End   string
}

//冲撞警报参数
type CollisionAlarmParametersData struct {
	Collision             byte
	CollisionAcceleration byte
}

//定时拍照参数
type TimingCameraControlData struct {
	OtherSign []byte
	Interval  int64
}

//定距拍照参数
type DistanceCameraControlData struct {
	OtherSign []byte
	Interval  int64
}

//收到的透行消息
type PassThroughMes struct {
	Header HeaderReceive
	Data   PassThroughReceive
}

//收到的透行消息体
type PassThroughReceive struct {
	CNSSData    string
	ICData      string
	SerialPort1 string
	SerialPort2 string
	Custom      CustomData
}

//收到的终端属性消息
type DevicePropertiesMes struct {
	Header HeaderReceive
	Data   DevicePropertiesData
}

//收到的终端属性消息体
type DevicePropertiesData struct {
	DeviceType           []byte
	ManufacturerID       string //制造商ID
	DeviceModel          string //终端型号
	DeviceID             string //终端ID
	DeviceSIM            string
	DeviceHardwareLength byte
	DeviceHardware       string
	DeviceFirmwareLength byte
	DeviceFirmware       string
	GNSSType             []byte
	CommunicationType    []byte
}

//音视频参数设置
type AudioAndVideoParametersData struct {
	RealTimeStreamCoding           byte //0：CBR(固定码率)； 1：VBR(可变码率)； 2：ABR(平均码率)；
	RealTimeStreamResolutionRatio  byte //0：QCIF； 1：CIF； 2：WCIF； 3：D1； 4：WD1； 5：720P； 6：1 080P； 100 ～127：自定义
	RealTimeStreamKeyFrameInterval uint16
	RealTimeStreamTargetFrame      byte
	RealTimeStreamTargetKbps       uint32
	StorageStreamCoding            byte
	StorageStreamResolutionRatio   byte
	StorageStreamKeyFrameInterval  uint16
	StorageStreamTargetFrame       byte
	StorageStreamTargetKbps        uint32
	OSDSubtitleOverlaySettings     []byte
	UseAudioOrNot                  byte
}

//音视频通道列表
type AudioAndVideoChannelListData struct {
	AudioAndVideoChannelCount byte
	AudioChannelCount         byte
	VideoChannelCount         byte
	ChannelComparisonTable    []ChannelComparisonTableData
}

//音视频通道对照表
type ChannelComparisonTableData struct {
	PhysicalChannelID byte
	LogicalChannelID  byte
	ChannelType       byte
	ConnectedToPTZ    byte
}

//单独通道视频参数
type SingleChannelVideoParameterData struct {
	ChannelCount          byte
	SingleVideoParameters []SingleVideoParametersData
}

//单独通道视频参数设置
type SingleVideoParametersData struct {
	ChannelID                      byte
	RealTimeStreamCoding           byte //0：CBR(固定码率)； 1：VBR(可变码率)； 2：ABR(平均码率)；
	RealTimeStreamResolutionRatio  byte //0：QCIF； 1：CIF； 2：WCIF； 3：D1； 4：WD1； 5：720P； 6：1 080P； 100 ～127：自定义
	RealTimeStreamKeyFrameInterval uint16
	RealTimeStreamTargetFrame      byte
	RealTimeStreamTargetKbps       uint32
	StorageStreamCoding            byte
	StorageStreamResolutionRatio   byte
	StorageStreamKeyFrameInterval  uint16
	StorageStreamTargetFrame       byte
	StorageStreamTargetKbps        uint32
	OSDSubtitleOverlaySettings     []byte
}

//特殊报警录像参数
type SpecialAlarmRecordingData struct {
	Threshold byte
	Duration  byte
	StartTime byte
}

//视频分析报警参数
type VideoAnalysisAlarmData struct {
	LoadNum          byte
	FatigueThreshold byte
}

//终端休眠唤醒模式设置数据
type DeviceWakeupModeData struct {
	WakeupMode      byte
	WakeupCondition byte
	TimingWakeupDay byte
	UseTimingWakeup byte
	TimeStart1      string
	TimeEnd1        string
	TimeStart2      string
	TimeEnd2        string
	TimeStart3      string
	TimeEnd3        string
	TimeStart4      string
	TimeEnd4        string
}

//设备上传音视频属性消息
type DeviceStreamPropertiesMes struct {
	Header HeaderReceive
	Data   DeviceStreamPropertiesData
}

//设备上传音视频属性消息体
type DeviceStreamPropertiesData struct {
	AudioCoding            byte
	AudioChannelCount      byte
	AudioSamplingFrequency byte
	AudioSamplingBits      byte
	AudioFrameLength       uint16
	AllowAudioOutput       byte
	VideoCoding            byte
	MaxAudioChannelCount   byte
	MaxVideoChannelCount   byte
}

//服务器收到的运营登记信息
type OperationRegistrationMes struct {
	Header HeaderReceive
	Data   OperationRegistrationData
}

//运营登记信息消息体
type OperationRegistrationData struct {
	LineNumber     uint32
	EmployeeNumber string
}

//服务器收到的多媒体上传数据消息
type MultimediaMes struct {
	Header HeaderReceive
	Data   MultimediaData
}

//多媒体上传数据消息体
type MultimediaData struct {
	ID             uint32
	MediaType      byte
	MediaCoding    byte
	IncidentCoding byte
	ChannelID      byte
	LocateData     LocateDataReceive
	DataPackage    []byte
}

//服务器收到的录像列表消息
type VideoListMes struct {
	Header HeaderReceive
	Data   VideoListData
}

//录像列表消息体
type VideoListData struct {
	SequenceSend uint16
	VideoCount   uint32
	VideoList    []VideoListInfo
}

type VideoListInfo struct {
	ChannelID   byte
	StartTime   string
	EndTime     string
	AlarmType   string
	MediaType   byte
	Bitstream   byte
	StorageType byte
	Size        float64
}

//通知是否切链接
type ConnectChangeData struct {
	PhoneNum string
	Change string
}

var ConnectChangeChan = make(chan ConnectChangeData, 500)