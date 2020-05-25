package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var PayType = map[string]bool{
	"paytm": true,
	"bank":  true,
	"upi":   true,
}

type UserLoginReq struct {
	Phone          int64  `json:"phone"`
	Platform       string `json:"platform"`
	Tuid           string `json:"tuid"`
	Photo          string `json:"photo"`
	Nickname       string `json:"nickname"`
	FatherID       int64  `json:"father_id"`
	FatherNickName string `json:"father_nickname"`
}

type PubgUserRecord struct {
	ID             int64  `gorm:"column:id" json:"id"`
	Phone          int64  `gorm:"column:phone" json:"phone"`
	Platform       string `gorm:"column:platform" json:"platform"`
	Tuid           string `gorm:"column:tuid" json:"tuid"`
	Photo          string `gorm:"column:photo" json:"photo"`
	Nickname       string `gorm:"column:nickname" json:"nickname"`
	FatherID       int64  `gorm:"column:father_id" json:"father_id"`
	FatherNickName string `gorm:"column:father_nickname" json:"father_nickname"`
	CreateTime     int64  `gorm:"column:create_time" json:"create_time"`
}

type UserMsg struct {
	PhoneNumber int64 `json:"phone"`
}

type MsgVerfiy struct {
	PhoneNumber int64 `json:"phone"`
	VerifyCode  int   `json:"verifyCode"`
}

type UserDepositGetReq struct {
	Uid         int64  `json:"uid"`
	Money       int64  `json:"money"`
	PhoneNumber int64  `json:"phone"`
	TradeBiilID string `json:"trade_bill_id"`
}

type UserTradeStatusGetReq struct {
	Uid         int64  `json:"uid"`
	Money       int64  `json:"money"`
	PhoneNumber int64  `json:"phone"`
	TradeBiilID string `json:"trade_bill_id"`
}

type UserTradeUpdateReq struct {
	Uid         string `json:"user_id"`
	Money       int64  `json:"amount"`
	PhoneNumber string `json:"phone"`
	Status      int    `json:"status"`
	TradeBiilID string `json:"trade_bill_id"`
}

type OrderCallbackData struct {
	PayChannel     string `json:"pay_channel"`
	AppOrderId     string `json:"app_order_id"`
	OrderId        string `json:"order_id"`
	Amount         int64  `json:"amount"`
	PaymentOrderId string `json:"payment_order_id"`
	PaymentId      string `json:"payment_id"`
	Status         int    `json:"status"`
	UserId         string `json:"user_id"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
}

//用户交易记录
type UserTransactionRecordReq struct {
	Uid         int64 `json:"uid"`
	PhoneNumber int64 `json:"phone"`
}

//用户中心数据
type UserCenterDataReq struct {
	Uid         int64 `json:"uid"`
	PhoneNumber int64 `json:"phone"`
}

//用户中心数据
type UserCenterDataResp struct {
	Uid      int64  `json:"uid"`
	Photo    string `json:"photo"`
	Nickname string `json:"nickname"`
	Money    int64  `json:"money"`
	Deposit  int64  `json:"deposit"`
	Bonus    int64  `json:"bonus"`
}

//获取房间信息
type RoomReq struct {
	Uid         int64 `json:"uid"`
	PhoneNumber int64 `json:"phone"`
}

type AddPubgNameReq struct {
	UID         int64  `json:"uid"`
	PhoneNumber int64  `json:"phone"`
	RoomID      int64  `json:"room_id"`
	NickName    string `json:"nick_name"`
}

type AddTeamNameReq struct {
	UID         int64  `json:"uid"`
	PhoneNumber int64  `json:"phone"`
	RoomID      int64  `json:"room_id"`
	TeamName    string `json:"team_name"`
	TeamMode    int    `json:"team_mode"`
}

type GetTeamReq struct {
	UID         int64 `json:"uid"`
	PhoneNumber int64 `json:"phone"`
	RoomID      int   `json:"room_id"`
}

//充值接入第三方，考虑参数问题
type UserRechargeReq struct {
	UID         int64  `json:"uid"`
	Money       int64  `json:"money"`
	PhoneNumber int64  `json:"phone"`
	TradeBiilID string `json:"trade_bill_id"`
}

//用户交易订单
type GetUserTradeBillReq struct {
	Uid         int64 `json:"uid"`
	PhoneNumber int64 `json:"phone"`
	Money       int64 `json:"money"`
}

type WithDepositsReq struct {
	UID      int64  `json:"uid"`
	Mail     string `json:"mail"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Phone    int64  `json:"phone"`
	PayType  string `json:"pay_type"`
	BankCard string `json:"bank_card"`
	Paytm    string `json:"paytm"`
	Vpa      string `json:"vpa"`
	Money    int64  `json:"money"`
}

