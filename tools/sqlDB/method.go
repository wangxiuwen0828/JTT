package sqlDB

import (
	"fmt"
	"gitee.com/ictt/JTTM/config"
	"gitee.com/ictt/JTTM/tools/logs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

//初始化创建设备配置信息表-getVehicleInfo
func InitGetVehicleInfo() {
	if !gormDB.Migrator().HasTable(&GetVehicleInfo{}) {
		//不存在，创建表格
		CreateTable(&GetVehicleInfo{})
	} else {
		//表格存在，则更新Status=OFF
		getVehicleInfo := map[string]interface{}{
			"Status": "OFF",
		}
		Updates(GetTableName(GetVehicleInfo{}), getVehicleInfo, "1 = 1")
	}
}

//初始化设备通道列表-GetChannelInfo
func InitGetChannelInfo() {
	if !gormDB.Migrator().HasTable(&GetChannelInfo{}) {
		//不存在，创建表格
		CreateTable(&GetChannelInfo{})
	}
	//else {
		//表格存在，则更新Status=OFF
	//	getTemporaryInfo := map[string]interface{}{
	//		"Status": "OFF",
	//	}
	//	Updates(GetChannelInfo{}, getTemporaryInfo)
	//}
}

//初始化用户信息表
func InitGetUserInfo() {
	if !gormDB.Migrator().HasTable(&GetUserInfo{}) {
		//不存在，创建表格
		CreateTable(&GetUserInfo{})
		CreateDefaultUsers(config.AdminUser)
		//创建一个默认游客账号
		CreateDefaultUsers(config.GuestUser)
	} else {

		CheckDefaultUser(config.AdminUser)
		//检查游客账号是否存在
		CheckDefaultUser(config.GuestUser)
	}
}

//初始化创建车辆位置信息表-GetLocationInfo
func InitGetLocationInfo(tableName string) {
	if !gormDB.Migrator().HasTable(tableName) {
		//不存在，创建表格
		CreateMyTable(tableName, &GetLocationInfo{})
	}
}

//初始化创建报警信息表-GetAlarmInfo
func InitGetAlarmInfo(tableName string) {
	if !gormDB.Migrator().HasTable(tableName) {
		//不存在，创建表格
		CreateMyTable(tableName, &GetAlarmInfo{})
	}
}

//初始化创建驾驶行为信息表-GetAlarmInfo
func InitGetDriverInfo(tableName string) {
	if !gormDB.Migrator().HasTable(tableName) {
		//不存在，创建表格
		CreateMyTable(tableName, &DrivingAction{})
	}
}

//初始化gps与phoneNum对照表
func InitGpsAndPhoneInfo() {
	if !gormDB.Migrator().HasTable(&GpsAndPhoneNum{}) {
		//不存在，创建表格
		CreateTable(&GpsAndPhoneNum{})
	}
}

//创建表格，如果表格不存在
func CreateTable(tab interface{}) {
	tableName := GetTableName(tab)
	switch dbType {
	case "mysql":
		if err := gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Migrator().CreateTable(tab); err != nil {
			logs.PanicLogger.Panicln(fmt.Sprintf("create %s table failed: %s", tableName, err))
		}
	case "sqlite3":
		if err := gormDB.Migrator().CreateTable(tab); err != nil {
			logs.PanicLogger.Panicln(fmt.Sprintf("create %s table failed: %s", tableName, err))
		}
	}
	logs.BeeLogger.Info("create %s table success", tableName)
}

//创建表格，如果表格不存在
func CreateMyTable(tableName string, tab interface{}) {
	switch dbType {
	case "mysql":
		if err := gormDB.Table(tableName).Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Migrator().CreateTable(tableName, tab); err != nil {
			logs.PanicLogger.Panicln(fmt.Sprintf("create %s table failed: %s", tableName, err))
		}
	case "sqlite3":
		if err := gormDB.Table(tableName).Migrator().CreateTable(tab); err != nil {
			logs.PanicLogger.Panicln(fmt.Sprintf("create %s table failed: %s", tableName, err))
		}
	}
	logs.BeeLogger.Info("create %s table success", tableName)
}

//创建默认用户账号
func CreateDefaultUsers(userName string) {
	user := new(GetUserInfo)
	switch userName {
	case config.AdminUser:
		//管理员admin
		user = &GetUserInfo{
			UserName:        config.AdminUser,
			Password:        config.AdminPwd,
			PermissionLevel: "admin",
			Remarks:         "",
		}

	case config.GuestUser:
		//游客guest
		user = &GetUserInfo{
			UserName:        config.GuestUser,
			Password:        config.GuestPwd,
			PermissionLevel: "guest",
			Remarks:         "",
		}
	}
	if err := gormDB.Create(user).Error; err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("create %s's user error in users'table: %s", userName, err))
	} else {
		logs.BeeLogger.Info("create %s's user successful in users'table", userName)
	}
}

