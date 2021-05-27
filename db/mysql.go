package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// 数据库配置结构
type DBConfig struct {
	Host        string // 地址
	Port        int    // 端口
	Name        string // 名称
	Pass        string // 密码
	DBName      string // 库名
	Charset     string // 编码
	Timezone    string // 时区
	MaxIdleConn int    // 最大空间连接
	MaxOpenConn int    // 最大连接数
}

//连接数据库
func InitDB(cfg DBConfig) (*sql.DB, error) {
	if strings.EqualFold(cfg.Timezone, "") {
		cfg.Timezone = "'Asia/Shanghai'"
	}
	if strings.EqualFold(cfg.Charset, "") {
		cfg.Charset = "utf8mb4"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&loc=Local",
		cfg.Name,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)
	connection, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = connection.Ping(); err != nil {
		return nil, err
	}
	connection.SetMaxIdleConns(cfg.MaxIdleConn)
	connection.SetMaxOpenConns(cfg.MaxOpenConn)
	return connection, nil
}

