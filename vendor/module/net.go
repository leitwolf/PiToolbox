package module

//
// 前后端连接处理
//
import (
	"encoding/json"
	"io/ioutil"
	"lib"
	"log"
	"net/http"
	"strconv"
)

// Receiver 客户端发来的请求信息
type Receiver struct {
	// 该次通信所需模块
	Module string `json:"module"`
	// 所执行的操作
	Action string `json:"action"`
	// 数据内容
	Data interface{} `json:"data"`
}

// Sender 发给客户端的响应信息
type Sender struct {
	// 该次通信所需模块
	Module string `json:"module"`
	// 所执行的操作
	Action string `json:"action"`
	// 数据内容
	Data interface{} `json:"data"`
	// 错误信息，没有则为空
	Err string `json:"err"`
}

// reqHandler 客户端请求处理
func reqHandler(res http.ResponseWriter, req *http.Request) {
	// defer req.Body.Close()
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.Write([]byte("bad request body"))
		return
	}
	// log.Println("receive: ", string(content))
	var receiver Receiver
	err = json.Unmarshal(content, &receiver)
	if err != nil {
		res.Write([]byte("bad request content"))
		return
	}
	sender := &Sender{}
	sender.Module = receiver.Module
	sender.Action = receiver.Action
	sender.Err = ""
	m := receiver.Module
	a := receiver.Action

	actionDispatch(m, a, receiver.Data, sender)

	b, err := json.Marshal(*sender)
	if err != nil {
		res.Write([]byte("bad response data"))
		return
	}
	// log.Println("send:", string(b))
	res.Header().Add("Content-type:", "application/json")
	res.Write(b)
}

// actionDispatch 操作分发
func actionDispatch(m string, a string, data interface{}, sender *Sender) {
	if m == "aria2" {
		if a == "getConfig" {
			C.Aria2.GetConfig(sender)
		} else if a == "saveConfig" {
			C.Aria2.SaveConfig(sender, data)
		} else if a == "getVersion" {
			C.Aria2.GetVersion(sender)
		} else if a == "getStat" {
			C.Aria2.GetStat(sender)
		} else if a == "start" {
			C.Aria2.Start(sender, data)
		} else if a == "pause" {
			C.Aria2.Pause(sender, data)
		} else if a == "remove" {
			C.Aria2.Remove(sender, data)
		} else if a == "startAll" {
			C.Aria2.StartAll(sender)
		} else if a == "pauseAll" {
			C.Aria2.PauseAll(sender)
		} else if a == "removeStoped" {
			C.Aria2.RemoveStoped(sender, data)
		} else if a == "removeAllStoped" {
			C.Aria2.RemoveAllStoped(sender)
		}
	} else if m == "xunlei" {
		if a == "getAccountList" {
			C.Xunlei.GetAccountList(sender)
		} else if a == "loadData" {
			C.Xunlei.LoadData(sender, data)
		} else if a == "download" {
			C.Xunlei.Download(sender, data)
		}
	} else if m == "yun360" {
		if a == "getAccountList" {
			C.Yun360.GetAccountList(sender)
		} else if a == "loadData" {
			C.Yun360.LoadData(sender, data)
		} else if a == "download" {
			C.Yun360.Download(sender, data)
		}
	} else if m == "xuanfeng" {
		if a == "getAccountList" {
			C.Xuanfeng.GetAccountList(sender)
		} else if a == "loadData" {
			C.Xuanfeng.LoadData(sender, data)
		} else if a == "download" {
			C.Xuanfeng.Download(sender, data)
		}
	} else if m == "cookies" {
		if a == "save" {
			// 保存cookies
			saveCookies(sender, data)
		}
	}
}

// 保存cookies
func saveCookies(sender *Sender, data interface{}) {
	data2 := data.(map[string]interface{})
	filename, ok := data2["filename"].(string)
	if !ok {
		sender.Err = "bad cookies filename"
		return
	}
	content, ok := data2["content"].(string)
	if !ok {
		sender.Err = "bad cookies content"
		return
	}
	err := lib.WriteFile("config/"+filename, []byte(content))
	if err != nil {
		sender.Err = err.Error()
		return
	}
	// 当前页面
	page, _ := data2["page"].(string)
	sender.Data = page
}

// StartServer 开启web服务
func StartServer(port int) {
	addr := ":" + strconv.Itoa(port)
	log.Println("Listening on port: " + strconv.Itoa(port))
	// 静态文件
	http.Handle("/", http.FileServer(http.Dir("./html")))
	// 客户端处理操作
	http.HandleFunc("/action", reqHandler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
