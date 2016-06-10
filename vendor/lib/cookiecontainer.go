package lib

//
// 一个cookie容器
//
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// CookieContainer cookies管理
type CookieContainer struct {
	cookies []http.Cookie
}

// Load 加载
func (c *CookieContainer) Load(path string) (err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &c.cookies)
	return
}

// GetValueByName 查找指定Name的cookie的值
func (c *CookieContainer) GetValueByName(name string) (value string) {
	for _, cookie := range c.cookies {
		if cookie.Name == name {
			value = cookie.Value
			return
		}
	}
	return
}

// AddToReqeust 添加cookie给request
func (c *CookieContainer) AddToReqeust(request *http.Request) {
	for _, cookie := range c.cookies {
		request.AddCookie(&cookie)
	}
}

// NewCookieContainer 新建一个CookieJar
func NewCookieContainer(path string) (c *CookieContainer, err error) {
	c = &CookieContainer{cookies: make([]http.Cookie, 0)}
	err = c.Load(path)
	return
}
