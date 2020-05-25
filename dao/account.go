package dao

import (
	"fmt"
	pubgsql "github.com/fanyanggang/BackendPlatform/fsql"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"log"
	"pubgserver/model"
	"time"
)

func AddUserFinancialDB(uid, phone int64, address, mail, name string, gtype int) error {
	record := model.UserFinanciaRecord{}

	r := pubgsql.Get("financial").Master().Table("pubg_user_financia").Where("uid = ? ", uid).Find(&record)
	if r.RecordNotFound() {
		logger.Debug("GetUserFinancialDB RecordNotFound uid:%v", uid)

		financiaRecord := &model.UserFinanciaRecord{
			UID:     uid,
			Mail:    mail,
			Address: address,
			Phone:   phone,
			Gtype:   gtype,
			Name:    name,
		}

		rInsert := pubgsql.Get("financial").Master().Table("pubg_user_financia").Create(financiaRecord)
		if rInsert.Error != nil {
			logger.Error("AddUserFinancialDB create err:%v, data:%v", rInsert.Error, record)
			return r.Error
		}
		logger.Debug("GetUserFinancialDB succ uid:%v", uid)
		return nil
	} else {
		logger.Error("GetUserFinancialDB err:%v, uid:%v", r.Error, uid)
		return r.Error
	}
}

func GetUserFinancialDB(uid int64) (model.UserFinanciaRecord, int, error) {
	record := model.UserFinanciaRecord{}

	r := pubgsql.Get("financial").Master().Table("pubg_user_financia").Where("uid = ? ", uid).Find(&record)
	if r.RecordNotFound() {
		logger.Debug("GetUserFinancialDB RecordNotFound uid:%v", uid)
		return record, model.STATUS_DB_DATA_INEXISTENCE, nil
	} else if r.Error != nil {
		logger.Error("GetUserFinancialDB err:%v, uid:%v", r.Error, uid)
		return record, model.STATUS_DB_ERROR, r.Error
	}

	logger.Debug("GetUserFinancialDB succ uid:%v, record:%v", uid, record)
	return record, model.STATUS_DB_SUCC, nil
}

func GetUserTransactionRecordDB(uid int64) ([]model.FinancialBillRecord, error) {
	record := make([]model.FinancialBillRecord, 0)

	tableName := GetFinancialBillTable(uid)
	r := pubgsql.Get("financial").Master().Table(tableName).Where("uid = ? and status = 1", uid).Order("create_time desc").Find(&record)
	if r.Error != nil {
		logger.Error("GetUserTransactionRecordDB err:%v, uid:%v", r.Error, uid)
		return record, r.Error
	}

	logger.Debug("GetUserTransactionRecordDB succ uid:%v, record:%v", uid, record)
	return record, nil
}

func CheckUserValiDB(uid, phoneNumber int64) error {
	accountRecord := &model.PubgUserRecord{}
	r := pubgsql.Get("pubg").Master().Table("pubg_user_info").Where("id = ? and phone = ? ", uid, phoneNumber).Find(&accountRecord)
	if r.Error != nil {
		logger.Error("CheckUserValiDB err:%v, uid:%v", r.Error, uid)
		return r.Error
	}
	return nil
}
func GetUserCenterData(uid, phoneNumber int64) (model.UserCenterDataResp, error) {

	resp := model.UserCenterDataResp{}
	accountRecord := &model.UserAccountRecord{}
	r := pubgsql.Get("financial").Master().Table("pubg_user_account").Where("uid = ? ", uid).Find(&accountRecord)
	if r.Error != nil {
		logger.Error("UserPayDB pubg_user_account get  err:%v, uid:%v", r.Error, uid)
		return resp, r.Error
	}

	resp.Money = accountRecord.Money
	resp.Deposit = accountRecord.Deposit
	resp.Bonus = accountRecord.Bonus

	record := new(model.PubgUserRecord)
	r = pubgsql.Get("pubg").Master().Table("pubg_user_info").Where(" id = ? and phone = ? ", uid, phoneNumber).Find(record)
	if r.Error != nil {
		log.Printf("GetUserData err:%v, tuid:%v", r.Error, uid)
		return resp, fmt.Errorf("GetUserData err:%v, tuid:%v", r.Error, uid)
	}

	resp.Uid = uid
	resp.Photo = record.Photo
	resp.Nickname = record.Nickname

	logger.Debug("GetUserCenterData succ:%v", resp)
	return resp, nil
}

