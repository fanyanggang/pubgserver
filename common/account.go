package game

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wonderivan/logger"
	"log"
	"math/rand"
	"pubgserver/dao"
	"pubgserver/model"
	//pubgrpc "pubgserver/rpc-go"
	pubgrpc "github.com/fanyanggang/BackendPlatform/frpc-go"
	"strconv"
	"strings"
	"time"
)

func AddUserFinancial(uid, phone int64, address, mail, name string, gtype int) error {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("AddUserFinancial panic.  err:%v", r)
			logger.Error(msg)
		}
	}()

	return dao.AddUserFinancialDB(uid, phone, address, mail, name, gtype)
}

func GetUserFinancial(uid int64) (model.UserFinanciaRecord, int, error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetUserFinancial panic.  err:%v", r)
			logger.Error(msg)
		}
	}()
	return dao.GetUserFinancialDB(uid)
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func GetMsgCode(phoneNum int64) int {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("GetMsgCode panic.  err:%v", r)
			logger.Error(msg)
		}
	}()

	code := GenValidateCode(4)

	redis := GetRedis()
	if redis == nil {
		log.Print("getRedisClient err")
		return model.STATUS_MSGCODE_ERROR
	}

	content := fmt.Sprintf("%s  is OTP for your Jeeto Games account. Please do not share it with anyone.", code)
	data := map[string]interface{}{
		"phone":   phoneNum,
		"content": content,
		"app_id":  "pubg",
		"type":    1,
	}
	//发送用户短信验证码
	body, _ := json.Marshal(data)
	url := "http://msg.jeeto.work/v1/send_sms_json"
	//body := fmt.Sprintf("phone=%d&content=%s", phoneNum, content)
	resp, err := pubgrpc.HttpPost(context.Background(), url, string(body))

	if err != nil {
		logger.Debug("GetMsgCode err:%v, phone:%v, data:%v", err, phoneNum, string(body))
		return model.STATUS_MSGCODE_ERROR
	}
	ret, err := redis.SetExSecond(phoneNum, code, 300)
	if err != nil {
		logger.Debug("GetMsgCode err:%v,ret:%v, resp:%v, phoneNum:%v, code:%v", err, ret, string(resp), phoneNum, code)
		return model.STATUS_MSGCODE_ERROR
	}

	logger.Debug("body: %v", string(body))
	logger.Debug("GetMsgCode succ:%v, phoneNum:%v, code:%v", string(resp), phoneNum, code)
	return model.STATUS_CODE_SUCC
}

func innerVerfiyCode(phoneNum int64, code int) error {
	redis := GetRedis()
	if redis == nil {
		log.Print("getRedisClient err")
		return fmt.Errorf("getRedisClient err")
	}

	str := strconv.FormatInt(phoneNum, 10)
	value, err := redis.Get(str)

	if err != nil {
		log.Printf("innerVerfiyCode redis get err:%v", err)
		return nil
	}

	intCode, err := strconv.Atoi(string(value))
	log.Printf("innerVerfiyCode:%v, code:%v, intcode:%v", phoneNum, code, intCode)

	if intCode != code {
		return fmt.Errorf("code illegal:%v", code)
	}

	redis.Del(str)
	return nil
}
