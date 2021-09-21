package sqlDB

import (
	"fmt"
	"gitee.com/ictt/JTTM/tools"
	"gitee.com/ictt/JTTM/tools/logs"
	"github.com/astaxie/beego"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	gormDB *gorm.DB
	dbType string
)

func InitDB() {
	var err error
	var dataSource string
	//获取数据库类型
	dbType = beego.AppConfig.String("dbType")
	if dbType == "" || dbType != "mysql" && dbType != "sqlite3" {
		dbType = "mysql"
	}
	//获取配置文件中的数据库名
	dbName := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbName"))
	switch dbType {
		case "mysql":
			dbUser := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbUser"))
			//fmt.Println(dbUser)
			dbPwd := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbPwd"))
			dbAddr := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbAddr"))
			dbCharset := beego.AppConfig.String(tools.StringsJoin(dbType, "::dbCharset"))
			dataSource = tools.StringsJoin(dbUser, ":", dbPwd, "@tcp(", dbAddr, ")/", "?charset=", dbCharset, "&parseTime=True&loc=Local")
			//dataSource = tools.StringsJoin(dbUser, ":", dbPwd, "@tcp(", dbAddr, ")/", "?charset=", dbCharset)

			gormDB, err = gorm.Open(mysql.Open(dataSource), &gorm.Config{})

			if err != nil {
				logs.PanicLogger.Panicln(fmt.Sprintf("failed to connect %s database: %s", dbType, err))
			}
			//判断MySQL是否创建该数据库，若无则创建
			if err = gormDB.Exec(fmt.Sprintf("CREATE DATABASE if not exists `%s`", dbName)).Error; err != nil {
				logs.PanicLogger.Panicln("failed to create database error: ", err)
			}

			dataSource = tools.StringsJoin(dbUser, ":", dbPwd, "@tcp(", dbAddr, ")/", dbName, "?charset=", dbCharset, "&parseTime=True&loc=Local")
			gormDB, err = gorm.Open(mysql.Open(dataSource), &gorm.Config{
				//禁止AutoMigrate自动创建数据库外键约束
				DisableForeignKeyConstraintWhenMigrating: true,
			})
			if err != nil {
				logs.PanicLogger.Panicln(fmt.Sprintf("failed to connect %s database: %s", dbType, err))
			}

		case "sqlite3":
			//作为服务启动时需切换到当前工作目录
			//作为服务启动时需切换到当前工作目录
			os.Chdir(tools.GetAbsPath())
			gormDB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{
				//禁止AutoMigrate自动创建数据库外键约束
				DisableForeignKeyConstraintWhenMigrating: true,
			})
			if err != nil {
				logs.PanicLogger.Panicln(fmt.Sprintf("failed to connect %s database: %s", dbType, err))
			}
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		logs.PanicLogger.Panicln(fmt.Sprintf("init gorm.DB error: %s", err))
	}


	//SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	//SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	//SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	//全局禁用表名复数
	//如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	//gormDB.SingularTable(true)
	//启用Logger，显示详细日志
	//gormDB.LogMode(true)

	logs.BeeLogger.Info(fmt.Sprintf("successful connection to %s database", dbType))
	fmt.Printf("%s successful connection to %s database\n", time.Now().Format("2006-01-02 15:04:05"), dbType)

	//初始化数据表格
	InitGetVehicleInfo()
	InitGetChannelInfo()
	InitGpsAndPhoneInfo()
	InitGetUserInfo()
}
