package main

import "module"

// 监听的端口
var port = 5000

// func test() {
// 	cc, err := lib.NewCookieContainer("config/cookies_xuanfeng.json")
// 	if err != nil {
// 		println(err.Error())
// 		return
// 	}
// 	urlStr := "http://lixian.qq.com/handler/lixian/get_http_url.php"
// 	str := "hash=2C0366287FD3B02F51603A015E56BA91195E5278&filename=%5B%E5%B1%B1%E6%B2%B3%E6%95%85%E4%BA%BA%5D.Mountains.May.Depart.2015.BluRay.720p.x264.AC3-CMCT.mkv&browser=other&g_tk=111"
// 	body := []byte(str)
// 	req, err := lib.MakeRequest("POST", urlStr, body, cc)
// 	if err != nil {
// 		println(err.Error())
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	// req.Header.Set("Content-Length", strconv.Itoa(len(str)))
// 	req.Header.Set("Referer", "http://lixian.qq.com/main.html")
// 	res, err := lib.FetchHTML(req, cc)
// 	if err != nil {
// 		println(err.Error())
// 		return
// 	}
// 	defer res.Body.Close()
// 	b, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		println(err.Error())
// 		return
// 	}
// 	content := string(b)
// 	println(content)
// }
// func test2() {
// 	lib.WriteFile("aa/bb/cc.txt", []byte("aaa"))
// }

func main() {
	module.Init()
	module.StartServer(port)
	// test2()
}