//检查账号是否存在
func CheckDefaultUser(userName string) {
	var userPwd string
	switch userName {
	case config.AdminUser:
		userPwd = config.AdminPwd
	case config.GuestUser:
		userPwd = config.GuestPwd
	}

	user := new(GetUserInfo)
	if err := gormDB.Where(map[string]interface{}{"Username": userName, "Password": userPwd}).First(user).Error; err != nil {
		if err.Error() != "record not found" {
			logs.PanicLogger.Panicln(fmt.Sprintf("query default %s user error: %s", userName, err))
		}
		//账号不存在，创建账号
		CreateDefaultUsers(userName)
	} else {
		//账号存在，更新时间
		gormDB.Model(user).Update("UpdatedAt", time.Now())
	}
}

//解析获得用户要操作的表名
func GetTableName(value interface{}) (tableName string) {
	if name, ok := value.(string); ok {
		tableName = name
	} else {
		//fmt.Println("类型打印：", reflect.TypeOf(gormDB.Migrator()))
		switch mg := gormDB.Migrator().(type) {
		case sqlite.Migrator:
			mg.RunWithValue(value, func(statement *gorm.Statement) error {
				tableName = statement.Table
				return nil
			})
		case mysql.Migrator:
			mg.RunWithValue(value, func(statement *gorm.Statement) error {
				tableName = statement.Table
				return nil
			})
		}
	}

	return
}

//查询某个表的总数
func TotalCount(tab interface{}) (totalCount int64, retBool bool) {
	tableName := GetTableName(tab)
	if err := gormDB.Table(tableName).Count(&totalCount).Error; err != nil {
		logs.BeeLogger.Error("query total count from %s's table error: %s", tableName, err)
		return
	}

	return totalCount, true
}

//添加记录到表格中，存在则更新，不存在则插入
func Save(tab interface{}) bool {
	tableName := GetTableName(tab)
	if err := gormDB.Save(tab).Error; err != nil {
		logs.BeeLogger.Error("%s's table save record error: %s", tableName, err)
		return false
	}
	return true
}

//分页查询
func Limit(tab interface{}, start, limit int, where ...interface{}) {
	gormDB.Limit(limit).Offset(start).Find(tab, where...)
}

//执行更新操作，更新更改字段
//func Updates(tbl interface{}, data interface{}) {
//	tableName := GetTableName(tbl)
//	if err := gormDB.Model(tbl).Updates(data).Error; err != nil {
//		logs.BeeLogger.Error("%s table run db.Updates() failed: %s", tableName, err)
//		return
//	}
//}

func Updates(tableName string, mapData map[string]interface{}, query string, args ...interface{}, ) {
	if err := gormDB.Table(tableName).Where(query, args...).Updates(mapData).Error; err != nil {
		logs.BeeLogger.Error("batch update in the %s table error: %s", tableName, err)
	}
}

