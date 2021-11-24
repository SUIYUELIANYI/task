package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type a struct {
	U        string `json:"username"`
	Password string `json:"password"`
}

func main() {
	requestUrl := "http://pass.muxi-tech.xyz/auth/api/signin"
	/*data := url.Values{}
	data.Set("username", "用户名")
	data.Set("password", "密码")*/
	var user a
	user = a{"方瑜诚", "UWlhbnJlbjQwMTY="}
	buf, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		panic(err)
	}
	payload := strings.NewReader(string(buf))
	req, err := http.NewRequest("POST", requestUrl, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", "54")
	req.Header.Add("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Add("Host", "pass.muxi-tech.xyz")
	req.Header.Add("Origin", "http://pass.muxi-tech.xyz")
	req.Header.Add("Referer", "http://pass.muxi-tech.xyz/intro")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
