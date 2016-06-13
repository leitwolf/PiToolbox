package module

//
// 处理与Aria2的通信
//
import (
	"encoding/json"
	"io/ioutil"
	"lib"
	"log"
	"strconv"
	"strings"
	"time"
)

// Aria2Task 下载任务
type Aria2Task struct {
	GID string `json:"gid"`
	// 文件名
	Filename string `json:"filename"`
	// 状态 active waiting paused error complete removed
	Status string `json:"status"`
	// 总大小
	Size int `json:"size"`
	// 已完成大小
	CompletedLength int `json:"completedLength"`
	// 进度0-100
	Progress float32 `json:"progress"`
	Speed    int     `json:"speed"`
	// 与服务器连接数
	Connections string `json:"connections"`
}

// Aria2Stat aria2状态信息
type Aria2Stat struct {
	Speed string `json:"speed"`
	// 活动的下载列表
	ActiveTasks []Aria2Task `json:"activeTasks"`
	// 等待中的下载列表
	WaitingTasks []Aria2Task `json:"waitingTasks"`
	// 已停止的下载列表
	StopedTasks []Aria2Task `json:"stopedTasks"`
}

// Aria2 与aria2相关的操作
type Aria2 struct {
	url     string
	version string
}

// ===start 交互相关==

// GetVersion 获取版本号
func (a *Aria2) GetVersion(sender *Sender) {
	sender.Data = a.getVersion()
}

// GetStat 获取实时信息
func (a *Aria2) GetStat(sender *Sender) {
	stat, err := a.getStat()
	if err != nil {
		sender.Err = err.Error()
	} else {
		sender.Data = *stat
	}
}

// Start 开始某些任务
func (a *Aria2) Start(sender *Sender, data interface{}) {
	gids, ok := data.([]interface{})
	if !ok {
		sender.Err = "start invalid gids"
		return
	}
	for i := 0; i < len(gids); i++ {
		params := []interface{}{}
		params = append(params, gids[i])
		req := a.getJSONRPCRequest()
		req.Method = "aria2.unpause"
		req.Params = params
		lib.CallJSONRPC(a.url, req)
	}
}

// Pause 暂停某些任务
func (a *Aria2) Pause(sender *Sender, data interface{}) {
	gids, ok := data.([]interface{})
	if !ok {
		sender.Err = "pause invalid gids"
		return
	}
	for i := 0; i < len(gids); i++ {
		params := []interface{}{}
		params = append(params, gids[i])
		req := a.getJSONRPCRequest()
		req.Method = "aria2.pause"
		req.Params = params
		lib.CallJSONRPC(a.url, req)
	}
}

// Remove 删除某些任务
func (a *Aria2) Remove(sender *Sender, data interface{}) {
	gids, ok := data.([]interface{})
	if !ok {
		sender.Err = "remove invalid gids"
		return
	}
	for i := 0; i < len(gids); i++ {
		params := []interface{}{}
		params = append(params, gids[i])
		req := a.getJSONRPCRequest()
		req.Method = "aria2.forceRemove"
		req.Params = params
		lib.CallJSONRPC(a.url, req)
		req.Method = "aria2.removeDownloadResult"
		lib.CallJSONRPC(a.url, req)
	}
}

// StartAll 开始所有任务
func (a *Aria2) StartAll(sender *Sender) {
	params := []interface{}{}
	req := a.getJSONRPCRequest()
	req.Method = "aria2.unpauseAll"
	req.Params = params
	lib.CallJSONRPC(a.url, req)
}

// PauseAll 暂停所有任务
func (a *Aria2) PauseAll(sender *Sender) {
	params := []interface{}{}
	req := a.getJSONRPCRequest()
	req.Method = "aria2.pauseAll"
	req.Params = params
	lib.CallJSONRPC(a.url, req)
}

// RemoveStoped 删除已停止的某些任务
func (a *Aria2) RemoveStoped(sender *Sender, data interface{}) {
	gids, ok := data.([]interface{})
	if !ok {
		sender.Err = "removeStoped invalid gids"
		return
	}
	for i := 0; i < len(gids); i++ {
		params := []interface{}{}
		params = append(params, gids[i])
		req := a.getJSONRPCRequest()
		req.Method = "aria2.removeDownloadResult"
		req.Params = params
		lib.CallJSONRPC(a.url, req)
	}
}

// RemoveAllStoped 暂停所有已停止的任务
func (a *Aria2) RemoveAllStoped(sender *Sender) {
	params := []interface{}{}
	req := a.getJSONRPCRequest()
	req.Method = "aria2.purgeDownloadResult"
	req.Params = params
	lib.CallJSONRPC(a.url, req)
}

// ===end 交互相关==

// AddDownload 添加下载
// downURL 下载地址
// filename 下载文件名
// header 需要设置的头部信息
func (a *Aria2) AddDownload(downURL string, filename string, header string) (gid string, err error) {
	params := []interface{}{}
	params = append(params, []string{downURL})
	options := map[string]string{}
	options["out"] = filename
	if header != "" {
		options["header"] = header
	}
	params = append(params, options)

	req := a.getJSONRPCRequest()
	req.Method = "aria2.addUri"
	req.Params = params
	res, err := lib.CallJSONRPC(a.url, req)
	if err != nil {
		return
	}
	str, ok := res.Result.(string)
	if ok {
		gid = str
	}
	return
}

