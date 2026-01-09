# gotoken
```
go实现登录态token存储，有效性判断等

```

## t_user_token_* 表结构参考
```
DROP TABLE IF EXISTS `t_user_token_1`;
CREATE TABLE IF NOT EXISTS `t_user_token_1` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID,自增主键',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `saas_seq` varchar(32) NOT NULL DEFAULT '' COMMENT '账套ID',
  `user_token` varchar(50) NOT NULL DEFAULT '' COMMENT '登录态',
  `expire_at` int(10) NOT NULL DEFAULT '0' COMMENT '登录态到期时间,时间戳 0 永不过期',
  `device_id` varchar(50) NOT NULL DEFAULT '' COMMENT '设备唯一ID',
  `active_time` int(10) NOT NULL DEFAULT '0' COMMENT '最后活跃时间',
  `app_platform` varchar(25) NOT NULL DEFAULT '' COMMENT '平台 p-PC i-IOS a-Android h5-H5 mp-小程序',
  `app_type` varchar(50) NOT NULL DEFAULT 'mqj' COMMENT 'APP类型',
  `ip` varchar(32) NOT NULL DEFAULT '' COMMENT 'login ip',
  `out_system` varchar(255) NOT NULL DEFAULT '' COMMENT '外部系统调用生成token，可选',
  `invalid_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '失效类型，0-未失效，1-自动失效 2-相同平台踢下线失效，3-用户主动退出,4-人工处理退出',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除 0-正常 1-删除',
  `create_time` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) NOT NULL DEFAULT '0' COMMENT '修改时间',
  `delete_time` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
  `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'mysql更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unq_token` (`user_token`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_modify_time` (`modify_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户登录态表' ROW_FORMAT=COMPACT;

```

## 新增token记录
```
package main

import (
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/jellycheng/gosupport/dbutils"
	"github.com/jellycheng/gotoken"
)

func main() {
	dsn := dbutils.GetDsn(map[string]interface{}{
		"host":     "localhost",
		"port":     "3306",
		"dbname":   "db_mall",
		"username": "root",
		"password": "88888888",
	})
	// 或者 dsn := "root:88888888@tcp(localhost:3306)/db_mall?charset=utf8"
	con, err := dbutils.GetDbConnect("db_mall", dsn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data := map[string]interface{}{
		"user_token":  "token-xxx123",
		"user_id":     789,
		"create_time": gosupport.TimeNow(),
	}
	if lastId, err := gotoken.AddToken(con, data, "t_user_token_1"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(lastId)
	}

}

```

## 查询token信息
```
package main

import (
	"fmt"
	"github.com/jellycheng/gosupport/dbutils"
	"github.com/jellycheng/gotoken"
)

func main() {
	dsn := dbutils.GetDsn(map[string]interface{}{
		"host":     "localhost",
		"port":     "3306",
		"dbname":   "db_mall",
		"username": "root",
		"password": "88888888",
	})
	con, err := dbutils.GetDbConnect("db_mall", dsn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 查询token信息
	token := "4cbca931c92015902b7483877e42f5d3"
	if tokenInfo, err := gotoken.GetTokenInfo2Dto(con, token, "t_user_token_1"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(fmt.Sprintf("tokenInfo: %+v", tokenInfo))
		fmt.Println(tokenInfo.UserID)
	}

}

```

## 通过缓存查询token信息
```
package main

import (
	"fmt"
	"github.com/jellycheng/gosupport/dbutils"
	"github.com/jellycheng/gotoken"
)

func main() {
	dsn := dbutils.GetDsn(map[string]interface{}{
		"host":     "localhost",
		"port":     "3306",
		"dbname":   "db_mall",
		"username": "root",
		"password": "88888888",
	})
	con, err := dbutils.GetDbConnect("db_mall", dsn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	redisCfg := gotoken.RedisCfg{
		Host:     "127.0.0.1",
		Port:     "6379",
		Password: "", // redis密码
		Prefix:   "token:",         //key前缀
	}

	rdb := gotoken.NewRedisClient(redisCfg)
	// 查询token信息
	token := "4cbca931c92015902b7483877e42f5d3"
	if tokenInfo, err := gotoken.GetTokenInfo2Dto4Redis(rdb, con, token, "t_user_token_1"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(fmt.Sprintf("tokenInfo: %+v", tokenInfo))
		fmt.Println(tokenInfo.UserToken)
	}

}

```

## 删除token
```
package main

import (
	"fmt"
	"github.com/jellycheng/gosupport/dbutils"
	"github.com/jellycheng/gotoken"
)

func main() {
	dsn := dbutils.GetDsn(map[string]interface{}{
		"host":     "localhost",
		"port":     "3306",
		"dbname":   "db_mall",
		"username": "root",
		"password": "88888888",
	})
	con, err := dbutils.GetDbConnect("db_mall", dsn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	token := "4cbca931c92015902b7483877e42f5d3"
	if affectedNum, err := gotoken.DelToken(con, token, "t_user_token_1"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(affectedNum)
	}

}

```

## 更新token记录
```
package main

import (
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/jellycheng/gosupport/dbutils"
	"github.com/jellycheng/gotoken"
)

func main() {
	dsn := dbutils.GetDsn(map[string]interface{}{
		"host":     "localhost",
		"port":     "3306",
		"dbname":   "db_mall",
		"username": "root",
		"password": "88888888",
	})
	con, err := dbutils.GetDbConnect("db_mall", dsn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	token := "token-xxx123"
	data := map[string]interface{}{
		"invalid_type": 1,
		"active_time":  gosupport.TimeNow(),
		"update_time":  gosupport.TimeNow(),
	}
	if affectedNum, err := gotoken.UpdateToken(con, token, data, "t_user_token_1"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(affectedNum)
	}

}

```