type WithDepositResp struct {
	Code string          `json:"code"`
	Data WithDepositData `json:"data"`
}

type WithDepositData struct {
	ID             int    `json:"id"`
	AppID          string `json:"app_id"`
	PayChannel     string `json:"pay_channel"`
	AppOrderID     string `json:"app_order_id"`
	OrderID        string `json:"order_id"`
	PayType        string `json:"pay_type"`
	Amount         int64  `json:"amount"`
	Uid            string `json:"user_id"`
	Phone          string `json:"phone"`
	Mail           string `json:"email"`
	Paytm          string `json:"paytm"`
	BankCard       string `json:"bank_card"`
	Address        string `json:"address"`
	Username       string `json:"user_name"`
	PaymentOrderID string `json:"payment_order_id"`
	PaymentID      string `json:"payment_id"`
	Status         int    `json:"status"`
	ThirdCode      string `json:"third_code"`
	ThirdDesc      string `json:"third_desc"`
	Create         int64  `json:"created"`
	Update         int64  `json:"updated"`
}

type GetRechargeBillID struct {
	Code string     `json:"code"`
	Data BillIDDate `json:"data"`
}

type BillIDDate struct {
	ID             int    `json:"id"`
	AppID          string `json:"app_id"`
	PayChannel     string `json:"pay_channel"`
	AppOrderID     string `json:"app_order_id"`
	OrderID        string `json:"order_id"`
	PaymentOrderID string `json:"payment_order_id"`
	PaymentID      string `json:"payment_id"`
	Status         string `json:"status"`
	ThirdCode      string `json:"third_code"`
	ThirdDesc      string `json:"third_desc"`
	Create         int64  `json:"created"`
	Update         int64  `json:"updated"`
}

type UserGetStatusReq struct {
	UID    int64 `json:"uid"`
	Phone  int64 `json:"phone"`
	RoomID int64 `json:"room_id"`
}

type UserGameWin struct {
	RoomID int64     `json:"room_id"`
	List   []WinList `json:"list"`
}

type WinList struct {
	UID         int64  `json:"uid"`
	Money       int64  `json:"money"`
	TradeBiilID string `json:"trade_bill_id"`
}

type UserPayReq struct {
	UID         int64  `json:"uid"`
	PhoneNumber int64  `json:"phone"`
	Money       int64  `json:"money"`
	RoomID      int64  `json:"room_id"`
	AccountType int    `json:"account_type"`
	TradeBiilID string `json:"trade_bill_id"`
	TeamMode    int    `json:"team_mode"`
}

//获取prize列表
type RoomPrizeReq struct {
	UID         int64 `json:"uid"`
	SoloType    int   `json:"soloType"`
	RoomID      int64 `json:"room_id"`
	PhoneNumber int64 `json:"phone"`
}

