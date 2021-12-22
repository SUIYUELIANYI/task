package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //数据库驱动，因为使用sql包时必须注入至少一个数据库驱动
	"log"
	"net/http"
	"strconv" //提供将字符串转化为基础类型的功能
)

var (
	Db  *sql.DB
	err error //定义err为error类型，不同于bool类型的true或者false
)

func init() { //init 初始化
	Db, err = sql.Open("mysql", "root:123456@/ManageSystem") //数据库类型；数据库用户名:数据库密码@[tcp(localhost:3306)]/数据库名，@后面的在本地运行可以省略
	if err != nil {
		log.Println(err)
	}
	log.Println(Db)
}

type User struct { //定义一个结构体存储用户，用户具有用户名、昵称等基本信息
	user_name     string //用户名
	user_password string //密码
	user_role     int    //用户职位
	user_gender   string //用户性别
	user_age      int    //用户年龄
}

var num = make([]User, 0, 100) //make 的使用方式是：func make([]T, len, cap)，其中 cap 是可选参数。
var infor User                 //设置infor为结构体类型

//实现输入的用户信息插入到数据库的功能
func CreatUser(name string, pwd string, role int, gender string, age int) error { //CreatUser为创建用户函数，其返回值为error类型
	sqlStr := "insert into user(user_name,user_password,user_role,user_gender,user_age)values(?,?,?,?,?)" //这里的"user_name,user_password,user_role"均为建表时的结构
	_, err := Db.Exec(sqlStr, name, pwd, role, gender, age)                                               //Exec里面为函数调用的五个值个值name，pwd，role，gender,age
	if err != nil {
		fmt.Println("error", err)
		return err
	} else {
		return nil
	}
}

//检查用户名是否存在功能
/*func IfExistUser(name string) bool {
	var user = make([]User, 1)
	if err := Db.Table("users").Where("user_name=?", name).Find(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	if len(user) != 1 {
		return true
	}
	return false
}*/

//注册功能
func Register(w http.ResponseWriter, r *http.Request) {
	infor.user_name = r.FormValue("username") //infor.Username是前面结构体中定义的，后面的是postman里表单填的
	infor.user_password = r.FormValue("password")
	infor.user_role, err = strconv.Atoi(r.FormValue("role")) //由于formvalue返回的是string类型，我们要通过strconv将其转化为int类型
	infor.user_gender = r.FormValue("gender")
	infor.user_age, err = strconv.Atoi(r.FormValue("age")) //我们需要注意Atoi返回两个参数，不能省略err

	var ok = true //定义ok为true类型
	for _, values := range num {
		if values.user_name == infor.user_name {
			ok = false //如果该用户名已经存在，则ok为false
		}
	}
	if ok {
		num = append(num, infor)
		/*cookies := http.Cookie{
			Name:     infor.user_name,
			Value:    infor.user_password,
			HttpOnly: true,
		}*/
		//w.Header().Set("Set-Cookie", cookies.String()) //注册应该不用设置cookie
		w.Write([]byte("注册成功~"))
	}
	if !ok {
		w.Write([]byte("用户已存在~"))
	}
	err := CreatUser(infor.user_name, infor.user_password, infor.user_role, infor.user_gender, infor.user_age)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, "success")
	} //发现如果用这段代码，在第一次注册账号后立即用相同的用户名注册确实会反馈用户已存在，但是如果之前注册账号即数据库里有数据，再注册相同的账号却显示"注册成功".因为我将user_name设置为唯一索引所以数据库里并不会添加数据，但我希望不是通过此方式
	//所以最好通过一个判断用户是否存在的函数，不过上述功能还算齐全，我这里先不修改了
}

