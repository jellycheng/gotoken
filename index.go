package gotoken

import (
	"database/sql"
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/jellycheng/gosupport/dbutils"
	"time"
)

func GetTokenInfo(connect *sql.DB, token string, tbl string) (map[string]string, error) {
	if tbl == "" {
		tbl = DefaultTokenTable
	}
	sqlStr := fmt.Sprintf("select * from %s where `%s`=? and is_delete=0 limit 1;", tbl, TokenFieldCfg.UserToken)
	ret, err := dbutils.SelectOne(connect, sqlStr, token)
	return ret, err
}

func GetTokenInfo2Dto(connect *sql.DB, token string, tbl string) (UserTokenDto, error) {
	ret := UserTokenDto{}
	if tokenInfo, err := GetTokenInfo(connect, token, tbl); err != nil {
		return ret, err
	} else {
		tmpJson := gosupport.ToJson(tokenInfo)
		_ = gosupport.JsonUnmarshal(tmpJson, &ret)
	}
	return ret, nil
}

func GetTokenInfo4Redis(myRedis *MyRedisClient, connect *sql.DB, token string, tbl string) (map[string]string, error) {
	ret := make(map[string]string)
	var e error
	var tmpRedisKey = fmt.Sprintf(RedisTokenKeyFormat, token)
	userTokenInfo4redis := GetKeyValue(myRedis, tmpRedisKey)
	if userTokenInfo4redis == "" {
		ret, e = GetTokenInfo(connect, token, tbl)
		if e == nil && ret[TokenFieldCfg.Id] != "" {
			_ = SetKeyValue(myRedis, tmpRedisKey, gosupport.ToJson(ret), 86400*time.Second)
		}
	} else {
		_ = gosupport.JsonUnmarshal(userTokenInfo4redis, &ret)
	}
	return ret, e
}

func GetTokenInfo2Dto4Redis(myRedis *MyRedisClient, connect *sql.DB, token string, tbl string) (UserTokenDto, error) {
	ret := UserTokenDto{}
	if tokenInfo, err := GetTokenInfo4Redis(myRedis, connect, token, tbl); err != nil {
		return ret, err
	} else {
		tmpJson := gosupport.ToJson(tokenInfo)
		_ = gosupport.JsonUnmarshal(tmpJson, &ret)
	}
	return ret, nil
}

func CleanTokenCache(myRedis *MyRedisClient, token string) {
	var tmpFormatKey = fmt.Sprintf(RedisTokenKeyFormat, token)
	tmpRedisKey := fmt.Sprintf("%s%s", myRedis.cfg.Prefix, tmpFormatKey)
	_ = myRedis.rdb.Del(ctx, tmpRedisKey)
}

// 删除token
func DelToken(connect *sql.DB, token string, tbl string) (int64, error) {
	if tbl == "" {
		tbl = DefaultTokenTable
	}
	curTime := gosupport.TimeNow()
	sqlStr := fmt.Sprintf("update %s set is_delete=1,delete_time=?,update_time=? where `%s`=? and is_delete=0 limit 1;", tbl, TokenFieldCfg.UserToken)
	return dbutils.UpdateSql(connect, sqlStr, curTime, curTime, token)
}

// 新增token记录
func AddToken(connect *sql.DB, data map[string]interface{}, tbl string) (int64, error) {
	if tbl == "" {
		tbl = DefaultTokenTable
	}
	var lastid int64 = 0
	if len(data) == 0 {
		return lastid, fmt.Errorf("data is empty")
	}
	fileds := make([]string, 0)
	values := make([]interface{}, 0)
	for k, v := range data {
		fileds = append(fileds, k)
		values = append(values, v)
	}
	sqlObj := dbutils.NewSQLBuilderInsert().SetTable(tbl).SetInsertData(fileds, values...)
	sql, _ := sqlObj.GetSql()
	lastid, _ = dbutils.InsertSql(connect, sql, sqlObj.GetParamValues()...)
	return lastid, nil
}

// 更新token记录
func UpdateToken(connect *sql.DB, token string, data map[string]interface{}, tbl string) (int64, error) {
	if tbl == "" {
		tbl = DefaultTokenTable
	}
	var affectedNum int64 = 0
	if token == "" {
		return affectedNum, fmt.Errorf("token is empty")
	}
	if len(data) == 0 {
		return affectedNum, fmt.Errorf("data is empty")
	}
	fileds := make([]string, 0)
	values := make([]interface{}, 0)
	for k, v := range data {
		fileds = append(fileds, k)
		values = append(values, v)
	}
	sqlObj := dbutils.NewSQLBuilderUpdate().SetTable(tbl).SetUpdateData(fileds, values...).Where(TokenFieldCfg.UserToken, "=", token).Where(TokenFieldCfg.IsDelete, "=", 0)
	sql, _ := sqlObj.GetSQL()
	affectedNum, _ = dbutils.UpdateSql(connect, sql, sqlObj.GetParamValues()...)
	return affectedNum, nil
}
