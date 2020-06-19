package game

import (
	pubghttp "github.com/fanyanggang/BackendPlatform/fhttp"
	"github.com/wonderivan/logger"
	"log"
	"net/http"
	"pubgserver/model"
	"strconv"
)

func GetRoomBonusHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserGetStatusReq{}
	model.ParsePostRequestParams(request, param)
	err := CheckUserVali(param.UID, param.Phone)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	data, err := GetRoomBonus(param.UID, param.RoomID)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	logger.Debug("GameWinHandler succc param:%v", param)
	response.Write(pubghttp.Success(data).Return())
	return
}

func UserPayStatusHandler(response http.ResponseWriter, request *http.Request) {

	param := &model.UserGetStatusReq{}
	model.ParsePostRequestParams(request, param)
	err := CheckUserVali(param.UID, param.Phone)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	stauts, err := UserPayStatus(param.UID, param.RoomID)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	logger.Debug("GameWinHandler succc param:%v", param)
	response.Write(pubghttp.Success(stauts).Return())
	return

}

func GameWinHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserGameWin{}
	model.ParsePostRequestParams(request, param)

	if param.RoomID < 1 || len(param.List) <= 0 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	//err := CheckUserVali(param.UID, param.PhoneNumber)
	//if err != nil {
	//	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
	//	return
	//}

	_, err := PubgGameWin(param.RoomID, param.List)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	logger.Debug("GameWinHandler succc param:%v", param)
	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return
}
func GetUserFinancialHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.GetUserFinancialReq{}
	model.ParsePostRequestParams(request, param)

	if param.UID < 1 || param.Phone < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.UID, param.Phone)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	data, status, err := GetUserFinancial(param.UID)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	} else if status == model.STATUS_DB_DATA_INEXISTENCE {
		logger.Debug("GetUserFinancialHandler data not exist param:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_DATA_NOT_EXIST).Return())
		return
	}

	logger.Debug("GetUserFinancialHandler succc param:%v, status:%v", param, data)
	response.Write(pubghttp.Success(data).Return())
	return
}