//pg_rooms
type RoomRecord struct {
	ID           int64  `gorm:"column:id" json:"id"`
	RoomName     string `gorm:"column:room_name" json:"room_name"`
	TeamMode     int    `gorm:"column:team_mode" json:"team_mode"`
	EntryFee     int64  `gorm:"column:entry_fee" json:"entry_fee"`
	PerKillPrize int64  `gorm:"column:per_kill_prize" json:"per_kill_prize"`
	RankPlayer   int    `gorm:"column:rank_player" json:"rank_player"`
	Map          int    `gorm:"column:map" json:"map"`
	TotalPrize   int    `gorm:"column:total_prize" json:"total_prize"`
	StartTime    int64  `gorm:"column:start_time" json:"start_time"`
	SignTime     int64  `gorm:"column:sign_time" json:"sign_time"`
	SignEndTime  int    `gorm:"column:sign_end_time" json:"sign_end_time"`
	Status       int    `gorm:"column:status" json:"status"`
	ResultStatus int    `gorm:"column:result_status" json:"result_status"`
	YouTubeURL   string `gorm:"column:youtube_url" json:"youtube_url"`
	ResultImage  string `gorm:"column:result_image" json:"result_image"`
	GameRooMID   int    `gorm:"column:game_room_id" json:"game_room_id"`
	GameRoomPWD  string `gorm:"column:game_room_pwd" json:"game_room_pwd"`
	PrizeTitle   string `gorm:"column:prize_title" json:"prize_title"`
	UserCount    int64  `gorm:"column:user_count" json:"user_count"`
	RreateTime   int64  `gorm:"column:created" json:"created"`
	UploadTime   int64  `gorm:"column:updated" json:"updated"`
}

type RoomConfigRecord struct {
	ID         int64  `gorm:"column:id" json:"id"`
	RoomID     int64  `gorm:"column:room_id" json:"room_id"`
	ConfKey    string `gorm:"column:conf_key" json:"conf_key"`
	Value      string `gorm:"column:value" json:"value"`
	CreateTime int64  `gorm:"column:created" json:"created"`
	UploadTime int64  `gorm:"column:updated" json:"updated"`
}

type RoomRecordResp struct {
	ID           int64  `json:"id"`
	RoomName     string `json:"room_name"`
	TeamMode     int    `json:"team_mode"`
	EntryFee     int64  `json:"entry_fee"`
	PerKillPrize int64  `json:"per_kill_prize"`
	RankPlayer   int    `json:"rank_player"`
	Map          int    `json:"map"`
	TotalPrize   int    `json:"total_prize"`
	StartTime    int64  `json:"start_time"`
	SignTime     int64  `json:"sign_time"`
	Uid          int64  `json:"uid"`
	SignEndTime  int    `json:"sign_end_time"`
	Status       int    `json:"status"`
	PayStatus    int    `json:"pay_status"`
	ResultStatus int    `json:"result_status"`
	YouTubeURL   string `json:"youtube_url"`
	ResultImage  string `json:"result_image"`
	GameRooMID   int    `json:"game_room_id"`
	GameRoomPWD  string `json:"game_room_pwd"`
	UserCount    int64  `json:"user_count"`
	PrizeTitle   string `json:"prize_title"`
	RreateTime   int64  `json:"created"`
	UploadTime   int64  `json:"updated"`
}

type GetRoomResp struct {
	Room []RoomRecordResp `json:"room_data"`
}

//pg_rooms_team
type RoomTeamRecord struct {
	ID       int64  `gorm:"column:id" json:"id"`
	RoomID   int64  `gorm:"column:room_id" json:"room_id"`
	TeamNum  int    `gorm:"column:team_num" json:"team_num"`
	TeamName string `gorm:"column:team_name" json:"team_name"`
	Create   int64  `gorm:"column:created" json:"created"`
}

type SoloResultRecord struct {
	ID       int64 `gorm:"column:id" json:"id"`
	RoomID   int64 `gorm:"column:room_id" json:"room_id"`
	PayerID  int64 `gorm:"column:payer_id" json:"payer_id"`
	TeamNum  int   `gorm:"column:team_num" json:"team_num"`
	TeamRank int   `gorm:"column:team_rank" json:"team_rank"`
	KillNum  int64 `gorm:"column:kill_num" json:"kill_num"`
	Create   int64 `gorm:"column:created" json:"created"`
	Update   int64 `gorm:"column:updated" json:"updated"`
}

