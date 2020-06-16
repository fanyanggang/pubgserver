package game

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fanyanggang/BackendPlatform/frpc-go"
	"github.com/wonderivan/logger"
	"pubgserver/dao"
	"pubgserver/model"
	"strconv"
)

func GetRoomBonus(uid, roomid int64) ([]model.RoomConfigRecord, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("UserWithDepositsStatusUpdate panic.  err:%v", r)
			logger.Error(msg)
		}
	}()

	return dao.GetRoomBonusByRoomID(uid, roomid)
}
func UserWithDepositsStatusUpdate(uid, money int64, tradeBillID string, status int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("UserWithDepositsStatusUpdate panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	redis := GetRedis()
	value, _ := redis.Get(tradeBillID)
	if string(value) != "1" {
		return model.STATUS_PAY_BILLID_NOT_EXIST, nil
	}

	status, err := dao.UserWithDepositsStatusUpdateDB(uid, money, tradeBillID, status)
	redis.Del(tradeBillID)
	return status, err
}

func UserTradeStatusGet(uid int64, tradeBillID string) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("UserTradeStatusGet panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.UserTradeStatusGetDB(uid, tradeBillID)
}

func UserRechargeStatusUpdate(uid, money int64, tradeBillID, orderBillID string, status int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("UserRechargeStatusUpdate panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	redis := GetRedis()
	value, _ := redis.Get(orderBillID)
	if string(value) != "1" {
		return model.STATUS_PAY_BILLID_NOT_EXIST, nil
	}

	status, err := dao.UserRechargeStatusUpdateDB(uid, money, tradeBillID, orderBillID, status)
	redis.Del(orderBillID)
	return status, err
}

func CheckUserVali(uid, phoneNumber int64) error {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("CheckUserVali panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.CheckUserValiDB(uid, phoneNumber)
}

func GetUserInfo(data *model.UserLoginReq) (*model.PubgUserRecord, int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetUserInfo panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.GetUserInfoDB(data)
}
func UserDepositGet(uid, money int64, tradeBillID string) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("UserDepositGet panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	redis := GetRedis()
	value, _ := redis.Get(tradeBillID)
	if string(value) == "1" {
		return model.STATUS_PAY_BILLID_EXIST, nil
	}

	status, err := dao.UserDepositGetDB(uid, money, tradeBillID)
	redis.SetExSecond(tradeBillID, "1", 24*60*60)
	return status, err
}

func GetUserTransactionRecord(uid int64) ([]model.FinancialBillRecord, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetUserTransactionRecord panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.GetUserTransactionRecordDB(uid)
}

func GetUserCenterData(uid, phoneNumber int64) (model.UserCenterDataResp, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetUserCenterData panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.GetUserCenterData(uid, phoneNumber)
}

func GetRoomInfo(phone, uid int64) (model.GetRoomResp, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetRoomInfo panic.  err:%v", r)
			logger.Error(msg)
		}
	}()

	resp := model.GetRoomResp{}
	resp.Room = make([]model.RoomRecordResp, 0)

	room, err := dao.GetRoomDB(phone, uid)
	if err != nil {
		return resp, err
	}

	usePlay, err := dao.GetRoomPlayerUserStatus(uid)
	if err != nil {
		return resp, err
	}

	//usePlay := make([]int64, 0)
	for _, n := range room {
		temp := model.RoomRecordResp{}
		temp.ID = n.ID
		temp.UserCount = n.UserCount
		temp.EntryFee = n.EntryFee
		temp.GameRooMID = n.GameRooMID
		temp.GameRoomPWD = n.GameRoomPWD
		temp.RoomName = n.RoomName
		temp.EntryFee = n.EntryFee
		temp.Map = n.Map
		temp.PerKillPrize = n.PerKillPrize
		temp.RankPlayer = n.RankPlayer
		temp.ResultImage = n.ResultImage
		temp.RreateTime = n.RreateTime
		temp.SignTime = n.SignTime
		temp.StartTime = n.StartTime
		temp.SignEndTime = n.SignEndTime
		temp.UploadTime = n.UploadTime
		temp.YouTubeURL = n.YouTubeURL
		temp.ResultStatus = n.ResultStatus
		temp.TotalPrize = n.TotalPrize
		temp.Status = n.Status
		temp.TeamMode = n.TeamMode
		temp.PrizeTitle = n.PrizeTitle
		temp.Uid = 0

		for _, v := range usePlay {
			if v.RoomID == n.ID {
				temp.Uid = v.UID
				temp.PayStatus = v.Status
			}
		}
		resp.Room = append(resp.Room, temp)
	}

	//count, err := dao.GetRoomUserNumByRoomID()
	//resp.Count = count
	return resp, nil
}

func GetRoomPrizeInfo(uid, phone, roomid int64, solo int) (model.SoleResultResp, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetRoomPrizeInfo panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.GetRoomPrizeDB(uid, phone, roomid, solo)
}

func AddPubgName(uid, roomID, phoneNumber int64, nickname string) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("AddPubgName panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.AddPubgNameDB(uid, roomID, phoneNumber, nickname)
}

func JoinTeam(uid, roomID int64, teamname string) error {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("JoinTeam panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.JoinTeamDB(uid, roomID, teamname)
}

func GetTeamDate(uid int64, roomID int) (model.TeamResp, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetTeamDate panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.GetTeamDB(uid, roomID)
}

func UserRecharge(uid, money int64, tradeBillID string) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("UserRecharge panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.UserRechargeDB(uid, money, tradeBillID)
}

func UserPayStatus(roomID, uid int64) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("UserPayStatus panic.  err:%v", r)
			logger.Error(msg)
		}
	}()

	status, err := dao.UserPayStatusDB(roomID, uid)
	logger.Debug("UserPayStatus uid:%v, roomid:%v, status:%v, err:%v", uid, roomID, status, err)
	return status, err
}