// GetVersion 获取版本
func (a *Aria2) getVersion() (version string) {
	if a.version != "" {
		version = a.version
		return
	}
	req := a.getJSONRPCRequest()
	req.Method = "aria2.getVersion"
	res, err := lib.CallJSONRPC(a.url, req)
	if err != nil || res.Error != nil {
		return
	}
	var m = res.Result.(map[string]interface{})
	version = m["version"].(string)
	// 缓存起来
	a.version = version
	return
}

// getStat 获取Aria2当前状态，包括下载速度，各任务情况
// 使用system.multicall返回多个查询结果
// 每个返回的结果都在原来的基础上加上了[]
func (a *Aria2) getStat() (stat *Aria2Stat, err error) {
	stat = &Aria2Stat{}
	methodList := []interface{}{}
	// 查询速度
	obj := make(map[string]interface{}, 0)
	obj["methodName"] = "aria2.getGlobalStat"
	methodList = append(methodList, obj)
	// 查询活动任务
	obj = make(map[string]interface{}, 0)
	obj["methodName"] = "aria2.tellActive"
	params := []interface{}{}
	params = append(params, a.getTaskKeys())
	obj["params"] = params
	methodList = append(methodList, obj)
	// 查询等待中的任务
	obj = make(map[string]interface{}, 0)
	obj["methodName"] = "aria2.tellWaiting"
	params = []interface{}{}
	params = append(params, 0)
	params = append(params, 1000)
	params = append(params, a.getTaskKeys())
	obj["params"] = params
	methodList = append(methodList, obj)
	// 查询已停止的任务
	obj = make(map[string]interface{}, 0)
	obj["methodName"] = "aria2.tellStopped"
	params = []interface{}{}
	params = append(params, 0)
	params = append(params, 1000)
	params = append(params, a.getTaskKeys())
	obj["params"] = params
	methodList = append(methodList, obj)

	req := a.getJSONRPCRequest()
	req.Method = "system.multicall"
	req.Params = []interface{}{methodList}
	res, err := lib.CallJSONRPC(a.url, req)
	if err != nil {
		return
	}
	list := res.Result.([]interface{})
	// 解析速度
	l := list[0].([]interface{})
	m := l[0].(map[string]interface{})
	stat.Speed = lib.GetReadableSize(m["downloadSpeed"].(string)) + "B/s"
	// 解析各任务
	l = list[1].([]interface{})
	stat.ActiveTasks = a.analyseTasks(l[0])
	l = list[2].([]interface{})
	stat.WaitingTasks = a.analyseTasks(l[0])
	l = list[3].([]interface{})
	stat.StopedTasks = a.analyseTasks(l[0])
	return
}

// analyseTasks 分析返回的任务信息，处理一些信息
// @return [{}]
func (a *Aria2) analyseTasks(data interface{}) (tasks []Aria2Task) {
	// fmt.Printf("%v\n", data)
	list := data.([]interface{})
	for _, item := range list {
		m := item.(map[string]interface{})
		task := Aria2Task{}
		task.GID = m["gid"].(string)
		task.Status = m["status"].(string)
		totalLength := m["totalLength"].(string)
		task.Size, _ = strconv.Atoi(totalLength)
		completedLength := m["completedLength"].(string)
		task.CompletedLength, _ = strconv.Atoi(completedLength)
		// 完成百分比
		if task.Size > 0 {
			p := float32(task.CompletedLength) * 100.0 / float32(task.Size)
			progress, _ := strconv.ParseFloat(strconv.FormatFloat(float64(p), 'f', 2, 32), 32)
			task.Progress = float32(progress)
		}
		speed := m["downloadSpeed"].(string)
		task.Speed, _ = strconv.Atoi(speed)
		task.Connections, _ = m["connections"].(string)
		// 文件名从files中取，只取第一个，去掉路径
		files := m["files"].([]interface{})
		file := files[0].(map[string]interface{})
		path1 := file["path"].(string)
		index := strings.LastIndex(path1, "/")
		filename := path1[(index + 1):]
		task.Filename = filename
		tasks = append(tasks, task)
	}
	return
}

// getTaskKeys 查询一个任务所需的字段
func (a *Aria2) getTaskKeys() (keys []string) {
	keys = append(keys, "gid")
	// 状态 active waiting paused error complete removed
	keys = append(keys, "status")
	// 文件大小 byte
	keys = append(keys, "totalLength")
	// 已完成大小 byte
	keys = append(keys, "completedLength")
	// 下载速度 bytes/sec
	keys = append(keys, "downloadSpeed")
	// 与服务器连接数
	keys = append(keys, "connections")
	// 包含的文件列表
	keys = append(keys, "files")
	return
}

// getJSONRPCRequest 获取一个请求实体
func (a *Aria2) getJSONRPCRequest() (req *lib.JSONRPCRequest) {
	req = &lib.JSONRPCRequest{}
	req.Jsonrpc = "2.0"
	req.ID = strconv.Itoa(int(time.Now().UnixNano()))
	req.Params = []interface{}{}
	return
}

// loadAria2Config 加载配置文件
func loadAria2Config() (url string) {
	b, err := ioutil.ReadFile("config/aria2.json")
	if err != nil {
		return
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return
	}
	url, _ = m["url"].(string)
	return
}

// NewAria2 新建
func NewAria2() (aria2 *Aria2) {
	url := loadAria2Config()
	log.Println("aria2 url:", url)
	if url == "" {
		// 默认的aria2连接地址
		url = "http://localhost:6800/jsonrpc"
	}
	aria2 = &Aria2{url: url}
	return
}
