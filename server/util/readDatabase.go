package util

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// ReadDatabaseAgentData 数据库读取agent数据
func ReadDatabaseAgentData(dbParameter, dbIp, dbUser, dbPass, dbName string, dbPort int) map[int][]string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", dbUser, dbPass, dbIp, dbPort, dbName, dbParameter)
	// Agent 数据库表模型
	type Agent struct {
		ID      int    `db:"id"`
		Status  int    `db:"status"`
		Name    string `db:"name"`
		Mount   string `db:"partition"`
		AuthKey string `db:"key"`
	}
	// 打开数据库连接（不配置日志）
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[Database]连接失败: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Println(err)
		}
	}()
	// 执行查询（SQLx默认不输出SQL日志，仅在出错时返回错误）
	var agents []Agent
	if err = db.Select(&agents, "SELECT * FROM vigilix_agents WHERE status = 1"); err != nil {
		log.Fatalf("[Database] agents 查询失败: %v", err)
	}
	// 转换为map
	result := make(map[int][]string)
	for _, agent := range agents {
		result[agent.ID] = []string{agent.Name, agent.Mount, agent.AuthKey}
	}
	log.Printf("[Database] agents 读取成功: %v", result)
	return result
}

// ReadDatabaseClientData 数据库读取client数据
func ReadDatabaseClientData(dbParameter, dbIp, dbUser, dbPass, dbName string, dbPort int) map[string]string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", dbUser, dbPass, dbIp, dbPort, dbName, dbParameter)
	// Agent 数据库表模型
	type client struct {
		ID     int    `db:"id"`
		Status int    `db:"status"`
		User   string `db:"user"`
		Pass   string `db:"pass"`
	}
	// 打开数据库连接（不配置日志）
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[Database]连接失败: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Println(err)
		}
	}()
	// 执行查询（SQLx默认不输出SQL日志，仅在出错时返回错误）
	var clients []client
	if err = db.Select(&clients, "SELECT * FROM vigilix_clients WHERE status = 1"); err != nil {
		log.Fatalf("[Database] clients 查询失败: %v", err)
	}
	// 转换为map
	result := make(map[string]string)
	for _, clientConfig := range clients {
		result[clientConfig.User] = clientConfig.Pass
	}
	log.Printf("[Database] clients 读取成功: %v", result)
	return result
}