func UserPayDB(uid, roomid, money int64, tradeBillID string, accountType, teamMode int) (int, error) {

	record, err := GetUserRoomPlayerStatus(uid, roomid)
	if err != nil {
		logger.Error("record:%v, err:%v", record, err)
		return model.STATUS_DB_DATA_INEXISTENCE, nil
	}

	if record.Status == 1 {
		return model.STATUS_DB_HAS_PAY, nil
	}

	if record.TeamNum == 0 && teamMode == model.ROOM_PRIZE_SOLO_SQUARD {
		return model.STATUS_DB_NO_TEAM_NUM, nil
	}

	if record.GameNiceName == "" {
		return model.STATUS_DB_NO_NIECKNAME, nil
	}

	count, err := GetRoomPlayUserCountByRoomid(uid, roomid)
	if err != nil {
		return model.STATUS_DB_ERROR, err
	}

	if count == model.AT_FULL_STRENGTH {
		return model.STATUS_DB_DATA_FULL, nil
	}

	status, err := innerUserPayDB(uid, roomid, money, tradeBillID, accountType, teamMode)

	if status == model.STATUS_DB_SUCC {
		return UpdateRoomPlayerUserStatus(uid, roomid)
	} else {
		return status, err
	}

}

func GetRoomPlayerUserStatus(uid int64) ([]model.RoomPlayerRecord, error) {
	record := make([]model.RoomPlayerRecord, 0)
	tx := pubgsql.Get("pubg").Master().Table("pg_rooms_player").Where("uid = ? ", uid).Find(&record)

	if tx.Error != nil {
		logger.Error("getRoomPlayerUserStatus err:%v, uid:%v", tx.Error, uid)
		return record, tx.Error
	}

	logger.Debug("getRoomPlayerUserStatus succ, uid:%v, record:%v", uid, record)
	return record, nil
}

func GetUserRoomPlayerStatus(uid, roomID int64) (model.RoomPlayerRecord, error) {
	record := model.RoomPlayerRecord{}
	r := pubgsql.Get("pubg").Master().Table("pg_rooms_player").Where(" room_id = ? and uid = ? ", roomID, uid).Find(&record)

	if r.Error != nil {
		logger.Error("GetUserRoomPlayerStatus err:%v, uid:%v, roomid:%v", r.Error, uid, roomID)
		return record, r.Error
	}
	return record, nil
}

func GetRoomPlayUserCountByRoomid(uid, roomID int64) (int64, error) {

	record := model.RoomRecord{}
	tx := pubgsql.Get("pubg").Master().Table("pg_rooms").Where("id = ? ", roomID).Find(&record)

	if tx.Error != nil {
		logger.Error("getRoomPlayerUserStatus err:%v, uid:%v", tx.Error, uid)
		return 0, tx.Error
	}

	logger.Debug("getRoomPlayerUserStatus succ, uid:%v", uid)
	return record.UserCount, nil
}