//根据通道名更新数据库中表格的字段
func UpdateTableFromChannel(tbl interface{}, locate, change string, updateAt string, key string) {
	sql := fmt.Sprintf(`UPDATE %s SET %s = '%s', UpdatedAt= '%v' WHERE PhoneNumAndChannel = '%s'`, GetTableName(tbl), locate, change, updateAt, key)
	if err := gormDB.Exec(sql).Error; err != nil {
		logs.BeeLogger.Error("update table error: %s", err)
	}
}

//插入gps信息
func InsertGPSInfo(tableName string, gpsList GPSList) {
	sql := fmt.Sprintf(`INSERT INTO %s (PhoneNum, InfoType, AlarmState, Latitude, Longitude, Altitude, Speed, Direction, Time, Mileage, Oil, SpeedRecode) VALUES `, tableName)
	for i := 0; i < 10; i++ {
		sql += fmt.Sprintf(`('%s', '%s', '%s', '%s', '%s', %v, %v, %v, '%s', %v, %v, %v),`, gpsList.GpsInfo[i].PhoneNum, gpsList.GpsInfo[i].InfoType,
			gpsList.GpsInfo[i].AlarmState, gpsList.GpsInfo[i].Latitude, gpsList.GpsInfo[i].Longitude, gpsList.GpsInfo[i].Altitude, gpsList.GpsInfo[i].Speed,
			gpsList.GpsInfo[i].Direction, gpsList.GpsInfo[i].Time, gpsList.GpsInfo[i].Mileage, gpsList.GpsInfo[i].Oil, gpsList.GpsInfo[i].SpeedRecode)
	}
	//fmt.Println(sql)
	sql = sql[:len(sql)-1]
	if err := gormDB.Exec(sql).Error; err != nil {
		logs.BeeLogger.Error("INSERT GPS table %s error: %s", tableName, err)
	}
}

//插入报警信息
func InsertAlarmInfo(tableName string, alarmList AlarmList) {
	sql := fmt.Sprintf(`INSERT INTO %s (PhoneNum, VehicleAlarm, StreamAlarm, SignLostChannel, Time) VALUES `, tableName)
	for i := 0; i < 10; i++ {
		sql += fmt.Sprintf(`('%s', '%s', '%s', '%s', '%s'),`, alarmList.AlarmInfo[i].PhoneNum, alarmList.AlarmInfo[i].VehicleAlarm,
			alarmList.AlarmInfo[i].StreamAlarm, alarmList.AlarmInfo[i].SignLostChannel, alarmList.AlarmInfo[i].Time)
	}
	//fmt.Println(sql)
	sql = sql[:len(sql)-1]
	if err := gormDB.Exec(sql).Error; err != nil {
		logs.BeeLogger.Error("INSERT Alarm table %s error: %s", tableName, err)
	}
}

//插入报警信息
func InsertDriverActionInfo(tableName string, driverList DrivingAction) {
	sql := fmt.Sprintf(`INSERT INTO %s (PhoneNum, AbnormalDrivingType, FatigueDegree, Time) VALUES `, tableName)
	//for i := 0; i < 10; i++ {
		sql += fmt.Sprintf(`('%s', '%s', %v, '%s')`, driverList.PhoneNum, driverList.AbnormalDrivingType,
			driverList.FatigueDegree, driverList.Time)
	//}
	//fmt.Println(sql)
	//sql = sql[:len(sql)-1]
	if err := gormDB.Exec(sql).Error; err != nil {
		logs.BeeLogger.Error("INSERT DriverAction table %s error: %s", tableName, err)
	}
}

//根据设备名更新数据库中表格的状态
func UpdateTableFromPhone(tbl interface{}, status string, updateAt string, key string) {
	sql := fmt.Sprintf(`UPDATE %s SET Status = '%s', UpdatedAt= '%v' WHERE PhoneNum = '%s'`, GetTableName(tbl), status, updateAt, key)
	if err := gormDB.Exec(sql).Error; err != nil {
		logs.BeeLogger.Error("update table error: %s", err)
	}
}

