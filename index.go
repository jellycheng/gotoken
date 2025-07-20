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