func AddUserFinancialHandler(response http.ResponseWriter, request *http.Request) {

	param := &model.AddUserFinancialReq{}
	model.ParsePostRequestParams(request, param)

	if !model.VerifyCheckPhone(param.Phone) || !model.VerifyEmailFormat(param.Mail) {
		log.Printf("param err a:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	if param.UID < 1 || param.Address == "" || param.Gtype != 0 {
		log.Printf("param err b:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := AddUserFinancial(param.UID, param.Phone, param.Address, param.Mail, param.Name, param.Gtype)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	logger.Debug("AddUserFinancialHandler succc param:%v", param)
	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return
}

func UserWithDepositsStatusGetHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserTradeStatusGetReq{}
	model.ParsePostRequestParams(request, param)

	if param.Uid < 1 || param.PhoneNumber < 1 || param.Money < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.Uid, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	data, err := UserTradeStatusGet(param.Uid, param.TradeBiilID)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}

	logger.Debug("UserWithDepositsStatusGetHandler succc param:%v, status:%v", param, data)
	response.Write(pubghttp.Success(data).Return())
	return
}

func UserWithDepositsStatusUpdateHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserTradeUpdateReq{}
	model.ParsePostRequestParams(request, param)

	phone, _ := strconv.ParseInt(param.PhoneNumber, 10, 64)

	uid, _ := strconv.ParseInt(param.Uid, 10, 64)

	//(param.Status != model.STATUS_TRADE_SUCC) && (param.Status != model.STATUS_TRADE_FAIL)
	if uid < 1 || param.Money < 1 {
		log.Printf("param err:%v, money:%v", uid, param.Money)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(uid, phone)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	_, err = UserWithDepositsStatusUpdate(uid, param.Money, param.TradeBiilID, param.Status)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	logger.Debug("UserWithDepositsStatusUpdateHandler succc param:%v", param)
	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return
}

func GetUserWithDepositsBillHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.WithDepositsReq{}
	model.ParsePostRequestParams(request, param)

	//if !model.VerifyCheckPhone(param.Phone) {
	//	log.Printf("phone param err:%v", param)
	//	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
	//	return
	//}

	if !model.VerifyEmailFormat(param.Mail) {
		logger.Error("mail param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	if _, ok := model.PayType[param.PayType]; !ok {
		logger.Error("paytype param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
	}
	if param.UID < 1 || param.Address == "" || param.Money < 1 {
		logger.Error("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.UID, param.Phone)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	status, billId, err := GetUserWithDepositsBill(param.UID, param.Phone, param.Money, param.Address, param.Mail, param.Name, param.PayType, param.BankCard, param.Paytm, param.Vpa)
	if err != nil {
		logger.Debug("GetUserWithDepositsBillHandler err:%v, param:%v", err, param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	} else if status == model.STATUS_ACCOUNT_NOT_ENOUGH {
		logger.Debug("GetUserWithDepositsBillHandler account not enough param:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ACCOUNT_NOT_ENOUGH).Return())
		return
	}
	logger.Debug("GetUserWithDepositsBillHandler succc param:%v, status:%v", param, status)
	response.Write(pubghttp.Success(billId).Return())
	return
}

func UserRechargeStatusGetHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserTradeStatusGetReq{}
	model.ParsePostRequestParams(request, param)

	if param.Uid < 1 || param.PhoneNumber < 1 || param.Money < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.Uid, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	data, err := UserTradeStatusGet(param.Uid, param.TradeBiilID)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}

	logger.Debug("UserRechargeStatusGetHandler succc param:%v, status:%v", param, data)
	response.Write(pubghttp.Success(data).Return())
	return
}

func UserRechargeStatusUpdateHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.OrderCallbackData{}
	model.ParsePostRequestParams(request, param)

	uid, err := strconv.ParseInt(param.UserId, 10, 64)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	phone, err := strconv.ParseInt(param.Phone, 10, 64)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	if uid < 1 || phone < 1 || param.Amount < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err = CheckUserVali(uid, phone)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	status, err := UserRechargeStatusUpdate(uid, param.Amount, param.OrderId, param.AppOrderId, param.Status)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	if status == model.STATUS_DB_DATA_INEXISTENCE {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	logger.Debug("UserRechargeStatusUpdateHandler succc param:%v", param)
	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return
}

func GetUserRechargeBillHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.GetUserTradeBillReq{}
	model.ParsePostRequestParams(request, param)

	if param.Uid < 1 || param.PhoneNumber < 1 || param.Money < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.Uid, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	bill, err := GetUserRechargeBill(param.Uid, param.PhoneNumber, param.Money)
	if err != nil {
		logger.Error("GetUserRechargeBillHandler err param:%v, err:%v", param, err)

		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}

	logger.Debug("GetUserRechargeBillHandler succc param:%v, data:%v", param, bill)
	response.Write(pubghttp.Success(bill).Return())
	return
}

func UserStatusGetHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.StatusReq{}
	model.ParsePostRequestParams(request, param)

	err := CheckUserVali(param.Uid, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	data := map[string]int{
		"code":1,
	}

	logger.Debug("GetUserTransactionRecordHandler succc param:%v, data:%v", param, data)
	response.Write(pubghttp.Success(data).Return())
	return
}

func UserDepositGetHandler(response http.ResponseWriter, request *http.Request) {
	//logger.Debug("JoinTeamHandler succc param:%v", param)

	param := &model.UserDepositGetReq{}
	model.ParsePostRequestParams(request, param)

	if param.Uid < 1 || param.TradeBiilID == "" || param.Money < 1 || param.PhoneNumber < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.Uid, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	status, err := UserDepositGet(param.Uid, param.Money, param.TradeBiilID)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	if status == model.STATUS_DB_DATA_INEXISTENCE {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return
}

func GetUserTransactionRecordHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserTransactionRecordReq{}
	model.ParsePostRequestParams(request, param)

	if param.Uid < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.Uid, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	data, err := GetUserTransactionRecord(param.Uid)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	logger.Debug("GetUserTransactionRecordHandler succc param:%v, data:%v", param, data)
	response.Write(pubghttp.Success(data).Return())
	return
}

func GetUserCenterDataHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserCenterDataReq{}
	model.ParsePostRequestParams(request, param)

	if param.Uid < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.Uid, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	data, err := GetUserCenterData(param.Uid, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	logger.Debug("GetUserCenterDataHandler succc param:%v, data:%v", param, data)
	response.Write(pubghttp.Success(data).Return())
	return
}

func UserPayHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserPayReq{}
	model.ParsePostRequestParams(request, param)

	if param.UID < 1 || param.Money < 1 || param.TradeBiilID == "" || param.AccountType < 0 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.UID, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	status, err := UserPay(param.UID, param.RoomID, param.Money, param.TradeBiilID, param.AccountType, param.TeamMode)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}

	logger.Debug("status:%v, err:%v", status, err)
	if status == model.STATUS_ACCOUNT_NOT_ENOUGH {
		logger.Debug("UserPayHandler account not enough param:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ACCOUNT_NOT_ENOUGH).Return())
		return
	} else if status == model.STATUS_PAY_BILLID_EXIST {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	} else if status == model.STATUS_DB_DATA_FULL {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_FUll_USER).Return())
		return
	} else if status == model.STATUS_DB_HAS_PAY {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_HAS_PAY).Return())
		return
	} else if status == model.STATUS_DB_DATA_INEXISTENCE {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_DATA_INEXISTENCE).Return())
		return
	} else if status == model.STATUS_DB_NO_TEAM_NUM {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_NO_TEAM_NUM).Return())
		return
	} else if status == model.STATUS_DB_NO_NIECKNAME {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_NO_NICKNAME).Return())
		return
	} else {
		logger.Debug("UserPayHandler succc param:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
		return
	}
}

func UserRechargeHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserRechargeReq{}
	model.ParsePostRequestParams(request, param)

	if param.UID < 1 || param.Money < 1 || param.TradeBiilID == "" {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	err := CheckUserVali(param.UID, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	_, err = UserRecharge(param.UID, param.Money, param.TradeBiilID)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}

	logger.Debug("JoinTeamHandler succc param:%v", param)
	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return

}

func GetTeamHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.GetTeamReq{}
	model.ParsePostRequestParams(request, param)

	if param.UID < 1 || param.RoomID < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	err := CheckUserVali(param.UID, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	data, err := GetTeamDate(param.UID, param.RoomID)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	logger.Debug("GetTeamHandler succc param:%v, data:%v", param, data)
	response.Write(pubghttp.Success(data).Return())
	return
}

func JoinTeamHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.AddTeamNameReq{}
	model.ParsePostRequestParams(request, param)

	if param.UID < 1 || param.RoomID < 1 || param.TeamName == "" {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	if param.TeamMode == model.ROOM_PRIZE_SOLO {
		logger.Debug("AddPubgNameHandle succc param:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
		return
	}

	err := CheckUserVali(param.UID, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	err = JoinTeam(param.UID, param.RoomID, param.TeamName)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}

	logger.Debug("JoinTeamHandler succc param:%v", param)
	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return
}

func AddPubgNameHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.AddPubgNameReq{}
	model.ParsePostRequestParams(request, param)

	if param.UID < 1 || param.RoomID < 1 || param.NickName == "" || param.PhoneNumber < 0 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := CheckUserVali(param.UID, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	status, err := AddPubgName(param.UID, param.RoomID, param.PhoneNumber, param.NickName)

	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}
	if status == model.STATUS_DB_HAS_ADDPUBG {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_HAS_ADDPUBG).Return())
		return
	}

	logger.Debug("AddPubgNameHandle succc param:%v", param)
	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return
}

func GetRoomPrizeHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.RoomPrizeReq{}
	model.ParsePostRequestParams(request, param)

	if param.PhoneNumber < 1 || param.UID < 1 || (param.SoloType < model.ROOM_PRIZE_SOLO || param.SoloType > model.ROOM_PRIZE_SOLO_SQUARD) {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	err := CheckUserVali(param.UID, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}
	data, err := GetRoomPrizeInfo(param.UID, param.PhoneNumber, param.RoomID, param.SoloType)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SERVER_ERROR).Return())
		return
	}

	response.Write(pubghttp.Success(data).Return())
	logger.Debug("UserLogin succ param:%v, soleResult:%v", param, data)
	return
}

func GetRoomConfigHandler(response http.ResponseWriter, request *http.Request) {

	param := &model.RoomReq{}
	model.ParsePostRequestParams(request, param)

	if param.PhoneNumber < 1 || param.Uid < 1 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	err := CheckUserVali(param.Uid, param.PhoneNumber)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_ILLE_USER).Return())
		return
	}

	data, err := GetRoomInfo(param.PhoneNumber, param.Uid)
	if err != nil {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	logger.Debug("UserLogin succ phone:%v, record:%v", param.PhoneNumber, data)
	response.Write(pubghttp.Success(data).Return())

}

func MsgCodeVerifyHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.MsgVerfiy{}
	model.ParsePostRequestParams(request, param)
	if param.PhoneNumber < 0 || param.VerifyCode < 0 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	err := innerVerfiyCode(param.PhoneNumber, param.VerifyCode)
	if err != nil {
		logger.Debug("MsgCodeVerify innerVerfiyCode err:%v, phone:%v, code:%v", err, param.PhoneNumber, param.VerifyCode)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	logger.Debug("MsgCodeVerify succ:%v, code:%v", param.PhoneNumber, param.VerifyCode)

	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return
}

func GetVerifyCodeHandler(response http.ResponseWriter, request *http.Request) {
	param := &model.UserMsg{}
	model.ParsePostRequestParams(request, param)

	if param.PhoneNumber < 0 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	status := GetMsgCode(param.PhoneNumber)
	if status != model.STATUS_CODE_SUCC {
		log.Printf("GetMsgCode err:%v", status)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	logger.Debug("GetVerifyCode succ:%v", param.PhoneNumber)
	response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_SUCCESS).Return())
	return
}

func UserLoginHandler(response http.ResponseWriter, request *http.Request) {

	param := &model.UserLoginReq{}
	model.ParsePostRequestParams(request, param)

	logger.Debug("UserLogin:%v", param)
	if param.Phone == 0 {
		log.Printf("param err:%v", param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}
	record, status, err := GetUserInfo(param)

	if err != nil {
		logger.Debug("UserLogin CreateUserInfo err:%v, param:%v", err, param)
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_NOT_IMPLEMENTED).Return())
		return
	}

	if status == model.STATUS_DB_DATA_INEXISTENCE {
		response.Write(pubghttp.DMError(pubghttp.ERROR_CODE_CLIENT_ERROR).Return())
		return
	}

	logger.Debug("UserLogin succ param:%v", param)
	response.Write(pubghttp.Success(record).Return())

	log.Printf("UserLogin succ:%v", param)
	return
}
