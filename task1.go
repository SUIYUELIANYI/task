package main

import (
	"fmt"
	"net/http"
)

type User struct{//定义一个结构体存储用户，用户具有用户名、昵称等基本信息
	Username string//用户名
	Password string//密码
	Nickname string//昵称
	Gender string//性别
	Age string//年龄
} 

var num = make([]User, 0, 100)//make 的使用方式是：func make([]T, len, cap)，其中 cap 是可选参数。
var infor User

//注册
func register(w http.ResponseWriter, r *http.Request) {
	infor.Username = r.FormValue("username")
	infor.Password = r.FormValue("password")
	infor.Nickname = r.FormValue("nickname")
	infor.Gender = r.FormValue("gender")
	infor.Age = r.FormValue("age")
	var ok = true
	for _, values := range num {
		if values.Username == infor.Username {
			ok = false
		}
	}
	if ok {
		num = append(num, infor)
		cookies := http.Cookie{
			Name:     infor.Username,
			Value:    infor.Password,
			HttpOnly: true,
		}
		w.Header().Set("Set-Cookie", cookies.String())
		w.Write([]byte("注册成功~"))
	}
	if !ok {
		w.Write([]byte("用户已存在~"))
	}
}
//登录
func login(w http.ResponseWriter, r *http.Request){
	for _, value := range num {
		if value.Username == r.FormValue("username") && value.Password == r.FormValue("password") {
			cookies := http.Cookie{
				Name:     infor.Username,
				Value:    infor.Password,
				HttpOnly: true,
			}
			w.Header().Set("Set-Cookie", cookies.String())
			w.Write([]byte("登陆成功！"))
		}
		if value.Username == r.FormValue("username") && value.Password != r.FormValue("password") {
			w.Write([]byte("密码错误！"))
		}
	}
}


func show(p *User, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%+v", *p)
	fmt.Printf("\n")
}

//修改
func edit(w http.ResponseWriter, r *http.Request){
	var name = r.FormValue("username")
	var value = r.FormValue("password")
	cookie, err := r.Cookie(name)
	if cookie != nil {
		for i := 0; i < len(num); i++ {
			if name == num[i].Username && num[i].Password == value {
				w.Write([]byte("你的信息如下：\n"))
				show(&num[i], w, r)
				num[i].Password = r.FormValue("password")
				num[i].Nickname = r.FormValue("nickname")
				num[i].Age = r.FormValue("age")
				num[i].Gender = r.FormValue("gender")
				w.Write([]byte("修改后的信息如下：\n"))
				show(&num[i], w, r)
			}
		}
	} else {
		fmt.Fprintln(w, err)
	}
}
//查看
func check(w http.ResponseWriter, r *http.Request){
	var tmp = r.FormValue("username")
	cookie, err := r.Cookie(tmp)
	if cookie != nil {
		if cookie.Name == tmp && cookie.Value == r.FormValue("password") {
			for index, value := range num {
				fmt.Fprintf(w, "已注册的第%d的人的信息:%+v\n", index+1, value)
			}
		}
	} else {
		fmt.Fprintln(w, err)
	}
}


func main() {
	servemux := http.NewServeMux()
	servemux.HandleFunc("/register", register)//注册
	servemux.HandleFunc("/login",login)//登录
	servemux.HandleFunc("/login/edit",edit)//登录后修改
	servemux.HandleFunc("/login/check",check)//查看
	http.ListenAndServe(":4016", servemux)
}