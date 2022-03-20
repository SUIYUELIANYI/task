package main

import (
	"Spider/model/get"
	//"fmt"

	_ "github.com/spf13/viper"
)

func main() {
	var id int

	id = 2021210001
	var ch1 chan string

	ch1 = make(chan string, 10000)
	// model.DB = model.Initdb() //好像这种方式连接失败

	cookie := "ASP.NET_SessionId=ltchuj55lz05wg35h4nhfc2k; _d_id=5bac9949032b9d0789855b1e7f0355" // 用于登录，最好是通过模拟登陆获取cookie,而不是在网页登录后直接复制过来，因为cookie是会过期的。
	// url1 := "http://kjyy.ccnu.edu.cn/clientweb/xcus/ic2/Default.aspx"                          // 图书馆登录后的url，如果没有携带cookie会自动跳转到身份验证界面
	for i := 0; i < 10; i++ { // 开10个协程，为了让主线程比协程后结束，我们要用channel,channel里的数据没有全部读出来之前主线程不会提前结束。
		go get.Getinfor(id+1000*i, id+1000*(i+1), cookie, ch1)
	}

	/* 	for i:= range ch1{
		fmt.Println(i)
	}  */
	var b int
	for {
		a := <-ch1
		if a == "1" {
			b++
			if b == 10 {
				break
			}
		}
	}
}

// [\u4E00-\u9FFF] 表示一个汉字
