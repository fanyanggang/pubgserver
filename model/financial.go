package model

const (
	//TABLE_FINANCIAL_BILLID_SEQ = "pubg_financial_billid_seq_%04d%02d%02d"
	//TABLE_FINANCIAL_BILL = "pubg_financial_bill_%04d%02d%02d"
	//TABLE_RECHARGE_BILL  = "pubg_recharge_bill_%04d%02d%02d"
	TABLE_FINANCIAL_BILL       = "pubg_financial_bill_%02d"
	TABLE_FINANCIAL_BILLID_SEQ = "pubg_financial_billid_seq_%02d"

	TABLE_FINANCIAL_BILLID_SEQ_BASE = "pubg_financial_billid_seq"
	TABLE_FINANCIAL_BILL_BASE       = "pubg_financial_bill"
	TABLE_RECHARGE_BILL_BASE        = "pubg_recharge_bill"
)

const (
	STATUS_TRADE_SUCC   = 1
	STATUS_TRADE_FAIL   = 2
	STATUS_TRADE_CANCEL = 3
)
const (
	STATUS_BILL_UNPROCESSED  = 0  // 交易开始
	STATUS_BILL_SUCC         = 1  // 交易成功
	STATUS_BILL_FAIL         = -1 // 交易失败
	STATUS_BILL_ACOUNT_NOT   = -2 // 账户余额不足
	STATUS_BILL_SYSTEM_ERROR = -3 // 支付系统返回异常
)

const (
	STATUS_PAY_BILLID_EXIST     = -1 //订单已经存在
	STATUS_PAY_BILLID_NOT_EXIST = -2 //订单不存在
	STATUS_PAY_SUCC             = 0  //扣费成功
	STATUS_PAY_FAIL             = 1  //扣费成功
)

const (
	STATUS_ACCOUNT_NOT_ENOUGH = -12 //余额不足
	STATUS_ACCOUNT_FAIL       = -11 //账户操作失败
	STATUS_ACCOUNT_SUCC       = 0   //账户操作成功
)

const (
	STATUS_TYPE_REGISTER      = 0
	STATUS_TYPE_INVITE        = 1
	STATUS_TYPE_RECHARGE      = 2
	STATUS_TYPE_GAME_TICKETS  = 3
	STATUS_TYPE_GAME_WIN      = 4
	STATUS_TYPE_WITHDEPOSISTS = 5
)

const (
	USER_REGISTER = 20
	USER_INVITE   = 30
)
const (
	ACCOUNT_TYPE_DEPOSIT = 0
	ACCOUNT_TYPE_BONUS   = 1
)
const (
	DATABASE = "financial"
)

//订单类型 0注册 1邀请 2充值 3购买 4游戏赢 5提现
type FinancialBillSeqRecord struct {
	ID          int64  `gorm:"column:id" json:"id"`
	UID         int64  `gorm:"column:uid" json:"uid"`
	BillID      string `gorm:"column:bill_id" json:"bill_id"`
	TradeBillID string `gorm:"column:trade_bill_id" json:"trade_bill_id"`
	BillType    int    `gorm:"column:bill_type" json:"bill_type"`
	CreateTime  int64  `gorm:"column:create_time" json:"create_time"`
}

type FinancialBillRecord struct {
	ID          int    `gorm:"column:id" json:"id"`
	UID         int64  `gorm:"column:uid" json:"uid"`
	BillID      string `gorm:"column:bill_id" json:"bill_id"`
	RoomID      int64  `gorm:"column:room_id" json:"room_id"`
	Money       int64  `gorm:"column:money" json:"money"`
	BillType    int    `gorm:"column:bill_type" json:"bill_type"`
	AccountType int    `gorm:"column:account_type" json:"account_type"`
	Status      int    `gorm:"column:status" json:"status"`
	TeamMode    int    `gorm:"column:team_mode" json:"team_mode"`
	CreateTime  int64  `gorm:"column:create_time" json:"create_time"`
}

type UserAccountRecord struct {
	ID            int   `gorm:"column:id" json:"id"`
	UID           int64 `gorm:"column:uid" json:"uid"`
	Money         int64 `gorm:"column:money" json:"money"`
	Deposit       int64 `gorm:"column:deposit" json:"deposit"`
	Bonus         int64 `gorm:"column:bonus" json:"bonus"`
	RechargeMoney int64 `gorm:"column:recharge_money" json:"recharge_money"`
	HistoryMoney  int64 `gorm:"column:history_money" json:"history_money"`
	CreateTime    int64 `gorm:"column:create_time" json:"create_time"`
	UpdateTime    int64 `gorm:"column:update_time" json:"update_time"`
}

type RechargeBillRecord struct {
	ID         int    `gorm:"column:id" json:"id"`
	UID        int64  `gorm:"column:uid" json:"uid"`
	BillID     string `gorm:"column:bill_id" json:"bill_id"`
	From       int    `gorm:"column:from" json:"from"`
	Money      int64  `gorm:"column:money" json:"money"`
	Status     int    `gorm:"column:status" json:"status"`
	CreateTime int64  `gorm:"column:create_time" json:"create_time"`
}

type UserFinanciaRecord struct {
	ID      int    `gorm:"column:id" json:"id"`
	UID     int64  `gorm:"column:uid" json:"uid"`
	Mail    string `gorm:"column:mail" json:"mail"`
	Address string `gorm:"column:address" json:"address"`
	Name    string `gorm:"column:name" json:"name"`
	Phone   int64  `gorm:"column:phone" json:"phone"`
	Gtype   int    `gorm:"column:gtype" json:"gtype"`
}

type GetUserFinancialReq struct {
	UID   int64 `json:"uid"`
	Phone int64 `json:"phone"`
}

type AddUserFinancialReq struct {
	UID     int64  `json:"uid"`
	Phone   int64  `json:"phone"`
	Mail    string `json:"mail"`
	Address string `json:"address"`
	Name    string `json:"name"`
	Gtype   int    `json:"gtype"`
}
