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
	"strings"
	"time"
)

// Xunlei 迅雷离线下载
type Xunlei struct {
	YunBase
}

// LoadData 加载数据
// @param data {account,id}
// @return 返回 {account:account,id:id,list:[{id,title,size,url}]}
func (xl *Xunlei) LoadData(sender *Sender, data interface{}) {
	data2, ok := data.(map[string]interface{})
	if !ok {
		sender.Err = "error data"
		return
	}
	accountName, _ := data2["account"].(string)
	id, _ := data2["id"].(string)
	cc := xl.getCookieContainer(accountName)
	if cc == nil {
		sender.Err = "No account name: " + accountName
		return
	}
	if id == "" {
		// 获取主页面
		resultList, err := xl.getMainList(cc)
		if err != nil {
			sender.Err = err.Error() + " | xunlei:39"
			return
		}
		sender.Data = map[string]interface{}{"account": accountName, "id": id, "list": resultList}
	} else {
		// 获取bt
		resultList, err := xl.getBt(cc, id)
		if err != nil {
			sender.Err = err.Error() + " | xunlei:47"
			return
		}
		sender.Data = map[string]interface{}{"account": accountName, "id": id, "list": resultList}
	}
}

// Download 下载
// @param data {account:xxx,list:[{title,url},xxx]}
// @return 如果成功返回ok
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
		urlStr, _ := obj["url"].(string)
		_, err := C.Aria2.AddDownload(urlStr, title, header)
		if err != nil {
			success = false
			sender.Err = err.Error() + " | xunlei:88"
		}
	}
	if success {
		// 添加成功
		sender.Data = "ok"
	}
}

// getMainList 获取主页面列表信息
// @return 返回 [{id,title,size,url,isdir}]
func (xl *Xunlei) getMainList(cc *lib.CookieContainer) (resultList []interface{}, err error) {
	// 需要随机
	ran := strconv.FormatInt(time.Now().UnixNano(), 10)
	callback := "jsonp" + ran
	urlStr := "http://dynamic.cloud.vip.xunlei.com/interface/showtask_unfresh?callback=" + callback + "&type_id=4&page=1&tasknum=300&p=1&interfrom=task"
	res, err := lib.GetHTML(urlStr, cc)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	str := string(b)
	// 获取jsonp中的内容
	str = str[len(callback)+1 : len(str)-1]
	result := map[string]interface{}{}
	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		return
	}
	// 检测返回代码
	rtcode, ok := result["rtcode"].(float64)
	if !ok {
		err = errors.New("bad response data")
		return
	}
	if rtcode != 0 {
		err = errors.New("wrong rtcode " + strconv.Itoa(int(rtcode)))
		return
	}
	info, ok := result["info"].(map[string]interface{})
	if !ok {
		err = errors.New("bad response data['info']")
		return
	}
	tasks, ok := info["tasks"].([]interface{})
	if !ok {
		err = errors.New("bad response data['info']['tasks']")
		return
	}
	resultList = []interface{}{}
	for i := 0; i < len(tasks); i++ {
		task, ok := tasks[i].(map[string]interface{})
		if !ok {
			continue
		}
		id := task["id"]
		title := task["taskname"]
		isdir := false
		urlStr, _ := task["lixian_url"].(string)
		if strings.HasPrefix(urlStr, "bt:") {
			// bt文件夹
			urlStr = ""
			isdir = true
		}
		status, _ := task["download_status"].(string)
		if status != "2" && !isdir {
			// 没有下载完成且非文件夹的，不计算
			continue
		}
		size, _ := task["file_size"].(string)
		size = lib.GetReadableSize(size) + "B"
		resultList = append(resultList, map[string]interface{}{"id": id, "title": title, "size": size, "url": urlStr, "isdir": isdir})
	}
	return
}

// getBt 获取bt列表
// @return 返回 [{id,title,size,url,isdir}]，url为空则是bt文件夹
func (xl *Xunlei) getBt(cc *lib.CookieContainer, taskID string) (resultList []interface{}, err error) {
	// 请求页面数据
	userid := cc.GetValueByName("userid")
	callback := "fill_bt_list"
	urlStr := "http://dynamic.cloud.vip.xunlei.com/interface/fill_bt_list?tid=" + taskID + "&g_net=1&p=1&uid=" + userid + "&callback=" + callback
	res, err := lib.GetHTML(urlStr, cc)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	str := string(b)
	// 获取jsonp中的内容
	str = str[len(callback)+1 : len(str)-1]
	result := map[string]interface{}{}
	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		return
	}
	result, ok := result["Result"].(map[string]interface{})
	if !ok {
		err = errors.New("bad response data['Result']")
		return
	}
	list, ok := result["Record"].([]interface{})
	if !ok {
		err = errors.New("bad response data['Result']['Record']")
		return
	}
	resultList = []interface{}{}
	for i := 0; i < len(list); i++ {
		obj, ok := list[i].(map[string]interface{})
		if !ok {
			err = errors.New("bad response data['Result']['Record'][...]")
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
		title := obj["title"]
		size, _ := obj["filesize"].(string)
		size = lib.GetReadableSize(size) + "B"
		urlStr := obj["downurl"]
		isdir := false
		if urlStr == "" {
			isdir = true
		}
		resultList = append(resultList, map[string]interface{}{"id": id, "title": title, "size": size, "url": urlStr, "isdir": isdir})
	}
	return
}

// NewXunlei 新建
func NewXunlei() (xunlei *Xunlei) {
	xunlei = &Xunlei{YunBase{accountType: "xunlei", accountList: []lib.Account{}}}
	xunlei.initAccountList()
	return
}
