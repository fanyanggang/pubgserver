package game

import (
	"fmt"
	pubgsql "github.com/fanyanggang/BackendPlatform/fsql"
	"github.com/wonderivan/logger"
	"time"
)

const (
	DB_COUNT = 4
)

func CheckTables() {
	for {
		CheckTablesImp()
		time.Sleep(time.Minute * 5)
	}
}

func CheckTablesImp() {
	//创建连续3天的表
	//for i := time.Duration(0); i < DB_COUNT; i++ {
	//	t := time.Now().Add(time.Hour * 24 * i)
	//	billSeqTable := dao.GetBillIdSeqTable(t)
	//	CreateTable(model.DATABASE, billSeqTable, model.TABLE_FINANCIAL_BILLID_SEQ_BASE)
	//
	//	financialBill := dao.GetFinancialBillTable(t)
	//	CreateTable(model.DATABASE, financialBill, model.TABLE_FINANCIAL_BILL_BASE)
	//
	//	rechargeBill := dao.GetRechargeBillTable(t)
	//	CreateTable(model.DATABASE, rechargeBill, model.TABLE_RECHARGE_BILL_BASE)
	//}
}

func CreateTable(db, newTableName, fromTableName string) bool {
	if pubgsql.Get(db).Master().HasTable(newTableName) {
		return true
	}

	createTableSql := fmt.Sprintf("CREATE TABLE `%s` LIKE `%s`", newTableName, fromTableName)
	ctDB := pubgsql.Get(db).Master().Exec(createTableSql)
	if ctDB.Error != nil {
		logger.Error("CreateTable create table error:%v, sql:%s", ctDB.Error, createTableSql)
		return false
	}
	return true
}
