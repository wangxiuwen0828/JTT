package models

import (
	"gitee.com/ictt/JTTM/config"
	"gitee.com/ictt/JTTM/tools/sqlDB"
)

func GetChannelData(start, limit int64) (replyData map[string]interface{}) {
	channelList, retBool := sqlDB.QueryFindChannel()
	if !retBool {
		replyData = map[string]interface{}{
			"ErrorCode": config.JTT_ERROR_DB_ERR,
			"ErrorMsg":  config.ErrCodeMap[config.JTT_ERROR_DB_ERR],
		}
		return
	}
	replyData = map[string]interface{}{
		"ErrorCode":   config.JTT_ERROR_SUCCESS_OK,
		"ErrorMsg":    config.ErrCodeMap[config.JTT_ERROR_SUCCESS_OK],
		"ChannelList": channelList,
	}
	return
}

func GetDeviceData() (replyData map[string]interface{}) {
	deviceList, retBool := sqlDB.QueryFindVehicle()
	if !retBool {
		replyData = map[string]interface{}{
			"ErrorCode": config.JTT_ERROR_DB_ERR,
			"ErrorMsg":  config.ErrCodeMap[config.JTT_ERROR_DB_ERR],
		}
		return
	}
	replyData = map[string]interface{}{
		"ErrorCode":   config.JTT_ERROR_SUCCESS_OK,
		"ErrorMsg":    config.ErrCodeMap[config.JTT_ERROR_SUCCESS_OK],
		"ChannelList": deviceList,
	}
	return
}
