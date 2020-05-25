package game

import (
	"log"
	"net/http"
)

func InitService() {

	http.HandleFunc("/user/login", UserLoginHandler)
	http.HandleFunc("/user/msgcode/get", GetVerifyCodeHandler)
	http.HandleFunc("/user/msgcode/verify", MsgCodeVerifyHandler)

	http.HandleFunc("/user/pay/status", UserPayStatusHandler)

	http.HandleFunc("/room/info/get", GetRoomConfigHandler)
	http.HandleFunc("/room/prize/get", GetRoomPrizeHandler)

	//添加游戏昵称
	http.HandleFunc("/pubg/name/add", AddPubgNameHandler)
	//添加对名
	http.HandleFunc("/pubg/team/join", JoinTeamHandler)
	//获取团队列表
	http.HandleFunc("/pubg/team/get", GetTeamHandler)

	//用户支付门票
	http.HandleFunc("/pubg/user/pay", UserPayHandler)

	//房间奖金配置
	http.HandleFunc("/pubg/room/bonus/get", GetRoomBonusHandler)

	//游戏获胜
	http.HandleFunc("/pubg/game/win", GameWinHandler)

	//用户充值交易订单获取
	http.HandleFunc("/pubg/recharge/bill/get", GetUserRechargeBillHandler)
	//用户充值
	//http.HandleFunc("/pubg/user/recharge", UserRechargeHandler)
	//用户充值状态更新
	http.HandleFunc("/pubg/recharge/status/update", UserRechargeStatusUpdateHandler)
	//充值订单状态查询
	http.HandleFunc("/pubg/recharge/status/get", UserRechargeStatusGetHandler)

	//用户提现交易订单获取
	http.HandleFunc("/pubg/withdeposits/bill/get", GetUserWithDepositsBillHandler)
	//用户提现状态更新
	http.HandleFunc("/pubg/withdeposits/status/update", UserWithDepositsStatusUpdateHandler)
	//充值订单状态查询
	http.HandleFunc("/pubg/withdeposits/status/get", UserWithDepositsStatusGetHandler)

	//用户金融信息查询
	http.HandleFunc("/pubg/user/financial/get", GetUserFinancialHandler)
	//用户金融信息更新
	http.HandleFunc("/pubg/user/financial/add", AddUserFinancialHandler)

	//用户中心数据
	http.HandleFunc("/pubg/user/get", GetUserCenterDataHandler)

	//用户交易记录查询
	http.HandleFunc("/pubg/transaction/get", GetUserTransactionRecordHandler)

	//用户提现
	http.HandleFunc("/pubg/deposit/get", UserDepositGetHandler)

	// 设置监听的端口
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