//1、solo 单人赛  按照kills > 0 排序  pg_rooms_player
//2、solo top  单人赛 按照rank 排序 pg_rooms_solo_top
//3、solo two/four 按照rank排序 ORDER BY rank,kills pg_rooms_duo_squad

//pg_rooms_player
type RoomPlayerRecord struct {
	ID           int64  `gorm:"column:id" json:"id"`
	UID          int64  `gorm:"column:uid" json:"uid"`
	GameNiceName string `gorm:"column:game_nick_name" json:"game_nick_name"`
	RoomID       int64  `gorm:"column:room_id" json:"room_id"`
	TeamNum      int    `gorm:"column:team_num" json:"team_num"`
	KillNum      int    `gorm:"column:kills" json:"kills"`
	Status       int    `gorm:"column:status" json:"status"`
	Prize        int64  `gorm:"column:prize" json:"prize"`
	Create       int64  `gorm:"column:created" json:"created"`
	Update       int64  `gorm:"column:updated" json:"updated"`
}

//pg_rooms_solo_top
type RoomsSoloTopRecord struct {
	ID      int64 `gorm:"column:id" json:"id"`
	RoomID  int64 `gorm:"column:room_id" json:"room_id"`
	Uid     int64 `gorm:"column:uid" json:"uid"`
	Rank    int   `gorm:"column:rank" json:"rank"`
	KillNum int   `gorm:"column:kills" json:"kills"`
	Prize   int64 `gorm:"column:prize" json:"prize"`
	Create  int64 `gorm:"column:created" json:"created"`
	Update  int64 `gorm:"column:updated" json:"updated"`
}

//pg_rooms_duo_squad
type RoomsDuoSquadRecord struct {
	ID       int64  `gorm:"column:id" json:"id"`
	UID      int64  `gorm:"column:uid" json:"uid"`
	GameID   int64  `gorm:"column:game_id" json:"game_id"`
	TeamNum  int    `gorm:"column:team_num" json:"team_num"`
	TeamName string `gorm:"column:team_name" json:"team_name"`
	Rank     int    `gorm:"column:rank" json:"rank"`
	MVP      int    `gorm:"column:is_mvp" json:"is_mvp"`
	KillNum  int    `gorm:"column:kills" json:"kills"`
	Prize    int64  `gorm:"column:prize" json:"prize"`
	Create   int64  `gorm:"column:created" json:"created"`
	Update   int64  `gorm:"column:updated" json:"updated"`
}

type SoleResultResp struct {
	Result map[int][]SolePerson
}

type SolePerson struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Kills int    `json:"kills"`
	Prize int64  `json:"prize"`
	Photo string `json:"photo"`
	Phone string `json:"phone"`
}

type TeamResp struct {
	TeamUser map[int][]TeamUserData `json:"team_user"`
}

type TeamUserData struct {
	Uid      int64  `json:"uid"`
	Photo    string `json:"photo"`
	Nickname string `json:"nickname"`
}

type RoomUserNumByRoomID struct {
	RoomID int64 `gorm:"column:room_id" json:"room_id"`
	Count  int64 `gorm:"column:count" json:"count"`
}

func ParsePostRequestParams(r *http.Request, param interface{}) {
	body, _ := ioutil.ReadAll(r.Body)
	if body == nil || strings.EqualFold(string(body), "") {
		return
	}

	log.Printf("ParsePostRequestParams|PostParams|body:data:%v", string(body))
	err := json.Unmarshal(body, param)
	if err != nil {
		log.Printf("ParsePostRequestParams|PostParams:error%v", err)
		return
	}
	log.Printf("ParsePostRequestParams|PostParams|Unmarshal|data:%v", string(body))
}

func VerifyCheckPhone(phone int64) bool {
	//if phone > PHONE_MIN && phone < PHONE_MAX {
	//	return true
	//}
	//return false
	return true
}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
