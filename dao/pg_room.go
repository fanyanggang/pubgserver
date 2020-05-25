package dao

import (
	"fmt"
	pubgsql "github.com/fanyanggang/BackendPlatform/fsql"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"pubgserver/model"
	"strconv"
	"time"
)

func GetRoomDB(phone, uid int64) ([]model.RoomRecord, error) {

	record := make([]model.RoomRecord, 0)

	//获取当前已上线以及下线的前三名的
	//sql := "select * from `pg_rooms` where status = 0 order by created"
	currentTime := time.Now().AddDate(0, 0, -2).Unix()
	r := pubgsql.Get("pubg").Master().Table("pg_rooms").Where("start_time > ? ", currentTime).Order("sign_time desc").Scan(&record)

	if r.Error != nil {
		return record, r.Error
	}

	logger.Debug("GetRoomConfDB succ: phone:%v,  record:%v", phone, record)
	return record, nil
}

func GetRoomBonusByRoomID(uid, roomid int64) ([]model.RoomConfigRecord, error) {
	result := make([]model.RoomConfigRecord, 0)
	r := pubgsql.Get("pubg").Master().Table("pg_rooms_config").Where("room_id = ? ", roomid).Find(&result)
	if r.Error != nil {
		logger.Error("GetRoomConfigByRoomID err:%v, uid:%v, roomid:%v", r.Error, uid, roomid)
		return result, r.Error
	}
	return result, nil
}

func GetRoomUserNumByRoomID() ([]model.RoomUserNumByRoomID, error) {

	result := make([]model.RoomUserNumByRoomID, 0)
	sql := "select room_id , COUNT(*) as count from pg_rooms_player where status = 1 group by room_id"

	r := pubgsql.Get("pubg").Master().Table("pg_rooms_player").Raw(sql).Count(&result)
	if r.Error != nil {
		logger.Error("GetRoomUserNumByRoomID err:%v", r.Error)
		return result, r.Error
	}

	logger.Debug("GetRoomUserNumByRoomID succ: roomid:%v ", result)
	return result, nil
}

func GetRoomPrizeDB(uid, phone, roomid int64, solo int) (model.SoleResultResp, error) {

	result := model.SoleResultResp{}

	if solo == model.ROOM_PRIZE_SOLO {
		return innerGetPrizeSolo(roomid)
	} else if solo == model.ROOM_PRIZE_SOLO_TOP {
		return innerGetPrizeSoloTop(roomid)
	} else if solo == model.ROOM_PRIZE_SOLO_SQUARD {
		return innerGetPrizeSoloSquard(roomid)
	} else {
		return result, nil
	}

	return result, nil
}

type rankData struct {
	ID    int
	Kills int
	Prize int64
}

func innerGetPrizeSoloSquard(roomid int64) (model.SoleResultResp, error) {
	result := model.SoleResultResp{}
	record := make([]model.RoomsDuoSquadRecord, 0)
	r := pubgsql.Get("pubg").Master().Table("pg_rooms_duo_squad").Where("room_id = ?", roomid).Order("rank, kills").Find(&record)
	if r.Error != nil {
		logger.Error("innerGetPrizeSoloSquard err:%v, roomid:%v", r.Error, roomid)
		return result, r.Error
	}

	userRank := make(map[int64]rankData)
	rankNum := 1
	uids := make([]int64, 0)
	for _, v := range record {
		uids = append(uids, v.UID)

		temp := rankData{
			Kills: v.KillNum,
			Prize: v.Prize,
			ID:    v.Rank,
		}
		userRank[v.UID] = temp
		rankNum++
	}
	logger.Debug("innerGetPrizeSoloSquard userRank:%v, uids:%v", userRank, uids)
	return getUserDataDB(userRank, uids)
}

