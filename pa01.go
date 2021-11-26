package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	requestUrl := "http://pass.muxi-tech.xyz/intro#/login"
	// 发送Get请求
	rsp, err := http.Get(requestUrl)
	if err != nil {
		log.Println(err.Error())
		return
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	content := string(body)
	defer rsp.Body.Close()

	buf := content
	//fmt.Println(buf)
	//reg := regexp.MustCompile()//应该是改这里来筛选
	reg := regexp.MustCompile(`<title>(?s:(.*?))</title>`)
	//reg1 := regexp.MustCompile(`<div class="span1">(?s:(.*?))</div>`) 这里是想爬点别的，无奈看不懂正则怎么用
	if reg == nil {
		fmt.Println("MustCompile err")
		return
	}
	// 提取关键信息
	result := reg.FindAllStringSubmatch(buf, -1)
	//result1 := reg1.FindAllStringSubmatch(buf, -1)
	// 过滤<> </>
	for i, text := range result {
		fmt.Println("text", i+1, "=", text[1])//这里如果text[1]改为text，则会读成text 1 = [<title>木犀内网门户</title> 木犀内网门户]
	}

	/*for i, text := range result1 {
		fmt.Println("text", i+1, "=", text[1])
	}*/
}