//查询符合条件的某一条用户账号信息，参数一为表，参数二为查询条件
func QueryUserTake(tab interface{}, query interface{}) bool {
	if err := gormDB.Where(query).Take(tab).Error; err != nil {
		if err.Error() != "record not found" {
			logs.BeeLogger.Error("query users error: %s", err)
		} else {
			logs.BeeLogger.Info("this user not found!!!")
		}
		return false
	}

	return true
}

//查询符合条件的第一条记录，有数据返回1，无数据返回0，其他查询出错返回-1
func First(tbl interface{}, data interface{}) int {
	tableName := GetTableName(tbl)
	if err := gormDB.Where(data).First(tbl).Error; err != nil {
		logs.BeeLogger.Error("%s table run db.First() failed: %s", tableName, err)
		//gormDB.Where(data).First(tbl).RecordNotFound()
		if err.Error() == "record not found" {
			//查询失败，无符合条件的数据
			return 0
		}
		return -1
	}

	return 1
}

////查询所有记录
//func Find(tbl interface{}, list interface{}) {
//	tableName := GetTableName(tbl)
//	if err := gormDB.Find(list).Error; err != nil {
//		logs.BeeLogger.Error("%s table run db.Find() failed: %s", tableName, err)
//	}
//}

//查询所有记录，参数一为数组名，参数二为table结构体
func Find(array, tab interface{}) {
	tableName := GetTableName(tab)
	if err := gormDB.Find(array).Error; err != nil {
		logs.BeeLogger.Error("%s table run db.Find() failed: %s", tableName, err)
	}
}

//创建记录
func Create(tbl interface{}) bool {
	tableName := GetTableName(tbl)
	if err := gormDB.Create(tbl).Error; err != nil {
		logs.BeeLogger.Error("error inserting records into the %s table: %s", tableName, err)
	} else {
		return true
	}
	return false
}

//获取所有通道信息
func QueryFindChannel() (ChannelList []GetChannelInfo, retBool bool) {
	tableName := GetTableName(&GetChannelInfo{})
	if err := gormDB.Find(&ChannelList).Error; err != nil {
		logs.BeeLogger.Error("queryFindAllChannels record from %s's table error: %s", tableName, err)
		return
	}

	return ChannelList, true
}

//获取所有车辆信息
func QueryFindVehicle() (vehicleList []GetVehicleInfo, retBool bool) {
	tableName := GetTableName(&GetVehicleInfo{})
	if err := gormDB.Find(&vehicleList).Error; err != nil {
		logs.BeeLogger.Error("queryFindAll vehicleList record from %s's table error: %s", tableName, err)
		return vehicleList, false
	}

	return vehicleList, true
}

//获取某表所有gps信息
func QueryFindGPS(tableName string) (gpsList []GetLocationInfo, retBool bool) {
	//tableName := GetTableName(&GetLocationInfo{})
	if err := gormDB.Table(tableName).Find(&gpsList).Error; err != nil {
		logs.BeeLogger.Error("queryFindAllChannels record from %s's table error: %s", tableName, err)
		return
	}

	return gpsList, true
}

//获取某表所有报警信息
func QueryFindAlarm(tableName string) (alarmList []GetAlarmInfo, retBool bool) {
	//tableName := GetTableName(&GetLocationInfo{})
	if err := gormDB.Table(tableName).Find(&alarmList).Error; err != nil {
		logs.BeeLogger.Error("queryFindAllChannels record from %s's table error: %s", tableName, err)
		return
	}

	return alarmList, true
}

//获取某表所有驾驶行为信息
func QueryFindDriverAction(tableName string) (driverList []DrivingAction, retBool bool) {
	//tableName := GetTableName(&GetLocationInfo{})
	if err := gormDB.Table(tableName).Find(&driverList).Error; err != nil {
		logs.BeeLogger.Error("queryFindAllChannels record from %s's table error: %s", tableName, err)
		return
	}

	return driverList, true
}
