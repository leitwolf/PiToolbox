package module

//
// Xunlei 迅雷离线下载
//
import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"lib"
	"strconv"
)

// Xunlei 迅雷离线下载
type Xunlei struct {
	// 账户列表
	accountList []lib.Account
}

// GetAccountList 获取账户列表
// 返回 [Account.Name] 列表
func (xl *Xunlei) GetAccountList(sender *Sender) {
	list, err := lib.LoadAccountList("xunlei")
	if err != nil {
		sender.Err = err.Error()
		return
	}
	xl.accountList = list
	var names []string
	for i := 0; i < len(list); i++ {
		names = append(names, list[i].Name)
	}
	sender.Data = names
}

// LoadData 加载数据
// 返回 {id:id,account:account,list:[{id,title,size,url}]}
func (xl *Xunlei) LoadData(sender *Sender, data interface{}) {
	data2, ok := data.(map[string]interface{})
	if !ok {
		sender.Err = "error data"
		return
	}
	accountName, _ := data2["account"].(string)
	id, _ := data2["id"].(string)
	// println("account", accountName, taskID)
	cc := xl.getCookieContainer(accountName)
	if cc == nil {
		sender.Err = "No account name: " + accountName
		return
	}
	if id == "" {
		// 获取主页面
		resultList, err := xl.getMainList(cc)
		if err != nil {
			sender.Err = err.Error()
			return
		}
		sender.Data = map[string]interface{}{"id": id, "account": accountName, "list": resultList}
	} else {
		// 获取bt
		resultList, err := xl.getBt(cc, id)
		if err != nil {
			sender.Err = err.Error()
			return
		}
		sender.Data = map[string]interface{}{"id": id, "account": accountName, "list": resultList}
	}
}

