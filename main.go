package main

import (
	"log"
	"module"
	"net/http"
	"strconv"
	// "github.com/VividCortex/godaemon"
)

// 监听的端口
var port = 5000

// startServer 开启web服务
func startServer(port int) {
	addr := ":" + strconv.Itoa(port)
	log.Println("Listening on port: " + strconv.Itoa(port))
	// 静态文件
	http.Handle("/", http.FileServer(http.Dir("html")))
	// http.Handle("/", http.FileServer(FS(false)))
	// 客户端处理操作
	http.HandleFunc("/action", module.ReqHandler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	// godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	module.Init()
	startServer(port)
}
