package model

const (
	STATUS_DB_SUCC             = 0  //DB操作成功
	STATUS_DB_ERROR            = -1 //DB错误
	STATUS_DB_NEW_ACCOUNT      = 1  //新建用户
	STATUS_DB_NORMAL_ACCOUNT   = 2  //正常登录用户
	STATUS_DB_HAS_ADDPUBG      = 3  //已经报名
	STATUS_DB_HAS_PAY          = 4  //已经支付
	STATUS_DB_NO_TEAM_NUM      = 5  //没有添加队名
	STATUS_DB_NO_NIECKNAME     = 6  //没有添加用户昵称
	STATUS_DB_DATA_INEXISTENCE = -2 //数据不存在
	STATUS_DB_DATA_FULL        = -3 //用户报名满员
)

const (
	STATUS_CODE_SUCC     = 0
	STATUS_MSGCODE_ERROR = -1
)
const (
	ROOM_PRIZE_SOLO        = 1
	ROOM_PRIZE_SOLO_TOP    = 2
	ROOM_PRIZE_SOLO_SQUARD = 3
)

const (
	AT_FULL_STRENGTH = 100
	PHONE_MIN        = 9100000000
	PHONE_MAX        = 9199999999
)

//测试环境开关

const CONTROL = false