func UpdateRoomPlayerUserStatus(uid, roomID int64) (int, error) {
	r := pubgsql.Get("pubg").Master().Table("pg_rooms_player").Where(" room_id = ? and uid = ?", roomID, uid).Updates(map[string]interface{}{"status": 1, "updated": time.Now().Unix()})

	if r.Error != nil {
		logger.Error("updateRoomUserStatus err:%v, uid:%v, roomid:%v", r.Error, uid, roomID)
		return model.STATUS_DB_ERROR, r.Error
	}

	record := model.RoomRecord{}
	tx := pubgsql.Get("pubg").Master().Table("pg_rooms").Where("id = ? ", roomID).Find(&record)
	if tx.Error != nil {
		logger.Error("get pg_rooms  err:%v, uid:%v, roomid:%v", r.Error, uid, roomID)
		return model.STATUS_DB_ERROR, r.Error
	}

	tx = pubgsql.Get("pubg").Master().Table("pg_rooms").Where(" id = ?", roomID).Update("user_count", record.UserCount+1)
	if tx.Error != nil {
		logger.Error("get pg_rooms update user_count  err:%v, uid:%v, roomid:%v, record:%v", tx.Error, uid, roomID, record)
		return model.STATUS_DB_ERROR, r.Error
	}

	logger.Error("updateRoomUserStatus succ, uid:%v, roomid:%v", uid, roomID)
	return model.STATUS_DB_SUCC, nil
}

func innerUserPayDB(uid, roomid, money int64, tradeBillID string, accountType, teamMode int) (int, error) {
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	return innerUpdate(tx, uid, money, roomid, model.STATUS_TYPE_GAME_TICKETS, accountType, teamMode, tradeBillID)
}

func PubgGameWinDB(uid, roomID, money int64, tradeBillID string) (int, error) {
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	return innerUpdate(tx, uid, money, roomID, model.STATUS_TYPE_GAME_WIN, model.ACCOUNT_TYPE_DEPOSIT, 0, tradeBillID)

}

//func innerUpdate(tx *gorm.DB, uid, money int64, aType, accountType int, tradeBillID string, payMode int) {

//tx := pubgsql.Get("pubg").Master().Begin()
//defer tx.Commit()
//t := time.Now()
//seqTable := GetBillIdSeqTable(t)
//
//billID := genBillId(uid)
//record := model.FinancialBillSeqRecord{
//	UID:         uid,
//	BillID:      billID,
//	TradeBillID: tradeBillID,
//	CreateTime:  time.Now().Unix(),
//}
//
//rInsert := tx.Table(seqTable).Create(record)
//if rInsert.Error != nil {
//	logger.Error("UserPayDB seq create err:%v, data:%v", rInsert.Error, record)
//	tx.Rollback()
//	return model.STATUS_ACCOUNT_FAIL, rInsert.Error
//}
//
//billTable := GetFinancialBillTable(t)
//billRecord := model.FinancialBillRecord{
//	UID:        uid,
//	BillID:     billID,
//	RoomID:     roomid,
//	Money:      money,
//	Status:     model.STATUS_BILL_UNPROCESSED,
//	CreateTime: time.Now().Unix(),
//}
//
//rInsert = tx.Table(billTable).Create(billRecord)
//if rInsert.Error != nil {
//	tx.Rollback()
//	logger.Error("UserPayDB bill create err:%v, data:%v", rInsert.Error, record)
//	return model.STATUS_ACCOUNT_FAIL, rInsert.Error
//}
//
//accountRecord := &model.UserAccountRecord{}
//r := tx.Table("pubg_user_account").Where("uid = ? ", uid).Find(accountRecord)
//if r.Error != nil {
//	tx.Rollback()
//	logger.Error("UserPayDB pubg_user_account get  err:%v, uid:%v, tradebillid:%v", uid, tradeBillID)
//	return model.STATUS_ACCOUNT_FAIL, r.Error
//}
//
//if accountRecord.Money >= money {
//	r = tx.Table("pubg_user_account").Where("uid = ? ", uid).Update("mondy", accountRecord.Money-money)
//	if r.Error != nil {
//		tx.Rollback()
//		logger.Error("UserPayDB pubg_user_account update err:%v, uid:%v, tradebillid:%v", uid, tradeBillID)
//		return model.STATUS_ACCOUNT_FAIL, r.Error
//	}
//} else {
//	r = tx.Table("billTable").Where("bill_id = ? ", billID).Update("Status", model.STATUS_BILL_ACOUNT_NOT)
//	if r.Error != nil {
//		tx.Rollback()
//		logger.Error("UserPayDB pubg_user_account billTable  err:%v, uid:%v, billID:%v, tradebillid:%v", uid, billID, tradeBillID)
//		return model.STATUS_ACCOUNT_FAIL, r.Error
//	}
//	logger.Debug("UserPayDB succ account not enough, uid:%v, roomid:%v, money:%v, tradebillid:%v", uid, roomid, money, tradeBillID)
//	return model.STATUS_ACCOUNT_NOT_ENOUGH, nil
//}
//
//r = tx.Table("billTable").Where("bill_id = ? ", billID).Update("Status", model.STATUS_BILL_PROCESSED)
//if r.Error != nil {
//	tx.Rollback()
//	logger.Error("UserPayDB billTable update err:%v, uid:%v, tradebillid:%v", uid, tradeBillID)
//	return model.STATUS_ACCOUNT_FAIL, r.Error
//}
//logger.Debug("UserPayDB succ, uid:%v, roomid:%v, money:%v, tradebillid:%v", uid, roomid, money, tradeBillID)
//return model.STATUS_ACCOUNT_SUCC, nil