//通过查询数据库来判断用户是否存在，以及用户名和密码是否正确
func IsTrue(name string, pwd string) (string, string, error) {
	var username string
	var password string
	var err error
	sqlStr := "select user_name,user_password from user where user_name = ? and user_password = ?" //问号是表单里填写的东西，等下在执行语句中用name和pwd代入
	row, _ := Db.Query(sqlStr, name, pwd)                                                          //同理Query里面填的是语句sqlStr和IsTrue调用的两个值
	for row.Next() {
		err = row.Scan(&username, &password)
	}
	if err != nil {
		return "", "", err
	}
	return username, password, nil
}

/*func login(w http.ResponseWriter, r *http.Request){
	for _, value := range num {
		if value.user_name == r.FormValue("username") && value.user_password == r.FormValue("password") {
			cookies := http.Cookie{
				Name:     infor.user_name,
				Value:    infor.user_password,
				HttpOnly: true,
			}
			w.Header().Set("Set-Cookie", cookies.String())
			w.Write([]byte("登陆成功！"))
		}
		if value.user_name == r.FormValue("username") && value.user_password != r.FormValue("password") {
			w.Write([]byte("密码错误！"))
		}
	}
}*/

//登录功能
func Login(w http.ResponseWriter, r *http.Request) {
	Uname := r.FormValue("username")
	Upwd := r.FormValue("password")
	name, _, _ := IsTrue(Uname, Upwd)
	var id string
	//根据用户名从数据库中搜索id
	sql := "SELECT user_id FROM user where user_name = ?"
	rows, err := Db.Query(sql, Uname)

	if err != nil {
		fmt.Println("查找失败", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id) //这里应该是把读取到的数据存入id中
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println(id)
	if name == "" {
		fmt.Fprintln(w, "fail")
		//fmt.Println(err)
	} else {
		//fmt.Fprintln(w,"success")
		//fmt.Println(name)
		//尝试在这里添加cookie
		cookies := http.Cookie{
			Name:  "cookie", //因为Value不能为中文字符，打算用id来代替
			Value: id,       //今天发现在实例化一个cookie变量时，value中不能含有中文字符，否则添加会失败且不会报错！！！
		}
		w.Header().Set("Set-Cookie", cookies.String())
		w.Write([]byte("登录成功!"))
	}
}

//注册球赛功能
//var Team []string

type Match struct {
	name        string
	date        string
	place       string
	info        string // 详情
	appointment int    // 预约数
	teamA       string
	teamB       string
}

var game Match //定义game为结构体类型

//插入创建的球赛到数据库的函数,注意由于match是sql的保留字,我们需要用反引号
func InsertMatch(name string, date string, place string, info string, appointment int, teama string, teamb string) error { //CreatUser为创建用户函数，其返回值为error类型
	sqlStr := "insert into `match`(match_name,match_date,match_place,match_info,match_appointment,match_teama,match_teamb) values(?,?,?,?,?,?,?)" //这里的"user_name,user_password,user_role"均为建表时的结构
	_, err := Db.Exec(sqlStr, name, date, place, info, appointment, teama, teamb)                                                                 //Exec里面为函数调用的五个值个值name，pwd，role，gender,age
	if err != nil {
		fmt.Println("error", err)
		return err
	} else {
		return nil
	}
}

func CreateMatch(w http.ResponseWriter, r *http.Request) {
	//我们需要根据cookie中存储的id去判断登录者的权限
	/*cookie, err := r.Cookie("cookie")
	if err != nil {
		log.Println("获取cookie失败!", err)
	}
	id := cookie.Value*/

	//我们可以用函数来实现获取Cookie
	id := getcookie(w, r)
	sql := "SELECT user_role role FROM user where user_id = ?"
	rows, err := Db.Query(sql, id)

	if err != nil {
		fmt.Println("查找失败", err)
	}
	defer rows.Close()

	var role string
	for rows.Next() {
		err = rows.Scan(&role) //这里应该是把读取到的数据存入role中
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if role == "1" {
		w.Write([]byte("用户权限不够"))
	} else {
		game.name = r.FormValue("name") //game.name是前面结构体中定义的，后面的是postman里表单填的
		game.date = r.FormValue("date")
		game.place = r.FormValue("place")
		game.info = r.FormValue("information")
		game.appointment, err = strconv.Atoi(r.FormValue("appointment")) //由于formvalue返回的是string类型，我们要通过strconv将其转化为int类型
		game.teamA = r.FormValue("teama")
		game.teamB = r.FormValue("teamb")
		err = InsertMatch(game.name, game.date, game.place, game.info, game.appointment, game.teamA, game.teamB)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Fprintln(w, "成功创建球赛")
		}
	}
}

func GrantAuthority(w http.ResponseWriter, r *http.Request) { //Grant(授予)Authority(权力)
	//这里修改权限时依然要进行身份验证，和创建球赛时一样
	//我们可以用函数来实现获取Cookie
	id := getcookie(w, r)
	sql := "SELECT user_role role FROM user where user_id = ?"
	rows, err := Db.Query(sql, id)

	if err != nil {
		fmt.Println("查找失败", err)
	}
	defer rows.Close()

	var role string
	for rows.Next() {
		err = rows.Scan(&role) //这里应该是把读取到的数据即权限存入role中
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if role == "1" || role == "2" {
		w.Write([]byte("当前用户权限不够，无法赋予其他用户权限"))
	} else {
		username := r.FormValue("username")
		changerole := r.FormValue(("changerole"))
		sql = "UPDATE user SET user_role= ? where user_name= ?"
		_, err := Db.Exec(sql, changerole, username) //Exec里面为函数调用的五个值个值name，pwd，role，gender,age
		if err != nil {
			fmt.Println("update data fail,err", err)
			return
		} else {
			w.Write([]byte("成功赋予权限"))
		}
	}
}

//获取球赛列表，10条记录一页，默认按时间顺序
func GetInfor1(w http.ResponseWriter, r *http.Request) {
	//执行查询操作
	rows, err := Db.Query("SELECT match_name,match_date,match_place,match_info,match_appointment,match_teama,match_teamb FROM `match` ORDER BY match_date DESC LIMIT 10;") //每次取10条记录,按照时间顺序降序即最新的时间在最前面,注意分页放最后
	if err != nil {
		fmt.Println("select Db failed,err:", err)
		return
	}
	// rows.Next(),用于循环获取所有
	for rows.Next() {
		err = rows.Scan(&game.name, &game.date, &game.place, &game.info, &game.appointment, &game.teamA, &game.teamB)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "比赛名:%s时间:%s场地:%s球赛详情:%s预约数:%d 比赛队伍:%s%s \n", game.name, game.date, game.place, game.info, game.appointment, game.teamA, game.teamB)
	}
	rows.Close()
}

//获取球赛列表，10条记录一页，按热度排序
func GetInfor2(w http.ResponseWriter, r *http.Request) {
	//执行查询操作
	rows, err := Db.Query("SELECT match_name,match_date,match_place,match_info,match_appointment,match_teama,match_teamb FROM `match` ORDER BY match_appointment DESC LIMIT 10;") //每次取10条记录,按照时间顺序降序即最新的时间在最前面,注意分页放最后
	if err != nil {
		fmt.Println("select Db failed,err:", err)
		return
	}
	// rows.Next(),用于循环获取所有
	for rows.Next() {
		err = rows.Scan(&game.name, &game.date, &game.place, &game.info, &game.appointment, &game.teamA, &game.teamB)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "比赛名: %s  时间: %s  场地: %s  球赛详情: %s  预约数: %d  比赛队伍: %s %s \n", game.name, game.date, game.place, game.info, game.appointment, game.teamA, game.teamB)
	}
	rows.Close()
}

