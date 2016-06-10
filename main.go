package main

import "module"

// 监听的端口
var port = 5000

func main() {
	module.Init()
	module.StartServer(port)
}
