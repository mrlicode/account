package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/grpcbrick/account/dao"
	"github.com/yinxulai/goutils/config"
	"github.com/yinxulai/goutils/sqldb"
)

func TestMain(m *testing.M) {
	fmt.Println("准备测试环境")                                                                       // 日志
	config.Set("encrypt-password", "encrypt-password")                                          // 密码加密中会用到的
	sqldb.Init("mysql", "root:root@tcp(localhost:3306)/default?charset=utf8mb4&parseTime=true") // 测试数据库
	dao.InitTables()                                                                            // 初始化测试数据库
	fmt.Println("开始执行测试")                                                                       // 日志
	exitCode := m.Run()                                                                         // 执行测试
	fmt.Println("测试执行完成,清理测试数据")                                                                // 日志
	dao.TruncateTables()                                                                        // 重置测试数据库
	os.Exit(exitCode)                                                                           // 推出
}
