package module

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"lib"
	"net/url"
	"regexp"
	"strings"
)

// Yun360 360云盘下载
type Yun360 struct {
	YunBase
}

// LoadData 加载列表
// @data {account,id,path}
// @return 返回{account,id,list:{id,title,size,path,isdir}} isdir时size=""
func (y3 *Yun360) LoadData(sender *Sender, data interface{}) {
	data2, ok := data.(map[string]interface{})
	if !ok {
		sender.Err = "error data"
		return
	}
	accountName, _ := data2["account"].(string)
	id, _ := data2["id"].(string)
	pathStr, _ := data2["path"].(string)
	cc := y3.getCookieContainer(accountName)
	if cc == nil {
		sender.Err = "No account name: " + accountName
		return
	}
	urlStr := "http://c69.yunpan.360.cn/file/list"
	if pathStr == "" {
		pathStr = "/"
	}
	// 路径要urlencode
	pathStr = url.QueryEscape(pathStr)
	bodyStr := "type=2&t=0.01148906020119389&order=asc&field=file_name&path=" + pathStr + "&page=0&page_size=300&ajax=1"
	// println("body", bodyStr)
	body := []byte(bodyStr)
	req, err := lib.MakeRequest("POST", urlStr, body, cc)
	if err != nil {
		sender.Err = err.Error()
		return
	}
	// ***必须加***
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "http://c69.yunpan.360.cn/my")
	res, err := lib.FetchHTML(req, cc)
	if err != nil {
		sender.Err = err.Error()
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sender.Err = err.Error()
		return
	}
	content := string(b)
	// 转成json格式
	jsonStr := y3.transJSON(content)
	// test
	ioutil.WriteFile("test/3601.json", []byte(jsonStr), 777)
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		sender.Err = err.Error()
		return
	}
	// 开始解析
	errno := jsonData["errno"]
	if errno != "0" {
		msg, _ := jsonData["errmsg"].(string)
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
		obj, _ := datas[i].(map[string]interface{})
		id := obj["nid"]
		if id == nil {
			continue
		}
		title := obj["oriName"]
		size := ""
		isdir := false
		isDir, _ := obj["isDir"].(float64)
		if isDir == 1 {
			isdir = true
		} else {
			size = obj["oriSize"].(string)
			size = lib.GetReadableSize(size) + "B"
		}
		pathStr := obj["path"]
		resultList = append(resultList, map[string]interface{}{"id": id, "title": title, "size": size, "path": pathStr, "isdir": isdir})
	}
	sender.Data = map[string]interface{}{"account": accountName, "id": id, "list": resultList}
}

// Download 下载
// @param data {account:xxx,list:[{id,title,path},xxx]}
// @return 如果成功返回ok
func (y3 *Yun360) Download(sender *Sender, data interface{}) {
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
	cc := y3.getCookieContainer(accountName)
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
		pathStr, _ := obj["path"].(string)
		downURL, err := y3.getDownURL(id, pathStr, cc)
		if err != nil {
			success = false
			sender.Err = err.Error()
			continue
		}
		header := cc.GetHeaderStr()
		// println("header", header)
		_, err = C.Aria2.AddDownload(downURL, title, header)
		if err != nil {
			success = false
			sender.Err = err.Error()
		}
	}
	if success {
		// 添加成功
		sender.Data = "ok"
	}
}

// getDownURL 获取下载链接
func (y3 *Yun360) getDownURL(id string, pathStr string, cc *lib.CookieContainer) (downURL string, err error) {
	urlStr := "http://c69.yunpan.360.cn/file/download"
	// 路径要urlencode
	pathStr = url.QueryEscape(pathStr)
	bodyStr := "nid=" + id + "&fname=" + pathStr + "&ajax=1"
	// println("body", bodyStr)
	body := []byte(bodyStr)
	req, err := lib.MakeRequest("POST", urlStr, body, cc)
	if err != nil {
		return
	}
	// ***必须加***
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "http://c69.yunpan.360.cn/my")
	res, err := lib.FetchHTML(req, cc)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	// println(string(b))
	var jsonData map[string]interface{}
	err = json.Unmarshal(b, &jsonData)
	if err != nil {
		return
	}
	// 开始解析
	errno, _ := jsonData["errno"].(float64)
	// println("err", errno)
	if errno != 0 {
		msg, _ := jsonData["errmsg"].(string)
		err = errors.New(msg)
		return
	}
	data, ok := jsonData["data"].(map[string]interface{})
	if !ok {
		err = errors.New("bad response data")
		return
	}
	downURL, _ = data["download_url"].(string)
	return
}

// transJSON 把不规范的字符串转换成json
func (y3 *Yun360) transJSON(str string) (jsonStr string) {
	str = strings.Replace(str, "'", "\"", -1)
	// 以,分隔
	list := strings.Split(str, ",")
	for i := 0; i < len(list); i++ {
		list[i] = y3.handleDot(list[i], 0)
	}
	// 合并成新的
	jsonStr = strings.Join(list, ",")
	return
}

// handleDot 给冒号前的key加上双引号，只处理一次冒号
// @param start 开始查找的位置
func (y3 *Yun360) handleDot(str string, start int) (newStr string) {
	if start >= len(str) {
		newStr = str
		return
	}
	index2 := strings.Index(str[start:], ":")
	if index2 == -1 {
		newStr = str
		return
	}
	index2 += start
	// 第一个可用字符（字母，数字,_）
	index1 := -1
	for j := start; j < len(str); j++ {
		c := str[j]
		m, _ := regexp.Match(`[A-Za-z]`, []byte{c})
		if m {
			index1 = j
			break
		}
	}
	if index1 == -1 {
		newStr = str
		return
	}
	// 键名
	name := strings.TrimSpace(str[index1:index2])
	if name != "http" && name != "https" {
		str = strings.Replace(str, name, "\""+name+"\"", 1)
	}
	// 再处理一次，可能有两个冒号，要处理到后面没有冒号为止
	newStr = y3.handleDot(str, index2+3)
	return
}

// NewYun360 新建
func NewYun360() (yun360 *Yun360) {
	yun360 = &Yun360{YunBase{accountType: "360", accountList: []lib.Account{}}}
	return
}
