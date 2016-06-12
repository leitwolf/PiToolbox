package main

import "module"

// 监听的端口
var port = 5000

// func test() {
// 	cc, err := lib.NewCookieContainer("config/cookies_360.json")
// 	if err != nil {
// 		println(err.Error())
// 		return
// 	}
// 	urlStr := "https://c69.yunpan.360.cn/file/list"
// 	str := "type=2&t=0.7683906133348188&order=asc&field=file_name&path=%2Faa%2F&page=0&page_size=300&ajax=1"
// 	body := []byte(str)
// 	req, err := lib.MakeRequest("POST", urlStr, body, cc)
// 	if err != nil {
// 		println(err.Error())
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	// req.Header.Add("Content-Length", strconv.Itoa(len(str)))
// 	req.Header.Add("Referer", "https://c69.yunpan.360.cn/my")
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
// }

func main() {
	module.Init()
	module.StartServer(port)
	// test()
}