func innerGetPrizeSoloTop(roomid int64) (model.SoleResultResp, error) {
	result := model.SoleResultResp{}
	record := make([]model.RoomsSoloTopRecord, 0)
	r := pubgsql.Get("pubg").Master().Table("pg_rooms_solo_top").Where("room_id = ?", roomid).Order("rank").Find(&record)
	if r.Error != nil {
		logger.Error("innerGetPrizeSolo err:%v, roomid:%v", r.Error, roomid)
		return result, r.Error
	}

	userRank := make(map[int64]rankData)
	rankNum := 1
	uids := make([]int64, 0)
	for _, v := range record {
		uids = append(uids, v.Uid)

		temp := rankData{
			Kills: v.KillNum,
			Prize: v.Prize,
			ID:    v.Rank,
		}
		userRank[v.Uid] = temp
		rankNum++
	}
	return getUserDataDB(userRank, uids)

}

func innerGetPrizeSolo(roomid int64) (model.SoleResultResp, error) {
	result := model.SoleResultResp{}
	record := make([]model.RoomPlayerRecord, 0)
	r := pubgsql.Get("pubg").Master().Table("pg_rooms_player").Where("kills > 0 and room_id = ? ", roomid).Order("kills").Find(&record)
	if r.Error != nil {
		logger.Error("innerGetPrizeSolo err:%v, roomid:%v", r.Error, roomid)
		return result, r.Error
	}

	userRank := make(map[int64]rankData)
	rankNum := 1
	uids := make([]int64, 0)
	for _, v := range record {
		uids = append(uids, v.UID)

		temp := rankData{
			Kills: v.KillNum,
			Prize: v.Prize,
			ID:    rankNum,
		}
		userRank[v.UID] = temp
		rankNum++
	}
	logger.Debug("innerGetPrizeSolo userRank:%v", userRank)
	return getUserDataDB(userRank, uids)
}

func getUserDataDB(userRank map[int64]rankData, uids []int64) (model.SoleResultResp, error) {
	record := model.SoleResultResp{}

	record.Result = make(map[int][]model.SolePerson, 0)
	userRecord := make([]model.PubgUserRecord, 0)
	r := pubgsql.Get("pubg").Master().Table("pubg_user_info").Where("id IN (?)", uids).Scan(&userRecord)
	if r.Error != nil {
		logger.Error("getUserDataDB err:%v, uids", r.Error, uids)
		return record, r.Error
	}

	logger.Debug("getUserDataDB beign :%v, uids:%v, userRank:%v", userRecord, uids, userRank)
	for _, v := range userRecord {
		logger.Debug("getUserDataDB v:%v", v)
		rank, ok := userRank[v.ID]
		if !ok {
			continue
		}
		logger.Debug("getUserDataDB vrank:%v", rank)

		phone := strconv.FormatInt(v.Phone, 10)

		var strPhone string
		for i := 0; i < len(phone); i++ {
			if i > 2 && i < 7 {
				strPhone = strPhone + "*"
			} else {
				strPhone = strPhone + string(phone[i])
			}
		}

		temp := model.SolePerson{
			Uid:   v.ID,
			Name:  v.Nickname,
			Photo: v.Photo,
			Kills: rank.Kills,
			Prize: rank.Prize,
			Phone: strPhone,
		}
		logger.Debug("getUserDataDB vtemp:%v", temp)

		solo, ok := record.Result[rank.ID]
		logger.Debug("getUserDataDB solo a:%v", solo)

		if !ok {
			record.Result[rank.ID] = []model.SolePerson{temp}
		} else {
			record.Result[rank.ID] = append(record.Result[rank.ID], temp)
		}
	}

	logger.Debug("getUserDataDB succ record:%v", record)
	return record, nil
}