//用户预约操作
func InsertAppointment(username string, matchname string) error {
	var appointment int //用来存储查询到的比赛的预约数
	sql := "INSERT INTO usertomatch(user_name,match_name) VALUES(?,?)"
	_, err := Db.Exec(sql, username, matchname)
	if err != nil {
		fmt.Println("预约失败", err)
		return err
	} else {
		//预约成功，对应的预约人数增加一
		sql1 := "SELECT match_appointment FROM `match` where match_name = ?"
		rows, err1 := Db.Query(sql1, matchname)
		if err1 != nil {
			fmt.Println("查找比赛失败", err1)
		}
		
		for rows.Next() {
			err1 = rows.Scan(&appointment)
			if err1 != nil {
				fmt.Println("读取预约数失败", err)
			}
		}
		rows.Close()
		appointment ++ //预约数加一
		sql2 := "UPDATE `match` SET match_appointment = ? where match_name= ?"
		_, err = Db.Exec(sql2, appointment, matchname)
		if err != nil{
			fmt.Println("预约数没有改变",err)
		}
		fmt.Println("预约成功")
		return nil
	}
}

func MakeAnAppointment(w http.ResponseWriter, r *http.Request) {
	id := getcookie(w, r)
	sql1 := "SELECT user_role role FROM user where user_id = ?"
	sql2 := "SELECT user_name FROM user where user_id = ?" //同时找到用户的名称
	rows1, err1 := Db.Query(sql1, id)
	if err1 != nil {
		fmt.Println("查找失败", err)
	}
	defer rows1.Close()

	var role string
	for rows1.Next() {
		err = rows1.Scan(&role) //把读取到的数据即权限存入role中
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	rows2, err2 := Db.Query(sql2, id)
	if err2 != nil {
		fmt.Println("查找失败", err)
	}
	defer rows2.Close()

	var username string
	for rows2.Next() {
		err = rows2.Scan(&username) //把读取到的数据即权限存入role中
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if role == "1" { //这是为普通用户设置的
		game.name = r.FormValue("matchname")
		err := InsertAppointment(username, game.name)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Fprintln(w, "成功预约球赛!")
		}
	}
}

//查询某个用户预约情况
func GetAppointmentInfor(w http.ResponseWriter, r *http.Request){
	var matchname string
	username := r.FormValue("Name")
	sql :="SELECT match_name FROM usertomatch WHERE user_name=?"
	rows, err :=Db.Query(sql,username)
	if err != nil{
		fmt.Println("查找预约情况失败",err)
	}
	defer rows.Close()
	for rows.Next(){
		err = rows.Scan(&matchname)
		if err !=nil {
			fmt.Println("读取预约比赛信息失败",err)
			return
		}
		fmt.Fprintf(w, "预约的比赛名称: %s  \n", matchname)
	}
}

//输入球队信息
type Team struct {
	team_name        string
	team_logo        string
	team_member      string
}
var team Team

func InsertTeam(name string, logo string, member string) error {
	sqlStr := "INSERT INTO team(team_name,team_logo,team_member)"
	_, err := Db.Exec(sqlStr, name, logo, member) //Exec里面为函数调用的五个值个值name，pwd，role，gender,age
	if err != nil {
		fmt.Println(" insert team error", err)
		return err
	} else {
		return nil
	}
}

func RegisterTeam(w http.ResponseWriter, r *http.Request) {
	fmt.Println("录入球队信息")
	team.team_name = r.FormValue("teamname")
	team.team_logo = r.FormValue("teamlogo")
	team.team_member = r.FormValue("teammember")
	err = InsertTeam(team.team_name, team.team_logo, team.team_member)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, "成功登记球队")
	}
}