func PubgGameWin(roomID int64, list []model.WinList) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("PubgGameWin panic.  err:%v", r)
			logger.Error(msg)
		}
	}()

	redis := GetRedis()

	for _, v := range list {

		value, _ := redis.Get(v.TradeBiilID)
		if string(value) == "1" {
			return model.STATUS_PAY_BILLID_EXIST, nil
		}

		status, err := dao.PubgGameWinDB(v.UID, roomID, v.Money, v.TradeBiilID)
		if err != nil {
			return status, err
		}

		redis.SetExSecond(v.TradeBiilID, "1", 24*60*60)
	}

	return model.STATUS_DB_SUCC, nil
}

func UserPay(uid, roomID, money int64, tradeBillID string, accountType, teamMode int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("UserPay panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	redis := GetRedis()
	value, _ := redis.Get(tradeBillID)
	if string(value) == "1" {
		return model.STATUS_PAY_BILLID_EXIST, nil
	}

	logger.Debug("UserPay get uid :%v, value:%v, billID:%v", uid, string(value), tradeBillID)
	status, err := dao.UserPayDB(uid, roomID, money, tradeBillID, accountType, teamMode)
	if err != nil {
		logger.Error("status:%v, err:%v", status, err)
		return status, err
	}

	redis.SetExSecond(tradeBillID, "1", 24*60*60)
	return status, err
}

func GetUserWithDepositsBill(uid, phone, money int64, address, mail, name, payType, bankCard, paytm, vpa string) (int, string, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetUserWithDepositsBill panic.  err:%v", r)
			logger.Error(msg)
		}
	}()

	billID := dao.GetBillId(uid)

	status, err := dao.CreateUserWithDepositsBillRecord(uid, money, billID, "")
	if err != nil || status == model.STATUS_ACCOUNT_NOT_ENOUGH {
		return status, "", err
	}

	trade_billID, err := innerWithDeposit(uid, phone, money, billID, address, mail, name, payType, bankCard, paytm, vpa)
	if err != nil {
		logger.Error("GetUserWithDepositsBill inner err:%v, uid:%v, trade_billid:%v", err, uid, trade_billID)
		status, err := dao.UserWithDepositsErrRollBack(uid, money, billID)
		return status, trade_billID, fmt.Errorf("GetUserWithDepositsBill inner err:%v, uid:%v, trade_billid:%v", err, uid, trade_billID)
	}

	dao.UpdateBillStatus(uid, billID, trade_billID)
	return status, trade_billID, err
}

