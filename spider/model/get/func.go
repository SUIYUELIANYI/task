package get

import (
	"Spider/model/types"
	"Spider/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

// 将学生信息存入数据库
func CreatUser(studentId string, name string, grade string) error {
	sqlStr := "insert into users(student_id,name,grade)values(?,?,?)"
	_, err := utils.Db.Exec(sqlStr, studentId, name, grade)
	if err != nil {
		fmt.Println("error", err)
		return err
	} else {
		return nil
	}
}

type Form struct {
	Execution string `json:"execution"` // n.执行，实行
	Lt        string `json:"lt"`
	EventId   string `json:"_eventId"`
	Submit    string `json:"submit"` // vt.提交
}

// 爬取数据
func Getinfor(id1 int, id2 int,cookie string,ch1 chan string) {
	var p types.User
	for ;id1<id2;id1 ++{
	url2 := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx?type=logonname&ReservaApply=ReservaApply&term=" + strconv.Itoa(id1) + "&_="
	client := &http.Client{}

	req, err := http.NewRequest("GET", url2, nil)

	if err != nil {
		fmt.Println("错误!")
	}

	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("登录错误!")
	}

	respByte, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	respString := string(respByte)
	//fmt.Println(respString)
	if respString == "[]" {
		continue
		//return
	}

	//fmt.Println(respString)
	//1)解释规则，他会解析正则表达式，如果成功返回解释器
	reg1 := regexp.MustCompile("\"name\": \".+?\"")
	if reg1 == nil { //解释失败，返回nil
		fmt.Println("err = ")
		return
	}

	//2)根据规则提取关键信息,-1表示匹配所有项，正数为匹配项数
	//result := reg1.FindAllStringSubmatch(respString, -1)
	name_reg := regexp.MustCompile(`"name": "(.*?)"`)
	re := name_reg.FindAllStringSubmatch(respString, -1)[0][1] // 这里是yyj写的匹配姓名

	//fmt.Println(re)
	//fmt.Println(result[0])

	p.Name = re                    //获取学生姓名
	p.StudentID = strconv.Itoa(id1) //学号
	s := p.StudentID[0:4]
	p.Grade = s //年级对应学号的前四位
	//fmt.Println(p)
	err1 := CreatUser(p.StudentID, p.Name, p.Grade)
	if err1 != nil {
		fmt.Println(err1)
	} 
	}
	ch1 <- "1"
	fmt.Println("协程已完成")
}
