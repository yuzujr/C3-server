package database_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/yuzujr/C3/internal/database"
)

func TestPrintUsersTableSchema(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	os.Setenv("ENV", "test")
	database.InitDatabase()
	db := database.DB

	tables, err := db.Migrator().GetTables()
	if err != nil {
		t.Fatalf("无法获取表列表: %v", err)
	}

	fmt.Println("数据库中已有表：", tables)

	// 检查关键表是否存在
	requiredTables := []string{"users", "clients", "command_logs", "screenshots"}
	for _, table := range requiredTables {
		found := false
		for _, tname := range tables {
			if tname == table {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("缺少关键表: %s", table)
		}
	}

	// 打印每个表的字段结构
	for _, table := range tables {
		columns, err := db.Migrator().ColumnTypes(table)
		if err != nil {
			t.Errorf("无法获取表 %s 的字段信息: %v", table, err)
			continue
		}
		fmt.Printf("表 %s 字段:\n", table)
		for _, col := range columns {
			fmt.Printf("  - %s (%s)\n", col.Name(), col.DatabaseTypeName())
		}
	}
}