func JoinTeamDB(uid, roomID int64, teamname string) error {

	tx := pubgsql.Get("pubg").Master().Begin()

	defer tx.Commit()
	record := make([]model.RoomTeamRecord, 0)

	r := tx.Table("pg_rooms_team").Where(" room_id = ? ", roomID).Order("team_num").Find(&record)

	var teamNum int
	bflags := false
	//insert new
	if r.RecordNotFound() {
		data := &model.RoomTeamRecord{
			RoomID:   roomID,
			TeamNum:  1,
			TeamName: teamname,
			Create:   time.Now().Unix(),
		}

		teamNum = 1
		rInsert := pubgsql.Get("pubg").Master().Table("pg_rooms_team").Create(data)
		if rInsert.Error != nil {
			logger.Error("JoinTeamDB insert err:%v, data:%v", rInsert.Error, data)
			tx.Rollback()
			return rInsert.Error
		}

		if rInsert.RowsAffected != 1 {
			logger.Error("insert RowsAffected err:%d, data:%v", rInsert.RowsAffected, data)
			tx.Rollback()
			return fmt.Errorf(" JoinTeamDB insert RowsAffected err:%d , data:%v", rInsert.RowsAffected, data)
		}
		err := updateRoomPlayer(tx, uid, roomID, 1)
		if err != nil {
			logger.Error("updateRoomPlayer err:%d, data:%v", rInsert.RowsAffected, data)
			tx.Rollback()
			return fmt.Errorf(" updateRoomPlayer err:%d , data:%v", rInsert.RowsAffected, data)
		}

	} else if r.Error != nil {
		logger.Error("JoinTeamDB err:%v, uid:%v, roomid:%v, teamname:%v", r.Error, uid, roomID, teamname)
		return r.Error
	} else {
		for _, v := range record {
			teamNum = v.TeamNum
			if v.TeamName == teamname {
				bflags = true
				break
			}
		}

		//没有当前记录
		if !bflags {
			teamNum++
			data := &model.RoomTeamRecord{
				RoomID:   roomID,
				TeamNum:  teamNum,
				TeamName: teamname,
				Create:   time.Now().Unix(),
			}

			rInsert := pubgsql.Get("pubg").Master().Table("pg_rooms_team").Create(data)
			if rInsert.Error != nil {
				logger.Error("JoinTeamDB no insert err:%v, data:%v", rInsert.Error, data)
				tx.Rollback()
				return rInsert.Error
			}

			if rInsert.RowsAffected != 1 {
				logger.Error("insert no RowsAffected err:%d, data:%v", rInsert.RowsAffected, data)
				tx.Rollback()
				return fmt.Errorf(" JoinTeamDB no insert RowsAffected err:%d , data:%v", rInsert.RowsAffected, data)
			}
			err := updateRoomPlayer(tx, uid, roomID, teamNum)
			if err != nil {
				logger.Error("updateRoomPlayer no err:%d, data:%v", rInsert.RowsAffected, data)
				tx.Rollback()
				return fmt.Errorf(" updateRoomPlayer no err:%d , data:%v", rInsert.RowsAffected, data)
			}
		} else {
			err := updateRoomPlayer(tx, uid, roomID, teamNum)
			if err != nil {
				logger.Error("updateRoomPlayer no err:%v, uid:%v, roomid:%v, teamname:%v", r.Error, uid, roomID, teamname)
				tx.Rollback()
				return fmt.Errorf(" updateRoomPlayer no err:%v, uid:%v, roomid:%v, teamname:%v", r.Error, uid, roomID, teamname)
			}
		}
	}
	logger.Debug("JoinTeamDB succ uid:%v, roomid:%v, teamname:%v", uid, roomID, teamname)
	return nil
}

func updateRoomPlayer(tx *gorm.DB, uid, roomID int64, teamNum int) error {

	r := tx.Table("pg_rooms_player").Where(" room_id = ? and uid = ?", roomID, uid).Update("team_num", teamNum)
	if r.RowsAffected != 1 {
		record := &model.RoomPlayerRecord{
			UID:          uid,
			GameNiceName: "",
			RoomID:       roomID,
			Prize:        0.00,
			TeamNum:      teamNum,
			KillNum:      0,
			Create:       0,
			Update:       0,
			Status:       0,
		}

		rInsert := tx.Table("pg_rooms_player").Create(record)
		if rInsert.Error != nil {
			logger.Error("updateRoomPlayer no insert err:%v, data:%v", rInsert.Error, record)
			tx.Rollback()
			return rInsert.Error
		}

		if rInsert.RowsAffected != 1 {
			logger.Error("updateRoomPlayer insert no RowsAffected err:%d, data:%v", rInsert.RowsAffected, record)
			tx.Rollback()
			return fmt.Errorf(" updateRoomPlayer no insert RowsAffected err:%d , data:%v", rInsert.RowsAffected, record)
		}

		if r.Error != nil {
			logger.Error("updateRoomPlayer err:%v, uid:%v, roomid:%v, teamNum:%v", r.Error, uid, roomID, teamNum)
			return r.Error
		}
	}
	return nil
}