func AddPubgNameDB(uid, roomID, phoneNumber int64, nickname string) (int, error) {
	tx := pubgsql.Get("pubg").Master()

	precord := &model.RoomPlayerRecord{}
	r := tx.Table("pg_rooms_player").Where(" room_id = ? and uid = ?", roomID, uid).Find(precord)

	if r.RecordNotFound() {
		record := &model.RoomPlayerRecord{
			UID:          uid,
			GameNiceName: nickname,
			RoomID:       roomID,
			Prize:        0.00,
			TeamNum:      0,
			KillNum:      0,
			Create:       0,
			Update:       0,
			Status:       0,
		}

		rInsert := tx.Table("pg_rooms_player").Create(record)
		if rInsert.Error != nil {
			logger.Error("AddPubgNameDB no insert err:%v, data:%v", rInsert.Error, record)
			return model.STATUS_DB_ERROR, rInsert.Error
		}

		if rInsert.RowsAffected != 1 {
			logger.Error("AddPubgNameDB add no RowsAffected err:%d, data:%v", rInsert.RowsAffected, record)
			return model.STATUS_DB_ERROR, fmt.Errorf(" AddPubgNameDB no insert RowsAffected err:%d , data:%v", rInsert.RowsAffected, record)
		}

		logger.Debug("AddPubgNameDB add succ uid:%v, roomid:%v, nickname:%v", uid, roomID, nickname)
		return model.STATUS_CODE_SUCC, nil

	} else if r.RowsAffected == 1 {
		r = tx.Table("pg_rooms_player").Where(" room_id = ? and uid = ?", roomID, uid).Update("game_nick_name", nickname)
		if r.Error != nil {
			logger.Debug("AddPubgNameDB update err uid:%v, roomid:%v, nickname:%v, r.Error:%v", uid, roomID, nickname, r.Error)
			return model.STATUS_DB_ERROR, r.Error
		}
		logger.Debug("AddPubgNameDB update succ uid:%v, roomid:%v, nickname:%v", uid, roomID, nickname)
		return model.STATUS_CODE_SUCC, nil
	} else if r.Error != nil {
		logger.Debug("AddPubgNameDB err uid:%v, roomid:%v, nickname:%v, r.Error:%v", uid, roomID, nickname, r.Error)
		return model.STATUS_DB_ERROR, r.Error
	}

	logger.Debug("AddPubgNameDB succ uid:%v, roomid:%v, nickname:%v", uid, roomID, nickname)
	return model.STATUS_DB_HAS_ADDPUBG, nil
}

