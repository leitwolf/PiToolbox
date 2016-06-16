package module

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"lib"
	"net/url"
	"strconv"
)

// Xuanfeng 旋风离线下载
type Xuanfeng struct {
	YunBase
}

// LoadData 加载列表
// @data {account}
// @return 返回{account,id,list:{id,title,size}} id为空
func (xf *Xuanfeng) LoadData(sender *Sender, data interface{}) {
	data2, ok := data.(map[string]interface{})
	if !ok {
		sender.Err = "error data"
		return
	}
	accountName, _ := data2["account"].(string)
	cc := xf.getCookieContainer(accountName)
	if cc == nil {
		sender.Err = "No account name: " + accountName
		return
	}
	urlStr := "http://lixian.qq.com/handler/lixian/get_lixian_items.php"
	bodyStr := "page=0&limit=200"
	body := []byte(bodyStr)
	req, err := lib.MakeRequest("POST", urlStr, body, cc)
	if err != nil {
		sender.Err = err.Error()
		return
	}
	// ***必须加***
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "http://lixian.qq.com/main.html")
	res, err := lib.FetchHTML(req, cc)
	if err != nil {
		sender.Err = err.Error() + " | xuanfeng:44"
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sender.Err = err.Error() + " | xuanfeng:50"
		return
	}
	// content := string(b)
	var jsonData map[string]interface{}
	err = json.Unmarshal(b, &jsonData)
	if err != nil {
		// test
		ioutil.WriteFile("test/xf1.json", b, 777)
		sender.Err = err.Error() + " | xuanfeng:59"
		return
	}
	// 开始解析
	ret, _ := jsonData["ret"].(float64)
	if ret != 0 {
		msg, _ := jsonData["msg"].(string)
		err = errors.New(msg)
		return
	}
	datas, ok := jsonData["data"].([]interface{})
	if !ok {
		sender.Err = "can not parse data"
		return
	}
	// 取出数据
	l := len(datas)
	resultList := []interface{}{}
	for i := 0; i < l; i++ {
		obj, ok := datas[i].(map[string]interface{})
		if !ok {
			continue
		}
		// 是否已下载完
		status, _ := obj["dl_status"].(float64)
		if status != 12 {
			continue
		}
		// 以hash为id
		id1 := obj["hash"]
		title := obj["file_name"]
		size1, _ := obj["file_size"].(float64)
		size := strconv.Itoa(int(size1))
		size = lib.GetReadableSize(size) + "B"
		resultList = append(resultList, map[string]interface{}{"id": id1, "title": title, "size": size})
	}
	sender.Data = map[string]interface{}{"account": accountName, "id": "", "list": resultList}
}

// Download 下载
// @param data {account:xxx,list:[{id,title},xxx]}
// @return 如果成功返回ok
func (xf *Xuanfeng) Download(sender *Sender, data interface{}) {
	data2, ok := data.(map[string]interface{})
	if !ok {
		sender.Err = "error data"
		return
	}
	accountName, _ := data2["account"].(string)
	list, ok := data2["list"].([]interface{})
	if !ok {
		sender.Err = "convert list fail"
		return
	}
	cc := xf.getCookieContainer(accountName)
	if cc == nil {
		sender.Err = "No accountName name: " + accountName
		return
	}
	success := true
	for i := 0; i < len(list); i++ {
		obj, ok := list[i].(map[string]interface{})
		if !ok {
			sender.Err = "convert list item fail"
			return
		}
		id, _ := obj["id"].(string)
		title, _ := obj["title"].(string)
		downURL, err := xf.getDownURL(id, title, cc)
		if err != nil {
			success = false
			sender.Err = err.Error() + " | xuanfeng:131"
			continue
		}
		header := cc.GetHeaderStr()
		// println("header", header)
		_, err = C.Aria2.AddDownload(downURL, title, header)
		if err != nil {
			success = false
			sender.Err = err.Error() + " | xuanfeng:139"
		}
	}
	if success {
		// 添加成功
		sender.Data = "ok"
	}
}

// getDownURL 获取下载链接
func (xf *Xuanfeng) getDownURL(id string, title string, cc *lib.CookieContainer) (downURL string, err error) {
	urlStr := "http://lixian.qq.com/handler/lixian/get_http_url.php"
	// filename要urlencode
	title = url.QueryEscape(title)
	bodyStr := "hash=" + id + "&filename=" + title + "&browser=other&g_tk=1529673135"
	// println("body", bodyStr)
	body := []byte(bodyStr)
	req, err := lib.MakeRequest("POST", urlStr, body, cc)
	if err != nil {
		println(err.Error() + " xuanfeng 158")
		return
	}
	// ***必须加***
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "http://lixian.qq.com/main.html")
	res, err := lib.FetchHTML(req, cc)
	if err != nil {
		println(err.Error() + " xuanfeng 166")
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		println(err.Error() + " xuanfeng 172")
		return
	}
	var jsonData map[string]interface{}
	// 前三个字节为无效的。。。
	err = json.Unmarshal(b[3:], &jsonData)
	if err != nil {
		// test
		println(b)
		println(err.Error() + " xuanfeng 179")
		return
	}
	// 开始解析
	ret, _ := jsonData["ret"].(float64)
	if ret != 0 {
		msg, _ := jsonData["msg"].(string)
		err = errors.New(msg)
		return
	}
	data, ok := jsonData["data"].(map[string]interface{})
	if !ok {
		err = errors.New("bad response data")
		return
	}
	// 更新cookies字段FTN5K <- com_cookie
	FTN5K, _ := data["com_cookie"].(string)
	cc.UpdateValue("FTN5K", FTN5K)
	downURL, _ = data["com_url"].(string)
	return
}

// NewXuanfeng 新建
func NewXuanfeng() (xf *Xuanfeng) {
	xf = &Xuanfeng{YunBase{accountType: "xuanfeng", accountList: []lib.Account{}}}
	return
}