//获取团队列表
func GetTeamInfor(w http.ResponseWriter, r *http.Request) {
	//执行查询操作
	rows, err := Db.Query("SELECT team_name,team_logo,team_member FROM `team`;")
	if err != nil {
		fmt.Println("select Db failed,err:", err)
		return
	}
	// rows.Next(),用于循环获取所有
	for rows.Next() {
		err = rows.Scan(&team.team_name, &team.team_logo, &team.team_member)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "团队名称: %s  团队Logo: %s  团队成员: %s \n", team.team_name, team.team_logo, team.team_member)
	}
	rows.Close()
}

//输入球员信息
type Player struct {
	player_name        string
	player_team        string
	player_avatar      string
	player_infor		string
}
var player Player

func InsertPlayer(name string, team string, avatar string,infor string) error {
	sqlStr := "INSERT INTO player(player_name,player_team,player_avatar,player_information)"
	sqlStr1 :="INSERT INTO playertoteam(player_name,team_name)"
	_, err := Db.Exec(sqlStr, name, team, avatar, infor )
	if err != nil {
		fmt.Println(" insert player error", err)
		return err
	} else {
		_ , err = Db.Exec(sqlStr1, name, team ) 
		if err != nil {
			fmt.Println(" insert playertoteam error", err)
			return err
		} else {
			return nil
		}
	}
}

func RegisterPlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("录入球员信息")
	playername := r.FormValue("playername")
	playerteam := r.FormValue("playerteam")
	playeravatar := r.FormValue("playeravatar")
	playerinfor := r.FormValue("playinfor")
	err = InsertPlayer(playername, playerteam, playeravatar, playerinfor)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, "成功登记球员")
	}
}

//获取球员列表
func GetPlayerInfor(w http.ResponseWriter, r *http.Request) {
	//执行查询操作
	rows, err := Db.Query("SELECT player_name,player_team,player_avatar,player_information FROM `player`;")
	if err != nil {
		fmt.Println("select Db failed,err:", err)
		return
	}
	// rows.Next(),用于循环获取所有
	for rows.Next() {
		err = rows.Scan(&player.player_name, &player.player_team, &player.player_avatar, &player.player_infor)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "球员名称: %s  所在团队: %s  头像: %s 信息:%s \n",player.player_name, player.player_team, player.player_avatar, player.player_infor)
	}
	rows.Close()
}

func ModifyRelation(w http.ResponseWriter, r *http.Request){
	//读入球员名称和新的球队名称
	playername := r.FormValue("playername")
	teamname := r.FormValue("modifyname")
	sqlStr1 :="INSERT INTO playertoteam(player_name,team_name)"
	_, err = Db.Exec(sqlStr1, playername, teamname ) 
	if err != nil {
		fmt.Println(" insert playertoteam error", err)
		return
	} else {
		fmt.Fprintf(w ,"更改球员所在球队成功！")
	}
}


//获取cookie中的信息进行身份识别
func getcookie(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("cookie")
	if err != nil {
		log.Println("获取cookie失败!", err)
	}
	return cookie.Value
}

func main() {
	servemux := http.NewServeMux()
	servemux.HandleFunc("/register", Register)                         //1注册用户
	servemux.HandleFunc("/login", Login)                               //2登录
	servemux.HandleFunc("/login/grantauthority", GrantAuthority)       //3用户权限
	servemux.HandleFunc("/login/getinfortime", GetInfor1)              //4获取球赛列表(按时间)
	servemux.HandleFunc("/login/getinforhot", GetInfor2)               //4获取球赛列表(按热度)
	servemux.HandleFunc("/login/creatematch", CreateMatch)             //5创建单个球赛
	servemux.HandleFunc("/login/makeanappointment", MakeAnAppointment) //6预约球赛
	servemux.HandleFunc("/getappointmentinfor",GetAppointmentInfor)	   //7查看单个用户预约的所有球赛
	servemux.HandleFunc("/registerplayer", RegisterPlayer)             //8登记球员
	servemux.HandleFunc("/registerteam", RegisterTeam)                 //9登记团队
	servemux.HandleFunc("/modifyrelation",ModifyRelation)			   //10修改团队和球员关系
	servemux.HandleFunc("/getteaminfor",GetTeamInfor)				   //11获取团队列表
	servemux.HandleFunc("/getplayerinfor",GetPlayerInfor)			   //12获取球员列表
	http.ListenAndServe(":4016", servemux)
}
