package gotoken

type TokenFieldMapping struct {
	Id        string // 主键ID字段
	UserToken string // 用户token字段
	UserID    string // 用户ID字段
}

var TokenFieldCfg = TokenFieldMapping{
	Id:        "id",
	UserToken: "user_token",
	UserID:    "user_id",
}

type UserTokenDto struct {
	ID          string `json:"id"`           //主键ID
	UserID      string `json:"user_id"`      //用户ID
	UserToken   string `json:"user_token"`   //用户token
	SaasSeq     string `json:"saas_seq"`     //账套ID
	ExpireAt    string `json:"expire_at"`    //过期时间戳
	ActiveTime  string `json:"active_time"`  //最后活跃时间
	AppPlatform string `json:"app_platform"` //平台 p-PC i-IOS a-Android h5-H5 mp-小程序
	AppType     string `json:"app_type"`     //app类型，产品类型
	DeviceID    string `json:"device_id"`    //设备ID
	InvalidType string `json:"invalid_type"` //失效类型，0-未失效，1-自动失效 2-相同平台踢下线失效，3-用户主动退出,4-人工处理退出
	OutSystem   string `json:"out_system"`   //外部系统代号
	IP          string `json:"ip"`           //登录IP
	IsDelete    string `json:"is_delete"`    //是否删除，0-正常，1-已删除
	CreateTime  string `json:"create_time"`  //创建时间
	UpdateTime  string `json:"update_time"`  //更新时间
	DeleteTime  string `json:"delete_time"`  //删除时间
	ModifyTime  string `json:"modify_time"`  //修改时间
}