// Download 下载
// {account:xxx,list:[{title,url},xxx]}
func (xl *Xunlei) Download(sender *Sender, data interface{}) {
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
	cc := xl.getCookieContainer(accountName)
	if cc == nil {
		sender.Err = "No accountName name: " + accountName
		return
	}
	// 要加上头
	header := "Cookie: gdriveid=" + cc.GetValueByName("gdriveid")
	success := true
	for i := 0; i < len(list); i++ {
		obj, ok := list[i].(map[string]interface{})
		if !ok {
			sender.Err = "convert list item fail"
			return
		}
		title, _ := obj["title"].(string)
		url, _ := obj["url"].(string)
		_, err := C.Aria2.AddDownload(url, title, header)
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

// getMainList 获取主页面列表信息
// 返回 {id,title,size,url}，url为空则是bt文件夹
func (xl *Xunlei) getMainList(cc *lib.CookieContainer) (resultList []interface{}, err error) {
	// 请求页面数据
	url := "http://dynamic.cloud.vip.xunlei.com/user_task?userid=" + cc.GetValueByName("userid")
	res, err := lib.GetHTMLWithCookies(url, cc)
	if err != nil {
		// println("main aa", err.Error())
		return
	}
	defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	println("body", err.Error())
	// }
	// ioutil.WriteFile("res.html", body, 0777)
	q, err := lib.NewQueryFromReader(res.Body)
	if err != nil {
		// println("main bb", err.Error())
		return
	}
	// 获取列表
	// <input id="dflag175673410781696" name="dflag" type="hidden" value="0" />
	// <input id="dl_url175673410781696" type="hidden" value="http://gdl.lixian.vip.xunlei.com/download?xxx" />
	// <input id="taskname175673410781696" type="hidden" value="[一刀倾城].Blade.Of.Fury.1993.DVDRip.x264.AC3.2Audios-CMCT.mkv" />
	// <input id="ysfilesize175673410781696" type="hidden" value="4676679449" />
	// <input id="d_status1381822237122304" type="hidden" value="2" />
	// [{taskID,taskName,url}]
	list := q.GetNodesByName("dflag")
	resultList = []interface{}{}
	for i := 0; i < len(list); i++ {
		str := q.GetNodeAttr(list[i], "id")
		id := str[5:]
		// println("taskid", id)
		// 查看下载状态是否已完成
		node := q.GetNodeByID("d_status" + id)
		status := q.GetNodeAttr(node, "value")
		if status != "2" {
			continue
		}
		node = q.GetNodeByID("taskname" + id)
		title := q.GetNodeAttr(node, "value")
		node = q.GetNodeByID("ysfilesize" + id)
		size := q.GetNodeAttr(node, "value")
		size = lib.GetReadableSize(size) + "B"
		node = q.GetNodeByID("dl_url" + id)
		downurl := q.GetNodeAttr(node, "value")
		resultList = append(resultList, map[string]string{"id": id, "title": title, "size": size, "url": downurl})
	}
	return
}

// getBt 获取bt列表
func (xl *Xunlei) getBt(cc *lib.CookieContainer, taskID string) (resultList []interface{}, err error) {
	// 请求页面数据
	userid := cc.GetValueByName("userid")
	callback := "fill_bt_list"
	url := "http://dynamic.cloud.vip.xunlei.com/interface/fill_bt_list?tid=" + taskID + "&g_net=1&p=1&uid=" + userid + "&callback=" + callback
	// println("url", url)
	res, err := lib.GetHTMLWithCookies(url, cc)
	if err != nil {
		// println("bt bb", err.Error())
		return
	}
	defer res.Body.Close()
	// 返回jsonp格式：
	// fill_bt_list({
	// 	"Result": {
	// 		"Record": [
	// 			{
	// 				"id": 0,
	// 				"title": "[\u5175\u4e34\u57ce\u4e0b(\u56fd\u82f1\u53cc\u8bed)].Enemy.At.The.Gates.2001.BluRay.720p.x264.AC3-CMCT.mkv",
	// 				"download_status": "2",
	// 				"percent": 100,
	// 				"taskid": "1381832231362368",
	// 				"downurl": "xxx",
	//  			"filesize": "3061626294"
	// 			}
	// 		]
	// 	}
	// })
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// println("bt bb", err.Error())
		return
	}
	str := string(b)
	// fmt.Println(str)
	str = str[len(callback)+1 : len(str)-1]
	result := map[string]interface{}{}
	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		// println("bt cc", err.Error())
		return
	}
	result, ok := result["Result"].(map[string]interface{})
	if !ok {
		err = errors.New("convert Result fail")
		return
	}
	list, ok := result["Record"].([]interface{})
	if !ok {
		err = errors.New("convert tasks fail")
		return
	}
	resultList = []interface{}{}
	for i := 0; i < len(list); i++ {
		obj, ok := list[i].(map[string]interface{})
		if !ok {
			err = errors.New("convert obj fail")
			return
		}
		// 查看下载状态是否已完成
		status, _ := obj["download_status"].(string)
		if status != "2" {
			continue
		}
		id1, _ := obj["id"].(int)
		id := strconv.Itoa(id1)
		taskid, _ := obj["taskid"].(string)
		id = taskid + id
		title, _ := obj["title"].(string)
		size, _ := obj["filesize"].(string)
		size = lib.GetReadableSize(size) + "B"
		downurl, _ := obj["downurl"].(string)
		resultList = append(resultList, map[string]string{"id": id, "title": title, "size": size, "url": downurl})
	}
	return
}

// getCookieContainer 获取指定的cookie
func (xl *Xunlei) getCookieContainer(accountName string) (cc *lib.CookieContainer) {
	for i := 0; i < len(xl.accountList); i++ {
		if xl.accountList[i].Name == accountName {
			cc = xl.accountList[i].CookieContainer
			return
		}
	}
	return
}

// NewXunlei 新建
func NewXunlei() (xunlei *Xunlei) {
	xunlei = &Xunlei{accountList: []lib.Account{}}
	return
}
