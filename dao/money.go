package dao

import (
	"fmt"
	pubgsql "github.com/fanyanggang/BackendPlatform/fsql"
	"github.com/wonderivan/logger"
	"math/rand"
	"pubgserver/model"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GetBillId(uid int64) string {
	t := time.Now()

	billid := fmt.Sprintf("AD%d%d%04d%02d", t.UnixNano()/1000, uid%8192+1000, 1024, random.Intn(99))
	logger.Debug("genBillId billid:%v", billid)

	return billid
}

func UserDepositGetDB(uid, money int64, tradeBillID string) (int, error) {
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	billid := GetBillId(uid)
	record := model.FinancialBillSeqRecord{
		UID:         uid,
		BillID:      billid,
		TradeBillID: tradeBillID,
		BillType:    model.STATUS_TYPE_WITHDEPOSISTS,
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
		RoomID:      0,
		Money:       money,
		BillType:    model.STATUS_TYPE_WITHDEPOSISTS,
		AccountType: model.ACCOUNT_TYPE_DEPOSIT,
		Status:      model.STATUS_BILL_UNPROCESSED,
		CreateTime:  time.Now().Unix(),
	}

	rInsert := tx.Table(billTable).Create(billRecord)
	if rInsert.Error != nil {
		tx.Rollback()
		logger.Error("UserPayDB bill create err:%v, data:%v", rInsert.Error, record)
		return model.STATUS_DB_ERROR, r.Error
	}

	status, err := innerUpdateAccount(tx, uid, money, model.STATUS_TYPE_WITHDEPOSISTS, model.ACCOUNT_TYPE_DEPOSIT)
	if err != nil {
		r = tx.Table(billTable).Where("bill_id = ? ", billid).Update("Status", model.STATUS_BILL_FAIL)
		if r.Error != nil {
			tx.Rollback()
			logger.Error("innerUpdate update bill err:%v, uid:%v, money:%v, aType:%v, accountType:%v", rInsert.Error, uid, money, model.STATUS_TYPE_RECHARGE, model.ACCOUNT_TYPE_DEPOSIT)
			return model.STATUS_DB_ERROR, r.Error
		}

		return model.STATUS_ACCOUNT_FAIL, nil
	} else {
		if status == model.STATUS_ACCOUNT_NOT_ENOUGH {
			r = tx.Table(billTable).Where("bill_id = ? ", billid).Update("Status", model.STATUS_BILL_ACOUNT_NOT)
		} else {
			r = tx.Table(billTable).Where("bill_id = ? ", billid).Update("Status", model.STATUS_BILL_SUCC)
		}
		if r.Error != nil {
			tx.Rollback()
			logger.Error("innerUpdate update bill err:%v, uid:%v, money:%v, aType:%v, accountType:%v", rInsert.Error, uid, money, model.STATUS_TYPE_RECHARGE, model.ACCOUNT_TYPE_DEPOSIT)
			return model.STATUS_DB_ERROR, r.Error
		}
		return model.STATUS_ACCOUNT_SUCC, nil
	}
}

func UserTradeStatusGetDB(uid int64, tradeBillID string) (int, error) {

	seqTableName := GetBillIdSeqTable(uid)
	record := &model.FinancialBillSeqRecord{}
	tx := pubgsql.Get("financial").Master().Table(seqTableName).Where("trade_bill_id = ? ", tradeBillID).Find(record)
	if tx.Error != nil {
		logger.Error("UserRechargeStatusGetDB get billseq err:%v, uid:%v, tradeBillID：%v", tx.Error, uid, tradeBillID)
		return model.STATUS_DB_DATA_INEXISTENCE, tx.Error
	}

	financialRecord := &model.FinancialBillRecord{}
	TableName := GetFinancialBillTable(uid)
	r := pubgsql.Get("financial").Master().Table(TableName).Where("bill_id = ? ", record.BillID).Find(financialRecord)
	if r.Error != nil {
		logger.Error("UserRechargeStatusGetDB get bill err:%v, uid:%v, tradeBillID：%v", tx.Error, uid, tradeBillID)
		return model.STATUS_DB_DATA_INEXISTENCE, tx.Error
	}

	logger.Info("UserRechargeStatusGetDB succ, uid:%v, tradeBillID:%v, financialRecord:%v", uid, tradeBillID, financialRecord)
	return financialRecord.Status, nil
}

func UserWithDepositsStatusUpdateDB(uid, money int64, tradeBillID string, status int) (int, error) {
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	seqTableName := GetBillIdSeqTable(uid)
	record := &model.FinancialBillSeqRecord{}
	r := tx.Table(seqTableName).Where("trade_bill_id = ? ", tradeBillID).Find(record)
	if r.Error != nil {
		logger.Error("UserRechargeStatusGetDB get billseq err:%v, uid:%v, tradeBillID：%v", tx.Error, uid, tradeBillID)
		return model.STATUS_DB_DATA_INEXISTENCE, tx.Error
	}

	financialRecord := &model.FinancialBillRecord{}
	TableName := GetFinancialBillTable(uid)
	r = tx.Table(TableName).Where("bill_id = ? ", record.BillID).Find(financialRecord)
	if r.Error != nil {
		logger.Error("UserRechargeStatusGetDB get bill err:%v, uid:%v, tradeBillID：%v", tx.Error, uid, tradeBillID)
		return model.STATUS_DB_DATA_INEXISTENCE, tx.Error
	}

	if status == model.STATUS_TRADE_SUCC {
		r = tx.Table(TableName).Where("bill_id = ? ", record.BillID).Update("status", model.STATUS_BILL_SUCC)
		if r.Error != nil {
			logger.Info("UserWithDepositsStatusUpdateDB update err:%v, uid:%v, tradebillid:%v", r.Error, uid, tradeBillID)
			return model.STATUS_DB_ERROR, r.Error
		}
		logger.Info("UserWithDepositsStatusUpdateDB update succ, uid:%v, tradebillid:%v", uid, tradeBillID)
		return model.STATUS_ACCOUNT_SUCC, nil
	} else {
		account := &model.UserAccountRecord{}
		r = tx.Table("pubg_user_account").Where("uid = ? ", uid).Find(account)
		if r.Error != nil {
			logger.Info("UserWithDepositsStatusUpdateDB pubg_user_account err:%v, uid:%v, tradebillid:%v", r.Error, uid, tradeBillID)
			return model.STATUS_DB_ERROR, r.Error
		}

		rUpdate := tx.Table("pubg_user_account").Where("uid = ? ", uid).Updates(map[string]interface{}{"deposit": account.Deposit + money, "money": account.Money + money})
		if rUpdate.Error != nil {
			logger.Error("UserWithDepositsStatusUpdateDB WithDeposits  err:%v, uid:%v, money:%v, tradeBillID:%v", r.Error, uid, money, tradeBillID)
			return model.STATUS_DB_ERROR, rUpdate.Error
		}

		r = tx.Table(TableName).Where("bill_id = ? ", record.BillID).Update("status", model.STATUS_BILL_FAIL)
		if r.Error != nil {
			logger.Info("UserWithDepositsStatusUpdateDB WithDeposits bill err:%v, uid:%v, tradebillid:%v", r.Error, uid, tradeBillID)
			return model.STATUS_DB_ERROR, r.Error
		}

		logger.Error("UserWithDepositsStatusUpdateDB WithDeposits succ uid:%v, money:%v, tradeBillID:%v", uid, money, tradeBillID)
		return model.STATUS_DB_SUCC, nil
	}
}

func UpdateBillStatus(uid int64, billID, tradeBillID string) error {
	seqTableName := GetBillIdSeqTable(uid)
	r := pubgsql.Get("financial").Master().Table(seqTableName).Where("bill_id = ? ", billID).Update("trade_bill_id", tradeBillID)
	if r.Error != nil {
		logger.Error("UpdateBillStatus error:%v, uid:%d, billID:%d, tradeBillID:%d\n", uid, billID, tradeBillID)
		return r.Error
	}

	return nil
}

func UserWithDepositsErrRollBack(uid, money int64, billid string) (int, error) {
	financialRecord := &model.FinancialBillRecord{}
	TableName := GetFinancialBillTable(uid)
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	r := tx.Table(TableName).Where("bill_id = ? ", billid).Find(financialRecord)
	if r.Error != nil {
		logger.Error("UserRechargeStatusGetDB get bill err:%v, uid:%v, billid：%v", tx.Error, uid, billid)
		return model.STATUS_DB_DATA_INEXISTENCE, tx.Error
	}

	r = tx.Table(TableName).Where("bill_id = ? ", billid).Update("status", model.STATUS_BILL_SYSTEM_ERROR)
	if r.Error != nil {
		logger.Info("UserWithDepositsStatusUpdateDB update err:%v, uid:%v, billid:%v", r.Error, uid, billid)
		return model.STATUS_DB_ERROR, r.Error
	}

	account := &model.UserAccountRecord{}
	r = tx.Table("pubg_user_account").Where("uid = ? ", uid).Find(account)
	if r.Error != nil {
		logger.Info("UserWithDepositsStatusUpdateDB pubg_user_account err:%v, uid:%v, billid:%v", r.Error, uid, billid)
		return model.STATUS_DB_ERROR, r.Error
	}

	rUpdate := tx.Table("pubg_user_account").Where("uid = ? ", uid).Updates(map[string]interface{}{"deposit": account.Deposit + money, "money": account.Money + money})
	if rUpdate.Error != nil {
		logger.Error("UserWithDepositsStatusUpdateDB WithDeposits  err:%v, uid:%v, money:%v, billid:%v", r.Error, uid, money, billid)
		return model.STATUS_DB_ERROR, rUpdate.Error
	}

	return model.STATUS_DB_SUCC, nil
}

func CreateUserWithDepositsBillRecord(uid, money int64, billid, tradeBillID string) (int, error) {
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	record := model.FinancialBillSeqRecord{
		UID:         uid,
		BillID:      billid,
		TradeBillID: tradeBillID,
		BillType:    model.STATUS_TYPE_WITHDEPOSISTS,
		CreateTime:  time.Now().Unix(),
	}

	seqTableName := GetBillIdSeqTable(uid)
	r := tx.Table(seqTableName).Create(&record)
	if r.Error != nil {
		tx.Rollback()
		logger.Error("UserRechargeStatusUpdateDB create billseq err:%v, uid:%v", r.Error, uid)
		return model.STATUS_DB_ERROR, r.Error
	}

	billRecord := &model.FinancialBillRecord{
		UID:         uid,
		BillID:      billid,
		RoomID:      0,
		Money:       money,
		BillType:    model.STATUS_TYPE_WITHDEPOSISTS,
		AccountType: model.ACCOUNT_TYPE_DEPOSIT,
		Status:      model.STATUS_BILL_UNPROCESSED,
		CreateTime:  time.Now().Unix(),
	}

	billTable := GetFinancialBillTable(uid)
	rInsert := tx.Table(billTable).Create(billRecord)
	if rInsert.Error != nil {
		tx.Rollback()
		logger.Error("UserPayDB bill create err:%v, data:%v", rInsert.Error, record)
		return model.STATUS_DB_ERROR, r.Error
	}

	status, err := innerUpdateAccount(tx, uid, money, model.STATUS_TYPE_WITHDEPOSISTS, model.ACCOUNT_TYPE_DEPOSIT)

	if status == model.STATUS_DB_ERROR {
		tx.Rollback()
		return model.STATUS_DB_ERROR, err
	} else if status == model.STATUS_ACCOUNT_NOT_ENOUGH {
		r = tx.Table(billTable).Where("bill_id = ? ", billid).Update("Status", model.STATUS_BILL_ACOUNT_NOT)
		if r.Error != nil {
			tx.Rollback()
			logger.Error("innerUpdate update bill err:%v, uid:%v, money:%v", rInsert.Error, uid, money)
			return model.STATUS_DB_ERROR, r.Error
		}
		return model.STATUS_ACCOUNT_NOT_ENOUGH, nil
	}

	logger.Info("CreateUserWithDepositsBillRecord succ, uid:%v, money:%v, billid:%v, status:%v，err:%v", uid, money, billid, status, err)
	return status, err
}

func UserRechargeStatusUpdateDB(uid, money int64, tradeBillID, billid string, status int) (int, error) {
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	record := model.FinancialBillSeqRecord{
		UID:         uid,
		BillID:      billid,
		TradeBillID: tradeBillID,
		BillType:    model.STATUS_TYPE_RECHARGE,
		CreateTime:  time.Now().Unix(),
	}

	seqTableName := GetBillIdSeqTable(uid)
	r := tx.Table(seqTableName).Create(&record)
	if r.Error != nil {
		tx.Rollback()
		logger.Error("UserRechargeStatusUpdateDB create billseq err:%v, uid:%v", r.Error, uid)
		return model.STATUS_DB_ERROR, r.Error
	}

	billTable := GetFinancialBillTable(uid)
	if status == model.STATUS_TRADE_FAIL || status == model.STATUS_TRADE_CANCEL {
		billRecord := &model.FinancialBillRecord{
			UID:         uid,
			BillID:      tradeBillID,
			RoomID:      0,
			Money:       money,
			BillType:    model.STATUS_TYPE_RECHARGE,
			AccountType: model.ACCOUNT_TYPE_DEPOSIT,
			Status:      model.STATUS_BILL_FAIL,
			CreateTime:  time.Now().Unix(),
		}

		rInsert := tx.Table(billTable).Create(billRecord)
		if rInsert.Error != nil {
			tx.Rollback()
			logger.Error("UserPayDB bill create err:%v, data:%v", rInsert.Error, record)
			return model.STATUS_DB_ERROR, r.Error
		}

		logger.Info("UserRechargeUpdateDB succ uid:%v, money:%v, status:%v", uid, money, status)
		return model.STATUS_ACCOUNT_SUCC, nil

	} else if status == model.STATUS_TRADE_SUCC {
		_, err := innerUpdateAccount(tx, uid, money, model.STATUS_TYPE_RECHARGE, model.ACCOUNT_TYPE_DEPOSIT)

		if err != nil {
			billRecord := &model.FinancialBillRecord{
				UID:         uid,
				BillID:      billid,
				RoomID:      0,
				Money:       money,
				BillType:    model.STATUS_TYPE_RECHARGE,
				AccountType: model.ACCOUNT_TYPE_DEPOSIT,
				Status:      model.STATUS_BILL_FAIL,
				CreateTime:  time.Now().Unix(),
			}

			rInsert := tx.Table(billTable).Create(billRecord)
			if rInsert.Error != nil {
				tx.Rollback()
				logger.Error("UserPayDB bill create err:%v, data:%v", rInsert.Error, record)
				return model.STATUS_DB_ERROR, r.Error
			}
			logger.Error("UserRechargeUpdateDB err :%v, uid:%v, billid:%v", r.Error, uid, billid)
			return model.STATUS_DB_ERROR, r.Error
		} else {
			billRecord := &model.FinancialBillRecord{
				UID:         uid,
				BillID:      billid,
				RoomID:      0,
				Money:       money,
				BillType:    model.STATUS_TYPE_RECHARGE,
				AccountType: model.ACCOUNT_TYPE_DEPOSIT,
				Status:      model.STATUS_BILL_SUCC,
				CreateTime:  time.Now().Unix(),
			}

			rInsert := tx.Table(billTable).Create(billRecord)
			if rInsert.Error != nil {
				tx.Rollback()
				logger.Error("UserPayDB bill create err:%v, data:%v", rInsert.Error, record)
				return model.STATUS_DB_ERROR, r.Error
			}

			logger.Info("UserRechargeUpdateDB succ uid:%v, money:%v, status:%v", uid, money, status)
			return model.STATUS_ACCOUNT_SUCC, nil
		}
	} else {
		logger.Info("UserRechargeUpdateDB data err uid:%v, money:%v, status:%v", uid, money, status)
		return model.STATUS_DB_DATA_INEXISTENCE, nil
	}
}

func UserPayStatusDB(roomID, uid int64) (int, error) {
	billTable := GetFinancialBillTable(uid)
	billRecord := &model.FinancialBillRecord{}

	tx := pubgsql.Get("financial").Master().Table(billTable).Where("uid = ? and room_id = ? and status = 1", uid, roomID).Find(billRecord)
	if tx.Error != nil {
		logger.Error("UserPayStatusDB err:%v, uid:%v, roomid:%v", tx.Error, uid, roomID)
		return 0, tx.Error
	}

	logger.Error("UserPayStatusDB succ:%v, uid:%v, roomid:%v", billRecord, uid, roomID)
	return billRecord.Status, nil
}

func UserRechargeDB(uid, money int64, tradeBillID string) (int, error) {
	tx := pubgsql.Get("financial").Master().Begin()
	defer tx.Commit()

	billid := GetBillId(uid)
	record := model.FinancialBillSeqRecord{
		UID:         uid,
		BillID:      billid,
		TradeBillID: tradeBillID,
		BillType:    model.STATUS_TYPE_RECHARGE,
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
		RoomID:      0,
		Money:       money,
		BillType:    model.STATUS_TYPE_RECHARGE,
		AccountType: model.ACCOUNT_TYPE_DEPOSIT,
		Status:      model.STATUS_BILL_UNPROCESSED,
		CreateTime:  time.Now().Unix(),
	}

	rInsert := tx.Table(billTable).Create(billRecord)
	if rInsert.Error != nil {
		tx.Rollback()
		logger.Error("UserPayDB bill create err:%v, data:%v", rInsert.Error, record)
		return model.STATUS_DB_ERROR, r.Error
	}

	return model.STATUS_ACCOUNT_SUCC, nil

}

//tx := pubgsql.Get(model.DATABASE).Master().Begin()
//defer tx.Commit()
//
//billid := genBillId(uid)
//record := model.RechargeBillRecord{
//	UID:        uid,
//	BillID:     billid,
//	From:       from,
//	Money:      money,
//	Status:     model.STATUS_BILL_UNPROCESSED,
//	CreateTime: time.Now().Unix(),
//}
//
//TableName := GetRechargeBillTable(time.Now())
//result := tx.Table(TableName).Create(&record)
//
//if result.Error != nil {
//	tx.Rollback()
//	logger.Error("CreateRechargeBillTable error:%v, uid:%d, money:%d, time:%d\n", result.Error, uid, money, time.Now().Unix())
//	return result.Error
//}
//
////更新用户账户
//err := updateUserAccount(tx, uid, money)
//if err != nil {
//	r := tx.Table(TableName).Where("bill_id = ? ", billid).Update("status", model.STATUS_BILL_FAIL)
//	if r.Error != nil {
//		tx.Rollback()
//		logger.Error("UserRechargeDB updateUserAccount error:%v, uid:%d, money:%d, time:%d\n", result.Error, uid, money, time.Now().Unix())
//		return result.Error
//	}
//	return nil
//} else {
//	r := tx.Table(TableName).Where("bill_id = ? ", billid).Update("status", model.STATUS_BILL_PROCESSED)
//	if r.Error != nil {
//		tx.Rollback()
//		logger.Error("UserRechargeDB update bill status err:%v, uid:%d, money:%d, time:%d\n", result.Error, uid, money, time.Now().Unix())
//		return result.Error
//	}
//	return nil
//}

//func updateUserAccount(tx *gorm.DB, uid int64, money int) error {
//	//更新用户账户
//	account := &model.UserAccountRecord{}
//	r := tx.Table("pubg_user_account").Where("uid = ? ", uid).Set("gorm:query_option", "FOR UPDATE").Find(account)
//
//	if r.RecordNotFound() {
//		account = &model.UserAccountRecord{
//			UID:          uid,
//			Money:        money,
//			HistoryMoney: int64(money),
//			CreateTime:   time.Now().Unix(),
//			UpdateTime:   time.Now().Unix(),
//		}
//
//		rInsert := tx.Table("pubg_user_account").Create(account)
//
//		if rInsert.RowsAffected != 1 {
//			logger.Error("updateUserAccount create RowsAffected :%v, uid:%v, money:%v", rInsert.RowsAffected, uid, money)
//			return fmt.Errorf("updateUserAccount create RowsAffected :%v, uid:%v, money:%v", rInsert.RowsAffected, uid, money)
//
//		} else if rInsert.Error != nil {
//			tx.Rollback()
//			logger.Error("updateUserAccount create err:%v, uid:%v, money:%v", rInsert.Error, uid, money)
//			return rInsert.Error
//		}
//
//		return nil
//
//	} else if r.Error != nil {
//		tx.Rollback()
//		logger.Error("updateUserAccount find err:%v, uid:%v, money:%v", r.Error, uid, money)
//		return r.Error
//	} else {
//		rUpdateAccount := tx.Table("pubg_user_account").Where("uid = ? ", uid).Updates(map[string]interface{}{"money": account.Money + money, "history_money": account.HistoryMoney + int64(money)})
//		if rUpdateAccount.Error != nil {
//			tx.Rollback()
//			logger.Error("updateUserAccount update err:%v, uid:%v, money:%v", rUpdateAccount.Error, uid, money)
//			return rUpdateAccount.Error
//		} else if rUpdateAccount.RowsAffected != 1 {
//			logger.Error("updateUserAccount update RowsAffected :%v, uid:%v, money:%v", rUpdateAccount.RowsAffected, uid, money)
//			return fmt.Errorf("updateUserAccount update RowsAffected :%v, uid:%v, money:%v", rUpdateAccount.RowsAffected, uid, money)
//		}
//		return nil
//	}
//}