func GetUserInfoDB(data *model.UserLoginReq) (*model.PubgUserRecord, int, error) {

	record := new(model.PubgUserRecord)
	tx := pubgsql.Get("pubg").Master().Begin()
	defer tx.Commit()
	r := tx.Table("pubg_user_info").Where(" phone = ? ", data.Phone).Find(record)

	var photo string
	if data.Photo == "" {
		photo = "https://jeeto-file.s3.ap-south-1.amazonaws.com/images/5.jpg"
	} else {
		photo = data.Photo
	}

	if r.RowsAffected != 1 {
		record := &model.PubgUserRecord{
			Phone:          data.Phone,
			Platform:       data.Platform,
			Tuid:           data.Tuid,
			Photo:          photo,
			Nickname:       data.Nickname,
			FatherID:       data.FatherID,
			FatherNickName: data.FatherNickName,
			CreateTime:     time.Now().Unix(),
		}

		db := tx.Table("pubg_user_info").Create(record)

		if db.Error != nil {
			log.Printf("CreateUserInfo error:%v, data:%v", db.Error, data)
			return record, model.STATUS_DB_ERROR, db.Error
		}

		//收徒
		if data.FatherID > 0 {
			frecord := new(model.PubgUserRecord)
			r = tx.Table("pubg_user_info").Where(" id = ?  ", data.FatherID).Find(frecord)
			if r.Error != nil {
				tx.Rollback()
				return record, model.STATUS_DB_DATA_INEXISTENCE, nil
			}
			UserRegister(record.ID, data.FatherID)

			logger.Debug("User Register add father succ  :%v", data)
			return record, model.STATUS_DB_NEW_ACCOUNT, nil
		} else {
			UserRegister(record.ID, 0)
			logger.Debug("User Register succ :%v", data)
			return record, model.STATUS_DB_NEW_ACCOUNT, nil
		}
	} else if r.Error != nil {
		logger.Debug("CreateUserInfo err:%v, data:%v", r.Error, data)
		return record, model.STATUS_DB_DATA_INEXISTENCE, nil
	}

	logger.Debug("GetUserInfo :%v", data)
	return record, model.STATUS_DB_NORMAL_ACCOUNT, nil
}

func UserRegister(uid, fatherID int64) {
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	//用户注册
	innerUpdate(tx, uid, model.USER_REGISTER, 0, model.STATUS_TYPE_REGISTER, model.ACCOUNT_TYPE_BONUS, 0, "")
	if fatherID > 0 {
		//填写师父邀请码
		innerUpdate(tx, fatherID, model.USER_INVITE, 0, model.STATUS_TYPE_INVITE, model.ACCOUNT_TYPE_BONUS, 0, "")
	}
}

func UpdataAccountCreateBIll(uid, money int64, aType, accountType int) {
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	innerUpdate(tx, uid, money, 0, aType, accountType, 0, "")
}