type Team struct {
	Uid      int64  `gorm:"column:uid" json:"uid"`
	TeamNum  int    `gorm:"column:team_num" json:"team_num"`
	NickName string `gorm:"column:game_nick_name" json:"game_nick_name"`
}

func GetTeamDB(uid int64, roomid int) (model.TeamResp, error) {
	record := make([]Team, 0)
	teamData := model.TeamResp{}
	teamData.TeamUser = make(map[int][]model.TeamUserData, 0)
	r := pubgsql.Get("pubg").Master().Table("pg_rooms_player").Where("room_id = ? and status = 1 ", roomid).Order("team_num").Scan(&record)

	if r.RowsAffected == 0 {
		logger.Error("GetTeamDB RowsAffected uid:%v, roomid:%v", uid, roomid)
		return teamData, fmt.Errorf("GetTeamDB RowsAffected 0, uid:%v, roomid:%v", uid, roomid)
	} else if r.Error != nil {
		logger.Error("GetTeamDB RowsAffected 0, uid:%v, roomid:%v", r.Error, uid, roomid)
		return teamData, r.Error
	}

	type uid_nickname struct {
		TeamNum  int
		NickName string
	}
	userTeam := make(map[int64]uid_nickname)
	uids := make([]int64, 0)
	for _, v := range record {
		uids = append(uids, v.Uid)
		temp := uid_nickname{
			TeamNum:  v.TeamNum,
			NickName: v.NickName,
		}
		userTeam[v.Uid] = temp
	}
	userRecord := make([]model.PubgUserRecord, 0)
	logger.Debug("pg_rooms_player userTeam:%v ", userTeam)
	r = pubgsql.Get("pubg").Master().Table("pubg_user_info").Where("id IN (?)", uids).Scan(&userRecord)
	if r.RowsAffected == 0 {
		logger.Error("GetTeamDB pubg_user_info RowsAffected 0, uid:%v, roomid:%v", uid, roomid)
		return teamData, fmt.Errorf("GetTeamDB get user RowsAffected 0, uid:%v, roomid:%v", uid, roomid)
	} else if r.Error != nil {
		logger.Error("GetTeamDB get user err :%v, uid:%v, roomid:%v", r.Error, uid, roomid)
		return teamData, r.Error
	}

	//map[1078:{1 柚吸吸} 1086:{2 xxx} 1101:{1 呼叫}]
	for i, v := range userRecord {
		teamnum, ok := userTeam[v.ID]
		if !ok {
			continue
		}

		tempTeamUserData := model.TeamUserData{
			Uid:      v.ID,
			Photo:    v.Photo,
			Nickname: teamnum.NickName,
		}

		logger.Debug("tempTeamUserData  teamData begin i:%v, data:%v", i, teamData)
		_, ok = teamData.TeamUser[teamnum.TeamNum]
		if !ok {
			logger.Debug("tempTeamUserData  teamData mid i:%v, data:%v", i, teamData)
			teamData.TeamUser[teamnum.TeamNum] = []model.TeamUserData{tempTeamUserData}
		} else {
			teamData.TeamUser[teamnum.TeamNum] = append(teamData.TeamUser[teamnum.TeamNum], tempTeamUserData)
		}

		logger.Debug("tempTeamUserData  teamData i:%v, data:%v", i, teamData)
	}

	logger.Debug("GetTeamDB succ uid:%v, roomid:%v, teamData:%v", uid, roomid, teamData)
	return teamData, nil
}
