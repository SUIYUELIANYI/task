package init

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//DB 全局变量,这两个在func里第一次用
var DB *gorm.DB
var err error

func Initdb() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", viper.Get("mysql.user"), viper.Get("mysql.password"), viper.Get("mysql.host"), viper.Get("mysql.port"), viper.Get("mysql.db"))
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) //连接数据库的固定格式
	if err != nil {
		log.Fatal(err) //1.打印输s出内容 2.退出应用程序 3.defer函数不会执行
	}
	// 注意：port的端口号和main函数的端口号不能相同，不然会按占用处理
	return DB
}
