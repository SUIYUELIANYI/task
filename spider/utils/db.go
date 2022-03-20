package utils

import(
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"log"
)

var(
	Db *sql.DB
	err error
)

func init() {
	Db,err = sql.Open("mysql","root:123456@/student")//数据库类型；数据库用户名:数据库密码@[tcp(localhost:3306)]/数据库名
	if err !=nil{
		log.Println(err)
	}
	log.Println(Db)
}