func innerUpdate(tx *gorm.DB, uid, money, roomid int64, billType, accountType, teamMode int, tradeBillID string) (int, error) {
	logger.Debug("innerUpdate beigin")
	billid := GetBillId(uid)
	record := model.FinancialBillSeqRecord{
		UID:         uid,
		BillID:      billid,
		TradeBillID: tradeBillID,
		BillType:    billType,
		CreateTime:  time.Now().Unix(),
	}

	seqTableName := GetBillIdSeqTable(uid)
	r := tx.Table(seqTableName).Create(&record)
	if r.Error != nil {
		tx.Rollback()
		logger.Error("UserRegister create billseq err:%v, uid:%v", r.Error, uid)
		return model.STATUS_DB_ERROR, r.Error
	}

	billTable := GetFinancialBillTable(uid)
	billRecord := &model.FinancialBillRecord{
		UID:         uid,
		BillID:      billid,
		RoomID:      roomid,
		Money:       money,
		BillType:    billType,
		AccountType: accountType,
		TeamMode:    teamMode,
		Status:      model.STATUS_BILL_UNPROCESSED,
		CreateTime:  time.Now().Unix(),
	}

	rInsert := tx.Table(billTable).Create(billRecord)
	if rInsert.Error != nil {
		tx.Rollback()
		logger.Error("UserPayDB bill create err:%v, data:%v", rInsert.Error, record)
		return model.STATUS_DB_ERROR, r.Error
	}

	status, err := innerUpdateAccount(tx, uid, money, billType, accountType)
	if err != nil {
		if status == model.STATUS_DB_ERROR {
			tx.Rollback()
			return model.STATUS_DB_ERROR, nil
		} else if status == model.STATUS_ACCOUNT_NOT_ENOUGH {
			r = tx.Table(billTable).Where("bill_id = ? ", billid).Update("Status", model.STATUS_BILL_ACOUNT_NOT)
			if r.Error != nil {
				tx.Rollback()
				logger.Error("innerUpdate update bill err:%v, uid:%v, money:%v, aType:%v, accountType:%v", rInsert.Error, uid, money, billType, accountType)
				return model.STATUS_DB_ERROR, r.Error
			}
		} else {
			r = tx.Table(billTable).Where("bill_id = ? ", billid).Update("Status", model.STATUS_BILL_FAIL)
			if r.Error != nil {
				tx.Rollback()
				logger.Error("innerUpdate update bill err:%v, uid:%v, money:%v, aType:%v, accountType:%v", rInsert.Error, uid, money, billType, accountType)
				return model.STATUS_DB_ERROR, r.Error
			}
		}
		return model.STATUS_ACCOUNT_NOT_ENOUGH, nil
	} else {
		r = tx.Table(billTable).Where("bill_id = ? ", billid).Update("Status", model.STATUS_BILL_SUCC)
		if r.Error != nil {
			tx.Rollback()
			logger.Error("innerUpdate update bill err:%v, uid:%v, money:%v, aType:%v, accountType:%v", rInsert.Error, uid, money, billType, accountType)
			return model.STATUS_DB_ERROR, r.Error
		}
		return model.STATUS_DB_SUCC, nil
	}
}

