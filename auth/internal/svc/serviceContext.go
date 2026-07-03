// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"auth/internal/config"
)

type ServiceContext struct {
	Config config.Config
	// UserModel model.UserModel // 注入用户/车辆绑定相关的 Model
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. 初始化数据库连接
	// sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,
		// 2. 实例化 Model 并注入到上下文中
		// UserModel: model.NewUserModel(sqlConn, c.CacheRedis),
	}
}
