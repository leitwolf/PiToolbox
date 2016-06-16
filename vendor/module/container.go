package module

// Container 模块容器
type Container struct {
	Aria2    *Aria2
	Xunlei   *Xunlei
	Yun360   *Yun360
	Xuanfeng *Xuanfeng
}

// C 容器实例
var C Container

// Init 初始化各个模块
func Init() {
	C = Container{}
	C.Aria2 = NewAria2()
	C.Xunlei = NewXunlei()
	C.Yun360 = NewYun360()
	C.Xuanfeng = NewXuanfeng()
}