//订单类型 0注册 1邀请 2充值 3购买 4游戏赢 5提现
func innerUpdateAccount(tx *gorm.DB, uid, money int64, billType, accountType int) (int, error) {
	//更新用户账户
	account := &model.UserAccountRecord{}
	r := tx.Table("pubg_user_account").Where("uid = ? ", uid).Find(account)

	//用户注册
	if billType == model.STATUS_TYPE_REGISTER {
		if r.RecordNotFound() {
			account = &model.UserAccountRecord{
				UID:           uid,
				Money:         money,
				Deposit:       0,
				Bonus:         money,
				RechargeMoney: 0,
				HistoryMoney:  money,
				CreateTime:    time.Now().Unix(),
				UpdateTime:    time.Now().Unix(),
			}

			rInsert := tx.Table("pubg_user_account").Create(account)
			if rInsert.Error != nil {
				logger.Error("innerUpdateAccount create account err:%v, uid:%v, money:%v, aType:%v", rInsert.Error, uid, money, billType)
				return model.STATUS_DB_ERROR, rInsert.Error
			}
			logger.Error("innerUpdateAccount create account succ uid:%v, money:%v, aType:%v", uid, money, billType)
			return model.STATUS_DB_SUCC, nil
		} else {
			logger.Error("innerUpdateAccount find account err:%v, uid:%v, money:%v, aType:%v", r.Error, uid, money, billType)
			return model.STATUS_DB_ERROR, fmt.Errorf("innerUpdateAccount find account err:%v, uid:%v, money:%v, aType:%v", uid, money, billType)
		}
		//收徒
	} else if billType == model.STATUS_TYPE_INVITE {
		if r.Error != nil {
			logger.Error("innerUpdateAccount INVITE err:%v, uid:%v, money:%v, aType:%v", r.Error, uid, money, billType)
			return model.STATUS_DB_ERROR, fmt.Errorf("innerUpdateAccount INVITE err:%v, uid:%v, money:%v, aType:%v", uid, money, billType)
		}
		rUpdate := tx.Table("pubg_user_account").Where("uid = ? ", uid).Updates(map[string]interface{}{"bonus": account.Bonus + money, "history_money": account.HistoryMoney + money, "money": account.Money + money})
		if rUpdate.Error != nil {
			logger.Error("innerUpdateAccount INVITE update err:%v, uid:%v, money:%v, aType:%v", r.Error, uid, money, billType)
			return model.STATUS_DB_ERROR, rUpdate.Error
		} else {
			logger.Error("innerUpdateAccount INVITE update succ uid:%v, money:%v, aType:%v", uid, money, billType)
			return model.STATUS_DB_SUCC, nil
		}
		//充值
	} else if billType == model.STATUS_TYPE_RECHARGE {
		if r.Error != nil {
			logger.Error("innerUpdateAccount RECHARGE err:%v, uid:%v, money:%v, aType:%v", r.Error, uid, money, billType)
			return model.STATUS_DB_ERROR, fmt.Errorf("innerUpdateAccount RECHARGE err:%v, uid:%v, money:%v, aType:%v", uid, money, billType)
		}
		rUpdate := tx.Table("pubg_user_account").Where("uid = ? ", uid).Updates(map[string]interface{}{"deposit": account.Deposit + money, "history_money": account.HistoryMoney + money, "recharge_money": account.RechargeMoney + money, "money": account.Money + money})
		if rUpdate.Error != nil {
			logger.Error("innerUpdateAccount RECHARGE update err:%v, uid:%v, money:%v, aType:%v", r.Error, uid, money, billType)
			return model.STATUS_DB_ERROR, rUpdate.Error
		} else {
			logger.Error("innerUpdateAccount RECHARGE update succ uid:%v, money:%v, aType:%v", uid, money, billType)
			return model.STATUS_DB_SUCC, nil
		}
		//游戏门票
	} else if billType == model.STATUS_TYPE_GAME_TICKETS {
		if accountType == model.ACCOUNT_TYPE_BONUS {
			if account.Bonus < money || account.Money < money {
				logger.Debug("innerUpdateAccount TICKETS BONUS not enougth  uid:%v, money:%v, aType:%v, account:%v", uid, money, billType, account)
				return model.STATUS_ACCOUNT_NOT_ENOUGH, fmt.Errorf("innerUpdateAccount TICKETS BONUS not enougth  uid:%v, money:%v, aType:%v, account:%v", uid, money, billType, account)
			}
			rUpdate := tx.Table("pubg_user_account").Where("uid = ? ", uid).Updates(map[string]interface{}{"bonus": account.Bonus - money, "money": account.Money - money})
			if rUpdate.Error != nil {
				logger.Error("innerUpdateAccount TICKETS BONUS err:%v, uid:%v, money:%v, aType:%v", r.Error, uid, money, billType)
				return model.STATUS_DB_ERROR, rUpdate.Error
			} else {
				logger.Error("innerUpdateAccount TICKETS BONUS succ uid:%v, money:%v, aType:%v", uid, money, billType)
				return model.STATUS_DB_SUCC, nil
			}
		} else {
			if account.Money < money || account.Deposit < money {
				logger.Debug("innerUpdateAccount TICKETS BONUS not enougth  uid:%v, money:%v, aType:%v, account:%v", uid, money, billType, account)
				return model.STATUS_ACCOUNT_NOT_ENOUGH, fmt.Errorf("innerUpdateAccount TICKETS Deposit not enougth  uid:%v, money:%v, aType:%v, account:%v", uid, money, billType, account)
			}
			rUpdate := tx.Table("pubg_user_account").Where("uid = ? ", uid).Updates(map[string]interface{}{"deposit": account.Deposit - money, "money": account.Money - money})
			if rUpdate.Error != nil {
				logger.Error("innerUpdateAccount TICKETS DEPOSIT err:%v, uid:%v, money:%v, aType:%v", r.Error, uid, money, billType)
				return model.STATUS_DB_ERROR, rUpdate.Error
			} else {
				logger.Error("innerUpdateAccount TICKETS DEPOSIT succ uid:%v, money:%v, aType:%v", uid, money, billType)
				return model.STATUS_DB_SUCC, nil
			}
		}
		//游戏获胜
	} else if billType == model.STATUS_TYPE_GAME_WIN {
		rUpdate := tx.Table("pubg_user_account").Where("uid = ? ", uid).Updates(map[string]interface{}{"deposit": account.Deposit + money, "money": account.Money + money})
		if rUpdate.Error != nil {
			logger.Error("innerUpdateAccount GAME update err:%v, uid:%v, money:%v, aType:%v", r.Error, uid, money, billType)
			return model.STATUS_DB_ERROR, rUpdate.Error
		} else {
			logger.Error("innerUpdateAccount GAME update succ uid:%v, money:%v, aType:%v", uid, money, billType)
			return model.STATUS_DB_SUCC, nil
		}
		//用户提现
	} else if billType == model.STATUS_TYPE_WITHDEPOSISTS {
		if account.Deposit < money {
			logger.Debug("innerUpdateAccount WITHDEPOSISTS not enougth  uid:%v, money:%v, aType:%v, account:%v", uid, money, billType, account)
			return model.STATUS_ACCOUNT_NOT_ENOUGH, nil
		}
		rUpdate := tx.Table("pubg_user_account").Where("uid = ? ", uid).Updates(map[string]interface{}{"deposit": account.Deposit - money, "money": account.Money - money})
		if rUpdate.Error != nil {
			logger.Error("innerUpdateAccount WITHDEPOSISTS update err:%v, uid:%v, money:%v, aType:%v", r.Error, uid, money, billType)
			return model.STATUS_DB_ERROR, rUpdate.Error
		} else {
			logger.Error("innerUpdateAccount WITHDEPOSISTS update succ uid:%v, money:%v, aType:%v", uid, money, billType)
			return model.STATUS_DB_SUCC, nil
		}
	} else {
		return model.STATUS_DB_DATA_INEXISTENCE, nil
	}
}