func innerWithDeposit(uid, phone, money int64, billID, address, mail, name, payType, bankCard, paytm, vpa string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetUserRechargeBill panic.  err:%v", r)
			logger.Error(msg)
		}
	}()

	url := "http://msg.internal.rummybank.com/v1/create_out_order"
	//url := "http://192.168.1.157:6063/v1/create_in_order"
	strphone := strconv.FormatInt(phone, 10)
	struid := strconv.FormatInt(uid, 10)

	data := map[string]interface{}{
		"app_id":       "pubg",
		"app_order_id": billID,
		"pay_type":     payType,
		"user_id":      struid,
		"phone":        strphone,
		"user_name":    name,
		"email":        mail,
		"amount":       money,
		"bank_card":    bankCard,
		"vpa":          vpa,
		"paytm":        paytm,
		"address":      address,
	}

	content, _ := json.Marshal(data)
	body, err := frpc.HttpPost(context.Background(), url, string(content))
	if err != nil {
		logger.Error("innerWithDeposit http err:v%, data:%v", err, data)
		return "", err
	}

	resq := new(model.WithDepositResp)

	json.Unmarshal([]byte(body), resq)
	if resq.Code != "200" {
		logger.Error("get recharge billid err:%v, resq:%v body:%v", body, resq, string(content))
		return "", fmt.Errorf("get recharge billid err:%v, resq:%v data:%v", body, resq, body)
	}

	if resq.Code != "200" {
		logger.Error("innerWithDeposit err:%v, data:%v", resq, data)
		return "", fmt.Errorf("innerWithDeposit err resq:%v data:%v", resq, body)
	}

	redis := GetRedis()
	redis.SetExSecond(resq.Data.OrderID, "1", 24*60*60)

	logger.Info("GetUserRechargeBill succ, uid:%v, phone:%v, resq:%v, body:%v", uid, phone, resq, string(body))
	return resq.Data.OrderID, nil
}

func GetUserRechargeBill(uid, phone, money int64) (map[string]string, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetUserRechargeBill panic.  err:%v", r)
			logger.Error(msg)
		}
	}()

	billData := make(map[string]string)
	url := "http://payment.internal.rummybank.com/v1/create_in_order"
	//url := "http://192.168.1.157:6063/v1/create_in_order"
	billID := dao.GetBillId(uid)
	strphone := strconv.FormatInt(phone, 10)
	struid := strconv.FormatInt(uid, 10)

	data := map[string]interface{}{
		"app_id":       "pubg",
		"user_id":      struid,
		"app_order_id": billID,
		"pay_channel":  "cashfree",
		"amount":       money,
		"phone":        strphone,
	}

	content, _ := json.Marshal(data)
	body, err := frpc.HttpPost(context.Background(), url, string(content))
	if err != nil {
		fmt.Errorf("UploadVivo http err:%, data:%v", err, data)
		return billData, err
	}

	resq := new(model.GetRechargeBillID)

	json.Unmarshal([]byte(body), resq)
	if resq.Code != "200" {
		return billData, fmt.Errorf("get recharge billid err:%v, resq:%v data:%v", body, resq, body)
	}

	redis := GetRedis()
	redis.SetExSecond(billID, "1", 24*60*60)

	billData["bill_id"] = resq.Data.OrderID
	billData["token"] = resq.Data.PaymentOrderID
	logger.Info("GetUserRechargeBill succ, uid:%v, phone:%v, resq:%v, body:%v", uid, phone, resq, string(body))
	return billData, nil
}