func GetUserData(phone int64) (*model.PubgUserRecord, error) {
	record := new(model.PubgUserRecord)

	r := pubgsql.Get("pubg").Master().Table("pubg_user_info").Where(" phoneNumber = ? ", phone).Find(record)

	if r.Error != nil || r.RowsAffected != 1 {
		log.Printf("GetUserData err:%v, tuid:%v", r.Error, phone)
		return record, fmt.Errorf("GetUserData err:%v, tuid:%v", r.Error, phone)
	}

	return record, nil
}

//func GetBillIdSeqTable(t time.Time) string {
//	return fmt.Sprintf(model.TABLE_FINANCIAL_BILLID_SEQ, t.Year(), t.Month(), t.Day())
//}
//
//func GetFinancialBillTable(t time.Time) string {
//	return fmt.Sprintf(model.TABLE_FINANCIAL_BILL, t.Year(), t.Month(), t.Day())
//}
//
//func GetRechargeBillTable(t time.Time) string {
//	return fmt.Sprintf(model.TABLE_RECHARGE_BILL, t.Year(), t.Month(), t.Day())
//}

func GetBillIdSeqTable(uid int64) string {
	return fmt.Sprintf(model.TABLE_FINANCIAL_BILLID_SEQ, uid%8)
}

func GetFinancialBillTable(uid int64) string {
	return fmt.Sprintf(model.TABLE_FINANCIAL_BILL, uid%8)
}
